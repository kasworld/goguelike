// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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

import (
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
)

var gTile3D [tile.Tile_Count]*Tile3D
var gTile3DDark [tile.Tile_Count]*Tile3D

var gTileZShift = [tile.Tile_Count]float64{
	tile.Swamp:  1,
	tile.Soil:   2,
	tile.Stone:  2,
	tile.Sand:   2,
	tile.Sea:    0,
	tile.Magma:  0,
	tile.Ice:    3,
	tile.Grass:  2,
	tile.Tree:   2,
	tile.Road:   3,
	tile.Room:   2,
	tile.Wall:   2,
	tile.Window: 2,
	tile.Door:   2,
	tile.Fog:    4,
	tile.Smoke:  4,
}

func preMakeTileMatGeo() {
	var tlt tile.Tile

	tlt = tile.Swamp
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Soil
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Stone
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Sand
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Sea
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Magma
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Ice
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Road
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Room
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Fog
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Smoke
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt).MakeSrcDark()

	tlt = tile.Grass
	gTile3D[tlt] = NewTile3D_Grass()
	gTile3DDark[tlt] = NewTile3D_Grass().MakeSrcDark()

	tlt = tile.Tree
	gTile3D[tlt] = NewTile3D_Tree()
	gTile3DDark[tlt] = NewTile3D_Tree().MakeSrcDark()

	tlt = tile.Wall
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Window
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Door
	gTile3D[tlt] = NewTile3D_BoxTile(gClientTile.FloorTiles[tile.Door][0])
	gTile3DDark[tlt] = NewTile3D_BoxTile(gClientTile.FloorTiles[tile.Door][0])

	for i := 0; i < tile.Tile_Count; i++ {
		gTile3D[i].Shift[2] = gTileZShift[i]
		gTile3DDark[i].Shift[2] = gTileZShift[i]
	}
	// for i, tl := range gTile3D {
	// 	jslog.Infof("%v %v", tile.Tile(i), tl.GeoInfo)
	// }
}

var tileHeightCache [1 << uint(tile.Tile_Count)]float64

func init() {
	for i := range tileHeightCache {
		tileHeightCache[i] = Tile3DHeightMin
	}
}

func calcTile3DHeight(tl tile_flag.TileFlag) float64 {
	rtn := Tile3DHeightMin
	for i := 0; i < tile.Tile_Count; i++ {
		if !tl.TestByTile(tile.Tile(i)) {
			continue
		}
		z := gTile3D[i].GeoInfo.Len[2] + gTile3D[i].Shift[2]
		if z > rtn {
			rtn = z
		}
	}
	return rtn
}

func GetTile3DHeightByCache(tl tile_flag.TileFlag) float64 {
	z := tileHeightCache[tl]
	if z == Tile3DHeightMin {
		z = calcTile3DHeight(tl)
		if z == Tile3DHeightMin { // empty tile
			z = 0 // prevent recalc empty tile
		}
		tileHeightCache[tl] = z
	}
	return z
}
