package http

import (
	"github.com/bmizerany/pat"
	"github.com/searKing/tusd"
	"net/http"
)

const TusResumableDisabled = "1"

// ResumableUploadHandler is a ready to use handler with routing (using pat)
type ResumableUploadHandler struct {
	tusdHandler   *tusd.Handler
	uploadHandler *UploadHandler
}

// NewResumableUploadHandler creates a routed tus protocol handler. This is the simplest
// way to use tusd but may not be as configurable as you require. If you are
// integrating this into an existing app you may like to use tusd.NewUnroutedHandler
// instead. Using tusd.NewUnroutedHandler allows the tus handlers to be combined into
// your existing router (aka mux) directly. It also allows the GET and DELETE
// endpoints to be customized. These are not part of the protocol so can be
// changed depending on your needs.
func NewResumableUploadHandler(config tusd.Config) (*ResumableUploadHandler, error) {
	tusdHandler, err := tusd.NewHandler(config)
	if err != nil {
		return nil, err
	}

	uploadHandler, err := NewUploadHandler(config)
	if err != nil {
		return nil, err
	}

	return &ResumableUploadHandler{
		tusdHandler:   tusdHandler,
		uploadHandler: uploadHandler,
	}, nil
}

func (handler *ResumableUploadHandler) Handler(id string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		r.URL.Path = id
		defer func() {
			r.URL.Path = path
		}()

		m := pat.New()
		h := handler.Middleware(m)

		m.Post("", http.HandlerFunc(handler.PostFile))
		m.Post(":id", http.HandlerFunc(handler.PostFile))
		m.Put(":id", http.HandlerFunc(handler.PutFile))
		m.Head(":id", http.HandlerFunc(handler.HeadFile))
		m.Patch(":id", http.HandlerFunc(handler.PatchFile))
		m.Del(":id", http.HandlerFunc(handler.DelFile))
		m.Get(":id", http.HandlerFunc(handler.GetFile))
		h.ServeHTTP(w, r)
	})
}

// Middleware checks various aspects of the request and ensures that it
// conforms with the spec. Also handles method overriding for clients which
// cannot make PATCH AND DELETE requests. If you are using the tusd handlers
// directly you will need to wrap at least the POST and PATCH endpoints in
// this middleware.
func (handler *ResumableUploadHandler) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Tus-Resumable") == "" {
			// Set current version used by the server
			r.Header.Set("Tus-Resumable", "1.0.0")
			r.Header.Set("Tus-Resumable-Disable", TusResumableDisabled)
		}
		handler.tusdHandler.Middleware(h).ServeHTTP(w, r)
	})
}

// Verify that the Upload-Length and Upload-Defer-Length headers are acceptable for creating a
// new upload
func (handler *ResumableUploadHandler) validateTusResumableDisableHeader(tusResumableDisableHeader string) (tusResumableDisabled bool, err error) {
	haveInvalidDeferHeader := tusResumableDisableHeader != "" && tusResumableDisableHeader != TusResumableDisabled
	tusIsDisabled := tusResumableDisableHeader == TusResumableDisabled

	if haveInvalidDeferHeader {
		err = tusd.ErrInvalidUploadDeferLength
		return
	}
	if tusIsDisabled {
		tusResumableDisabled = true
		return
	}
	return
}

// PostFile creates a new file upload using the datastore after validating the
// length and parsing the metadata.
func (handler *ResumableUploadHandler) PostFile(w http.ResponseWriter, r *http.Request) {
	tusResumableDisabled, err := handler.validateTusResumableDisableHeader(r.Header.Get("Tus-Resumable-Disable"))
	if err != nil || !tusResumableDisabled {
		handler.tusdHandler.PostFile(w, r)
		return
	}
	handler.uploadHandler.PostFile(w, r)
	return
}

// PutFile upgrades a new file upload using the datastore after validating the
// length and parsing the metadata.
func (handler *ResumableUploadHandler) PutFile(w http.ResponseWriter, r *http.Request) {
	tusResumableDisabled, err := handler.validateTusResumableDisableHeader(r.Header.Get("Tus-Resumable-Disable"))
	if err != nil || !tusResumableDisabled {
		handler.tusdHandler.PutFile(w, r)
		return
	}
	handler.uploadHandler.PutFile(w, r)
	return
}

// HeadFile returns the length and offset for the HEAD request
func (handler *ResumableUploadHandler) HeadFile(w http.ResponseWriter, r *http.Request) {
	tusResumableDisabled, err := handler.validateTusResumableDisableHeader(r.Header.Get("Tus-Resumable-Disable"))
	if err != nil || !tusResumableDisabled {
		handler.tusdHandler.HeadFile(w, r)
		return
	}
	handler.uploadHandler.HeadFile(w, r)
	return
}

// PatchFile adds a chunk to an upload. This operation is only allowed
// if enough space in the upload is left.
func (handler *ResumableUploadHandler) PatchFile(w http.ResponseWriter, r *http.Request) {
	tusResumableDisabled, err := handler.validateTusResumableDisableHeader(r.Header.Get("Tus-Resumable-Disable"))
	if err != nil || !tusResumableDisabled {
		handler.tusdHandler.PatchFile(w, r)
		return
	}
	handler.uploadHandler.PatchFile(w, r)
	return
}

// DelFile terminates an upload permanently.
func (handler *ResumableUploadHandler) DelFile(w http.ResponseWriter, r *http.Request) {
	tusResumableDisabled, err := handler.validateTusResumableDisableHeader(r.Header.Get("Tus-Resumable-Disable"))
	if err != nil || !tusResumableDisabled {
		handler.tusdHandler.DelFile(w, r)
		return
	}
	handler.uploadHandler.DelFile(w, r)
	return
}

// GetFile handles requests to download a file using a GET request. This is not
// part of the specification.
func (handler *ResumableUploadHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	tusResumableDisabled, err := handler.validateTusResumableDisableHeader(r.Header.Get("Tus-Resumable-Disable"))
	if err != nil || !tusResumableDisabled {
		handler.tusdHandler.GetFile(w, r)
		return
	}
	handler.uploadHandler.GetFile(w, r)
	return
}
