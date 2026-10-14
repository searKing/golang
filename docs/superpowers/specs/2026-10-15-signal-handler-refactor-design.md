# Signal handler refactor

## Goal

Keep the supported purpose of `go/os/signal`: wrap signals registered through
`os/signal.Notify`, dump the native C/C++ stack, and preserve the Go runtime's
original signal behavior.

Make the native handler small, buildable, and async-signal-safe. Remove the
blocking signal-chain implementation while preserving source compatibility for
its public Go APIs.

## Existing behavior and original intent

`Notify` first calls the standard library's `signal.Notify`, then installs a
native wrapper for each non-runtime signal. The wrapper dumps a Boost native
stacktrace and invokes the Go handler that was installed before it.

The signal-chain API was added for fatal-signal handshakes:

1. receive `from`, such as `SIGSEGV`;
2. dump the native stack;
3. raise `to`, such as `SIGUSR1`, to wake a Go goroutine;
4. wait for an acknowledgement signal;
5. optionally sleep;
6. invoke the original Go handler for `from`.

This ordering inherently blocks the faulting thread. Moving the handshake to a
worker would lose the original `siginfo` and machine context needed by the Go
runtime for synchronous faults. There are no signal-chain users in this
repository other than the example, and the README only promises native
stacktraces around `signal.Notify`.

`PreviousStacktrace` reads the most recently dumped native stack in the current
process. The default dump path is process-specific, so it does not by itself
persist a crash across process restarts.

## Compatibility boundary

Keep these supported APIs and behavior:

- `Notify`
- `DumpSignalTo`
- `DumpStacktraceTo`
- `DumpPreviousStacktrace`
- `PreviousStacktrace`
- `Signum`

Keep `SetSigInvokeChain` and `SetSigInvokeChainContext` so existing callers
continue to compile, but mark both deprecated and make them non-blocking no-ops
on every platform. Remove their native ABI and C++ implementation.

Do not add the unfinished `SetSigInvokeChainWithTimeout` API.

Keep the existing no-CGO behavior: standard Go notification works, while native
stacktrace operations are no-ops.

## Handler safety boundary

Configuration code may allocate and lock. The asynchronous handler may only:

1. load lock-free atomically published immutable state;
2. write fixed diagnostic text and integers to the configured file descriptor;
3. call Boost `safe_dump_to`;
4. invoke the previously captured signal disposition.

The handler must not use maps, mutexes, condition variables, `std::function`,
dynamic allocation, streams, Go callbacks, waits, sleeps, or mutable strings.

The dump target is published as an immutable record containing a raw,
null-terminated path and its length. Replaced records are intentionally retained
because a signal handler may still hold a pointer to an older record. Updating a
dump path is rare, making this bounded-by-configuration leak preferable to
unsafe reclamation.

Signal-number tables publish immutable records through lock-free atomic
pointers. Registration is serialized outside the handler.

## Unix design

`Notify` preserves its ordering:

1. register the Go runtime through `os/signal.Notify`;
2. install the native wrapper.

Installing the wrapper:

- queries the current `sigaction`;
- returns successfully if the wrapper is already installed;
- publishes the complete previous disposition before installing the wrapper;
- preserves the previous mask and compatible flags;
- always enables `SA_ONSTACK`;
- serializes concurrent installations.

When invoked, the wrapper dumps the native stack and forwards to the captured
disposition using the original `signum`, `siginfo`, and context. `SIG_IGN`
returns. `SIG_DFL` restores the default disposition and raises the signal again.

Runtime-owned signals such as `SIGURG` and `SIGPROF` remain excluded by the Go
layer.

Calls to standard-library `signal.Stop`, `Reset`, or `Ignore` may replace the
wrapper. A later `Notify` captures the then-current Go disposition and installs
the wrapper again.

## Windows design

Windows keeps its existing limited signal support. Installation performs one
`signal(signum, wrapper)` call and records the returned previous handler,
eliminating the window where the disposition was temporarily `SIG_DFL`.

The wrapper dumps the stack, restores the previous handler, and raises the
signal. Existing documentation continues to require another `Notify` call after
delivery because C `signal` handlers reset.

## Internal structure

- `BaseSignalHandler` owns only dump configuration and dump decoding.
- `SignalTable<T>` remains the lock-free immutable publication primitive.
- Unix and Windows handlers own platform-specific previous dispositions.
- Remove unused native callback registration and all signal-chain state.
- Remove unused CRTP and `enable_shared_from_this` inheritance.
- Keep the C ABI limited to operations used by the Go CGO bridge.

## Tests

Establish the current clean baseline before editing, then verify:

- CGO and non-CGO package builds;
- `Notify` still relays an asynchronous signal;
- native diagnostics and stack dump are produced before forwarding;
- runtime-owned signals are not wrapped;
- repeated and concurrent `Notify` calls do not recurse or race registration;
- changing dump configuration while signals arrive does not corrupt state;
- deprecated chain APIs return immediately and do not delay signal delivery;
- a subprocess receiving a fatal signal still follows the Go runtime's fatal
  behavior;
- Unix tests pass with the race detector;
- Windows CGO sources build in CI.

Run repository formatting, lint, short tests, full `go/os/signal` tests, and
module verification after the focused tests pass.

## Migration and scope

The change is internal except for deprecation notices and the documented
non-blocking behavior of the signal-chain APIs. No new dependencies or public
APIs are introduced.

Work is performed in an isolated worktree based on `origin/master`; the user's
dirty main checkout remains untouched.
