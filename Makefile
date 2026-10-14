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

MODSET ?= stable-v1
COMMIT ?= HEAD

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
	@echo "  crosslink     add replace directives of intra-repo requires and prune unused ones"
	@echo "  verify-mods   verify versions.yaml lists every module of the repo"
	@echo "  prerelease    update intra-repo requires of MODSET to its version in versions.yaml,"
	@echo "                run 'go mod tidy' and commit to a new branch prerelease_<MODSET>_<VERSION>"
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

.PHONY: crosslink
crosslink:
	@echo "Updating intra-repo dependencies in all go modules" \
		&& $(CROSSLINK) --root=$(CURDIR) --prune

.PHONY: verify-mods
verify-mods:
	$(MULTIMOD) verify

# multimod requires a clean working tree without untracked files, so prerelease in a clone of HEAD.
.PHONY: prerelease
prerelease:
	@git diff --quiet HEAD || { echo "commit or stash changes of tracked files first" >&2; exit 1; }
	@set -e; \
	head=$$(git rev-parse HEAD); \
	tmp=$$(mktemp -d); \
	trap 'rm -rf "$$tmp"' EXIT; \
	git clone -q --no-hardlinks "$(CURDIR)" "$$tmp/repo"; \
	cd "$$tmp/repo"; \
	git checkout -q -B release_base "$$head"; \
	$(MULTIMOD) verify; \
	$(MULTIMOD) prerelease -m $(MODSET); \
	cd "$(CURDIR)"; \
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
