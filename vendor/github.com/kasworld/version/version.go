// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// goguelike2 version module
package version

import (
	"strings"
	"time"
)

var oriString string
var ver string
var buildDate time.Time
var isRelease bool

func GetVersion() string {
	return ver
}

func GetBuildDate() time.Time {
	return buildDate
}

func IsRelease() bool {
	return isRelease
}

func GetFullString() string {
	return oriString
}

func Set(s string) {
	if s == "" {
		ver = "NoVersion"
	} else {
		oriString = s
		verList := []string{}
		for _, v := range strings.Split(s, "_") {
			if t, err := time.Parse("2006-01-02T15:04:05Z07:00", v); err == nil {
				buildDate = t
			} else if v == "release" {
				isRelease = true
			} else {
				verList = append(verList, v)
			}
		}
		ver = strings.Join(verList, "_")
	}
}

func IsSame(argVer string) bool {
	return ver == "NoVersion" ||
		argVer == "NoVersion" ||
		ver == argVer
}
