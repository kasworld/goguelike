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

	"github.com/kasworld/gowasmlib/jslog"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/lib/imagecanvas"
	"github.com/kasworld/wrapper"
)

type Viewport struct {
	rnd *rand.Rand

	ViewWidth  int
	ViewHeight int

	// for tile draw
	textureTileList         [tile.Tile_Count]*imagecanvas.ImageCanvas
	textureTileWrapInfoList [tile.Tile_Count]textureTileWrapInfo

	DarkerTileImgCnv *imagecanvas.ImageCanvas

	CanvasGL      js.Value
	threejs       js.Value
	camera        js.Value
	light         js.Value
	renderer      js.Value
	fontLoader    js.Value
	textureLoader js.Value

	scene       js.Value
	jsSceneObjs map[g2id.G2ID]js.Value

	// title
	font_helvetiker_regular js.Value
	jsoTitle                js.Value

	// terrain
	floorG2ID2ClientField map[g2id.G2ID]*ClientField

	// cache
	colorMaterialCache map[uint32]js.Value
	textGeometryCache  map[textGeoKey]js.Value
}

func NewViewport() *Viewport {
	vp := &Viewport{
		rnd:                   rand.New(rand.NewSource(time.Now().UnixNano())),
		jsSceneObjs:           make(map[g2id.G2ID]js.Value),
		floorG2ID2ClientField: make(map[g2id.G2ID]*ClientField),
		colorMaterialCache:    make(map[uint32]js.Value),
		textGeometryCache:     make(map[textGeoKey]js.Value),
	}

	vp.threejs = js.Global().Get("THREE")
	vp.renderer = vp.ThreeJsNew("WebGLRenderer")
	vp.CanvasGL = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvasglholder").Call("appendChild", vp.CanvasGL)
	vp.CanvasGL.Set("tabindex", "1")

	vp.scene = vp.ThreeJsNew("Scene")
	vp.camera = vp.ThreeJsNew("PerspectiveCamera", 60, 1, 1, HelperSize*2)
	vp.textureLoader = vp.ThreeJsNew("TextureLoader")
	vp.fontLoader = vp.ThreeJsNew("FontLoader")

	// for tile draw
	srcCellSize := 64
	for i, v := range tile.TileScrollAttrib {
		if v.Texture {
			idstr := fmt.Sprintf("%vPng", tile.Tile(i))
			vp.textureTileList[i] = imagecanvas.NewByID(idstr)
			vp.textureTileWrapInfoList[i] = textureTileWrapInfo{
				CellSize: srcCellSize,
				Xcount:   vp.textureTileList[i].W / srcCellSize,
				Ycount:   vp.textureTileList[i].H / srcCellSize,
				WrapX:    wrapper.New(vp.textureTileList[i].W - srcCellSize).WrapSafe,
				WrapY:    wrapper.New(vp.textureTileList[i].H - srcCellSize).WrapSafe,
			}
		}
	}
	vp.DarkerTileImgCnv = imagecanvas.NewByID("DarkerPng")

	vp.initHelpers()
	vp.initTitle()
	return vp
}

func (vp *Viewport) Hide() {
	vp.CanvasGL.Get("style").Set("display", "none")
}
func (vp *Viewport) Show() {
	vp.CanvasGL.Get("style").Set("display", "initial")
}

func (vp *Viewport) ResizeCanvas(title bool) {
	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Int()
	winH := win.Get("innerHeight").Int()
	if title {
		winH /= 3
	}
	vp.ViewWidth = winW
	vp.ViewHeight = winH

	vp.CanvasGL.Call("setAttribute", "width", winW)
	vp.CanvasGL.Call("setAttribute", "height", winH)

	vp.camera.Set("aspect", float64(winW)/float64(winH))
	vp.camera.Call("updateProjectionMatrix")

	vp.renderer.Call("setSize", winW, winH)
}

func (vp *Viewport) Focus() {
	vp.CanvasGL.Call("focus")
}

func (vp *Viewport) Zoom(state int) {
	fov := [3]float64{
		30, 50, 70,
	}
	vp.camera.Set("fov", fov[state])
	vp.camera.Call("updateProjectionMatrix")
	jslog.Infof("fov %v", fov[state])
}

func (vp *Viewport) AddEventListener(evt string, fn func(this js.Value, args []js.Value) interface{}) {
	vp.CanvasGL.Call("addEventListener", evt, js.FuncOf(fn))
}

func (vp *Viewport) Draw() {
	vp.renderer.Call("render", vp.scene, vp.camera)
}

func (vp *Viewport) ThreeJsNew(name string, args ...interface{}) js.Value {
	return vp.threejs.Get(name).New(args...)
}
