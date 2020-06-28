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
	"fmt"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/clientinitdata"
	"github.com/kasworld/goguelike/lib/clienttile"
)

const (
	DisplayLineLimit = 3*gameconst.ViewPortH - gameconst.ViewPortH/2

	DstCellSize = 32
	// TextYShift  = 24
	HelperSize = 1024.0
)

var gInitData *clientinitdata.InitData = clientinitdata.New()
var gClientTile *clienttile.ClientTile = clienttile.New()
var gTextureTileList [tile.Tile_Count]*TextureTile = LoadTextureTileList()

var gRnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

var gTextureLoader js.Value = ThreeJsNew("TextureLoader")

var gColorMaterialCache map[uint32]js.Value = make(map[uint32]js.Value)

func getColorMaterial(co uint32) js.Value {
	mat, exist := gColorMaterialCache[co]
	if !exist {
		mat = ThreeJsNew("MeshPhongMaterial",
			map[string]interface{}{
				"color": co,
			},
		)
		gColorMaterialCache[co] = mat
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

func getTextGeometry(str string, size float64) js.Value {
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

func calcGeoMinMaxX(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("x").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("x").Float()
	return geoMin, geoMax
}

func calcGeoMinMaxY(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("y").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("y").Float()
	return geoMin, geoMax
}

func CalcCurrentFrame(difftick int64, fps float64) int {
	diffsec := float64(difftick) / float64(time.Second)
	frame := fps * diffsec
	return int(frame)
}

func getImgWH(srcImageID string) (js.Value, float64, float64) {
	img := js.Global().Get("document").Call("getElementById", srcImageID)
	if !img.Truthy() {
		fmt.Printf("fail to get %v", srcImageID)
		return js.Null(), 0, 0
	}
	srcw := img.Get("naturalWidth").Float()
	srch := img.Get("naturalHeight").Float()
	return img, srcw, srch
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
