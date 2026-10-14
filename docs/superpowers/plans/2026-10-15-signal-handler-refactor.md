# Signal Handler Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reduce `go/os/signal` to an async-signal-safe native stack dumper that faithfully forwards to Go's original signal handler.

**Architecture:** The Go layer keeps the supported public API and installs the native wrapper after `os/signal.Notify`. Native handlers read immutable state through lock-free atomic pointers, dump the native stack, and immediately forward the original disposition; deprecated signal-chain APIs remain compile-compatible no-ops.

**Tech Stack:** Go, CGO, C++11 atomics, POSIX `sigaction`, Windows C `signal`, Boost.Stacktrace

**Spec:** `docs/superpowers/specs/2026-10-15-signal-handler-refactor-design.md`

## Global Constraints

- Preserve `Notify`, stack dump/read, `Signum`, and no-CGO behavior.
- Keep the deprecated signal-chain APIs in both importable Go packages source-compatible as non-blocking no-ops.
- Do not add `SetSigInvokeChainWithTimeout`.
- Never allocate, lock, wait, sleep, invoke Go, or access mutable C++ state in an asynchronous handler.
- Unix wrappers must use `SA_ONSTACK` and forward synchronous faults with their original `siginfo` and context.
- Do not add dependencies.
- Work only in `/Users/chenhaixin/workspace/src/github.com/searKing/golang-signal-refactor`.

---

### Task 1: Deprecate blocking signal-chain behavior at the Go boundary

**Files:**
- Modify: `go/os/signal/signal.go`
- Modify: `go/os/signal/signal_cgo.go`
- Modify: `go/os/signal/signal_not_cgo.go`
- Modify: `go/os/signal/cgo/signal.go`
- Modify: `go/os/signal/signal_test.go`

**Interfaces:**
- Produces: compile-compatible `SetSigInvokeChain(...)` and `SetSigInvokeChainContext(...)` that return immediately.
- Removes: private Go and CGO chain bridge functions, which Task 2 removes from the C ABI.

- [ ] **Step 1: Replace the old timing assertion with subprocess coverage for both deprecated APIs**

  Add a parent/child test to `signal_test.go`. Each child configures one API with a missing acknowledgement signal, sends `SIGWINCH`, and must still receive it immediately:

  ```go
  func TestDeprecatedSigInvokeChainDoesNotBlock(t *testing.T) {
    if mode := os.Getenv("SEARKING_SIGNAL_CHAIN_TEST"); mode != "" {
      switch mode {
      case "plain":
        signal_.SetSigInvokeChain(syscall.SIGUSR1, syscall.SIGTTOU, 10, syscall.SIGWINCH)
      case "context":
        signal_.SetSigInvokeChainContext(context.Background(), syscall.SIGUSR1, syscall.SIGTTOU, 10*time.Second, syscall.SIGWINCH)
      default:
        t.Fatalf("unknown child mode %q", mode)
      }
      c := make(chan os.Signal, 1)
      signal_.Notify(c, syscall.SIGWINCH)
      if err := syscall.Kill(syscall.Getpid(), syscall.SIGWINCH); err != nil {
        t.Fatal(err)
      }
      select {
      case <-c:
      case <-time.After(500 * time.Millisecond):
        t.Fatal("deprecated signal chain blocked delivery")
      }
      return
    }

    for _, mode := range []string{"plain", "context"} {
      t.Run(mode, func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestDeprecatedSigInvokeChainDoesNotBlock$")
        cmd.Env = append(os.Environ(), "SEARKING_SIGNAL_CHAIN_TEST="+mode)
        if output, err := cmd.CombinedOutput(); err != nil {
          t.Fatalf("child failed: %v\n%s", err, output)
        }
      })
    }
  }
  ```

- [ ] **Step 2: Run the new test and verify the old implementation fails**

  Run:

  ```sh
  cd go
  GOWORK=off go test -count=1 -run '^TestDeprecatedSigInvokeChainDoesNotBlock$' ./os/signal
  ```

  Expected: both child processes exceed their internal delivery deadline or the parent's two-second timeout because the old implementation waits for `SIGTTOU`.

