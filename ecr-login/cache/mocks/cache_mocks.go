// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package mock_cache

import (
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/cache"
)

type MockCredentialsCache struct {
	GetFn       func(registry string) *cache.AuthEntry
	GetPublicFn func() *cache.AuthEntry
	SetFn       func(registry string, entry *cache.AuthEntry)
	ListFn      func() []*cache.AuthEntry
	ClearFn     func()
}

var _ cache.CredentialsCache = (*MockCredentialsCache)(nil)

func (m MockCredentialsCache) Get(registry string) *cache.AuthEntry {
	return m.GetFn(registry)
}

func (m MockCredentialsCache) GetPublic() *cache.AuthEntry {
	return m.GetPublicFn()
}

func (m MockCredentialsCache) Set(registry string, entry *cache.AuthEntry) {
	m.SetFn(registry, entry)
}

func (m MockCredentialsCache) List() []*cache.AuthEntry {
	return m.ListFn()
}

func (m MockCredentialsCache) Clear() {
	m.ClearFn()
}
