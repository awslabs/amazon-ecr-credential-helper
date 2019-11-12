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

# Taken from https://github.com/aws/amazon-ecs-agent/blob/ecda8a686200643081fe7de498b61b1c023b2c25/misc/certs/Dockerfile
FROM debian:latest as certs

RUN apt-get update &&  \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# If anyone has a better idea for how to trim undesired certs or a better ca list to use, I'm all ears
RUN cp /etc/ca-certificates.conf /tmp/caconf && \
  cat /tmp/caconf | grep -v "mozilla/CNNIC_ROOT\.crt" > /etc/ca-certificates.conf && \
  update-ca-certificates --fresh

FROM scratch

COPY --from=builder /go/bin/docker-credential-ecr-login /usr/local/bin/docker-credential-ecr-login
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT [ "/usr/local/bin/docker-credential-ecr-login" ]
CMD [ "eval" ]
