// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"slices"

	"github.com/searKing/golang/go/runtime"
	"github.com/searKing/golang/pkg/webserver/healthz"
	"golang.org/x/sync/errgroup"
)

// PostStartHookFunc is a function that is called after the server has started.
// It must properly handle cases like:
//  1. asynchronous start in multiple API server processes
//  2. conflicts between the different processes all trying to perform the same action
//  3. partially complete work (API server crashes while running your hook)
//  4. API server access **BEFORE** your hook has completed
//
// Think of it like a mini-controller that is super privileged and gets to run in-process
// If you use this feature, tag @deads2k on github who has promised to review code for anyone's PostStartHook
// until it becomes easier to use.
// ctx will be cancelled when WebServer is Closed or any other PostStartHookFunc failed.
type PostStartHookFunc func(ctx context.Context) error

// PreShutdownHookFunc is a function that can be added to the shutdown logic.
type PreShutdownHookFunc func() error

// PostStartHookProvider is an interface in addition to provide a post start hook for the api server
type PostStartHookProvider interface {
	PostStartHook() (string, PostStartHookFunc, error)
}

type postStartHookEntry struct {
	hook PostStartHookFunc
	// originatingStack holds the stack that registered postStartHooks. This allows us to show a more helpful message
	// for duplicate registration.
	originatingStack string

	// done will be closed when the postHook is finished
	done chan struct{}
}

type preShutdownHookEntry struct {
	hook PreShutdownHookFunc
}

// AddBootSequencePostStartHook allows you to add a PostStartHook in order.
func (s *WebServer) AddBootSequencePostStartHook(name string, hook PostStartHookFunc) error {
	return s.addPostStartHook(name, hook, true)
}

// AddPostStartHook allows you to add a PostStartHook.
func (s *WebServer) AddPostStartHook(name string, hook PostStartHookFunc) error {
	return s.addPostStartHook(name, hook, false)
}

// AddPostStartHookOrDie allows you to add a PostStartHook, but dies on failure
func (s *WebServer) AddPostStartHookOrDie(name string, hook PostStartHookFunc) {
	if err := s.AddPostStartHook(name, hook); err != nil {
		slog.Error(fmt.Sprintf("Error registering PostStartHook %q: %s", name, err.Error()))
		os.Exit(1)
	}
}

// AddBootSequencePreShutdownHook allows you to add a PreShutdownHook in reverse order.
func (s *WebServer) AddBootSequencePreShutdownHook(name string, hook PreShutdownHookFunc) error {
	return s.addPreShutdownHook(name, hook, true)
}

// AddPreShutdownHook allows you to add a PreShutdownHook.
func (s *WebServer) AddPreShutdownHook(name string, hook PreShutdownHookFunc) error {
	return s.addPreShutdownHook(name, hook, false)
}

// AddPreShutdownHookOrDie allows you to add a PostStartHook, but dies on failure
func (s *WebServer) AddPreShutdownHookOrDie(name string, hook PreShutdownHookFunc) {
	if err := s.AddPreShutdownHook(name, hook); err != nil {
		slog.Error(fmt.Sprintf("Error registering PreShutdownHook %q: %s", name, err.Error()))
		os.Exit(1)
	}
}

// RunPostStartHooks runs the PostStartHooks for the server
func (s *WebServer) RunPostStartHooks(ctx context.Context) error {
	s.postStartHookLock.Lock()
	defer s.postStartHookLock.Unlock()
	s.postStartHooksCalled = true

	g, gCtx := errgroup.WithContext(ctx)
	var keys = s.postStartHookOrderedKeys
	for k := range s.postStartHooks {
		if !slices.Contains(s.postStartHookOrderedKeys, k) {
			keys = append(keys, k)
		}
	}

	for i, k := range keys {
		if v, has := s.postStartHooks[k]; has {
			hookName, hookEntry := k, v
			if i < len(s.postStartHookOrderedKeys) {
				if err := runPostStartHook(gCtx, hookName, hookEntry); err != nil {
					return err
				}
				continue
			}

			g.Go(func() error {
				return runPostStartHook(gCtx, hookName, hookEntry)
			})
		} else { // never happen
			hookName := k
			slog.Warn(fmt.Sprintf("unknown PostStartHook %q", hookName))
		}
	}
	return g.Wait()
}

