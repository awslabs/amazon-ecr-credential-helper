# Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
BINARY_NAME=docker-credential-ecr-login
LOCAL_BINARY=bin/local/$(BINARY_NAME)

LINUX_AMD64_BINARY=bin/linux-amd64/$(BINARY_NAME)
DARWIN_AMD64_BINARY=bin/darwin-amd64/$(BINARY_NAME)
WINDOWS_AMD64_BINARY=bin/windows-amd64/$(BINARY_NAME).exe

CONTAINER_NAME=build-$(BINARY_NAME)

.PHONY: docker
docker: Dockerfile
	docker rm $(CONTAINER_NAME) 2>/dev/null ||:
	docker run \
		--name $(CONTAINER_NAME) \
		-e TARGET_GOOS=$(TARGET_GOOS) \
		-e TARGET_GOARCH=$(TARGET_GOARCH) \
		$(shell docker build -q .)
	mkdir -p bin/local
	docker cp $(CONTAINER_NAME):/go/src/github.com/awslabs/amazon-ecr-credential-helper/bin/local/$(BINARY_NAME) $(LOCAL_BINARY)

.PHONY: build
build: $(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES)
	. ./scripts/shared_env && ./scripts/build_binary.sh ./bin/local
	@echo "Built ecr-login"

.PHONY: test
test:
	. ./scripts/shared_env && go test -v -timeout 30s -short -cover $(shell go list ./ecr-login/... | grep -v /vendor/)

.PHONY: all-variants
all-variants: linux-amd64 darwin-amd64 windows-amd64

.PHONY: linux-amd64
linux-amd64: $(LINUX_AMD64_BINARY)
$(LINUX_AMD64_BINARY): $(SOURCES)
	./scripts/build_variant.sh linux amd64

.PHONY: darwin-amd64
darwin-amd64: $(DARWIN_AMD64_BINARY)
$(DARWIN_AMD64_BINARY): $(SOURCES)
	./scripts/build_variant.sh darwin amd64

.PHONY: windows-amd64
windows-amd64: $(WINDOWS_AMD64_BINARY)
$(WINDOWS_AMD64_BINARY): $(SOURCES)
	./scripts/build_variant.sh windows amd64
	@mv ./bin/windows-amd64/docker-credential-ecr-login ./bin/windows-amd64/docker-credential-ecr-login.exe

.PHONY: gogenerate
gogenerate:
	./scripts/gogenerate

.PHONY: get-deps
get-deps:
	go get github.com/tools/godep
	go get golang.org/x/tools/cmd/cover
	go get github.com/golang/mock/mockgen
	go get golang.org/x/tools/cmd/goimports

.PHONY: clean
clean:
	rm -rf ./bin ||:
	docker rm $(CONTAINER_NAME) 2>/dev/null ||:
