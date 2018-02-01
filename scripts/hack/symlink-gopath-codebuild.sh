#!/bin/bash
# Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

# This script is meant to make all package available in the 'canonical' location
# for execution in AWS CodeBuild in the golang container image.


if [[ ! -d "/go/src/github.com/awslabs/amazon-ecr-credential-helper" ]]; then
	mkdir -p "/go/src/github.com/awslabs"
	ln -s "$(pwd)" "/go/src/github.com/awslabs/amazon-ecr-credential-helper"
fi