// RunPreShutdownHooks runs the PreShutdownHooks for the server
func (s *WebServer) RunPreShutdownHooks() error {
	s.preShutdownHookLock.Lock()
	defer s.preShutdownHookLock.Unlock()
	s.preShutdownHooksCalled = true

	var keys = s.preShutdownHookOrderedKeys
	for k := range s.preShutdownHooks {
		if !slices.Contains(s.preShutdownHookOrderedKeys, k) {
			keys = append(keys, k)
		}
	}
	slices.Reverse(keys)
	var errs []error
	for _, k := range keys {
		if v, has := s.preShutdownHooks[k]; has {
			hookName, hookEntry := k, v
			errs = append(errs, runPreShutdownHook(hookName, hookEntry))
		} else {
			hookName := k
			slog.Warn(fmt.Sprintf("unknown PreShutdownHook %q", hookName))
		}
	}
	return errors.Join(errs...)
}

// isPostStartHookRegistered checks whether a given PostStartHook is registered
func (s *WebServer) isPostStartHookRegistered(name string) bool {
	s.postStartHookLock.Lock()
	defer s.postStartHookLock.Unlock()
	_, exists := s.postStartHooks[name]
	return exists
}

func (s *WebServer) addPostStartHook(name string, hook PostStartHookFunc, order bool) error {
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}
	if hook == nil {
		return fmt.Errorf("hook func may not be nil: %q", name)
	}

	s.postStartHookLock.Lock()
	defer s.postStartHookLock.Unlock()

	if s.postStartHooksCalled {
		return fmt.Errorf("unable to add %q because PostStartHooks have already been called", name)
	}
	if postStartHook, exists := s.postStartHooks[name]; exists {
		// this is programmer error, but it can be hard to debug
		return fmt.Errorf("unable to add %q because it was already registered by: %s", name, postStartHook.originatingStack)
	}

	// done is closed when the poststarthook is finished.  This is used by the health check to be able to indicate
	// that the poststarthook is finished
	done := make(chan struct{})
	if err := s.AddBootSequenceHealthChecks(postStartHookHealthz{name: "poststarthook/" + name, done: done}); err != nil {
		return err
	}
	if order {
		s.postStartHookOrderedKeys = append(s.postStartHookOrderedKeys, name)
	}
	s.postStartHooks[name] = postStartHookEntry{hook: hook, originatingStack: string(debug.Stack()), done: done}
	return nil
}

func (s *WebServer) addPreShutdownHook(name string, hook PreShutdownHookFunc, order bool) error {
	if len(name) == 0 {
		return fmt.Errorf("missing name")
	}
	if hook == nil {
		return nil
	}

	s.preShutdownHookLock.Lock()
	defer s.preShutdownHookLock.Unlock()

	if s.preShutdownHooksCalled {
		return fmt.Errorf("unable to add %q because PreShutdownHooks have already been called", name)
	}
	if _, exists := s.preShutdownHooks[name]; exists {
		return fmt.Errorf("unable to add %q because it is already registered", name)
	}

	if order {
		s.preShutdownHookOrderedKeys = append(s.preShutdownHookOrderedKeys, name)
	}
	s.preShutdownHooks[name] = preShutdownHookEntry{hook: hook}

	return nil
}

func runPostStartHook(ctx context.Context, name string, entry postStartHookEntry) error {
	var err error
	func() {
		// don't let the hook *accidentally* panic and kill the server
		defer runtime.NeverPanicButLog.Recover()
		err = entry.hook(ctx)
	}()
	// if the hook intentionally wants to kill server, let it.
	if err != nil {
		return fmt.Errorf("PostStartHook %q failed: %w", name, err)
	}
	close(entry.done)
	return nil
}

func runPreShutdownHook(name string, entry preShutdownHookEntry) error {
	var err error
	func() {
		// don't let the hook *accidentally* panic and kill the server
		defer runtime.NeverPanicButLog.Recover()
		err = entry.hook()
	}()
	if err != nil {
		return fmt.Errorf("PreShutdownHook %q failed: %w", name, err)
	}
	return nil
}

// postStartHookHealthz implements a healthz check for poststarthooks.  It will return a "hookNotFinished"
// error until the poststarthook is finished.
type postStartHookHealthz struct {
	name string

	// done will be closed when the postStartHook is finished
	done chan struct{}
}

var _ healthz.HealthChecker = postStartHookHealthz{}

func (h postStartHookHealthz) Name() string {
	return h.name
}

var hookNotFinished = errors.New("not finished")

func (h postStartHookHealthz) Check(req *http.Request) error {
	select {
	case <-h.done:
		return nil
	default:
		return hookNotFinished
	}
}
