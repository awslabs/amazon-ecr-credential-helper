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
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/stretchr/testify/mock"
)

var (
	_ ECRClient = &MockECRClient{}
)

type ECRClient interface {
	GetAuthorizationToken(*ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error)
}

type MockECRClient struct {
	mock.Mock
	GetAuthorizationTokenFn func(i *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error)
}

func (m *MockECRClient) GetAuthorizationToken(i *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
	if m.GetAuthorizationTokenFn != nil {
		return m.GetAuthorizationTokenFn(i)
	}

	args := m.Called(i)
	return args.Get(0).(*ecr.GetAuthorizationTokenOutput), args.Error(1)
}
