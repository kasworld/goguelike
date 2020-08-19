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

package terrain

import (
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/tileoptype"
)

func (tr *Terrain) resource2View() {
	for x, xv := range tr.resourceTileArea {
		for y, yv := range xv {
			tr.tileArea[x][y] = yv.Resource2View()
		}
	}
}

func (tr *Terrain) drawCorridors() {
	for _, v := range tr.corridorList {
		for _, w := range v.P {
			tr.tileArea[w[0]][w[1]].OverrideBits(v.Tile)
		}
	}
}

func (tr *Terrain) openBlockedDoor() {
	for _, r := range tr.roomManager.GetRoomList() {
		for _, connPos := range r.ConnectPos {
			x, y := tr.XWrap(connPos[0]), tr.YWrap(connPos[1])
			t := tr.tileArea[x][y]
			if t.CannotPlaceObj() {
				if tr.roomManager.GetRoomByPos(x, y) != nil {
					tr.tileArea[x][y].Op(
						tile_flag.TileTypeValue{
							Op:  tileoptype.OverrideBits,
							Arg: tile.Door},
					)
					tr.log.Warn("wall blocked door found %v [%v %v], change to door", tr, x, y)
				}
			}
		}
	}
}