- [ ] **Step 3: Make the public chain APIs deprecated no-ops**

  In `signal.go`, preserve signatures and replace bodies:

  ```go
  // SetSigInvokeChain is retained for source compatibility.
  //
  // Deprecated: signal handlers no longer block, raise auxiliary signals, or
  // sleep before forwarding to the Go runtime.
  func SetSigInvokeChain(to os.Signal, wait os.Signal, sleepInSeconds int, froms ...os.Signal) {
  }

  // SetSigInvokeChainContext is retained for source compatibility.
  //
  // Deprecated: signal handlers no longer block, raise auxiliary signals, or
  // sleep before forwarding to the Go runtime.
  func SetSigInvokeChainContext(ctx context.Context, to os.Signal, wait os.Signal, sleep time.Duration, froms ...os.Signal) {
  }
  ```

  Remove `setSigInvokeChain` and `setSigInvokeChainCancelable` from `signal_cgo.go` and `signal_not_cgo.go`. Keep deprecated no-op `SetSigInvokeChain` and `SetSigInvokeChainCancelable` functions in the public `cgo` package without restoring the C ABI.

- [ ] **Step 4: Run Go-layer tests in both build modes**

  ```sh
  cd go
  GOWORK=off go test -count=1 ./os/signal
  CGO_ENABLED=0 GOWORK=off go test -count=1 ./os/signal
  ```

  Expected: both pass; CGO still links because native chain symbols remain until Task 2.

- [ ] **Step 5: Commit the API migration**

  ```sh
  git add go/os/signal/signal.go go/os/signal/signal_cgo.go \
    go/os/signal/signal_not_cgo.go go/os/signal/cgo/signal.go \
    go/os/signal/signal_test.go
  git commit -m "refactor(os/signal): deprecate blocking signal chains"
  ```

### Task 2: Reduce native shared state to immutable dump configuration

**Files:**
- Modify: `go/os/signal/cgo/base_signal_handler.hpp`
- Modify: `go/os/signal/cgo/signal.cgo.h`
- Modify: `go/os/signal/cgo/signal.cgo.cpp`
- Modify: `go/os/signal/cgo/signal_handler_unix.hpp`
- Modify: `go/os/signal/cgo/signal_handler_unix.cpp`
- Modify: `go/os/signal/cgo/signal_handler_windows.hpp`
- Modify: `go/os/signal/cgo/signal_handler_windows.cpp`

**Interfaces:**
- Produces: non-template `BaseSignalHandler`, `SignalTable<T>`, and immutable `DumpPath`.
- Removes: native signal-chain ABI/state and unused native `RegisterOnSignal` callbacks.

- [ ] **Step 1: Add a dump-path replacement stress test**

  Add a Unix CGO test which alternates two paths while delivering `SIGUSR2`; keep the signal count modest and accept either complete path in diagnostics:

  ```go
  func TestDumpStacktracePathConcurrentUpdate(t *testing.T) {
    dir := t.TempDir()
    paths := []string{filepath.Join(dir, "a.dump"), filepath.Join(dir, "b.dump")}
    c := make(chan os.Signal, 32)
    signal_.Notify(c, syscall.SIGUSR2)
    defer signal.Stop(c)

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
      defer wg.Done()
      for i := 0; i < 32; i++ {
        signal_.DumpStacktraceTo(paths[i%len(paths)])
      }
    }()
    for i := 0; i < 32; i++ {
      if err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR2); err != nil {
        t.Fatal(err)
      }
      select {
      case <-c:
      case <-time.After(5 * time.Second):
        t.Fatal("SIGUSR2 not relayed")
      }
    }
    wg.Wait()
  }
  ```

- [ ] **Step 2: Verify the stress test passes before structural changes**

  ```sh
  cd go
  GOWORK=off go test -count=5 -run '^TestDumpStacktracePathConcurrentUpdate$' ./os/signal
  ```

  Expected: PASS. This is a behavior-preservation test, not a red test.

