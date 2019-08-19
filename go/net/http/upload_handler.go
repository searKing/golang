package http

import (
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/pkg/errors"
	"github.com/searKing/golang/thirdparty/github.com/sirupsen/logrus"
	"github.com/searKing/tusd"
	"github.com/searKing/tusd/filestore"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	reExtractFileID  = regexp.MustCompile(`([^/]+)\/?$`)
	reForwardedHost  = regexp.MustCompile(`host=([^,]+)`)
	reForwardedProto = regexp.MustCompile(`proto=(https?)`)
	reMimeType       = regexp.MustCompile(`^[a-z]+\/[a-z\-\+0-9]+$`)
)

// UploadHandler exposes methods to handle requests as part of the tus protocol,
// such as PostFile, HeadFile, PatchFile and DelFile. In addition the GetFile method
// is provided which is, however, not part of the specification.
type UploadHandler struct {
	composer *tusd.StoreComposer
	http.Handler

	MaxSize int64
	// BasePath defines the URL path used for handling uploads, e.g. "/files/".
	// If no trailing slash is presented it will be added. You may specify an
	// absolute URL containing a scheme, e.g. "http://tus.io"
	BasePath *url.URL

	*logrus.FieldLogger

	// CompleteUploads is used to send notifications whenever an upload is
	// completed by a user. The FileInfo will contain information about this
	// upload after it is completed. Sending to this channel will only
	// happen if the NotifyCompleteUploads field is set to true in the Config
	// structure. Notifications will also be sent for completions using the
	// Concatenation extension.
	NotifyCompleteUploads bool
	CompleteUploads       chan tusd.FileInfo
	// TerminatedUploads is used to send notifications whenever an upload is
	// terminated by a user. The FileInfo will contain information about this
	// upload gathered before the termination. Sending to this channel will only
	// happen if the NotifyTerminatedUploads field is set to true in the Config
	// structure.
	NotifyTerminatedUploads bool
	TerminatedUploads       chan tusd.FileInfo
	// UploadProgress is used to send notifications about the progress of the
	// currently running uploads. For each open PATCH request, every second
	// a FileInfo instance will be send over this channel with the Offset field
	// being set to the number of bytes which have been transfered to the server.
	// Please be aware that this number may be higher than the number of bytes
	// which have been stored by the data store! Sending to this channel will only
	// happen if the NotifyUploadProgress field is set to true in the Config
	// structure.
	NotifyUploadProgress bool
	UploadProgress       chan tusd.FileInfo
	// CreatedUploads is used to send notifications about the uploads having been
	// created. It triggers post creation and therefore has all the FileInfo incl.
	// the ID available already. It facilitates the post-create hook. Sending to
	// this channel will only happen if the NotifyCreatedUploads field is set to
	// true in the Config structure.
	NotifyCreatedUploads bool
	CreatedUploads       chan tusd.FileInfo
}

// NewUploadHandler creates a new handler without routing using the given
// configuration. It exposes the http handlers which need to be combined with
// a router (aka mux) of your choice. If you are looking for preconfigured
// handler see NewHandler.
func NewUploadHandler(config tusd.Config) (*UploadHandler, error) {
	l := logrus.New(nil)
	l.SetStdLogger(config.Logger)

	base := config.BasePath
	// Ensure base path ends with slash to remove logic from absFileURL
	if base != "" && string(base[len(base)-1]) != "/" {
		base += "/"
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return nil, errors.WithMessage(err, "malformed base path")

	}
	handler := &UploadHandler{
		composer:                config.StoreComposer,
		CompleteUploads:         make(chan tusd.FileInfo),
		TerminatedUploads:       make(chan tusd.FileInfo),
		UploadProgress:          make(chan tusd.FileInfo),
		CreatedUploads:          make(chan tusd.FileInfo),
		FieldLogger:             l,
		NotifyCompleteUploads:   config.NotifyCompleteUploads,
		NotifyTerminatedUploads: config.NotifyTerminatedUploads,
		NotifyUploadProgress:    config.NotifyUploadProgress,
		NotifyCreatedUploads:    config.NotifyCreatedUploads,
		MaxSize:                 config.MaxSize,
		BasePath:                baseUrl,
	}

	handler.Handler = handler.newHandler()
	return handler, nil
}
func (handler *UploadHandler) newHandler() http.Handler {
	m := pat.New()
	h := handler.Middleware(m)

	m.Post("", http.HandlerFunc(handler.PostFile))
	m.Post(":id", http.HandlerFunc(handler.PostFile))
	m.Put(":id", http.HandlerFunc(handler.PutFile))
	m.Head(":id", http.HandlerFunc(handler.HeadFile))
	m.Patch(":id", http.HandlerFunc(handler.PatchFile))
	m.Del(":id", http.HandlerFunc(handler.DelFile))
	m.Get(":id", http.HandlerFunc(handler.GetFile))
	return h
}

