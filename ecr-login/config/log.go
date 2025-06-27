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
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func SetupLogger() {
	logdir, err := homedir.Expand(GetCacheDir() + "/log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: failed to find directory: %v", err)
		logdir = os.TempDir()
	}
	// Clean the path to replace with OS-specific separators
	logdir = filepath.Clean(logdir)
	err = os.MkdirAll(logdir, os.ModeDir|0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log: failed to create directory: %v", err)
		logdir = os.TempDir()
	}
	file, err := os.OpenFile(filepath.Join(logdir, "ecr-login.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug})))
}
