# Copyright 2026 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Release modules by module sets defined in versions.yaml, as opentelemetry-go does.
# See RELEASING.md.

SHELL := /bin/bash

GO := GOWORK=off go
TOOLS_MOD := $(CURDIR)/internal/tools/go.mod
MULTIMOD := $(GO) tool -modfile=$(TOOLS_MOD) multimod
CROSSLINK := $(GO) tool -modfile=$(TOOLS_MOD) crosslink

ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -not -path './.*' -not -path '*/testdata/*' -exec dirname {} \; | sort)
# Modules of packages, except internal/tools, which requires tools only.
GO_MOD_DIRS := $(filter-out ./internal/tools,$(ALL_GO_MOD_DIRS))

MODSET ?= stable-v1
COMMIT ?= HEAD

# Modules of all module sets in versions.yaml.
MODSET_MODULES = $(shell awk '/^[^ ]/ { found = ($$0 == "module-sets:") } found && $$1 == "-" { print $$2 }' versions.yaml)

# Tags of the module set MODSET with its version in versions.yaml, as <module-dir>/<version>, or <version> for the root module.
MODSET_TAGS = $(shell awk -v set='  $(MODSET):' ' \
	$$0 == set { found = 1; next } \
	found && /^[^ ]|^  [^ ]/ { exit } \
	found && $$1 == "version:" { version = $$2 } \
	found && $$1 == "-" { mods[n++] = $$2 } \
	END { for (i = 0; i < n; i++) { dir = mods[i]; sub("^github.com/searKing/golang/?", "", dir); print (dir == "" ? version : dir "/" version) } }' \
	versions.yaml)

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "Usage: make TARGET [MODSET=stable-v1] [COMMIT=HEAD]"
	@echo
	@echo "Targets:"
	@echo "  go-mod-tidy   run 'go mod tidy' in every module"
	@echo "  test          run 'go test ./...' in every module, with GOTEST_FLAGS"
	@echo "  lint          run 'go vet ./...' in every module, and verify Go files are formatted by gofmt"
	@echo "  crosslink     add replace directives of intra-repo requires and prune unused ones"
	@echo "  verify-mods   verify versions.yaml lists every module of the repo, and verify-root-requires"
	@echo "  verify-root-requires"
	@echo "                verify the root module requires every module of module sets, against ambiguous imports"
	@echo "  prerelease    update intra-repo requires of MODSET to its version in versions.yaml,"
	@echo "                run 'go mod tidy' and commit to a new branch prerelease_<MODSET>_<VERSION>"
	@echo "  prerelease-clone"
	@echo "                prerelease in a clone of HEAD, then fetch the branch, if the working tree is not clean"
	@echo "  add-tags      tag modules of MODSET at COMMIT by its version in versions.yaml"
	@echo "  delete-tags   delete local tags of MODSET with its version in versions.yaml"
	@echo "  push-tags     push tags of MODSET with its version in versions.yaml to origin"

.PHONY: go-mod-tidy
go-mod-tidy: $(ALL_GO_MOD_DIRS:%=go-mod-tidy/%)
go-mod-tidy/%: DIR=$*
go-mod-tidy/%:
	@echo "$(GO) mod tidy in $(DIR)" \
		&& cd $(DIR) \
		&& $(GO) mod tidy

.PHONY: test
test: $(GO_MOD_DIRS:%=test/%)
test/%: DIR=$*
test/%:
	@echo "$(GO) test $(GOTEST_FLAGS) ./... in $(DIR)" \
		&& cd $(DIR) \
		&& $(GO) test $(GOTEST_FLAGS) ./...

.PHONY: lint
lint: $(GO_MOD_DIRS:%=vet/%) fmt-check
vet/%: DIR=$*
vet/%:
	@echo "$(GO) vet $(VET_FLAGS) ./... in $(DIR)" \
		&& cd $(DIR) \
		&& $(GO) vet $(VET_FLAGS) ./...

# Legacy generated protos copied from protobuf for tests use unkeyed fields of protobuf structs.
vet/./third_party/github.com/spf13/viper: VET_FLAGS = -composites=false

# Go files of testdata may be malformed on purpose.
.PHONY: fmt-check
fmt-check:
	@unformatted=$$(git ls-files -z '*.go' ':!:**/testdata/**' | xargs -0 gofmt -l); \
	if [ -n "$$unformatted" ]; then \
		echo "Go files not formatted by gofmt, run 'gofmt -w' on them:" >&2; \
		printf '  %s\n' $$unformatted >&2; \
		exit 1; \
	fi

