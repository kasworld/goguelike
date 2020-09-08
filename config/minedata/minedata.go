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

package minedata

import (
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/config/gameconst"
)

var MineData [gameconst.ViewPortW]findnear.XYLenList

func init() {
	fullData := findnear.NewXYLenList(
		gameconst.ClientViewPortW, gameconst.ClientViewPortH)[:gameconst.ViewPortWH]

	st := 0
	curRadius := 1.0
	for i, v := range fullData {
		if v.L > curRadius {
			MineData[int(curRadius)-1] = fullData[st:i]
			st = i
			curRadius++
		}
	}
	MineData[int(curRadius)-1] = fullData[st:]
	// fmt.Printf("%v", MineData)
}