- [ ] **Step 3: Replace mutable handler-visible strings with immutable raw records**

  In `base_signal_handler.hpp`:

  ```cpp
  struct DumpPath {
    explicit DumpPath(const std::string &name)
        : length(name.size()), path(new char[length + 1]) {
      memcpy(path, name.data(), length);
      path[length] = '\0';
    }

    const size_t length;
    char *const path;
  };

  class BaseSignalHandler : private boost::noncopyable {
   protected:
    BaseSignalHandler()
        : signal_dump_to_fd_(-1), stacktrace_dump_to_file_(nullptr) {}

   public:
    void SetStacktraceDumpToFile(const std::string &name) {
      stacktrace_dump_to_file_.store(new DumpPath(name),
                                     std::memory_order_release);
    }
    // Keep SetSignalDumpToFd, WriteSignalStacktrace,
    // DumpPreviousStacktrace and PreviousStacktrace.

   protected:
    const DumpPath *StacktraceDumpToFile() const {
      return stacktrace_dump_to_file_.load(std::memory_order_acquire);
    }

   private:
    std::atomic<int> signal_dump_to_fd_;
    std::atomic<const DumpPath *> stacktrace_dump_to_file_;
  };
  ```

  Make handler code use only `file->path` and `file->length`. Normal-context decoding opens `file->path`.

  Remove `SigChain`, `RunSigChain`, cancellation, polling, CRTP, `enable_shared_from_this`, and no-longer-used headers. Keep `SignalTable<T>::Load` and `Store`; remove `RemoveIf` and `StoreIfEmpty` only after platform code no longer uses them.

- [ ] **Step 4: Remove chain symbols from the C ABI**

  Delete `CGO_SetSigInvokeChain`, `CGO_NewSigInvokeChain`, and `CGO_CancelSigInvokeChain` from `signal.cgo.h` and `signal.cgo.cpp`. Remove `<stdint.h>` if unused.

- [ ] **Step 5: Remove unused native callbacks and chain execution**

  Make both platform classes inherit `BaseSignalHandler` without a template argument. Remove `OnSignal`, `on_signal_`, `RegisterOnSignal`, `std::function`, and their invocation. Remove `RunSigChain` calls from both platform operators.

- [ ] **Step 6: Format C++ and run focused tests**

  ```sh
  clang-format -i \
    go/os/signal/cgo/base_signal_handler.hpp \
    go/os/signal/cgo/signal.cgo.h \
    go/os/signal/cgo/signal.cgo.cpp \
    go/os/signal/cgo/signal_handler_unix.hpp \
    go/os/signal/cgo/signal_handler_unix.cpp \
    go/os/signal/cgo/signal_handler_windows.hpp \
    go/os/signal/cgo/signal_handler_windows.cpp
  cd go
  GOWORK=off go test -count=3 ./os/signal
  CGO_ENABLED=0 GOWORK=off go test -count=1 ./os/signal
  ```

  Expected: PASS.

- [ ] **Step 7: Commit the native-state simplification**

  ```sh
  git add go/os/signal/cgo
  git commit -m "refactor(os/signal): minimize native handler state"
  ```

### Task 3: Make Unix installation idempotent and race-free

**Files:**
- Modify: `go/os/signal/cgo/signal_handler_unix.hpp`
- Modify: `go/os/signal/cgo/signal_handler_unix.cpp`
- Modify: `go/os/signal/signal_test.go`

**Interfaces:**
- Produces: atomically published immutable `SignalDisposition` records and serialized handler installation.
- Consumes: `SignalTable<T>` and `BaseSignalHandler` from Task 2.

- [ ] **Step 1: Add repeated concurrent registration coverage**

  ```go
  func TestNotifyConcurrent(t *testing.T) {
    const workers = 8
    channels := make([]chan os.Signal, workers)
    var wg sync.WaitGroup
    for i := range channels {
      channels[i] = make(chan os.Signal, 1)
      wg.Add(1)
      go func(c chan os.Signal) {
        defer wg.Done()
        signal_.Notify(c, syscall.SIGUSR2)
      }(channels[i])
    }
    wg.Wait()
    for _, c := range channels {
      defer signal.Stop(c)
    }
    if err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR2); err != nil {
      t.Fatal(err)
    }
    for i, c := range channels {
      select {
      case <-c:
      case <-time.After(5 * time.Second):
        t.Fatalf("channel %d did not receive SIGUSR2", i)
      }
    }
  }
  ```

