// Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

const loadConfigEnvVar = "AWS_SDK_LOAD_CONFIG"

func TestSharedConfigState(t *testing.T) {
	defer func(value string) {
		os.Setenv(loadConfigEnvVar, value)
	}(os.Getenv(loadConfigEnvVar))

	cases := []struct {
		envValue string
		expected session.SharedConfigState
	}{
		{"", session.SharedConfigEnable},
		{"true", session.SharedConfigEnable},
		{"false", session.SharedConfigDisable},
	}

	for _, testCase := range cases {
		t.Run(testCase.envValue, func(t *testing.T) {
			os.Setenv(loadConfigEnvVar, testCase.envValue)
			state := loadSharedConfigState()
			assert.NotNil(t, state)
			assert.Equal(t, testCase.expected, state)
		})
	}
}
