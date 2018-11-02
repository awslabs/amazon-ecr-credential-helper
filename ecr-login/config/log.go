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

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	logrusConfig()
}

func logrusConfig() {
	logfile, err := homedir.Expand(GetCacheDir() + "/log/ecr-login.log")
	if err != nil {
		fmt.Errorf("%v", err)
		logfile = "/tmp/.ecr/log/ecr-login.log"
	}
	// Clean the path to replace with OS-specific separators
	logfile = filepath.Clean(logfile)
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return
	}
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(file)
}