.PHONY: crosslink
crosslink:
	@echo "Updating intra-repo dependencies in all go modules" \
		&& $(CROSSLINK) --root=$(CURDIR) --prune

.PHONY: verify-mods
verify-mods: verify-root-requires
	$(MULTIMOD) verify

# The root module requires every module of module sets, so that a package moved into a module of its own
# is not found in both the module and its parent module of an old version, as an ambiguous import.
# go mod tidy keeps a require only if its package is imported, by doc.go, or a tool directive for main packages.
.PHONY: verify-root-requires
verify-root-requires:
	@set -e; \
	required=$$(awk '/^require \(/ { found = 1; next } found && /^\)/ { found = 0 } found { print $$1 } /^require [^(]/ { print $$2 }' go.mod); \
	missing=$$(for mod in $(MODSET_MODULES); do \
		[ "$$mod" = "github.com/searKing/golang" ] || grep -qx "$$mod" <<< "$$required" || echo "$$mod"; \
	done); \
	if [ -n "$$missing" ]; then \
		echo "modules not required by the root module, import them in doc.go, or by tool directives for main packages:" >&2; \
		printf '  %s\n' $$missing >&2; \
		exit 1; \
	fi; \
	$(GO) mod tidy -diff >/dev/null || { echo "go.mod of the root module is not tidy, run 'make go-mod-tidy/.'" >&2; exit 1; }; \
	echo "PASS: Root module requires all modules of module sets."

.PHONY: prerelease
prerelease: verify-mods
	$(MULTIMOD) prerelease -m $(MODSET)

# prerelease in a clone of HEAD, for multimod requires a clean working tree without untracked files.
.PHONY: prerelease-clone
prerelease-clone:
	@git diff --quiet HEAD || { echo "commit or stash changes of tracked files first" >&2; exit 1; }
	@set -e; \
	head=$$(git rev-parse HEAD); \
	tmp=$$(mktemp -d); \
	trap 'rm -rf "$$tmp"' EXIT; \
	git clone -q --no-hardlinks "$(CURDIR)" "$$tmp/repo"; \
	git -C "$$tmp/repo" checkout -q -B release_base "$$head"; \
	$(MAKE) -C "$$tmp/repo" prerelease MODSET=$(MODSET); \
	git fetch -q "$$tmp/repo" 'refs/heads/prerelease_*:refs/heads/prerelease_*'; \
	git branch --list 'prerelease_*' --contains "$$head"

# Lightweight tags as the released ones, for 'multimod tag' signs tags by gpg always.
.PHONY: add-tags
add-tags:
	@set -e; \
	tags="$(MODSET_TAGS)"; \
	[ -n "$$tags" ] || { echo "module set $(MODSET) not found in versions.yaml" >&2; exit 1; }; \
	commit=$$(git rev-parse --verify "$(COMMIT)^{commit}"); \
	for tag in $$tags; do \
		if git rev-parse -q --verify "refs/tags/$$tag" >/dev/null; then echo "tag $$tag already exists" >&2; exit 1; fi; \
	done; \
	for tag in $$tags; do git tag "$$tag" "$$commit"; echo "$$tag"; done; \
	echo "$$(wc -w <<< "$$tags" | tr -d ' ') tags created at $$commit. Run 'make push-tags MODSET=$(MODSET)' to push them."

.PHONY: delete-tags
delete-tags:
	@for tag in $(MODSET_TAGS); do \
		if git rev-parse -q --verify "refs/tags/$$tag" >/dev/null; then git tag -d "$$tag"; fi; \
	done

.PHONY: push-tags
push-tags:
	@set -e; \
	tags="$(MODSET_TAGS)"; \
	[ -n "$$tags" ] || { echo "module set $(MODSET) not found in versions.yaml" >&2; exit 1; }; \
	printf '%s\n' $$tags; \
	read -r -p "Push the above $$(wc -w <<< "$$tags" | tr -d ' ') tags to origin? Enter y to continue: " reply; \
	if [ "$$reply" = "y" ]; then \
		git push origin $$(printf 'refs/tags/%s ' $$tags); \
	else \
		echo "Tag push cancelled."; \
	fi
