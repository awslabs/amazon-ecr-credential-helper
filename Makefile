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
LOCAL_BINARY=bin/local/docker-credential-ecr-login

.PHONY: docker
docker: Dockerfile
	docker run --rm \
	-e TARGET_GOOS=$(TARGET_GOOS) \
	-e TARGET_GOARCH=$(TARGET_GOARCH) \
	-v $(shell pwd)/bin:/go/src/github.com/awslabs/amazon-ecr-credential-helper/bin \
	$(shell docker build -q .)

.PHONY: build
build: $(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES)
	. ./scripts/shared_env && ./scripts/build_binary.sh ./bin/local
	@echo "Built ecr-login"

.PHONY: test
test:
	. ./scripts/shared_env && go test -v -timeout 30s -short -cover $(shell go list ./ecr-login/... | grep -v /vendor/)

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
