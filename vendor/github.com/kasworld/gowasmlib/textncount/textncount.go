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

// Package textncount handle ordered text with repeating count
package textncount

import (
	"bytes"
	"fmt"
)

type TextNCount struct {
	Text  string
	Count int
}

func (v TextNCount) String() string {
	if v.Count != 1 {
		return fmt.Sprintf(`%v (%v)`, v.Text, v.Count)
	} else {
		return fmt.Sprintf(`%v`, v.Text)
	}
}

type TextNCountList []*TextNCount

func (til *TextNCountList) Appendf(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	til.Append(s)
}

func (til *TextNCountList) Append(txt string) {
	ti := &TextNCount{txt, 1}
	if last := til.Last(); last != nil && last.Text == ti.Text {
		last.Count++
		return
	}
	*til = append(*til, ti)
}

func (til *TextNCountList) Clear() {
	*til = (*til)[:0]
}

func (til TextNCountList) Last() *TextNCount {
	if l := len(til); l > 0 {
		return til[l-1]
	}
	return nil
}

func (til TextNCountList) GetLastN(n int) TextNCountList {
	l := len(til)
	if l <= n {
		return til
	}
	return til[l-n:]
}

func (til TextNCountList) ToHtmlString() string {
	var buf bytes.Buffer
	for _, v := range til {
		if v.Count != 1 {
			fmt.Fprintf(&buf, `%v (%v)<br/>`,
				v.Text, v.Count)
		} else {
			fmt.Fprintf(&buf, `%v<br/>`,
				v.Text)
		}
	}
	return buf.String()
}

func (til TextNCountList) ToHtmlStringRev() string {
	var buf bytes.Buffer
	for i := len(til) - 1; i >= 0; i-- {
		if til[i].Count != 1 {
			fmt.Fprintf(&buf, `%v (%v)<br/>`,
				til[i].Text, til[i].Count)
		} else {
			fmt.Fprintf(&buf, `%v<br/>`,
				til[i].Text)
		}
	}
	return buf.String()
}
