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

package terrain

import (
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/game/terrain/resourcetilearea"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/goguelike/game/terrain/viewportcache"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/goguelike/lib/uuidposmani"
	"github.com/kasworld/wrapper"
)

// imp terraini
func (tr *Terrain) GetTiles() tilearea.TileArea {
	return tr.serviceTileArea
}

func (tr *Terrain) GetTileWrapped(x, y int) tile_flag.TileFlag {
	return tr.serviceTileArea[tr.XWrap(x)][tr.YWrap(y)]
}

func (tr *Terrain) GetName() string {
	return tr.Name
}

func (tr *Terrain) GetFieldObjPosMan() uuidposmani.UUIDPosManI {
	return tr.foPosMan
}

func (tr *Terrain) GetRoomList() []*room.Room {
	return tr.roomManager.GetRoomList()
}
func (tr *Terrain) GetResetAfterNAgeing() int64 {
	return tr.ResetAfterNAgeing
}
func (tr *Terrain) GetMSPerAgeing() int64 {
	return tr.MSPerAgeing
}
func (tr *Terrain) AgeingCount() int64 {
	return tr.ageingCount
}

func (tr *Terrain) GetViewportCache() *viewportcache.ViewportCache {
	return tr.viewportCache
}

func (tr *Terrain) GetXYLen() (int, int) {
	return tr.Xlen, tr.Ylen
}
func (tr *Terrain) GetXLen() int {
	return tr.Xlen
}
func (tr *Terrain) GetYLen() int {
	return tr.Ylen
}

func (tr *Terrain) GetXWrapper() *wrapper.Wrapper {
	return tr.XWrapper
}
func (tr *Terrain) GetYWrapper() *wrapper.Wrapper {
	return tr.YWrapper
}

func (tr *Terrain) WrapXY(x, y int) (int, int) {
	return tr.XWrap(x), tr.YWrap(y)
}

func (tr *Terrain) GetScript() []string {
	return tr.terrainScript
}

func (tr *Terrain) GetTile2Discover() int {
	return tr.Tile2Discover
}

func (tr *Terrain) GetRcsTiles() resourcetilearea.ResourceTileArea {
	return tr.resourceTileArea
}

func (tr *Terrain) GetOriRcsTiles() resourcetilearea.ResourceTileArea {
	return tr.oriTiles
}

func (tr *Terrain) GetActiveObjCount() int {
	return tr.ActiveObjCount
}
func (tr *Terrain) GetCarryObjCount() int {
	return tr.CarryObjCount
}

func (tr *Terrain) FindPath(dstx, dsty, srcx, srcy int, trylimit int) [][2]int {
	return tr.ta4ff.FindPath(dstx, dsty, srcx, srcy, trylimit)
}