// Middleware checks various aspects of the request and ensures that it
// conforms with the spec. Also handles method overriding for clients which
// cannot make PATCH AND DELETE requests. If you are using the tusd handlers
// directly you will need to wrap at least the POST and PATCH endpoints in
// this middleware.
func (handler *UploadHandler) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow overriding the HTTP method. The reason for this is
		// that some libraries/environments to not support PATCH and
		// DELETE requests, e.g. Flash in a browser and parts of Java
		if newMethod := r.Header.Get("X-HTTP-Method-Override"); newMethod != "" {
			r.Method = newMethod
		}

		handler.GetLogger().WithField("method", r.Method).WithField("path", r.URL.Path).Debug("RequestIncoming")

		header := w.Header()

		if origin := r.Header.Get("Origin"); origin != "" {
			header.Set("Access-Control-Allow-Origin", origin)

			if r.Method == "OPTIONS" {
				// Preflight request
				header.Add("Access-Control-Allow-Methods", "POST, GET, HEAD, PATCH, DELETE, OPTIONS")
				header.Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Content-Length, Upload-Offset, Upload-Metadata")
				header.Set("Access-Control-Max-Age", "86400")
			} else {
				// Actual request
				header.Add("Access-Control-Expose-Headers", "Upload-Offset, Location, Content-Length, Upload-Metadata")
			}
		}

		// Set appropriated headers in case of OPTIONS method allowing protocol
		// discovery and end with an 204 No Content
		if r.Method == "OPTIONS" {
			if handler.MaxSize > 0 {
				header.Set("Tus-Max-Size", strconv.FormatInt(handler.MaxSize, 10))
			}

			// Although the 204 No Content status code is a better fit in this case,
			// since we do not have a response body included, we cannot use it here
			// as some browsers only accept 200 OK as successful response to a
			// preflight request. If we send them the 204 No Content the response
			// will be ignored or interpreted as a rejection.
			// For example, the Presto engine, which is used in older versions of
			// Opera, Opera Mobile and Opera Mini, handles CORS this way.
			handler.sendResp(w, r, http.StatusOK)
			return
		}

		// Proceed with routing the request
		h.ServeHTTP(w, r)
	})
}

// PostFile creates a new file upload using the datastore after validating the
// length and parsing the metadata.
func (handler *UploadHandler) PostFile(w http.ResponseWriter, r *http.Request) {
	uploadId, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		uploadId = "" // generate by file store
	}

	if uploadId != "" {
		info, err := handler.composer.Core.GetInfo(uploadId)
		if err == nil {
			w.Header().Set("Upload-Offset", strconv.FormatInt(info.Offset, 10))

			// If a resource has been created on the origin server, the response SHOULD be 201 (Created)
			// and contain an entity which describes
			// the status of the request and refers to the new resource, and a Location header
			// Add the Location header directly after creating the new resource to even
			// include it in cases of failure when an error is returned
			url, err := handler.absFileURL(r, uploadId)
			if err != nil {
				handler.sendError(w, r, err)
				return
			}
			w.Header().Set("Location", url)
			handler.sendResp(w, r, http.StatusCreated)
			return
		}
	}

	info := tusd.FileInfo{
		ID:             uploadId,
		SizeIsDeferred: true,
		IsFinal:        true,
	}

	id, err := handler.composer.Core.NewUpload(info)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	info.ID = id

	// Add the Location header directly after creating the new resource to even
	// include it in cases of failure when an error is returned
	url, err := handler.absFileURL(r, id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}
	w.Header().Set("Location", url)

	if handler.NotifyCreatedUploads {
		handler.CreatedUploads <- info
	}

	if handler.NotifyCompleteUploads {
		handler.CompleteUploads <- info
	}

	if handler.composer.UsesLocker {
		locker := handler.composer.Locker
		if err := locker.LockUpload(id); err != nil {
			handler.sendError(w, r, err)
			return
		}

		defer locker.UnlockUpload(id)
	}
	// Get Content-Length if possible
	length := r.ContentLength
	uploadLength, err := handler.writeChunk(id, info, length, w, r)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	if err := handler.composer.LengthDeferrer.DeclareLength(id, uploadLength); err != nil {
		handler.sendError(w, r, err)
		return
	}

	info, err = handler.composer.Core.GetInfo(id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}
	// Directly finish the upload if the upload is empty (i.e. has a size of 0).
	// This statement is in an else-if block to avoid causing duplicate calls
	// to finishUploadIfComplete if an upload is empty and contains a chunk.
	if _, err := handler.writeChunk(id, info, 0, w, r); err != nil {
		handler.sendError(w, r, err)
		return
	}

	handler.sendResp(w, r, http.StatusCreated)
}

