# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You
# may not use this file except in compliance with the License. A copy of
# the License is located at
#
# 	http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is
# distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF
# ANY KIND, either express or implied. See the License for the specific
# language governing permissions and limitations under the License.

ROOT := $(shell pwd)

all: build

SOURCEDIR=./ecr-login
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
VERSION := $(shell cat VERSION)
GITFILES := $(shell test -d .git && find ".git/" -type f)
UID:=$(shell id -u)
GID:=$(shell id -g)

BINPATH:=$(abspath ./bin)
BINARY_NAME=docker-credential-ecr-login
LOCAL_BINARY=$(BINPATH)/local/$(BINARY_NAME)

LINUX_AMD64_BINARY=$(BINPATH)/linux-amd64/$(BINARY_NAME)
LINUX_ARM64_BINARY=$(BINPATH)/linux-arm64/$(BINARY_NAME)
DARWIN_AMD64_BINARY=$(BINPATH)/darwin-amd64/$(BINARY_NAME)
DARWIN_ARM64_BINARY=$(BINPATH)/darwin-arm64/$(BINARY_NAME)
WINDOWS_AMD64_BINARY=$(BINPATH)/windows-amd64/$(BINARY_NAME).exe
WINDOWS_ARM64_BINARY=$(BINPATH)/windows-arm64/$(BINARY_NAME).exe

.PHONY: docker
docker: build-in-docker

%-in-docker: GITCOMMIT_SHA
	docker run --rm \
		--user $(UID):$(GID) \
		--env TARGET_GOOS=$(TARGET_GOOS) \
		--env TARGET_GOARCH=$(TARGET_GOARCH) \
		--volume $(ROOT):/go/src/github.com/awslabs/amazon-ecr-credential-helper \
		$(shell docker build -q .) \
		make $(subst -in-docker,,$@)

.PHONY: build
build: $(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_binary.sh $(BINPATH)/local $(VERSION) $(shell cat GITCOMMIT_SHA)
	@echo "Built ecr-login"

.PHONY: test
test:
	cd $(SOURCEDIR) && go test -v -timeout 30s -short -cover ./...

.PHONY: all-variants
all-variants: linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64

.PHONY: linux-amd64
linux-amd64: $(LINUX_AMD64_BINARY)
$(LINUX_AMD64_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_variant.sh linux amd64 $(VERSION) $(shell cat GITCOMMIT_SHA)

.PHONY: linux-arm64
linux-arm64: $(LINUX_ARM64_BINARY)
$(LINUX_ARM64_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_variant.sh linux arm64 $(VERSION) $(shell cat GITCOMMIT_SHA)

.PHONY: darwin-amd64
darwin-amd64: $(DARWIN_AMD64_BINARY)
$(DARWIN_AMD64_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_variant.sh darwin amd64 $(VERSION) $(shell cat GITCOMMIT_SHA)

.PHONY: darwin-arm64
darwin-arm64: $(DARWIN_ARM64_BINARY)
$(DARWIN_ARM64_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_variant.sh darwin arm64 $(VERSION) $(shell cat GITCOMMIT_SHA)

.PHONY: windows-amd64
windows-amd64: $(WINDOWS_AMD64_BINARY)
$(WINDOWS_AMD64_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_variant.sh windows amd64 $(VERSION) $(shell cat GITCOMMIT_SHA)
	@mv $(BINPATH)/windows-amd64/$(BINARY_NAME) $(WINDOWS_AMD64_BINARY)

.PHONY: windows-arm64
windows-arm64: $(WINDOWS_ARM64_BINARY)
$(WINDOWS_ARM64_BINARY): $(SOURCES) GITCOMMIT_SHA
	./scripts/build_variant.sh windows arm64 $(VERSION) $(shell cat GITCOMMIT_SHA)
	@mv $(BINPATH)/windows-arm64/$(BINARY_NAME) $(WINDOWS_ARM64_BINARY)

GITCOMMIT_SHA: $(GITFILES)
	git rev-parse --short=7 HEAD > GITCOMMIT_SHA

release.tar: GITCOMMIT_SHA
	git archive -o release.tar HEAD
	tar -f release.tar --append GITCOMMIT_SHA --owner 0 --group 0

.PHONY: release-tarball
release-tarball: release.tar.gz
release.tar.gz: release.tar
	gzip release.tar

.PHONY: release-tarball-no-vendor
release-tarball-no-vendor: release-novendor.tar.gz
release-novendor.tar.gz: release.tar
	mv release.tar release-novendor.tar
	tar -f release-novendor.tar --wildcards --delete 'ecr-login/vendor/*'
	gzip release-novendor.tar

.PHONY: gogenerate
gogenerate:
	./scripts/gogenerate

.PHONY: get-deps
get-deps:
	go install golang.org/x/tools/cmd/goimports@698251aaa532d49ac69d2c416b0241afb2f65ea5

.PHONY: clean
clean:
	- rm -rf ./bin
	- rm -f GITCOMMIT_SHA
	- rm -f release.tar.gz
	- rm -f release-novendor.tar.gz
