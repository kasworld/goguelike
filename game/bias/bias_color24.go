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

package bias

import (
	"fmt"

	"github.com/kasworld/htmlcolors"
)

func (bs Bias) FactionColor() htmlcolors.Color24 {
	return bs.NearFaction().Color24()
}

func (bs Bias) ToColor24() htmlcolors.Color24 {
	val := bs.ToRGBList()
	return htmlcolors.NewColor24(val[0], val[1], val[2])
}

func (bs Bias) ToRGB() (uint8, uint8, uint8) {
	val := bs.ToRGBList()
	return val[0], val[1], val[2]
}

func (bs Bias) ToRGBString() string {
	val := bs.ToRGBList()
	return fmt.Sprintf("%02x%02x%02x", val[0], val[1], val[2])
}

func (bs Bias) ToHTMLColorString() string {
	val := bs.ToRGBList()
	return fmt.Sprintf("#%02x%02x%02x", val[0], val[1], val[2])
}

func (bs Bias) ToRGBList() [3]uint8 {
	max := 0.0
	for _, v := range bs {
		if v > max {
			max = v
		}
		if -v > max {
			max = -v
		}
	}
	var val [3]uint8
	for i, v := range bs {
		tv := 128 + v/max*128
		if tv > 255 {
			tv = 255
		}
		if tv < 0 {
			tv = 0
		}
		val[i] = uint8(tv)

	}
	return val
}
