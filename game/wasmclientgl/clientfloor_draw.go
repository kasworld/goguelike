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

package wasmclientgl

import "github.com/kasworld/goguelike/enum/tile"

func (cf *ClientFloorGL) calcWallTileDiff(flx, fly int) int {
	rtn := 0
	if cf.checkWallAt(flx, fly-1) {
		rtn |= 1
	}
	if cf.checkWallAt(flx+1, fly) {
		rtn |= 1 << 1
	}
	if cf.checkWallAt(flx, fly+1) {
		rtn |= 1 << 2
	}
	if cf.checkWallAt(flx-1, fly) {
		rtn |= 1 << 3
	}
	return rtn
}

func (cf *ClientFloorGL) checkWallAt(flx, fly int) bool {
	flx = cf.XWrapSafe(flx)
	fly = cf.YWrapSafe(fly)
	tl := cf.Tiles[flx][fly]
	return tl.TestByTile(tile.Wall) ||
		tl.TestByTile(tile.Door) ||
		tl.TestByTile(tile.Window)
}
