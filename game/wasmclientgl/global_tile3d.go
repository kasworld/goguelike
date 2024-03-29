// Copyright 2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
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

func CalcTile3DVisibleTop(tl tile_flag.TileFlag) float64 {
	rtn := 0.0
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

func CalcTile3DStepOn(tl tile_flag.TileFlag) float64 {
	rtn := 0.0
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
	tile.Tree:   {UnitTileZ * 16.0, UnitTileZ * 3.0, UnitTileZ * 3.0},
	tile.Road:   {UnitTileZ * 3.0, UnitTileZ * 3.0, UnitTileZ * 6.0},
	tile.Room:   {UnitTileZ * 3.0, UnitTileZ * 3.0, UnitTileZ * 6.0},
	tile.Wall:   {UnitTileZ * 16.0, UnitTileZ * 3.0, UnitTileZ * 19.0},
	tile.Window: {UnitTileZ * 16.0, UnitTileZ * 3.0, UnitTileZ * 19.0},
	tile.Door:   {UnitTileZ * 16.0, UnitTileZ * 3.0, UnitTileZ * 3.0},
	tile.Fog:    {UnitTileZ * 1.0, UnitTileZ * 6.0, UnitTileZ * 3.0},
	tile.Smoke:  {UnitTileZ * 1.0, UnitTileZ * 6.0, UnitTileZ * 3.0},
}

func preMakeTileMatGeo() {
	var tlt tile.Tile

	tlt = tile.Swamp
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 0.5)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 0.5).MakeSrcDark()

	tlt = tile.Soil
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 1.0)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 1.0).MakeSrcDark()

	tlt = tile.Stone
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 1.0)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 1.0).MakeSrcDark()

	tlt = tile.Sand
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 1.0)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 1.0).MakeSrcDark()

	tlt = tile.Sea
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 0.5)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 0.5).MakeSrcDark()

	tlt = tile.Magma
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 1.0)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 1.0).MakeSrcDark()

	tlt = tile.Ice
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 0.5)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 0.5).MakeSrcDark()

	tlt = tile.Road
	gTile3D[tlt] = NewTile3D_OctCylinderTexture(tlt, 1.0)
	gTile3DDark[tlt] = NewTile3D_OctCylinderTexture(tlt, 1.0).MakeSrcDark()

	tlt = tile.Room
	gTile3D[tlt] = NewTile3D_OctCylinderTexture(tlt, 1.0)
	gTile3DDark[tlt] = NewTile3D_OctCylinderTexture(tlt, 1.0).MakeSrcDark()

	tlt = tile.Fog
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 0.5)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 0.5).MakeSrcDark()

	tlt = tile.Smoke
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 0.5)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 0.5).MakeSrcDark()

	tlt = tile.Grass
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 1.0)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 1.0).MakeSrcDark()

	tlt = tile.Tree
	gTile3D[tlt] = NewTile3D_Tree(1.0)
	gTile3DDark[tlt] = NewTile3D_Tree(1.0).MakeSrcDark()

	tlt = tile.Wall
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 1.0)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 1.0).MakeSrcDark()

	tlt = tile.Window
	gTile3D[tlt] = NewTile3D_BoxTexture(tlt, 0.5)
	gTile3DDark[tlt] = NewTile3D_BoxTexture(tlt, 0.5).MakeSrcDark()

	tlt = tile.Door
	gTile3D[tlt] = NewTile3D_Door(0.5)
	gTile3DDark[tlt] = NewTile3D_Door(0.5)

	// for i, tl := range gTile3D {
	// 	jslog.Infof("%v %v %v", tile.Tile(i), tl.GeoInfo, gTileZInfo[i])
	// }
}
