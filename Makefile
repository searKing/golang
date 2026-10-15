# Copyright 2026 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Release modules by module sets defined in versions.yaml, as opentelemetry-go does.
# See RELEASING.md.

SHELL := /bin/bash

GO := GOWORK=off go
TOOLS_DIR := $(CURDIR)/internal/tools
TOOLS_BIN := $(CURDIR)/.tools/bin
MULTIMOD := $(TOOLS_BIN)/multimod
CROSSLINK := $(TOOLS_BIN)/crosslink

ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -not -path './.*' -not -path '*/testdata/*' -exec dirname {} \; | sort)
# Modules of packages, except internal/tools, which requires tools only.
GO_MOD_DIRS := $(filter-out ./internal/tools,$(ALL_GO_MOD_DIRS))
# Local tools tracked by tool directives. Modules providing commands installed by
# 'go install module@version' must not contain replace directives.
LOCAL_TOOL_MODULES := $(shell awk '/^tool \(/ { found = 1; next } found && /^\)/ { found = 0 } found && $$1 ~ /^github.com\/searKing\/golang\/tools\// { print $$1 }' $(TOOLS_DIR)/go.mod)
LOCAL_TOOL_MOD_DIRS := $(patsubst github.com/searKing/golang/%,./%,$(LOCAL_TOOL_MODULES))

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
	@echo "  tools         install tools of internal/tools into .tools/bin, with tools of this repo built from local paths"
	@echo "  generate      run 'go generate ./...' in every module with tools of .tools/bin"
	@echo "  go-mod-tidy   run 'go mod tidy' in every module"
	@echo "  go-get        upgrade DEPS, as arguments of 'go get', in every module requiring them directly,"
	@echo "                or indirectly too if INDIRECT=1, then run 'go mod tidy', prompting for them line by line,"
	@echo "                until an empty line, if DEPS is empty"
	@echo "  test          run 'go test ./...' in every module, with GOTEST_FLAGS"
	@echo "  lint          run 'go vet ./...' in every module, and verify Go files are formatted by gofmt"
	@echo "  crosslink     add replace directives of intra-repo requires and prune unused ones,"
	@echo "                except modules providing tools installed by 'go install module@version'"
	@echo "  verify-mods   verify versions.yaml lists every module, root requires, and publishable tools"
	@echo "  verify-publishable-tools"
	@echo "                verify modules providing installable tools contain no replace directives"
	@echo "  verify-root-requires"
	@echo "                verify the root module requires every module of module sets, against ambiguous imports"
	@echo "  prerelease    update intra-repo requires of MODSET to its version in versions.yaml,"
	@echo "                run 'go mod tidy' and commit to a new branch prerelease_<MODSET>_<VERSION>"
	@echo "  prerelease-clone"
	@echo "                prerelease in a clone of HEAD, then fetch the branch, if the working tree is not clean"
	@echo "  add-tags      tag modules of MODSET at COMMIT by its version in versions.yaml"
	@echo "  delete-tags   delete local tags of MODSET with its version in versions.yaml"
	@echo "  push-tags     push tags of MODSET with its version in versions.yaml to origin"

# Install in internal/tools, for replace directives are relative to the main module, even with -modfile.
.PHONY: tools
tools:
	@echo "Installing tools of internal/tools into $(TOOLS_BIN)" \
		&& cd $(TOOLS_DIR) \
		&& GOBIN=$(TOOLS_BIN) $(GO) install tool

# protoc-gen-go-tag generates by protoc, which is not a Go tool, with PROTOC_INCLUDE for protos of google/protobuf.
.PHONY: generate
generate: tools $(GO_MOD_DIRS:%=generate/%)
generate/%: DIR=$*
generate/%: tools
	@echo "$(GO) generate ./... in $(DIR)" \
		&& cd $(DIR) \
		&& PATH="$(TOOLS_BIN):$$PATH" $(GO) generate ./...

