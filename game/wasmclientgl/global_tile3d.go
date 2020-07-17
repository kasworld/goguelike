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

var gTile3D [tile.Tile_Count]Tile3D
var gTile3DDark [tile.Tile_Count]Tile3D

func preMakeTileMatGeo() {
	var tlt tile.Tile

	tlt = tile.Swamp
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 1)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 1)

	tlt = tile.Soil
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 2)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 2)

	tlt = tile.Stone
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 2)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 2)

	tlt = tile.Sand
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 2)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 2)

	tlt = tile.Sea
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 0)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 0)

	tlt = tile.Magma
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 0)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 0)

	tlt = tile.Ice
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 3)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 3)

	tlt = tile.Road
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 3)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 3)

	tlt = tile.Room
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 2)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 2)

	tlt = tile.Fog
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 4)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 4)

	tlt = tile.Smoke
	gTile3D[tlt] = NewTile3D_PlaneGeo(tlt, 4)
	gTile3DDark[tlt] = NewTile3D_PlaneGeo(tlt, 4)

	tlt = tile.Wall
	gTile3D[tlt] = NewTile3D_Wall(2)
	gTile3DDark[tlt] = NewTile3D_Wall(2)

	tlt = tile.Grass
	gTile3D[tlt] = NewTile3D_Grass(2)
	gTile3DDark[tlt] = NewTile3D_Grass(2)

	tlt = tile.Tree
	gTile3D[tlt] = NewTile3D_Tree(2)
	gTile3DDark[tlt] = NewTile3D_Tree(2)

	tlt = tile.Window
	gTile3D[tlt] = NewTile3D_BoxTile(gClientTile.CursorTiles[2], 2)
	gTile3DDark[tlt] = NewTile3D_BoxTile(gClientTile.CursorTiles[2], 2)

	tlt = tile.Door
	gTile3D[tlt] = NewTile3D_BoxTile(gClientTile.FloorTiles[tile.Door][0], 2)
	gTile3DDark[tlt] = NewTile3D_BoxTile(gClientTile.FloorTiles[tile.Door][0], 2)

	for i := 0; i < tile.Tile_Count; i++ {
		gTile3DDark[i].Mat.Set("metalness", 1.0)
		// gTile3DDark[i].Mat.Set("roughness", 0.0)
	}
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
