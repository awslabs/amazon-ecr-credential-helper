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
	"fmt"
	"os"
)

const (
	TokenTypeRegional = "regional-token"
	TokenTypeGlobal   = "global-token"

	tokenTypeEnvironmentVariable = "AWS_ECR_PRIVATE_TOKEN_TYPE"
)

// GetTokenType returns the configured token type. An unset value preserves the
// existing regional token behavior.
func GetTokenType() string {
	if tokenType := os.Getenv(tokenTypeEnvironmentVariable); tokenType != "" {
		return tokenType
	}
	return TokenTypeRegional
}

// ValidateTokenType rejects values that the helper does not understand.
func ValidateTokenType() error {
	switch GetTokenType() {
	case TokenTypeRegional, TokenTypeGlobal:
		return nil
	default:
		return fmt.Errorf(
			"unsupported %s %q; expected %q, %q, or unset",
			tokenTypeEnvironmentVariable,
			GetTokenType(),
			TokenTypeRegional,
			TokenTypeGlobal,
		)
	}
}

// IsGlobalToken returns true when the customer explicitly selects a global
// token.
func IsGlobalToken() bool {
	return GetTokenType() == TokenTypeGlobal
}
