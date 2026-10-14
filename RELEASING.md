# Release Process

Modules of this repo are released by module sets defined in [versions.yaml](versions.yaml),
with [multimod](https://github.com/open-telemetry/opentelemetry-go-build-tools/tree/main/multimod),
as [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go/blob/main/RELEASING.md) does.
All modules of a module set are released together with the same version.

Tools are pinned by the module [internal/tools](internal/tools/go.mod), and run by [devops.sh](devops.sh).

## Module sets

Every module of the repo must be listed in exactly one module set, or in `excluded-modules`
if it is never released, such as tests, examples and testdata.
Run `./devops.sh verify` after adding or removing a module.

Modules require each other with replace directives to their paths, so that they are built against
local changes rather than the released versions. Run `./devops.sh crosslink` after requiring a
module of this repo, to add the missing replace directives and prune the unused ones.

## Release

1. Bump the version of the module set in [versions.yaml](versions.yaml), and commit it.
2. Update requires of modules of this repo to the new version, and run `go mod tidy`.

   ```sh
   ./devops.sh prerelease -m stable-v1
   ```

   It commits to a new branch `prerelease_stable-v1_<VERSION>`. Review and merge it.
3. Tag modules of the module set at the merged commit.

   ```sh
   ./devops.sh tag -m stable-v1 -c <COMMIT>
   ```

   Run `./devops.sh untag -m stable-v1` to delete the local tags created by mistake.
4. Push the tags.

   ```sh
   ./devops.sh push -m stable-v1
   ```

Tags are lightweight, in the form `<module-dir>/<VERSION>`, or `<VERSION>` for the root module.
