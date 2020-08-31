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

package terraini

import (
	"net/http"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/game/terrain/resourcetilearea"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/goguelike/game/terrain/viewportcache"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/wrapper"
)

type TerrainI interface {
	Init() error
	GetOriRcsTiles() resourcetilearea.ResourceTileArea
	GetTiles() tilearea.TileArea
	GetName() string
	GetFieldObjPosMan() *uuidposman.UUIDPosMan
	GetRoomList() []*room.Room
	GetViewportCache() *viewportcache.ViewportCache
	FindPath(dstx, dsty, srcx, srcy int, trylimit int) [][2]int

	GetResetAfterNAgeing() int64
	GetMSPerAgeing() int64
	Ageing() error
	AgeingCount() int64

	GetXYLen() (int, int)
	GetXLen() int
	GetYLen() int
	GetXWrapper() *wrapper.Wrapper
	GetYWrapper() *wrapper.Wrapper
	WrapXY(x, y int) (int, int)
	// GetTileWrapped(x, y int) tile_flag.TileFlag
	GetTile2Discover() int
	GetRcsTiles() resourcetilearea.ResourceTileArea

	GetActiveObjCount() int
	GetCarryObjCount() int
	GetScript() []string

	Search1stByXYLenList(
		xylenlist findnear.XYLenList,
		sx, sy int,
		filterfn func(t *tile_flag.TileFlag) bool) (int, int, bool)

	Web_TerrainInfo(w http.ResponseWriter, r *http.Request)
	Web_TerrainImageZoom(w http.ResponseWriter, r *http.Request)
	Web_TerrainImageAutoZoom(w http.ResponseWriter, r *http.Request)
	Web_TileInfo(w http.ResponseWriter, r *http.Request)
}
