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

package tilearea

import (
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_stats"
)

func (ta TileArea) CalcNotEmptyTileCount() int {
	rtn := 0
	for _, xv := range ta {
		for _, yv := range xv {
			if !yv.Empty() {
				rtn++
			}
		}
	}
	return rtn
}

func (ta TileArea) CalcStat() tile_stats.TileStat {
	var rtn tile_stats.TileStat
	for _, xv := range ta {
		for _, yv := range xv {
			for i := range rtn {
				if yv.TestByTile(tile.Tile(i)) {
					rtn[i]++
				}
			}
		}
	}
	return rtn
}
