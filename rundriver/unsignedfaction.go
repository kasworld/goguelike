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

package main

import (
	"fmt"
	"sort"

	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/htmlcolors"
)

func main() {
	// str := "가냐더려모뵤수유즈치캐테퍠혜꾀똬쀠쒀쯰"
	bsl := make(bias.BiasList, 0)
	st := -1
	ed := st + 3
	for i := st; i < ed; i++ {
		for j := st; j < ed; j++ {
			for k := st; k < ed; k++ {
				bs := bias.Bias{float64(i), float64(j), float64(k)}
				bsl = append(bsl, bs.MakeAbsSumTo(100))
			}
		}
	}
	sort.Sort(bias.BiasList(bsl))
	ftL := make(bias.BiasList, 0)
	for i, v := range bsl {
		if v[0]+v[1]+v[2] != 100 {
			continue
		}
		if i == 0 {
			ftL = append(ftL, v)
		} else {
			if v != bsl[i-1] {
				ftL = append(ftL, v)
			}
		}
	}
	for i, v := range ftL {
		// co := v.MakeAbsSumTo(255)
		// co2 := htmlcolors.NewColor24(uint8(co[0]), uint8(co[1]), uint8(co[2]))
		fmt.Printf("%v %v %v %v\n",
			v.IntString(), i, ToColor24_u(v),
			ToColor24_u(v).NearNamedColor24())
	}
}

func ToColor24_u(bs bias.Bias) htmlcolors.Color24 {
	val := ToRGBList_u(bs)
	return htmlcolors.NewColor24(val[0], val[1], val[2])
}

func ToRGBList_u(bs bias.Bias) [3]uint8 {
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
		tv := v / max * 255
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

/*
[0 0 100] 0 Blue Blue
[0 33 67] 1 Color24[0x00007fff] DodgerBlue
[0 50 50] 2 Cyan Cyan
[0 67 33] 3 SpringGreen SpringGreen
[0 100 0] 4 Lime Lime
[20 40 40] 5 Color24[0x007fffff] Aquamarine
[25 25 50] 6 Color24[0x007f7fff] MediumSlateBlue
[25 50 25] 7 Color24[0x007fff7f] LightGreen
[33 0 67] 8 Color24[0x007f00ff] DarkViolet
[33 33 33] 9 White White
[33 67 0] 10 Chartreuse Chartreuse
[40 20 40] 11 Color24[0x00ff7fff] Violet
[40 40 20] 12 Color24[0x00ffff7f] Khaki
[50 0 50] 13 Magenta Magenta
[50 25 25] 14 Color24[0x00ff7f7f] Salmon
[50 50 0] 15 Yellow Yellow
[67 0 33] 16 Color24[0x00ff007f] DeepPink
[67 33 0] 17 Color24[0x00ff7f00] DarkOrange
[100 0 0] 18 Red Red

*/
