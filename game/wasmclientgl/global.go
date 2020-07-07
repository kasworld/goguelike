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
	"math/rand"
	"syscall/js"
	"time"

	"github.com/kasworld/findnear"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/clientinitdata"
	"github.com/kasworld/goguelike/lib/clienttile"
	"github.com/kasworld/goguelike/lib/webtilegroup"
)

const (
	DisplayLineLimit = 3*gameconst.ViewPortH - gameconst.ViewPortH/2
	DstCellSize      = 32
	HelperSize       = DstCellSize * 32
	ClientViewLen    = 40
)

var gRnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

var gInitData *clientinitdata.InitData = clientinitdata.New()
var gClientTile *clienttile.ClientTile = clienttile.New()
var gTextureTileList [tile.Tile_Count]*TextureTile = LoadTextureTileList()

var gXYLenListView = findnear.NewXYLenList(ClientViewLen, ClientViewLen)

var gTileMaterial [tile.Tile_Count]js.Value
var gTileGeometry [tile.Tile_Count]js.Value
var gTileShift [tile.Tile_Count][3]float64

func preMakeTileMatGeo() {
	var tlt tile.Tile

	tlt = tile.Swamp
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	gTileShift[tlt] = [3]float64{0, 0, -1}

	tlt = tile.Soil
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)

	tlt = tile.Stone
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)

	tlt = tile.Sand
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)

	tlt = tile.Sea
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	gTileShift[tlt] = [3]float64{0, 0, -2}

	tlt = tile.Magma
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	gTileShift[tlt] = [3]float64{0, 0, -2}

	tlt = tile.Ice
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)

	tlt = tile.Grass
	gTileMaterial[tlt] = NewTextureTileMaterial(tile.Grass)
	gTileGeometry[tlt] = ThreeJsNew("BoxGeometry", DstCellSize, DstCellSize, DstCellSize/8)

	tlt = tile.Tree
	gTileMaterial[tlt] = NewTextureTileMaterial(tile.Grass)
	gTileGeometry[tlt] = ThreeJsNew("ConeGeometry", DstCellSize/2-1, DstCellSize-1)
	gTileGeometry[tlt].Call("rotateX", math.Pi/2)

	tlt = tile.Road
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	gTileShift[tlt] = [3]float64{0, 0, 1}

	tlt = tile.Room
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	gTileShift[tlt] = [3]float64{0, 0, 1}

	tlt = tile.Wall
	gTileMaterial[tlt] = NewTextureTileMaterial(tile.Stone)
	gTileGeometry[tlt] = ThreeJsNew("BoxGeometry", DstCellSize, DstCellSize, DstCellSize)

	tlt = tile.Window
	gTileMaterial[tlt] = NewTileMaterial(gClientTile.CursorTiles[2])
	gTileGeometry[tlt] = ThreeJsNew("BoxGeometry", DstCellSize, DstCellSize, DstCellSize)

	tlt = tile.Door
	gTileMaterial[tlt] = NewTileMaterial(gClientTile.FloorTiles[tile.Door][0])
	gTileGeometry[tlt] = ThreeJsNew("BoxGeometry", DstCellSize, DstCellSize, DstCellSize)

	tlt = tile.Fog
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	gTileShift[tlt] = [3]float64{0, 0, DstCellSize/8 + 1}

	tlt = tile.Smoke
	gTileMaterial[tlt] = NewTextureTileMaterial(tlt)
	gTileGeometry[tlt] = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	gTileShift[tlt] = [3]float64{0, 0, DstCellSize/8 + 1}

}

var gTextureLoader js.Value = ThreeJsNew("TextureLoader")

var gColorMaterialCache map[string]js.Value = make(map[string]js.Value)

func GetColorMaterialByCache(co string) js.Value {
	mat, exist := gColorMaterialCache[co]
	if !exist {
		mat = ThreeJsNew("MeshPhongMaterial",
			map[string]interface{}{
				"color": co,
			},
		)
		mat.Set("transparent", true)
		gColorMaterialCache[co] = mat
	}
	return mat
}

func NewTileMaterial(ti webtilegroup.TileInfo) js.Value {
	Cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", DstCellSize)
	Cnv.Set("height", DstCellSize)
	Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)

	Tex := ThreeJsNew("CanvasTexture", Cnv)
	mat := ThreeJsNew("MeshPhongMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	mat.Set("transparent", true)
	// mat.Set("side", ThreeJs().Get("DoubleSide"))
	return mat
}

var gTileMaterialCache map[webtilegroup.TileInfo]js.Value = make(map[webtilegroup.TileInfo]js.Value)

func GetTileMaterialByCache(ti webtilegroup.TileInfo) js.Value {
	mat, exist := gTileMaterialCache[ti]
	if !exist {
		mat = NewTileMaterial(ti)
		gTileMaterialCache[ti] = mat
	}
	return mat
}

