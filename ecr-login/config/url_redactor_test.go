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
	"errors"
	"net/url"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRedactURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "URL without sensitive data",
			input:    "https://example.com/path",
			expected: "https://example.com/path",
		},
		{
			name:     "URL with query parameters - values redacted",
			input:    "https://example.com/path?token=secret&id=123",
			expected: "https://example.com/path?id=redacted&token=redacted",
		},
		{
			name:     "URL with username only",
			input:    "https://user@example.com/path",
			expected: "https://user@example.com/path",
		},
		{
			name:     "ECR URL format with scheme",
			input:    "https://123456789012.dkr.ecr.us-west-2.amazonaws.com",
			expected: "https://123456789012.dkr.ecr.us-west-2.amazonaws.com",
		},
		{
			name:     "ECR public URL",
			input:    "https://public.ecr.aws",
			expected: "https://public.ecr.aws",
		},
		{
			name:     "HTTP URL",
			input:    "http://example.com/path",
			expected: "http://example.com/path",
		},
		{
			name:     "HTTP URL with credentials",
			input:    "http://user:pass@example.com/path",
			expected: "http://user:xxxxx@example.com/path",
		},
		{
			name:     "plain string - unchanged",
			input:    "not a url",
			expected: "not a url",
		},
		{
			name:     "URL with only query params - values redacted",
			input:    "https://example.com?secret=value",
			expected: "https://example.com?secret=redacted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestURLRedactorHook_Levels(t *testing.T) {
	hook := &URLRedactorHook{}
	levels := hook.Levels()

	assert.Equal(t, logrus.AllLevels, levels)
}

func TestURLRedactorHook_Fire(t *testing.T) {
	hook := &URLRedactorHook{}

	tests := []struct {
		name           string
		fields         logrus.Fields
		expectedFields logrus.Fields
	}{
		{
			name: "sanitizes serverURL field with password",
			fields: logrus.Fields{
				"serverURL": "https://user:pass@example.com",
			},
			expectedFields: logrus.Fields{
				"serverURL": "https://user:xxxxx@example.com",
			},
		},
		{
			name: "sanitizes serverURL field with query params",
			fields: logrus.Fields{
				"serverURL": "https://example.com?token=secret",
			},
			expectedFields: logrus.Fields{
				"serverURL": "https://example.com?token=redacted",
			},
		},
		{
			name: "serverURL without credentials unchanged",
			fields: logrus.Fields{
				"serverURL": "https://123456789012.dkr.ecr.us-west-2.amazonaws.com",
			},
			expectedFields: logrus.Fields{
				"serverURL": "https://123456789012.dkr.ecr.us-west-2.amazonaws.com",
			},
		},
		{
			name: "does not modify other fields",
			fields: logrus.Fields{
				"registry": "123456789012",
				"region":   "us-west-2",
				"service":  "ecr",
			},
			expectedFields: logrus.Fields{
				"registry": "123456789012",
				"region":   "us-west-2",
				"service":  "ecr",
			},
		},
		{
			name: "handles serverURL with other fields",
			fields: logrus.Fields{
				"serverURL": "https://user:secret@example.com",
				"registry":  "123456789012",
				"region":    "us-west-2",
			},
			expectedFields: logrus.Fields{
				"serverURL": "https://user:xxxxx@example.com",
				"registry":  "123456789012",
				"region":    "us-west-2",
			},
		},
		{
			name: "handles non-string values",
			fields: logrus.Fields{
				"serverURL": "https://user:pass@example.com",
				"count":     42,
				"enabled":   true,
			},
			expectedFields: logrus.Fields{
				"serverURL": "https://user:xxxxx@example.com",
				"count":     42,
				"enabled":   true,
			},
		},
		{
			name:           "handles empty fields",
			fields:         logrus.Fields{},
			expectedFields: logrus.Fields{},
		},
		{
			name: "handles serverURL with non-string value",
			fields: logrus.Fields{
				"serverURL": 12345,
			},
			expectedFields: logrus.Fields{
				"serverURL": 12345,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := &logrus.Entry{
				Data: tt.fields,
			}

			err := hook.Fire(entry)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedFields, entry.Data)
		})
	}
}

func TestRedactURLFromError(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedURL string
		unchanged   bool
	}{
		{
			name:      "nil error returns nil",
			inputErr:  nil,
			unchanged: true,
		},
		{
			name:      "non-url.Error returns unchanged",
			inputErr:  errors.New("some error"),
			unchanged: true,
		},
		{
			name: "url.Error with query params - redacted",
			inputErr: &url.Error{
				Op:  "Get",
				URL: "https://example.com/path?token=secret&id=123",
				Err: errors.New("connection refused"),
			},
			expectedURL: "https://example.com/path?id=redacted&token=redacted",
		},
		{
			name: "url.Error without query params - unchanged",
			inputErr: &url.Error{
				Op:  "Get",
				URL: "https://example.com/path",
				Err: errors.New("connection refused"),
			},
			expectedURL: "https://example.com/path",
		},
		{
			name: "url.Error with password and query params - both redacted",
			inputErr: &url.Error{
				Op:  "Get",
				URL: "https://user:pass@example.com/path?token=secret",
				Err: errors.New("connection refused"),
			},
			expectedURL: "https://user:xxxxx@example.com/path?token=redacted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactURLFromError(tt.inputErr)

			if tt.unchanged {
				assert.Equal(t, tt.inputErr, result)
				return
			}

			var urlErr *url.Error
			assert.True(t, errors.As(result, &urlErr))
			assert.Equal(t, tt.expectedURL, urlErr.URL)
		})
	}
}

func TestURLRedactorHook_Fire_WithError(t *testing.T) {
	hook := &URLRedactorHook{}

	t.Run("redacts url.Error with query params", func(t *testing.T) {
		urlErr := &url.Error{
			Op:  "Get",
			URL: "https://example.com/path?token=secret",
			Err: errors.New("connection refused"),
		}

		entry := &logrus.Entry{
			Data: logrus.Fields{
				logrus.ErrorKey: urlErr,
			},
		}

		err := hook.Fire(entry)
		assert.NoError(t, err)

		resultErr, ok := entry.Data[logrus.ErrorKey].(*url.Error)
		assert.True(t, ok)
		assert.Equal(t, "https://example.com/path?token=redacted", resultErr.URL)
	})

	t.Run("non-url.Error unchanged", func(t *testing.T) {
		plainErr := errors.New("some error")

		entry := &logrus.Entry{
			Data: logrus.Fields{
				logrus.ErrorKey: plainErr,
			},
		}

		err := hook.Fire(entry)
		assert.NoError(t, err)
		assert.Equal(t, plainErr, entry.Data[logrus.ErrorKey])
	})
}
