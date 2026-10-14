# Signal handler refactor final fix report

Date: 2026-10-15

## Status

All four final-review findings were addressed in one fix wave. No C signal-chain
ABI was restored, and no changes were pushed.

## Finding resolutions

1. Windows re-registration state
   - `SignalHandler::SetSig` now leaves the previously published handler record
     visible while `signal(signum, HandleSignal)` runs.
   - A successful installation release-publishes the newly returned previous
     handler through `SignalTable::Store`.
   - `HandleSignal` and `SIG_ERR` returns preserve the old record.
   - The standalone state helper covers an empty first-install window, an old
     record visible during re-registration, successful replacement, wrapper
     return, and `SIG_ERR`.

2. Public `go/os/signal/cgo` compatibility
   - Restored deprecated no-op
     `SetSigInvokeChain(from, to, wait, sleepInSeconds int)`.
   - Restored deprecated no-op
     `SetSigInvokeChainCancelable(...) func()` with a non-nil empty cancel
     function.
   - Added direct package API coverage.
   - Confirmed `CGO_SetSigInvokeChain`, `CGO_NewSigInvokeChain`, and
     `CGO_CancelSigInvokeChain` remain absent from production Go/C/C++ code.

3. Stacktrace documentation
   - The design and implementation plan now describe `PreviousStacktrace` as
     reading the most recent dump at the currently configured path.
   - They explicitly note that a caller-selected persistent path may contain a
     dump from a previous process.

4. Markdown diff hygiene
   - Replaced space-before-tab indentation in fenced Go examples without
     changing the examples' semantics.

## Test evidence

Regression tests were observed failing before implementation:

- Windows helper: `observed_record == saved` failed because the table was
  cleared during re-registration.
- `go test -count=1 ./os/signal/cgo`: failed to compile because both public
  compatibility functions were undefined.

Fresh passing verification after implementation:

- `GOWORK=off CGO_ENABLED=1 go test -count=1 ./os/signal ./os/signal/cgo`
- `GOWORK=off CGO_ENABLED=0 go test -count=1 ./os/signal`
- `GOWORK=off CGO_ENABLED=1 go test -race -count=1 ./os/signal ./os/signal/cgo`
- `GOWORK=off CGO_ENABLED=1 go vet ./os/signal ./os/signal/cgo`
- `GOWORK=off CGO_ENABLED=0 go vet ./os/signal`
- `GOWORK=off CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go test -c ./os/signal`
- Standalone Windows state helper compiled with
  `clang++ -std=c++11 -Wall -Wextra -Werror` and ran successfully.
- `git diff --check $(git merge-base origin/master HEAD)` passed for the full
  working diff before commit.

## Self-review

- The Windows table is untouched until `signal` returns, so an in-flight
  wrapper sees the old immutable record. First installation naturally remains
  empty until publication.
- Publication still uses the existing release-store/acquire-load
  `SignalTable` contract.
- Compatibility functions are inert and do not call removed native symbols.
- Scope is limited to the four findings and their regression evidence.

## Concerns

- The Darwin environment has no MinGW cross compiler and no `clang-format`
  executable. Windows CGO was therefore checked through the strict standalone
  C++ helper; the available Windows Go cross-check used `CGO_ENABLED=0`.
