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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValid_NewEntry(t *testing.T) {
	authEntry := &AuthEntry{
		RequestedAt: time.Now(),
		ExpiresAt:   time.Now().Add(12 * time.Hour),
	}
	assert.True(t, authEntry.IsValid(time.Now()))
}

func TestIsValid_OldEntry(t *testing.T) {
	authEntry := &AuthEntry{
		RequestedAt: time.Now().Add(-12 * time.Hour),
		ExpiresAt:   time.Now(),
	}
	assert.False(t, authEntry.IsValid(time.Now()))
}

func TestIsValid_BeforeRefreshTime(t *testing.T) {
	now := time.Now()
	authEntry := &AuthEntry{
		RequestedAt: now.Add(-6 * time.Hour),
		ExpiresAt:   now.Add(6 * time.Hour),
	}
	assert.True(t, authEntry.IsValid(now.Add(-1*time.Second)))
}

func TestIsValid_AtRefreshTime(t *testing.T) {
	now := time.Now()
	authEntry := &AuthEntry{
		RequestedAt: now.Add(-6 * time.Hour),
		ExpiresAt:   now.Add(6 * time.Hour),
	}
	assert.False(t, authEntry.IsValid(now))
}

func TestIsValid_AfterRefreshTime(t *testing.T) {
	now := time.Now()
	authEntry := &AuthEntry{
		RequestedAt: now.Add(-6 * time.Hour),
		ExpiresAt:   now.Add(6 * time.Hour),
	}
	assert.False(t, authEntry.IsValid(now.Add(time.Second)))
}