// It may return an tusd.ErrNotFound which will be interpreted as a
// 404 Not Found.
func (handler *UploadHandler) deleteFile(id string) error {
	if handler.composer.UsesLocker {
		locker := handler.composer.Locker
		if err := locker.LockUpload(id); err != nil {
			return err
		}

		defer locker.UnlockUpload(id)
	}

	info, err := handler.composer.Core.GetInfo(id)
	// Interpret os.ErrNotExist as 404 Not Found
	if os.IsNotExist(err) {
		err = tusd.ErrNotFound
	}

	if err != nil {
		return err
	}

	// Abort the request handling if the required interface is not implemented
	if !handler.composer.UsesTerminater {
		return tusd.ErrNotImplemented
	}

	err = handler.composer.Terminater.Terminate(id)
	if err != nil {
		return err
	}

	if handler.NotifyTerminatedUploads {
		handler.TerminatedUploads <- info
	}

	return nil
}

// https://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html 9.6 PUT
// PutFile upgrades a new file upload using the datastore after validating the
// length and parsing the metadata.
// If a new resource is created, the origin server MUST inform the user agent via the 201 (Created) response.
// If an existing resource is modified, either the 200 (OK) or 204 (No Content) response codes SHOULD be sent to indicate successful completion of the request.
// If the resource could not be created or modified with the Request-URI, an appropriate error response SHOULD be given that reflects the nature of the problem.
func (handler *UploadHandler) PutFile(w http.ResponseWriter, r *http.Request) {
	var created bool
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	rangeReq := r.Header.Get("Content-Range")
	if rangeReq != "" {
		http.Error(w, "range is not support", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	_, err = handler.composer.Core.GetInfo(id)
	if err != nil {
		// Interpret os.ErrNotExist as 404 Not Found
		// Ignore the error if the upload could not be found. In this case, the upload
		// has likely already been removed by another service (e.g. a cron job).
		if err != tusd.ErrNotFound && !os.IsNotExist(err) {
			handler.sendError(w, r, err)
			return
		}

		// Handle 404 Not Found
		created = true
	}

	// POST new data

	info := tusd.FileInfo{
		ID:             id,
		SizeIsDeferred: true,
		IsFinal:        true,
	}

	id, err = handler.composer.Core.NewUpload(info)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	info.ID = id

	// Add the Location header directly after creating the new resource to even
	// include it in cases of failure when an error is returned
	url, err := handler.absFileURL(r, id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}
	w.Header().Set("Location", url)

	if handler.NotifyCreatedUploads {
		handler.CreatedUploads <- info
	}

	if handler.NotifyCompleteUploads {
		handler.CompleteUploads <- info
	}

	// Get Content-Length if possible
	length := r.ContentLength
	uploadLength, err := handler.writeChunk(id, info, length, w, r)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	if err := handler.composer.LengthDeferrer.DeclareLength(id, uploadLength); err != nil {
		handler.sendError(w, r, err)
		return
	}

	info, err = handler.composer.Core.GetInfo(id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}
	// Directly finish the upload if the upload is empty (i.e. has a size of 0).
	// This statement is in an else-if block to avoid causing duplicate calls
	// to finishUploadIfComplete if an upload is empty and contains a chunk.
	if _, err := handler.writeChunk(id, info, 0, w, r); err != nil {
		handler.sendError(w, r, err)
		return
	}

	if created {
		handler.sendResp(w, r, http.StatusCreated)
		return
	}
	handler.sendResp(w, r, http.StatusNoContent)
}

// HeadFile returns the length and offset for the HEAD request
func (handler *UploadHandler) HeadFile(w http.ResponseWriter, r *http.Request) {

	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	if handler.composer.UsesLocker {
		locker := handler.composer.Locker
		if err := locker.LockUpload(id); err != nil {
			handler.sendError(w, r, err)
			return
		}

		defer locker.UnlockUpload(id)
	}

	info, err := handler.composer.Core.GetInfo(id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	if len(info.MetaData) != 0 {
		w.Header().Set("Upload-Metadata", tusd.SerializeMetadataHeader(info.MetaData))
	}

	w.Header().Set("Cache-Control", "no-store")
	cr := httpContentRange{
		firstBytePos:   0,
		lastBytePos:    info.Offset - 1,
		completeLength: info.Size,
	}
	w.Header().Set("Content-Range", cr.String())
	handler.sendResp(w, r, http.StatusOK)
}

// PatchFile adds a chunk to an upload. This operation is only allowed
// if enough space in the upload is left.
func (handler *UploadHandler) PatchFile(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	if handler.composer.UsesLocker {
		locker := handler.composer.Locker
		if err := locker.LockUpload(id); err != nil {
			handler.sendError(w, r, err)
			return
		}

		defer locker.UnlockUpload(id)
	}

	info, err := handler.composer.Core.GetInfo(id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	if r.ContentLength == 0 && info.SizeIsDeferred {
		if err := handler.composer.LengthDeferrer.DeclareLength(id, info.Offset); err != nil {
			handler.sendError(w, r, err)
			return
		}
		handler.sendResp(w, r, http.StatusNoContent)
	}

	var size = info.Size
	if info.SizeIsDeferred {
		size = -1
	}

	// Test whether the size is still allowed
	if handler.MaxSize > 0 && size > handler.MaxSize {
		handler.sendError(w, r, tusd.ErrMaxSizeExceeded)
		return
	}

	ranges, err := parseContentRanges(r.Header["Content-Range"])
	if err != nil {
		if err == errNoOverlap {
			if size < 0 {
				w.Header().Set("Content-Range", "bytes */*")
			} else {
				w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", size))
			}
		}
		handler.sendResp(w, r, http.StatusRequestedRangeNotSatisfiable)
		return
	}

	if len(ranges) == 0 {
		ranges = append(ranges, httpContentRange{
			firstBytePos:   0,
			lastBytePos:    r.ContentLength - 1,
			completeLength: r.ContentLength,
		})
	}

	if len(ranges) > 0 && ranges[0].completeLength > size {
		// The total number of bytes in all the ranges
		// is larger than the size of the file by
		// itself, so this is probably an attack, or a
		// dumb client. Ignore the range request.
		ranges = nil
	}

	// check if seek is needed, else only append mode can be done
	if _, ok := handler.composer.Core.(filestore.FileStore); !ok {
		for _, ra := range ranges {
			if ra.firstBytePos != info.Offset {
				handler.sendResp(w, r, http.StatusRequestedRangeNotSatisfiable)
				return
			}
		}
	}

	for _, ra := range ranges {
		info.Offset = ra.firstBytePos

		_, err := handler.writeChunk(id, info, ra.lastBytePos-ra.firstBytePos+1, w, r)
		if err != nil {
			handler.sendError(w, r, err)
			return
		}
	}
	writeContentRanges(w, ranges)

	w.Header().Set("Accept-Ranges", "bytes")
	handler.sendResp(w, r, http.StatusNoContent)
}

// GetFile handles requests to download a file using a GET request. This is not
// part of the specification.
func (handler *UploadHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	if !handler.composer.UsesGetReader {
		handler.sendError(w, r, tusd.ErrNotImplemented)
		return
	}

	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	//if handler.composer.UsesLocker {
	//	locker := handler.composer.Locker
	//	if err := locker.LockUpload(id); err != nil {
	//		handler.sendError(w, r, err)
	//		return
	//	}
	//
	//	defer locker.UnlockUpload(id)
	//}

	info, err := handler.composer.Core.GetInfo(id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	// Set headers before sending responses
	w.Header().Set("Content-Length", strconv.FormatInt(info.Offset, 10))

	fileType, fileName := extractFileInfo(info)
	if fileType != "" {
		w.Header().Set("Content-Type", fileType)
	}

	// If no data has been uploaded yet, respond with an empty "204 No Content" status.
	if info.Offset == 0 {
		handler.sendResp(w, r, http.StatusNoContent)
		return
	}

	src, err := handler.composer.GetReader.GetReader(id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	ServeContent(w, r, fileName, time.Time{}, src, info.Offset)

	// Try to close the reader if the io.Closer interface is implemented
	if closer, ok := src.(io.Closer); ok {
		closer.Close()
	}

}

// DelFile terminates an upload permanently.
func (handler *UploadHandler) DelFile(w http.ResponseWriter, r *http.Request) {
	// Abort the request handling if the required interface is not implemented
	if !handler.composer.UsesTerminater {
		handler.sendError(w, r, tusd.ErrNotImplemented)
		return
	}

	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	if handler.composer.UsesLocker {
		locker := handler.composer.Locker
		if err := locker.LockUpload(id); err != nil {
			handler.sendError(w, r, err)
			return
		}

		defer locker.UnlockUpload(id)
	}

	var info tusd.FileInfo
	if handler.NotifyTerminatedUploads {
		info, err = handler.composer.Core.GetInfo(id)
		if err != nil {
			handler.sendError(w, r, err)
			return
		}
	}

	err = handler.composer.Terminater.Terminate(id)
	if err != nil {
		handler.sendError(w, r, err)
		return
	}

	handler.sendResp(w, r, http.StatusNoContent)

	if handler.NotifyTerminatedUploads {
		handler.TerminatedUploads <- info
	}

}

// extractFileInfo returns the values for the Content-Type and
// Content-Disposition headers for a given upload. These values should be used
// in responses for GET requests to ensure that only non-malicious file types
// are shown directly in the browser. It will extract the file name and type
// from the "fileame" and "filetype".
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
func extractFileInfo(info tusd.FileInfo) (filetype string, filename string) {
	// Add a filename to Content-Disposition if one is available in the metadata
	filename = info.MetaData["filename"]

	filetype = strings.TrimSpace(info.MetaData["filetype"])

	return filetype, filename
}

// writeChunk reads the body from the requests r and appends it to the upload
// with the corresponding id. Afterwards, it will set the necessary response
// headers but will not send the response.
func (handler *UploadHandler) writeChunk(id string, info tusd.FileInfo, length int64, w http.ResponseWriter, r *http.Request) (int64, error) {
	offset := info.Offset

	// Test if this upload fits into the file's size
	if !info.SizeIsDeferred && offset+length > info.Size {
		return 0, tusd.ErrSizeExceeded
	}

	maxSize := info.Size - offset
	// If the upload's length is deferred and the PATCH request does not contain the Content-Length
	// header (which is allowed if 'Transfer-Encoding: chunked' is used), we still need to set limits for
	// the body size.
	if info.SizeIsDeferred {
		if handler.MaxSize > 0 {
			// Ensure that the upload does not exceed the maximum upload size
			maxSize = handler.MaxSize - offset
		} else {
			// If no upload limit is given, we allow arbitrary sizes
			maxSize = math.MaxInt64
		}
	}
	if length > 0 {
		maxSize = length
	}

	var bytesWritten int64
	// Prevent a nil pointer dereference when accessing the body which may not be
	// available in the case of a malicious request.
	if r.Body != nil {
		var uploadFile io.ReadCloser
		uploadFile, _, err := r.FormFile("file")
		if err != nil {
			if err != http.ErrMissingFile && err != http.ErrNotMultipart {
				return 0, err
			}
			uploadFile = r.Body
		}
		defer uploadFile.Close()

		// Limit the data read from the request's body to the allowed maximum
		reader := io.LimitReader(uploadFile, maxSize)

		if handler.NotifyUploadProgress {
			var stop chan<- struct{}
			reader, stop = handler.sendProgressMessages(info, reader)
			defer close(stop)
		}

		bytesWritten, err = handler.composer.Core.WriteChunk(id, offset, reader)
		if err != nil {
			return 0, err
		}
	}

	// Send new offset to client
	newOffset := offset + bytesWritten
	w.Header().Set("Upload-Offset", strconv.FormatInt(newOffset, 10))
	info.Offset = newOffset

	return info.Offset, handler.finishUploadIfComplete(info)
}

// finishUploadIfComplete checks whether an upload is completed (i.e. upload offset
// matches upload size) and if so, it will call the data store's FinishUpload
// function and send the necessary message on the CompleteUpload channel.
func (handler *UploadHandler) finishUploadIfComplete(info tusd.FileInfo) error {
	// If the upload is completed, ...
	if !info.SizeIsDeferred && info.Offset == info.Size {
		// ... allow custom mechanism to finish and cleanup the upload
		if handler.composer.UsesFinisher {
			if err := handler.composer.Finisher.FinishUpload(info.ID); err != nil {
				return err
			}
		}

		// ... send the info out to the channel
		if handler.NotifyCompleteUploads {
			handler.CompleteUploads <- info
		}
	}

	return nil
}

// Send the error in the response body. The status code will be looked up in
// ErrStatusCodes. If none is found 500 Internal Error will be used.
func (handler *UploadHandler) sendError(w http.ResponseWriter, r *http.Request, err error) {
	// Interpret os.ErrNotExist as 404 Not Found
	if os.IsNotExist(err) {
		err = tusd.ErrNotFound
	}

	// Errors for read timeouts contain too much information which is not
	// necessary for us and makes grouping for the metrics harder. The error
	// message looks like: read tcp 127.0.0.1:1080->127.0.0.1:53673: i/o timeout
	// Therefore, we use a common error message for all of them.
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		err = errors.New("read tcp: i/o timeout")
	}

	statusErr, ok := err.(tusd.HTTPError)
	if !ok {
		statusErr = tusd.NewHTTPError(err, http.StatusInternalServerError)
	}

	reason := append(statusErr.Body(), '\n')
	if r.Method == "HEAD" {
		reason = nil
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(reason)))
	w.WriteHeader(statusErr.StatusCode())
	w.Write(reason)
}

// sendResp writes the header to w with the specified status code.
func (handler *UploadHandler) sendResp(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}

// Make an absolute URLs to the given upload id. If the base path is absolute
// it will be prepended else the host and protocol from the request is used.
func (handler *UploadHandler) absFileURL(r *http.Request, id string) (string, error) {
	idUrl, err := url.Parse(id)
	if err != nil {
		return "", err
	}

	fileUrl := idUrl
	if handler.BasePath != nil {
		fileUrl = handler.BasePath.ResolveReference(idUrl)

		if handler.BasePath.IsAbs() {
			return fileUrl.String(), nil
		}
	}
	fileUrl = r.URL.ResolveReference(fileUrl)

	fileProxyUrl := ResolveProxyUrl(fileUrl, r, true)

	return fileProxyUrl.String(), nil
}

type progressWriter struct {
	Offset int64
}

func (w *progressWriter) Write(b []byte) (int, error) {
	atomic.AddInt64(&w.Offset, int64(len(b)))
	return len(b), nil
}

// sendProgressMessage will send a notification over the UploadProgress channel
// every second, indicating how much data has been transfered to the server.
// It will stop sending these instances once the returned channel has been
// closed. The returned reader should be used to read the request body.
func (handler *UploadHandler) sendProgressMessages(info tusd.FileInfo, reader io.Reader) (io.Reader, chan<- struct{}) {
	previousOffset := int64(0)
	progress := &progressWriter{
		Offset: info.Offset,
	}
	stop := make(chan struct{}, 1)
	reader = io.TeeReader(reader, progress)

	go func() {
		for {
			select {
			case <-stop:
				info.Offset = atomic.LoadInt64(&progress.Offset)
				if info.Offset != previousOffset {
					handler.UploadProgress <- info
					previousOffset = info.Offset
				}
				return
			case <-time.After(1 * time.Second):
				info.Offset = atomic.LoadInt64(&progress.Offset)
				if info.Offset != previousOffset {
					handler.UploadProgress <- info
					previousOffset = info.Offset
				}
			}
		}
	}()

	return reader, stop
}

// The get sum of all sizes for a list of upload ids while checking whether
// all of these uploads are finished yet. This is used to calculate the size
// of a final resource.
func (handler *UploadHandler) sizeOfUploads(ids []string) (size int64, err error) {
	for _, id := range ids {
		info, err := handler.composer.Core.GetInfo(id)
		if err != nil {
			return size, err
		}

		if info.SizeIsDeferred || info.Offset != info.Size {
			err = tusd.ErrUploadNotFinished
			return size, err
		}

		size += info.Size
	}

	return
}

// extractIDFromPath pulls the last segment from the url provided
func extractIDFromPath(url string) (string, error) {
	result := reExtractFileID.FindStringSubmatch(url)
	if len(result) != 2 {
		return "", tusd.ErrNotFound
	}
	return result[1], nil
}
