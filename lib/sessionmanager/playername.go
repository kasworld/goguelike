// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sessionmanager

import (
	"strings"
)

func ValidatePlayername(s string, padding string) string {
	validPlayername := strings.NewReplacer(
		" ", "_",
		"<", "_",
		">", "_",
		"\r", "_",
		"\t", "_",
		"\n", "_",
	)
	s = validPlayername.Replace(s)
	if l := len(s); l > 20 {
		s = s[:20]
	}
	if len(s) < 4 {
		s = s + padding
	}
	return s
}
