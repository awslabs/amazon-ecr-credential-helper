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

FROM public.ecr.aws/docker/library/golang:1.24-alpine

WORKDIR /go/src/github.com/awslabs/amazon-ecr-credential-helper

COPY ./scripts/container_init.sh /setup/container_init.sh

RUN /setup/container_init.sh