- [ ] **Step 2: Run registration coverage repeatedly**

  ```sh
  cd go
  GOWORK=off go test -count=20 -run '^TestNotifyConcurrent$' ./os/signal
  ```

  Expected: PASS on the baseline; this guards the refactor.

- [ ] **Step 3: Store the complete prior Unix disposition**

  In `signal_handler_unix.hpp`:

  ```cpp
  struct SignalDisposition {
    explicit SignalDisposition(const struct sigaction &action)
        : action(action) {}
    const struct sigaction action;
  };

  class SignalHandler : public BaseSignalHandler {
   private:
    static void HandleSignal(int signum, siginfo_t *info, void *context);
    void InvokeGoSignalHandler(int signum, siginfo_t *info, void *context);

    std::mutex install_mutex_;
    SignalTable<SignalDisposition> go_registered_handlers_;
  };
  ```

  Include `<mutex>`. Delete the old pair typedefs and callback declarations.

- [ ] **Step 4: Serialize installation and avoid wrapping the wrapper**

  Implement `SetSig` around the singleton:

  ```cpp
  void SignalHandler::HandleSignal(int signum, siginfo_t *info, void *context) {
    GetInstance()(signum, info, context);
  }

  int SignalHandler::SetSig(int signum) {
    SignalHandler &instance = GetInstance();
    std::lock_guard<std::mutex> lock(instance.install_mutex_);

    struct sigaction previous;
    memset(&previous, 0, sizeof(previous));
    if (sigaction(signum, nullptr, &previous) != 0) {
      return -1;
    }
    if ((previous.sa_flags & SA_SIGINFO) != 0 &&
        previous.sa_sigaction == &SignalHandler::HandleSignal) {
      return 0;
    }

    instance.go_registered_handlers_.Store(
        signum, new SignalDisposition(previous));

    struct sigaction wrapper = previous;
    wrapper.sa_flags |= SA_SIGINFO | SA_ONSTACK | SA_RESTART;
    wrapper.sa_sigaction = &SignalHandler::HandleSignal;
    return sigaction(signum, &wrapper, nullptr);
  }
  ```

  Preserve the existing alternate-stack check before installation. It may allocate only in `SetSig`, never in the handler.

- [ ] **Step 5: Forward the complete prior disposition**

  Load `SignalDisposition` once. Invoke `sa_sigaction` when `SA_SIGINFO` is set, return for `SIG_IGN`, and for `SIG_DFL` restore the captured action before `raise(signum)`. Otherwise invoke `sa_handler(signum)`. Guard against accidentally forwarding to `HandleSignal`.

- [ ] **Step 6: Run focused tests and race detector**

  ```sh
  cd go
  GOWORK=off go test -count=10 ./os/signal
  GOWORK=off go test -race -count=3 ./os/signal
  ```

  Expected: PASS without recursion, timeout, or race reports.

- [ ] **Step 7: Commit Unix registration**

  ```sh
  git add go/os/signal/cgo/signal_handler_unix.hpp \
    go/os/signal/cgo/signal_handler_unix.cpp \
    go/os/signal/signal_test.go
  git commit -m "fix(os/signal): serialize unix handler installation"
  ```

### Task 4: Remove the Windows default-disposition window

**Files:**
- Modify: `go/os/signal/cgo/signal_handler_windows.hpp`
- Modify: `go/os/signal/cgo/signal_handler_windows.cpp`

**Interfaces:**
- Produces: one-call Windows handler replacement with atomically published previous handlers.
- Consumes: `SignalTable<T>` and `BaseSignalHandler` from Task 2.

- [ ] **Step 1: Simplify the Windows class**

  Remove callback state and duplicate deleted-copy declarations already provided by `BaseSignalHandler`. Add an installation mutex:

  ```cpp
  class SignalHandler : public BaseSignalHandler {
   private:
    static void HandleSignal(int signum);
    void InvokeGoSignalHandler(int signum);

    std::mutex install_mutex_;
    SignalTable<SignalHandlerSignalHandler> go_registered_handlers_;
  };
  ```

