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

// Package wrapspan make html span wrapped color text with/without title
package wrapspan

import (
	"fmt"
)

// THCS is interface of html color string conversion
type THCS interface {
	ToHTMLColorString() string
}

func TitleTHCSText(title string, bi THCS, txt string) string {
	return TitleColorText(title, bi.ToHTMLColorString(), txt)
}

func TitleTHCSTextf(title string, bi THCS, format string, arg ...interface{}) string {
	txt := fmt.Sprintf(format, arg...)
	return TitleColorText(title, bi.ToHTMLColorString(), txt)
}

func TitleColorTextf(title string, co string, format string, arg ...interface{}) string {
	txt := fmt.Sprintf(format, arg...)
	return TitleColorText(title, co, txt)
}

func TitleColorText(title string, co string, txt string) string {
	return fmt.Sprintf(`<span style="color:%s" title="%s">%s</span>`,
		co, title, txt)
}

func THCSText(bi THCS, txt string) string {
	return ColorText(bi.ToHTMLColorString(), txt)
}

func THCSTextf(bi THCS, format string, arg ...interface{}) string {
	txt := fmt.Sprintf(format, arg...)
	return ColorText(bi.ToHTMLColorString(), txt)
}

func ColorTextf(co string, format string, arg ...interface{}) string {
	txt := fmt.Sprintf(format, arg...)
	return ColorText(co, txt)
}

func ColorText(co string, txt string) string {
	return fmt.Sprintf(`<span style="color:%s">%s</span>`,
		co, txt)
}
