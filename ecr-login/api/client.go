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

package api

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	log "github.com/cihub/seelog"
)

const proxyEndpointScheme = "https://"

type Client struct {
	ecrClient *ecr.ECR
}

func NewClient(region string) *Client {
	return &Client{
		ecrClient: ecr.New(session.New(), &aws.Config{Region: aws.String(region)}),
	}
}

func (self *Client) GetCredentials(registry, image string) (string, string, error) {
	log.Debugf("Calling ECR.GetAuthorizationToken for %s", registry)
	input := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(registry)},
	}

	output, err := self.ecrClient.GetAuthorizationToken(input)
	if err != nil {
		return "", "", err
	}
	if output == nil {
		return "", "", fmt.Errorf("Missing AuthorizationData in ECR response for %s", registry)
	}
	for _, authData := range output.AuthorizationData {
		if authData.ProxyEndpoint != nil &&
			strings.HasPrefix(proxyEndpointScheme+image, aws.StringValue(authData.ProxyEndpoint)) &&
			authData.AuthorizationToken != nil {
			return extractToken(authData)
		}
	}
	return "", "", fmt.Errorf("No AuthorizationToken found for %s", registry)
}

func extractToken(authData *ecr.AuthorizationData) (string, string, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(aws.StringValue(authData.AuthorizationToken))
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(string(decodedToken), ":", 2)
	return parts[0], parts[1], nil
}
