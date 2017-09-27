// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"fmt"
	"os"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/config"
	log "github.com/cihub/seelog"
	"github.com/docker/docker-credential-helpers/credentials"
)

func main() {
	defer log.Flush()
	config.SetupLogger()
	helper := ecr.ECRHelper{ClientFactory: api.DefaultClientFactory{}}
	if len(os.Args) > 1 && os.Args[1] == "eval" {
		evalCommand(helper)
	} else {
		credentials.Serve(helper)
	}
}

func evalCommand(helper credentials.Helper) {
	server, err := parseArgs(helper)
	if err == nil {
		var user, token string
		user, token, err = helper.Get(server)
		if err == nil {
			fmt.Printf("docker login -e none -u %s -p %s %s\n", user, token, server)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
}

func parseArgs(helper credentials.Helper) (string, error) {
	if len(os.Args) > 2 {
		return os.Args[2], nil
	}
	servers, err := helper.List()
	if err == nil {
		for k := range servers {
			return k, nil
		}
		return "", fmt.Errorf("No default ECR servers found")
	}
	return "", err
}
