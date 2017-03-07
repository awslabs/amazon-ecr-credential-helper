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
	"path/filepath"

	log "github.com/cihub/seelog"
	homedir "github.com/mitchellh/go-homedir"
)

func SetupLogger() {
	SetupLoggerWithConfig(loggerConfig())
}

func SetupLoggerWithConfig(config string) {
	logger, err := log.LoggerFromConfigAsString(config)
	if err == nil {
		log.ReplaceLogger(logger)
	} else {
		log.Error(err)
	}
}

func loggerConfig() string {
	logfile, err := homedir.Expand(GetCacheDir() + "/log/ecr-login.log")
	if err != nil {
		fmt.Errorf("%v", err)
		logfile = "/tmp/.ecr/log/ecr-login.log"
	}
	// Clean the path to replace with OS-specific separators
	logfile = filepath.Clean(logfile)
	config := `
	<seelog type="asyncloop" minlevel="debug">
		<outputs formatid="main">
			<rollingfile filename="` + logfile + `" type="date"
			 datepattern="2006-01-02-15" archivetype="none" maxrolls="2" />
			<filter levels="warn,error,critical">
				<console />
			</filter>
		</outputs>
		<formats>
			<format id="main" format="%UTCDate(2006-01-02T15:04:05Z07:00) [%LEVEL] %Msg%n" />
		</formats>
	</seelog>
`
	return config
}
