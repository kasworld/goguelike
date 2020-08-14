// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package canvastext

import (
	"fmt"
	"time"
)

type TextTemplate struct {
	Color    string
	SizeRate float64
	Dur      time.Duration
}

type CanvasText struct {
	Color     string
	SizeRate  float64
	Dur       time.Duration
	Text      string
	StartTime time.Time `prettystring:"simple"`
}

func (ct CanvasText) IsEnded() bool {
	return ct.StartTime.Add(ct.Dur).Before(time.Now())
}

func (ct CanvasText) ProgressRate() float64 {
	return time.Now().Sub(ct.StartTime).Seconds() / ct.Dur.Seconds()
}

func (ct CanvasText) CanMerge(dst CanvasText) bool {
	return ct.Color == dst.Color && ct.SizeRate == dst.SizeRate && ct.Text == dst.Text
}

type CanvasTextList []*CanvasText

func (til *CanvasTextList) AppendTf(
	t TextTemplate,
	format string, arg ...interface{}) {
	ti := &CanvasText{
		Color:     t.Color,
		SizeRate:  t.SizeRate,
		Dur:       t.Dur,
		Text:      fmt.Sprintf(format, arg...),
		StartTime: time.Now(),
	}
	*til = append(*til, ti)
}

func (til *CanvasTextList) Appendf(
	color string, sizeRate float64, dur time.Duration,
	format string, arg ...interface{}) {
	ti := &CanvasText{
		Color:     color,
		SizeRate:  sizeRate,
		StartTime: time.Now(),
		Dur:       dur,
		Text:      fmt.Sprintf(format, arg...),
	}
	*til = append(*til, ti)
}

func (til *CanvasTextList) Clear() {
	*til = (*til)[:0]
}

func (til CanvasTextList) Compact() CanvasTextList {
	var rtn CanvasTextList
	for _, v := range til {
		if v.IsEnded() {
			continue
		}
		rtn = append(rtn, v)
	}

	return rtn
}

func (til CanvasTextList) Last() *CanvasText {
	if l := len(til); l > 0 {
		return til[l-1]
	}
	return nil
}

func (til CanvasTextList) GetLastN(n int) CanvasTextList {
	l := len(til)
	if l <= n {
		return til
	}
	return til[l-n:]
}
