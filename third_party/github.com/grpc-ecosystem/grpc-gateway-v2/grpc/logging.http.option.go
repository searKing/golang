// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/codes"
)

// HttpLoggingOption customizes HttpInterceptor and HttpRoundTripDecorator,
// as logging.Option does for logging interceptors of gRPC, which is not applicable to http as opaque.
type HttpLoggingOption func(*httpLoggingOptions)

type httpLoggingOptions struct {
	levelFunc            logging.CodeToLevel
	loggableEvents       []logging.LoggableEvent
	errorToFieldsFunc    logging.ErrorToFields
	codeFunc             HttpToCode
	durationFieldFunc    logging.DurationToFields
	timestampFormat      string
	fieldsFromRequestFn  func(ctx context.Context, r *http.Request) logging.Fields
	disableHttpLogFields []string
}

func evaluateHttpLoggingOptions(levelFunc logging.CodeToLevel, opts []HttpLoggingOption) *httpLoggingOptions {
	o := &httpLoggingOptions{
		levelFunc:         levelFunc,
		loggableEvents:    []logging.LoggableEvent{logging.StartCall, logging.FinishCall},
		codeFunc:          DefaultHttpToCode,
		durationFieldFunc: DefaultHttpDurationToFields,
		timestampFormat:   time.RFC3339,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *httpLoggingOptions) has(event logging.LoggableEvent) bool {
	return slices.Contains(o.loggableEvents, event)
}

// HttpToCode maps the http status code, or the error of the http client if not nil, to the gRPC code.
type HttpToCode func(status int, err error) codes.Code

// DefaultHttpToCode maps the http status code to the gRPC code by reversing runtime.HTTPStatusFromCode of
// grpc-gateway, with the less severe code if ambiguous, or the error of the http client to the gRPC code
// reported if the transport fails, such as Unavailable.
func DefaultHttpToCode(status int, err error) codes.Code {
	if err != nil {
		return httpErrorToCode(err)
	}
	return httpStatusToCode(status)
}

// DefaultHttpDurationToFields is the default implementation of converting request duration to a field.
var DefaultHttpDurationToFields = HttpDurationToTimeMillisFields

// HttpDurationToTimeMillisFields converts the duration to milliseconds and uses the key `http.time_ms`,
// as logging.DurationToTimeMillisFields does with the key `grpc.time_ms`.
func HttpDurationToTimeMillisFields(duration time.Duration) logging.Fields {
	return logging.Fields{"http.time_ms", fmt.Sprintf("%v", float32(duration.Nanoseconds()/1000)/1000)}
}

// HttpDurationToDurationField uses a Duration field to log the request duration with the key `http.duration`,
// as logging.DurationToDurationField does with the key `grpc.duration`.
func HttpDurationToDurationField(duration time.Duration) logging.Fields {
	return logging.Fields{"http.duration", duration.String()}
}

// HttpFieldsFromRequestWithForward fill "X-Forwarded-For" and "X-Forwarded-Host" to record http callers,
// as FieldsFromContextWithForward does for gRPC.
func HttpFieldsFromRequestWithForward(_ context.Context, r *http.Request) logging.Fields {
	fields := logging.Fields{}
	for _, key := range []string{"X-Forwarded-For", "X-Forwarded-Host"} {
		if fwd := r.Header.Values(key); len(fwd) > 0 {
			fields = fields.AppendUnique(logging.Fields{"http." + strings.ToLower(key), fwd})
		}
	}
	return fields
}

// WithHttpFieldsFromContext allows overriding existing or adding extra fields to all log messages per given context,
// as logging.WithFieldsFromContext does.
// Only one of WithHttpFieldsFromContext or WithHttpFieldsFromContextAndRequest should be used,
// using both will result in the last one overwriting the previous.
func WithHttpFieldsFromContext(f func(ctx context.Context) logging.Fields) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.fieldsFromRequestFn = func(ctx context.Context, _ *http.Request) logging.Fields {
			return f(ctx)
		}
	}
}

// WithHttpFieldsFromContextAndRequest allows overriding existing or adding extra fields to all log messages
// per given context and http request, as logging.WithFieldsFromContextAndCallMeta does.
// Only one of WithHttpFieldsFromContext or WithHttpFieldsFromContextAndRequest should be used,
// using both will result in the last one overwriting the previous.
func WithHttpFieldsFromContextAndRequest(f func(ctx context.Context, r *http.Request) logging.Fields) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.fieldsFromRequestFn = f
	}
}

// WithHttpLogOnEvents customizes on what events the http interceptor should log on,
// as logging.WithLogOnEvents does, where only logging.StartCall and logging.FinishCall are supported.
func WithHttpLogOnEvents(events ...logging.LoggableEvent) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.loggableEvents = events
	}
}

// WithHttpErrorFields allows to extract logging fields from an error of the http client,
// as logging.WithErrorFields does.
func WithHttpErrorFields(f logging.ErrorToFields) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.errorToFieldsFunc = f
	}
}

// WithHttpLevels customizes the function for mapping gRPC codes and interceptor log level statements,
// as logging.WithLevels does, where the gRPC code is mapped by the function customized by WithHttpCodes.
func WithHttpLevels(f logging.CodeToLevel) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.levelFunc = f
	}
}

// WithHttpCodes customizes the function for mapping http status codes and errors to gRPC codes,
// as logging.WithCodes does.
func WithHttpCodes(f HttpToCode) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.codeFunc = f
	}
}

// WithHttpDurationField customizes the function for mapping request durations to log fields,
// as logging.WithDurationField does.
func WithHttpDurationField(f logging.DurationToFields) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.durationFieldFunc = f
	}
}

// WithHttpTimestampFormat customizes the timestamps emitted in the log fields,
// as logging.WithTimestampFormat does.
func WithHttpTimestampFormat(format string) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.timestampFormat = format
	}
}

// WithHttpDisableLoggingFields disables logging of http fields provided, as logging.WithDisableLoggingFields does.
// The following are the default logging fields:
//   - SystemTag[0]
//   - ComponentFieldKey
//   - http.method
//   - http.request.uri
//   - http.host
//   - http.remote_addr, of servers only
//   - http.client_ip, of servers only
//
// Usage example - WithHttpDisableLoggingFields("http.remote_addr", "http.client_ip")
func WithHttpDisableLoggingFields(disableHttpLogFields ...string) HttpLoggingOption {
	return func(o *httpLoggingOptions) {
		o.disableHttpLogFields = disableHttpLogFields
	}
}
