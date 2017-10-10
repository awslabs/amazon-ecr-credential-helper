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

package ecr

import (
	"errors"
	"fmt"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	log "github.com/cihub/seelog"
	"github.com/docker/docker-credential-helpers/credentials"
)

var notImplemented = errors.New("not implemented")

type ECRHelper struct {
	ClientFactory api.ClientFactory
}

// ensure ECRHelper adheres to the credentials.Helper interface
var _ credentials.Helper = (*ECRHelper)(nil)

func (ECRHelper) Add(creds *credentials.Credentials) error {
	// This does not seem to get called
	return notImplemented
}

func (ECRHelper) Delete(serverURL string) error {
	// This does not seem to get called
	return notImplemented
}

func (self ECRHelper) Get(serverURL string) (string, string, error) {
	defer log.Flush()

	registry, err := api.ExtractRegistry(serverURL)
	if err != nil {
		log.Errorf("Error parsing the serverURL: %s, error: %v", serverURL, err)
		return "", "", credentials.NewErrCredentialsNotFound()
	}

	client := self.ClientFactory.NewClientFromRegion(registry.Region)
	auth, err := client.GetCredentials(serverURL)
	if err != nil {
		log.Errorf("Error retrieving credentials: %v", err)
		return "", "", credentials.NewErrCredentialsNotFound()
	}
	return auth.Username, auth.Password, nil
}

func (self ECRHelper) List() (map[string]string, error) {
	log.Debug("Listing credentials")
	client := self.ClientFactory.NewClientWithDefaults()

	auths, err := client.ListCredentials()
	if err != nil {
		log.Errorf("Error listing credentials: %v", err)
		return nil, fmt.Errorf("Could not list credentials: %v:", err)
	}

	result := map[string]string{}

	for _, auth := range auths {
		serverURL := auth.ProxyEndpoint
		result[serverURL] = auth.Username
	}
	return result, nil
}
