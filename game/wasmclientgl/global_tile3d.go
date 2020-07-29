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

var gTileZInfo = [tile.Tile_Count]struct {
	Size     float64
	Shift    float64
	OnHeight float64 // for step on can be not visible
}{
	tile.Swamp:  {UnitTileZ * 2.0, UnitTileZ * 0.0, UnitTileZ * 2.0},
	tile.Soil:   {UnitTileZ * 3.0, UnitTileZ * 0.0, UnitTileZ * 3.0},
	tile.Stone:  {UnitTileZ * 3.0, UnitTileZ * 0.0, UnitTileZ * 3.0},
	tile.Sand:   {UnitTileZ * 3.0, UnitTileZ * 0.0, UnitTileZ * 3.0},
	tile.Sea:    {UnitTileZ * 1.0, UnitTileZ * 0.0, UnitTileZ * 1.0},
	tile.Magma:  {UnitTileZ * 1.0, UnitTileZ * 0.0, UnitTileZ * 1.0},
	tile.Ice:    {UnitTileZ * 4.0, UnitTileZ * 0.0, UnitTileZ * 4.0},
	tile.Grass:  {UnitTileZ * 5.0, UnitTileZ * 0.0, UnitTileZ * 3.0},
	tile.Tree:   {UnitTileZ * 0.0, UnitTileZ * 3.0, UnitTileZ * 3.0},
	tile.Road:   {UnitTileZ * 3.0, UnitTileZ * 3.0, UnitTileZ * 6.0},
	tile.Room:   {UnitTileZ * 3.0, UnitTileZ * 3.0, UnitTileZ * 6.0},
	tile.Wall:   {UnitTileZ * 16.0, UnitTileZ * 3.0, UnitTileZ * 18.0},
	tile.Window: {UnitTileZ * 16.0, UnitTileZ * 3.0, UnitTileZ * 18.0},
	tile.Door:   {UnitTileZ * 16.0, UnitTileZ * 3.0, UnitTileZ * 3.0},
	tile.Fog:    {UnitTileZ * 1.0, UnitTileZ * 6.0, UnitTileZ * 3.0},
	tile.Smoke:  {UnitTileZ * 1.0, UnitTileZ * 6.0, UnitTileZ * 3.0},
}

func preMakeTileMatGeo() {
	var tlt tile.Tile

	tlt = tile.Swamp
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Soil
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Stone
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Sand
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Sea
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Magma
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Ice
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Road
	gTile3D[tlt] = NewTile3D_OctCylinderTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_OctCylinderTexture(tlt).MakeSrcDark()

	tlt = tile.Room
	gTile3D[tlt] = NewTile3D_OctCylinderTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_OctCylinderTexture(tlt).MakeSrcDark()

	tlt = tile.Fog
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Smoke
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

	tlt = tile.Grass
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt).MakeSrcDark()

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
	gTile3D[tlt] = NewTile3D_Door()
	gTile3DDark[tlt] = NewTile3D_Door()

	for i := 0; i < tile.Tile_Count; i++ {
		gTile3D[i].Shift[2] = gTileZInfo[i].Shift
		gTile3DDark[i].Shift[2] = gTileZInfo[i].Shift
	}
	// for i, tl := range gTile3D {
	// 	jslog.Infof("%v %v", tile.Tile(i), tl.GeoInfo)
	// }
}

// for view
var tileflagTopCache [1 << uint(tile.Tile_Count)]float64

// for place obj step
var tileflagOnCache [1 << uint(tile.Tile_Count)]float64

func init() {
	for i := range tileflagTopCache {
		tileflagTopCache[i] = Tile3DHeightMin
	}
	for i := range tileflagOnCache {
		tileflagOnCache[i] = Tile3DHeightMin
	}
}

func calcTile3DTop(tl tile_flag.TileFlag) float64 {
	rtn := Tile3DHeightMin
	for i := 0; i < tile.Tile_Count; i++ {
		if !tl.TestByTile(tile.Tile(i)) {
			continue
		}
		z := gTileZInfo[i].Size + gTileZInfo[i].Shift
		if z > rtn {
			rtn = z
		}
	}
	return rtn
}

// top height for view, use for cursor
func GetTile3DTopByCache(tl tile_flag.TileFlag) float64 {
	z := tileflagTopCache[tl]
	if z == Tile3DHeightMin {
		z = calcTile3DTop(tl)
		if z == Tile3DHeightMin { // empty tile
			z = 0 // prevent recalc empty tile
		}
		tileflagTopCache[tl] = z
	}
	return z
}

func calcTile3DOn(tl tile_flag.TileFlag) float64 {
	rtn := Tile3DHeightMin
	for i := 0; i < tile.Tile_Count; i++ {
		if !tl.TestByTile(tile.Tile(i)) {
			continue
		}
		z := gTileZInfo[i].OnHeight
		if z > rtn {
			rtn = z
		}
	}
	return rtn
}

// height for step on
func GetTile3DOnByCache(tl tile_flag.TileFlag) float64 {
	z := tileflagTopCache[tl]
	if z == Tile3DHeightMin {
		z = calcTile3DOn(tl)
		if z == Tile3DHeightMin { // empty tile
			z = 0 // prevent recalc empty tile
		}
		tileflagTopCache[tl] = z
	}
	return z
}