.PHONY: go-mod-tidy
go-mod-tidy: $(ALL_GO_MOD_DIRS:%=go-mod-tidy/%)
go-mod-tidy/%: DIR=$*
go-mod-tidy/%:
	@echo "$(GO) mod tidy in $(DIR)" \
		&& cd $(DIR) \
		&& $(GO) mod tidy

# Upgrade only modules requiring DEPS, directly by default, and leave versions of the others to go mod tidy by MVS.
# INDIRECT=1 upgrades modules requiring them indirectly too, for fixes of vulnerabilities of indirect dependencies..PHONY: go-get
go-get:
	@set -e; \
	deps="$(DEPS)"; indirect="$(INDIRECT)"; \
	if [ -z "$$deps" ]; then \
		echo "Dependencies to upgrade, as arguments of go get, one or more separated by spaces a line, such as:"; \
		echo "  google.golang.org/grpc@v1.70.0 golang.org/x/net"; \
		echo "Enter an empty line to finish."; \
		while read -r -p "DEPS: " line && [ -n "$$line" ]; do deps="$$deps $$line"; done; \
		read -r -p "Upgrade modules requiring them indirectly too? Enter y for INDIRECT=1: " reply || true; \
		if [ "$$reply" = "y" ]; then indirect=1; fi; \
	fi; \
	[ -n "$$deps" ] || { echo "no dependencies to upgrade" >&2; exit 1; }; \
	updated=0; \
	for dir in $(ALL_GO_MOD_DIRS); do \
		requires=$$(awk '/^require \(/ { found = 1; next } found && /^\)/ { found = 0; next } \
			found { print $$1, (/\/\/ indirect/ ? "indirect" : "direct") } \
			/^require [^(]/ { print $$2, (/\/\/ indirect/ ? "indirect" : "direct") }' "$$dir/go.mod"); \
		args=; \
		for dep in $$deps; do \
			if grep -qx "$${dep%@*} direct" <<< "$$requires" || \
				{ [ "$$indirect" = "1" ] && grep -qx "$${dep%@*} indirect" <<< "$$requires"; }; then \
				args="$$args $$dep"; \
			fi; \
		done; \
		[ -n "$$args" ] || continue; \
		echo "$(GO) get$$args && $(GO) mod tidy in $$dir"; \
		(cd "$$dir" && $(GO) get$$args && $(GO) mod tidy); \
		updated=$$((updated + 1)); \
	done; \
	echo "$$updated modules updated. Run 'make verify-mods test' to verify them."

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
crosslink: tools
	@echo "Updating intra-repo dependencies in all go modules" \
		&& $(CROSSLINK) --root=$(CURDIR) --prune
	@$(MAKE) remove-publishable-tool-replaces

# Crosslink adds local replaces for development, but a module providing a command
# installed by 'go install module@version' is rejected if its go.mod has any replace.
.PHONY: remove-publishable-tool-replaces
remove-publishable-tool-replaces:
	@set -e; \
	for dir in $(LOCAL_TOOL_MOD_DIRS); do \
		(cd "$$dir" && $(GO) mod edit -json) \
			| awk '/"Old": \{/ { old = 1; next } old && /"Path":/ { gsub(/[",]/, "", $$2); print $$2; old = 0 }' \
			| while IFS= read -r replaced; do \
				(cd "$$dir" && $(GO) mod edit -dropreplace="$$replaced"); \
			done; \
	done

.PHONY: verify-mods
verify-mods: tools verify-root-requires verify-publishable-tools
	$(MULTIMOD) verify

.PHONY: verify-publishable-tools
verify-publishable-tools:
	@failed=; \
	for dir in $(LOCAL_TOOL_MOD_DIRS); do \
		if (cd "$$dir" && $(GO) mod edit -json) | grep -q '"Replace":'; then \
			echo "$$dir/go.mod contains replace directives, which make 'go install module@version' fail" >&2; \
			failed=1; \
		fi; \
	done; \
	test -z "$$failed"; \
	echo "PASS: Modules providing installable tools contain no replace directives."

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
