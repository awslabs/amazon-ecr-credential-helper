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

package cache

import (
	"github.com/stretchr/testify/mock"
)

type MockCredentialsCache struct {
	mock.Mock

	SetFn func(registry string, entry *AuthEntry)
}

func (c *MockCredentialsCache) Get(registry string) *AuthEntry {
	args := c.Called(registry)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(*AuthEntry)
}

func (c *MockCredentialsCache) Set(registry string, entry *AuthEntry) {
	if c.SetFn != nil {
		c.SetFn(registry, entry)
		return
	}
	c.Called(registry, entry)
}

func (c *MockCredentialsCache) List() []*AuthEntry {
	args := c.Called()
	return args.Get(0).([]*AuthEntry)
}

func (c *MockCredentialsCache) Clear() {
	c.Called()
}
