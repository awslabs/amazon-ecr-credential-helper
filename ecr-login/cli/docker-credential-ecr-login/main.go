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
	"flag"
	"fmt"
	"os"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/config"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/version"
	log "github.com/cihub/seelog"
	"github.com/docker/docker-credential-helpers/credentials"
)

const banner = `amazon-ecr-credential-helper
Version:    %s
Git commit: %s
`

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "v", false, "print version and exit")
	flag.Parse()

	// Exit safely when version is used
	if versionFlag {
		fmt.Printf(banner, version.Version, version.GitCommitSHA)
		os.Exit(0)
	}

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
	server, email, passwordStdin, err := parseArgs(helper)
	if err == nil {
		var user, token string
		user, token, err = helper.Get(server)
		if err == nil {
			var emailOpt string
			if email {
				emailOpt = " -e none"
			}
			var echoPassword, passwordOpt string
			if passwordStdin {
				echoPassword = fmt.Sprintf("echo %s|", token)
				passwordOpt = "--password-stdin"
			} else {
				passwordOpt = fmt.Sprintf("-p %s", token)
			}

			fmt.Printf("%sdocker login%s -u %s %s %s\n",
				echoPassword, emailOpt, user, passwordOpt, server)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}
}

func parseArgs(helper credentials.Helper) (string, bool, bool, error) {
	var err error
	var email, passwordStdin bool
	for i := 2; i < len(os.Args); i += 1 {
		switch os.Args[i] {
		case "-e":
			email = true
		case "--password-stdin":
			passwordStdin = true
		default:
			err = fmt.Errorf("Usage: %s [-e] [--password-stdin]", os.Args[0])
			break
		}
	}
	if err == nil {
		var servers map[string]string
		servers, err = helper.List()
		if err == nil {
			// Return any server in the map
			for k := range servers {
				return k, email, passwordStdin, nil
			}
			err = fmt.Errorf("No default ECR servers found")
		}
	}
	return "", false, false, err
}
