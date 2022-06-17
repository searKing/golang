// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

// InstallHandler registers handlers for health checking on the path
// "/healthz" to Muxer. *All handlers* for Muxer must be specified in
// exactly one call to InstallHandler. Calling InstallHandler more
// than once for the same Muxer will result in a panic.
func InstallHandler(mux Muxer, checks ...HealthChecker) {
	InstallPathHandler(mux, "/healthz", checks...)
}

// InstallReadyzHandler registers handlers for health checking on the path
// "/readyz" to Muxer. *All handlers* for Muxer must be specified in
// exactly one call to InstallHandler. Calling InstallHandler more
// than once for the same Muxer will result in a panic.
func InstallReadyzHandler(mux Muxer, checks ...HealthChecker) {
	InstallPathHandler(mux, "/readyz", checks...)
}

// InstallLivezHandler registers handlers for liveness checking on the path
// "/livez" to Muxer. *All handlers* for Muxer must be specified in
// exactly one call to InstallHandler. Calling InstallHandler more
// than once for the same Muxer will result in a panic.
func InstallLivezHandler(mux Muxer, checks ...HealthChecker) {
	InstallPathHandler(mux, "/livez", checks...)
}

// InstallReadyzHandlerWithHealthyFunc is like InstallReadyzHandler, but in addition call firstTimeReady
// the first time /readyz succeeds.
func InstallReadyzHandlerWithHealthyFunc(mux Muxer, firstTimeReady func(), checks ...HealthChecker) {
	InstallPathHandlerWithHealthyFunc(mux, "/readyz", firstTimeReady, checks...)
}

// InstallPathHandler registers handlers for health checking on
// a specific path to Muxer. *All handlers* for the path must be
// specified in exactly one call to InstallPathHandler. Calling
// InstallPathHandler more than once for the same path and Muxer will
// result in a panic.
func InstallPathHandler(mux Muxer, path string, checks ...HealthChecker) {
	InstallPathHandlerWithHealthyFunc(mux, path, nil, checks...)
}

// InstallPathHandlerWithHealthyFunc is like InstallPathHandler, but calls firstTimeHealthy exactly once
// when the handler succeeds for the first time.
func InstallPathHandlerWithHealthyFunc(mux Muxer, path string, firstTimeHealthy func(), checks ...HealthChecker) {
	if len(checks) == 0 {
		logrus.Infof("No default health checks specified. Installing the ping handler.")
		checks = []HealthChecker{PingHealthzCheck}
	}

	logrus.Infof("Installing health checkers for (%v): %v", path, formatQuoted(checkerNames(checks...)...))

	name := strings.Split(strings.TrimPrefix(path, "/"), "/")[0]
	mux.Handle(path,
		handleRootHealth(name, firstTimeHealthy, checks...))
	for _, check := range checks {
		mux.Handle(fmt.Sprintf("%s/%v", path, check.Name()), adaptCheckToHandler(check.Check))
	}
}

// Muxer is an interface describing the methods InstallHandler requires.
type Muxer interface {
	Handle(pattern string, handler http.Handler)
}

// handleRootHealth returns an http.HandlerFunc that serves the provided checks.
func handleRootHealth(name string, firstTimeHealthy func(), checks ...HealthChecker) http.HandlerFunc {
	var notifyOnce sync.Once
	return func(w http.ResponseWriter, r *http.Request) {
		excluded := getExcludedChecks(r)
		// failedVerboseLogOutput is for output to the log.  It indicates detailed failed output information for the log.
		var failedVerboseLogOutput bytes.Buffer
		var failedChecks []string
		var individualCheckOutput bytes.Buffer
		for _, check := range checks {
			// no-op the check if we've specified we want to exclude the check
			if _, has := excluded[check.Name()]; has {
				delete(excluded, check.Name())
				_, _ = fmt.Fprintf(&individualCheckOutput, "[+]%s excluded: ok\n", check.Name())
				continue
			}
			if err := check.Check(r); err != nil {
				// don't include the error since this endpoint is public.  If someone wants more detail
				// they should have explicit permission to the detailed checks.
				_, _ = fmt.Fprintf(&individualCheckOutput, "[-]%s failed: reason withheld\n", check.Name())
				// but we do want detailed information for our log
				_, _ = fmt.Fprintf(&failedVerboseLogOutput, "[-]%s failed: %v\n", check.Name(), err)
				failedChecks = append(failedChecks, check.Name())
			} else {
				_, _ = fmt.Fprintf(&individualCheckOutput, "[+]%s ok\n", check.Name())
			}
		}
		if len(excluded) > 0 {
			_, _ = fmt.Fprintf(&individualCheckOutput, "warn: some health checks cannot be excluded: no matches for %s\n", formatQuoted(maps.Keys(excluded)...))
			logrus.Warningf("cannot exclude some health checks, no health checks are installed matching %s",
				formatQuoted(maps.Keys(excluded)...))
		}
		// always be verbose on failure
		if len(failedChecks) > 0 {
			logrus.Errorf("%s check failed: %s\n%v", strings.Join(failedChecks, ","), name, failedVerboseLogOutput.String())
			http.Error(w, fmt.Sprintf("%s%s check failed", individualCheckOutput.String(), name), http.StatusInternalServerError)
			return
		}

		// signal first time this is healthy
		if firstTimeHealthy != nil {
			notifyOnce.Do(firstTimeHealthy)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if _, found := r.URL.Query()["verbose"]; !found {
			_, _ = fmt.Fprint(w, "ok")
			return
		}

		_, _ = individualCheckOutput.WriteTo(w)
		_, _ = fmt.Fprintf(w, "%s check passed\n", name)
	}
}

// adaptCheckToHandler returns an http.HandlerFunc that serves the provided checks.
func adaptCheckToHandler(c func(r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		} else {
			_, _ = fmt.Fprint(w, "ok")
		}
	}
}

// checkerNames returns the names of the checks in the same order as passed in.
func checkerNames(checks ...HealthChecker) []string {
	// accumulate the names of checks for printing them out.
	checkerNames := make([]string, 0, len(checks))
	for _, check := range checks {
		checkerNames = append(checkerNames, check.Name())
	}
	return checkerNames
}

// formatQuoted returns a formatted string of the health check names,
// preserving the order passed in.
func formatQuoted(names ...string) string {
	quoted := make([]string, 0, len(names))
	for _, name := range names {
		quoted = append(quoted, fmt.Sprintf("%q", name))
	}
	return strings.Join(quoted, ",")
}
