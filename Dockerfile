# Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You
# may not use this file except in compliance with the License. A copy of
# the License is located at
#
#       http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is
# distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF
# ANY KIND, either express or implied. See the License for the specific
# language governing permissions and limitations under the License.

FROM golang:1.9-alpine3.6 AS builder

COPY ecr-login /go/src/github.com/awslabs/amazon-ecr-credential-helper/ecr-login

RUN env CGO_ENABLED=0 go install -installsuffix "static" \
    github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cli/docker-credential-ecr-login

FROM alpine:3.6

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/docker-credential-ecr-login /usr/local/bin/docker-credential-ecr-login

ENTRYPOINT [ "/usr/local/bin/docker-credential-ecr-login" ]
CMD [ "eval" ]