func NewTextureTileMaterial(ti tile.Tile) js.Value {
	Cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", DstCellSize)
	Cnv.Set("height", DstCellSize)

	tlic := gTextureTileList[ti]
	Ctx.Call("drawImage", tlic.Cnv,
		0, 0, DstCellSize, DstCellSize,
		0, 0, DstCellSize, DstCellSize)

	Tex := ThreeJsNew("CanvasTexture", Cnv)
	mat := ThreeJsNew("MeshPhongMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	mat.Set("transparent", true)
	// mat.Set("side", ThreeJs().Get("DoubleSide"))
	return mat
}

var gTextureTileMaterialCache map[tile.Tile]js.Value = make(map[tile.Tile]js.Value)

func GetTextureTileMaterialByCache(ti tile.Tile) js.Value {
	mat, exist := gTextureTileMaterialCache[ti]
	if !exist {
		mat = NewTextureTileMaterial(ti)
		gTextureTileMaterialCache[ti] = mat
	}
	return mat
}

type textGeoKey struct {
	Str  string
	Size float64
}

var gFontLoader js.Value = ThreeJsNew("FontLoader")
var gFont_helvetiker_regular js.Value

var gTextGeometryCache map[textGeoKey]js.Value = make(map[textGeoKey]js.Value)

func GetTextGeometryByCache(str string, size float64) js.Value {
	geo, exist := gTextGeometryCache[textGeoKey{str, size}]
	curveSegments := size / 3
	if curveSegments < 1 {
		curveSegments = 1
	}
	bevelEnabled := true
	if size < 16 {
		bevelEnabled = false
	}
	bevelThickness := size / 8
	if bevelThickness < 1 {
		bevelThickness = 1
	}
	bevelSize := size / 16
	if bevelSize < 1 {
		bevelSize = 1
	}
	bevelSegments := size / 8
	if bevelSegments < 1 {
		bevelSegments = 1
	}
	if !exist {
		geo = ThreeJsNew("TextGeometry", str,
			map[string]interface{}{
				"font":           gFont_helvetiker_regular,
				"size":           size,
				"height":         5,
				"curveSegments":  curveSegments,
				"bevelEnabled":   bevelEnabled,
				"bevelThickness": bevelThickness,
				"bevelSize":      bevelSize,
				"bevelOffset":    0,
				"bevelSegments":  bevelSegments,
			})
		gTextGeometryCache[textGeoKey{str, size}] = geo
	}
	return geo
}

var gBoxGeometryCache map[[3]int]js.Value = make(map[[3]int]js.Value)

func GetBoxGeometryByCache(x, y, z int) js.Value {
	geo, exist := gBoxGeometryCache[[3]int{x, y, z}]
	if !exist {
		geo = ThreeJsNew("BoxGeometry", x, y, z)
		gBoxGeometryCache[[3]int{x, y, z}] = geo
	}
	return geo
}

var gPlaneGeometryCache map[[2]int]js.Value = make(map[[2]int]js.Value)

func GetPlaneGeometryByCache(w, h int) js.Value {
	geo, exist := gPlaneGeometryCache[[2]int{w, h}]
	if !exist {
		geo = ThreeJsNew("PlaneGeometry", w, h)
		gPlaneGeometryCache[[2]int{w, h}] = geo
	}
	return geo
}

var gConeGeometryCache map[[2]int]js.Value = make(map[[2]int]js.Value)

func GetConeGeometryByCache(r, h int) js.Value {
	geo, exist := gConeGeometryCache[[2]int{r, h}]
	if !exist {
		geo = ThreeJsNew("ConeGeometry", r, h)
		geo.Call("rotateX", math.Pi/2)
		gConeGeometryCache[[2]int{r, h}] = geo
	}
	return geo
}

var gSphereGeometryCache map[[3]int]js.Value = make(map[[3]int]js.Value)

func GetSphereGeometryByCache(r, ws, hs int) js.Value {
	geo, exist := gSphereGeometryCache[[3]int{r, ws, hs}]
	if !exist {
		geo = ThreeJsNew("SphereGeometry", r, ws, hs)
		// geo.Call("rotateX", math.Pi)
		gSphereGeometryCache[[3]int{r, ws, hs}] = geo
	}
	return geo
}

func CalcGeoMinMaxX(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("x").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("x").Float()
	return geoMin, geoMax
}

func CalcGeoMinMaxY(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("y").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("y").Float()
	return geoMin, geoMax
}

func CalcGeoMinMaxZ(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("z").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("z").Float()
	return geoMin, geoMax
}

func CalcCurrentFrame(difftick int64, fps float64) int {
	diffsec := float64(difftick) / float64(time.Second)
	frame := fps * diffsec
	return int(frame)
}

func SetPosition(jso js.Value, pos ...interface{}) {
	po := jso.Get("position")
	if len(pos) >= 1 {
		po.Set("x", pos[0])
	}
	if len(pos) >= 2 {
		po.Set("y", pos[1])
	}
	if len(pos) >= 3 {
		po.Set("z", pos[2])
	}
}

func ThreeJsNew(name string, args ...interface{}) js.Value {
	return js.Global().Get("THREE").Get(name).New(args...)
}

func ThreeJs() js.Value {
	return js.Global().Get("THREE")
}

func GetElementById(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}
