# CI test optimization

## Goal

Keep cross-platform CI useful while moving stress and long-running coverage to local full test runs.
Reduce pull-request feedback time without hiding correctness failures.

## Test policy

- All CI operating systems run `go test -short`.
- CI runs two independent modules at a time. Modules are isolated by `GOWORK=off`, and two-way
  parallelism reduces orchestration overhead without heavily oversubscribing hosted runners.
- Local `make test` remains the complete suite, including stress and wall-clock tests.
- Tests that already scale work under `testing.Short` keep their full local coverage.
- Tests that use real sleeps or wall-clock boundaries should either:
  - use deterministic synchronization instead of sleeping, when that tests the same behavior; or
  - skip in short mode when the wall-clock behavior itself is under test.
- Run and profile the complete local suite after each optimization. Only change tests with measured
  cost, and preserve their assertions.

## Pull request scope

PR #79 will contain:

1. The publishable-tool fix already on the branch.
2. All CI test jobs using `-short`.
3. Two-way module parallelism in CI.
4. Focused optimizations for measured slow tests.

It will not contain incremental module selection.

## Follow-up incremental CI

A separate change will compute changed Go modules from the pull-request base, add their recursive
reverse dependencies, and run only that set on pull requests. Changes to global build and module
metadata will fall back to all modules. Pushes to master will continue to test every module.

Generate and verify-mods remain full checks because they are fast and their dependency boundaries
are broader than source module boundaries. Lint remains full until changed-module selection is proven.

## Verification

- Run the complete local suite and record module and individual-test timings.
- Run the complete short suite.
- Run lint and verify-mods.
- Confirm all tests pass and compare wall time before pushing.
