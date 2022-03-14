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

package viewportdata

import (
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/lib/lineofsight"
)

var ViewportXYLenList = findnear.NewXYLenList(
	gameconst.ClientViewPortW, gameconst.ClientViewPortH)[:gameconst.ViewPortWH]

// same order with ViewportXYLenList
type ViewportSight2 [gameconst.ViewPortWH]float32
type ViewportTileArea2 [gameconst.ViewPortWH]tile_flag.TileFlag

var SightlinesByXYLenList = makeSightlinesByXYLenList(ViewportXYLenList)

// makeSightlinesByXYLenList make sight lines 0,0 to all xyLenList dst
func makeSightlinesByXYLenList(xyLenList findnear.XYLenList) []findnear.XYLenList {
	rtn := make([]findnear.XYLenList, len(xyLenList))
	for i, v := range xyLenList {
		// calc to dst near point
		shiftx := 0.5
		shifty := 0.5
		if v.X > 0 {
			shiftx = 0
		} else if v.X < 0 {
			shiftx = 1
		}
		if v.Y > 0 {
			shifty = 0
		} else if v.Y < 0 {
			shifty = 1
		}
		rtn[i] = lineofsight.MakePosLenList(
			0+0.5, 0+0.5, // from src center
			float64(v.X)+shiftx,
			float64(v.Y)+shifty,
		).ToCellLenList()
	}
	return rtn
}
