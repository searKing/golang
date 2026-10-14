# CI Test Optimization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make every CI platform run the short suite and remove avoidable wall-clock and filesystem cost from short tests while preserving the full local suite.

**Architecture:** Existing tests retain full behavior under normal `go test`. Short mode either scales the workload, uses deterministic synchronization, or skips tests whose subject is a real wall-clock boundary. CI passes `-short` on every operating system.

**Tech Stack:** Go `testing`, GitHub Actions, GNU Make

**Spec:** `docs/superpowers/specs/2026-10-14-ci-test-optimization-design.md`

## Global Constraints

- All CI operating systems run `go test -short`.
- Local `make test` remains the complete suite.
- Preserve assertions and full-test stress intensity.
- Do not implement incremental module selection in this PR.

---

### Task 1: Run the short suite on every CI platform

**Files:**
- Modify: `.github/workflows/go.yml`

**Interfaces:**
- Consumes: `Makefile` variable `GOTEST_FLAGS`
- Produces: every matrix entry invokes `make test GOTEST_FLAGS="-short"`

- [ ] **Step 1:** Change the Ubuntu matrix entry from an empty `gotest_flags` value to `-short`, and update the matrix comment to say stress tests run locally.
- [ ] **Step 2:** Parse the workflow with Ruby YAML:

  ```sh
  ruby -e 'require "yaml"; YAML.load_file(".github/workflows/go.yml")'
  ```

- [ ] **Step 3:** Commit:

  ```sh
  git add .github/workflows/go.yml
  git commit -m "ci: run short tests on every platform"
  ```

### Task 2: Reduce avoidable short-test work

**Files:**
- Modify: `go/io/copy_test.go`
- Modify: `go/testing/leakcheck/leakcheck_test.go`
- Modify: `go/time/rate/rate_test.go`
- Modify: `go/time/rate/example_test.go`
- Modify: `go/os/file_rotate_test.go`
- Modify: `go/net/mux/server_test.go`

**Interfaces:**
- Consumes: `testing.Short()`
- Produces: the same assertions with less short-mode work

- [ ] **Step 1:** Record focused short baselines:

  ```sh
  cd go
  GOWORK=off go test -short -count=1 ./io ./testing/leakcheck ./time/rate ./os ./net/mux
  ```

- [ ] **Step 2:** In `TestCopyDir`, use recursion depth 2 in short mode and depth 3 otherwise. This keeps nested-directory coverage while reducing generated files by about 10x.
- [ ] **Step 3:** In leakcheck tests, define durations from `testing.Short()`: leaked goroutines live 200ms, the detecting check waits 100ms, ignored goroutines live 300ms, and the cleanup check waits 500ms. Keep the existing 2s/1s/3s values in full mode.
- [ ] **Step 4:** In `TestWaitBlock`, use a 100ms context timeout in short mode and 1s otherwise.
- [ ] **Step 5:** Rewrite `TestSimultaneousLongRequests` to hold acquired tokens on a release channel. Wait for all acquisition attempts, close the release channel, wait for completion, then perform the final request. Remove the 1s sleeps without changing the expected count.
- [ ] **Step 6:** In `ExampleNewFullBurstLimiter`, cancel a background context immediately before the final empty-bucket wait and update expected output to `context canceled`. Do not wait for a real one-second deadline.
- [ ] **Step 7:** Skip short mode in the four top-level file-rotation tests that call `sleepToNextSecond`, and in the slow subtest `log file removed before rotate`.
- [ ] **Step 8:** Skip `net/mux.TestTimeout` in short mode because its assertion depends on real 500ms timeout boundaries.
- [ ] **Step 9:** Run focused tests in both modes:

  ```sh
  cd go
  GOWORK=off go test -short -count=3 ./io ./testing/leakcheck ./time/rate ./os ./net/mux
  GOWORK=off go test -count=1 ./io ./testing/leakcheck ./time/rate ./os ./net/mux
  ```

- [ ] **Step 10:** Commit:

  ```sh
  git add go/io/copy_test.go go/testing/leakcheck/leakcheck_test.go \
    go/time/rate/rate_test.go go/time/rate/example_test.go \
    go/os/file_rotate_test.go go/net/mux/server_test.go
  git commit -m "test: shorten short-mode timing and stress cases"
  ```

### Task 3: Profile and verify the repository

**Files:**
- No production files

**Interfaces:**
- Produces: measured full and short timing evidence

- [ ] **Step 1:** Profile the full suite module by module with `go test -count=1 ./...`, recording wall time and reporting the 20 slowest modules.
- [ ] **Step 2:** Profile the short suite the same way with `go test -short -count=1 ./...`.
- [ ] **Step 3:** Run repository gates:

  ```sh
  make verify-mods
  make lint
  make test GOTEST_FLAGS=-short
  make test
  git diff --check
  ```

- [ ] **Step 4:** Do not push. Report full and short timings and wait for approval.