- [ ] **Step 2: Replace the two-call default window**

  ```cpp
  void SignalHandler::HandleSignal(int signum) { GetInstance()(signum); }

  int SignalHandler::SetSig(int signum) {
    SignalHandler &instance = GetInstance();
    std::lock_guard<std::mutex> lock(instance.install_mutex_);
    SignalHandlerSignalHandler previous = ::signal(signum, &HandleSignal);
    if (previous == SIG_ERR) {
      return -1;
    }
    if (previous != &HandleSignal) {
      instance.go_registered_handlers_.Store(
          signum, new SignalHandlerSignalHandler(previous));
    }
    return 0;
  }
  ```

  `InvokeGoSignalHandler` restores the stored handler and raises the signal. If no record has been published yet, restore `SIG_DFL` before raising; this covers the narrow interval between `signal` returning and atomic publication without dereferencing null.

- [ ] **Step 3: Format and perform available cross-platform checks**

  ```sh
  clang-format -i \
    go/os/signal/cgo/signal_handler_windows.hpp \
    go/os/signal/cgo/signal_handler_windows.cpp
  cd go
  CGO_ENABLED=0 GOOS=windows GOWORK=off go test -c ./os/signal
  rm -f signal.test.exe
  ```

  Expected: the non-CGO Windows test binary compiles locally. Windows CGO compilation and execution is required from GitHub Actions because the local Darwin toolchain cannot link Windows CGO.

- [ ] **Step 4: Commit Windows installation**

  ```sh
  git add go/os/signal/cgo/signal_handler_windows.hpp \
    go/os/signal/cgo/signal_handler_windows.cpp
  git commit -m "fix(os/signal): replace windows handler atomically"
  ```

### Task 5: Update examples, docs, and integration verification

**Files:**
- Modify: `go/os/signal/example_test.go`
- Modify: `go/os/signal/signal.go`
- Modify: `README.md`
- Test: `go/os/signal/signal_test.go`

**Interfaces:**
- Produces: documentation matching the supported native-stack-dump behavior.

- [ ] **Step 1: Replace the example's chain handshake**

  Remove `SetSigInvokeChain`, `SIGUSR1`, `SIGUSR2`, `LockOSThread`, and acknowledgement logic. Keep a minimal `SIGINT` example that calls `signal_.Notify`, waits for the channel, and prints `Got signal: interrupt`.

- [ ] **Step 2: Correct stacktrace API documentation**

  Document `PreviousStacktrace` as the most recent native stack dump at the currently configured path. Explain that a caller-selected persistent path may contain a dump from a previous process; the default path is process-specific.

  Keep the README promise narrow:

  ```markdown
  * [signal](https://pkg.go.dev/github.com/searKing/golang/go/os/signal#Notify) — wraps `signal.Notify`
    to dump the native C/C++ stack before preserving the Go runtime's signal behavior.
  ```

- [ ] **Step 3: Add subprocess fatal-forwarding coverage**

  Add a Unix CGO child mode that calls `signal_.Notify` for `SIGSEGV`, then performs an invalid memory access through `unsafe.Pointer(uintptr(0x123))` and must not reach a sentinel print. The parent asserts a non-zero exit and absence of the sentinel. Do not send `SIGSEGV` with `kill`: that is asynchronous delivery and does not carry a real faulting PC, so `os/signal.Notify` may legitimately deliver it to the channel instead of exercising the runtime fault path.

- [ ] **Step 4: Run focused verification**

  ```sh
  gofmt -w go/os/signal/*.go
  cd go
  GOWORK=off go test -count=10 ./os/signal
  GOWORK=off go test -race -count=3 ./os/signal
  CGO_ENABLED=0 GOWORK=off go test -count=3 ./os/signal
  GOWORK=off go vet ./os/signal
  ```

  Expected: all pass.

- [ ] **Step 5: Run repository gates**

  ```sh
  make verify-mods
  make lint
  make test GOTEST_FLAGS=-short
  git diff --check
  git status --short
  ```

  Expected: all commands pass; only intended signal/docs changes remain committed.

- [ ] **Step 6: Commit docs and integration tests**

  ```sh
  git add README.md go/os/signal
  git commit -m "docs(os/signal): describe native stack forwarding"
  ```

- [ ] **Step 7: Push and require all platform CI before merge**

  Push only after local verification and user approval. Confirm Linux, macOS, and Windows short tests pass, with Windows CGO providing the required native build/link coverage.
