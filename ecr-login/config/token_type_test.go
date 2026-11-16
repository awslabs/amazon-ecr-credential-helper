// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenType(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "unset", expected: TokenTypeRegional},
		{name: "regional", value: TokenTypeRegional, expected: TokenTypeRegional},
		{name: "global", value: TokenTypeGlobal, expected: TokenTypeGlobal},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv(tokenTypeEnvironmentVariable, testCase.value)
			assert.Equal(t, testCase.expected, GetTokenType())
		})
	}
}

func TestValidateTokenType(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		wantError bool
	}{
		{name: "unset"},
		{name: "regional", value: TokenTypeRegional},
		{name: "global", value: TokenTypeGlobal},
		{name: "invalid", value: "signin", wantError: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv(tokenTypeEnvironmentVariable, testCase.value)
			err := ValidateTokenType()
			if testCase.wantError {
				assert.ErrorContains(t, err, "unsupported AWS_ECR_PRIVATE_TOKEN_TYPE")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
