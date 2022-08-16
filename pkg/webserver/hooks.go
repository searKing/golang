// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	errors_ "github.com/searKing/golang/go/errors"
	"github.com/searKing/golang/go/runtime"
	"github.com/searKing/golang/pkg/webserver/healthz"
	"github.com/sirupsen/logrus"
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

// AddPostStartHook allows you to add a PostStartHook.
func (s *WebServer) AddPostStartHook(name string, hook PostStartHookFunc) error {
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

	done := make(chan struct{})
	s.postStartHooks[name] = postStartHookEntry{hook: hook, originatingStack: string(debug.Stack()), done: done}

	return nil
}

// AddPostStartHookOrDie allows you to add a PostStartHook, but dies on failure
func (s *WebServer) AddPostStartHookOrDie(name string, hook PostStartHookFunc) {
	if err := s.AddPostStartHook(name, hook); err != nil {
		logrus.Fatalf("Error registering PostStartHook %q: %v", name, err)
	}
}

// AddPreShutdownHook allows you to add a PreShutdownHook.
func (s *WebServer) AddPreShutdownHook(name string, hook PreShutdownHookFunc) error {
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

	s.preShutdownHooks[name] = preShutdownHookEntry{hook: hook}

	return nil
}

// AddPreShutdownHookOrDie allows you to add a PostStartHook, but dies on failure
func (s *WebServer) AddPreShutdownHookOrDie(name string, hook PreShutdownHookFunc) {
	if err := s.AddPreShutdownHook(name, hook); err != nil {
		logrus.Fatalf("Error registering PreShutdownHook %q: %v", name, err)
	}
}

// RunPostStartHooks runs the PostStartHooks for the server
func (s *WebServer) RunPostStartHooks(ctx context.Context) error {
	var errs []error
	s.postStartHookLock.Lock()
	defer s.postStartHookLock.Unlock()
	s.postStartHooksCalled = true

	for hookName, hookEntry := range s.postStartHooks {
		if err := runPostStartHook(ctx, hookName, hookEntry); err != nil {
			errs = append(errs, err)
		}
	}
	return errors_.Multi(errs...)
}

// RunPreShutdownHooks runs the PreShutdownHooks for the server
func (s *WebServer) RunPreShutdownHooks() error {
	var errs []error

	s.preShutdownHookLock.Lock()
	defer s.preShutdownHookLock.Unlock()
	s.preShutdownHooksCalled = true

	for hookName, hookEntry := range s.preShutdownHooks {
		if err := runPreShutdownHook(hookName, hookEntry); err != nil {
			errs = append(errs, err)
		}
	}
	return errors_.Multi(errs...)
}

// isPostStartHookRegistered checks whether a given PostStartHook is registered
func (s *WebServer) isPostStartHookRegistered(name string) bool {
	s.postStartHookLock.Lock()
	defer s.postStartHookLock.Unlock()
	_, exists := s.postStartHooks[name]
	return exists
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
