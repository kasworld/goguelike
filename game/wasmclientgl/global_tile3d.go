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
	"math"
	"syscall/js"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
)

var gTile3D [tile.Tile_Count]Tile3D
var gTile3DDark [tile.Tile_Count]Tile3D

func preMakeTileMatGeo() {
	var tlt tile.Tile

	tlt = tile.Swamp
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 1},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 1},
	}

	tlt = tile.Soil
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Stone
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Sand
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Sea
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 0},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 0},
	}

	tlt = tile.Magma
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 0},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 0},
	}

	tlt = tile.Ice
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 3},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 3},
	}

	tlt = tile.Grass
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tile.Grass),
		Geo:   ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize/8),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tile.Grass),
		Geo:   ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize/8),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Tree
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tile.Grass),
		Geo:   MakeTreeGeo(),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tile.Grass),
		Geo:   MakeTreeGeo(),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Road
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 3},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 3},
	}

	tlt = tile.Room
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Wall
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tile.Stone),
		Geo:   ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tile.Stone),
		Geo:   ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Window
	gTile3D[tlt] = Tile3D{
		Mat:   NewTileMaterial(gClientTile.CursorTiles[2]),
		Geo:   ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTileMaterial(gClientTile.CursorTiles[2]),
		Geo:   ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Door
	gTile3D[tlt] = Tile3D{
		Mat:   NewTileMaterial(gClientTile.FloorTiles[tile.Door][0]),
		Geo:   ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize-1),
		Shift: [3]float64{0, 0, 2},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTileMaterial(gClientTile.FloorTiles[tile.Door][0]),
		Geo:   ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize-1),
		Shift: [3]float64{0, 0, 2},
	}

	tlt = tile.Fog
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 4},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 4},
	}

	tlt = tile.Smoke
	gTile3D[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 4},
	}
	gTile3DDark[tlt] = Tile3D{
		Mat:   NewTextureTileMaterial(tlt),
		Geo:   ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize),
		Shift: [3]float64{0, 0, 4},
	}

	for i := 0; i < tile.Tile_Count; i++ {
		gTile3D[i].GeoInfo = GetGeoInfo(gTile3D[i].Geo)
		// gTile3D[i].Mat.Set("metalness", 1.0)
		// gTile3D[i].Mat.Set("roughness", 0.0)

		gTile3DDark[i].GeoInfo = GetGeoInfo(gTile3DDark[i].Geo)
		gTile3DDark[i].Mat.Set("metalness", 1.0)
		// gTile3DDark[i].Mat.Set("roughness", 0.0)
	}
}

func MakeTreeGeo() js.Value {
	matrix := ThreeJsNew("Matrix4")
	geo := ThreeJsNew("CylinderGeometry", 2, 4, DstCellSize-1)
	geo.Call("rotateX", math.Pi/2)
	geo1 := ThreeJsNew("ConeGeometry", DstCellSize/3-2, DstCellSize/2-2)
	geo1.Call("rotateX", math.Pi/2)
	geo2 := ThreeJsNew("ConeGeometry", DstCellSize/2-2, DstCellSize/2-2)
	geo2.Call("rotateX", math.Pi/2)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, 0, DstCellSize/2-2,
	))
	geo.Call("merge", geo1, matrix)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, 0, 5,
	))
	geo.Call("merge", geo2, matrix)
	geo1.Call("dispose")
	geo2.Call("dispose")
	return geo
}

func MakeWallGeo() js.Value {
	shape := ThreeJsNew("Shape")
	shape.Call("moveTo", 0, 0)
	shape.Call("lineTo", 0, DstCellSize-1)
	shape.Call("lineTo", DstCellSize-1, DstCellSize-1)
	shape.Call("lineTo", DstCellSize-1, 0)
	shape.Call("lineTo", 0, 0)
	extrudeSettings := map[string]interface{}{
		"steps":          1,
		"depth":          DstCellSize - 1,
		"bevelEnabled":   true,
		"bevelThickness": 1,
		"bevelSize":      1,
		"bevelOffset":    0,
		"bevelSegments":  1,
	}
	geo := ThreeJsNew("ExtrudeGeometry", shape, extrudeSettings)
	geo.Call("translate", -DstCellSize/2, -DstCellSize/2, -DstCellSize/2)
	return geo
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
