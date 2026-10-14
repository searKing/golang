# Release Process

Modules of this repo are released by module sets defined in [versions.yaml](versions.yaml),
with [multimod](https://github.com/open-telemetry/opentelemetry-go-build-tools/tree/main/multimod),
as [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go/blob/main/RELEASING.md) does.
All modules of a module set are released together with the same version.

Tools are tracked by tool directives of the module [internal/tools](internal/tools/go.mod), with tools of this repo
replaced by their local paths, installed into `.tools/bin` by `make tools`, and run by targets of [Makefile](Makefile).
Run `make help` for all targets.

`make generate` runs `go generate ./...` in every module with these tools, and CI fails if generated files are out of date.
protoc-gen-go-tag also requires protoc, of the version in the CI workflow, as its version is written into generated files.
To track a new tool, run `cd internal/tools && go get -tool <TOOL>`, and `make crosslink` if it is of this repo.

## Module sets

Every module of the repo must be listed in exactly one module set, or in `excluded-modules`
if it is never released, such as tests, examples and testdata.
Run `make verify-mods` after adding or removing a module.

The root module requires every module of module sets, so that a package moved from a module into a
module of its own is not found in both of them, which fails as an ambiguous import, as the root module
selects versions of them after the move. As go mod tidy keeps a require only if its package is imported,
import the package of a new module in [doc.go](doc.go), or require it by a tool directive if it is main:

```sh
go get -tool github.com/searKing/golang/tools/<TOOL>
```

`make verify-mods` checks it by `make verify-root-requires`.

Modules require each other with replace directives to their paths, so that they are built against
local changes rather than the released versions. Run `make crosslink` after requiring a
module of this repo, to add the missing replace directives and prune the unused ones.

Modules providing tools installed by `go install module@version` are the exception: the Go command
rejects such modules if their go.mod contains any replace directive. `make crosslink` removes replaces
from the local tool modules tracked by internal/tools/go.mod, and `make verify-mods` checks this.

## Release

1. Bump the version of the module set in [versions.yaml](versions.yaml), and commit it.
2. Update requires of modules of this repo to the new version, and run `go mod tidy`.

   ```sh
   make prerelease MODSET=stable-v1
   ```

   It commits to a new branch `prerelease_stable-v1_<VERSION>`. Review and merge it.
   multimod requires a clean working tree without untracked files, run `make prerelease-clone MODSET=stable-v1`
   instead to prerelease in a clone of HEAD and fetch the branch back.
3. Tag modules of the module set at the merged commit.

   ```sh
   make add-tags MODSET=stable-v1 COMMIT=<COMMIT>
   ```

   Run `make delete-tags MODSET=stable-v1` to delete the local tags created by mistake.
4. Push the tags.

   ```sh
   make push-tags MODSET=stable-v1
   ```

Tags are lightweight, in the form `<module-dir>/<VERSION>`, or `<VERSION>` for the root module.
