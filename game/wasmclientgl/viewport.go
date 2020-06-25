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
	"math/rand"
	"syscall/js"
	"time"

	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/lib/imagecanvas"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type Viewport struct {
	rnd *rand.Rand

	zoom int

	ViewWidth  int
	ViewHeight int

	DarkerTileImgCnv *imagecanvas.ImageCanvas

	CanvasGL      js.Value
	threejs       js.Value
	camera        js.Value
	light         js.Value
	renderer      js.Value
	fontLoader    js.Value
	textureLoader js.Value

	scene       js.Value
	jsSceneObjs map[string]js.Value

	// title
	font_helvetiker_regular js.Value
	jsoTitle                js.Value

	// cache
	colorMaterialCache map[uint32]js.Value
	textGeometryCache  map[textGeoKey]js.Value
}

func NewViewport() *Viewport {
	vp := &Viewport{
		rnd:                rand.New(rand.NewSource(time.Now().UnixNano())),
		jsSceneObjs:        make(map[string]js.Value),
		colorMaterialCache: make(map[uint32]js.Value),
		textGeometryCache:  make(map[textGeoKey]js.Value),
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
	gTextureTileList = LoadTextureTileList()
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
	vp.zoom = state
}

func (vp *Viewport) AddEventListener(evt string, fn func(this js.Value, args []js.Value) interface{}) {
	vp.CanvasGL.Call("addEventListener", evt, js.FuncOf(fn))
}

func (vp *Viewport) Draw(
	frameProgress float64,
	scrollDir way9type.Way9Type,
	taNoti *c2t_obj.NotiVPTiles_data,
) {

	sx, sy := vp.calcShiftDxDy(frameProgress)
	scrollDx := -scrollDir.Dx() * sx
	scrollDy := scrollDir.Dy() * sy

	// move camera, light
	cameraX := taNoti.VPX*DstCellSize + scrollDx
	cameraY := -taNoti.VPY*DstCellSize + scrollDy
	SetPosition(vp.light,
		cameraX, cameraY, DstCellSize*2,
	)
	cameraZ := HelperSize / (vp.zoom + 1)
	SetPosition(vp.camera,
		cameraX, cameraY, cameraZ,
	)
	vp.camera.Call("lookAt",
		vp.ThreeJsNew("Vector3",
			cameraX, cameraY, 0,
		),
	)

	vp.renderer.Call("render", vp.scene, vp.camera)
}

func (vp *Viewport) ThreeJsNew(name string, args ...interface{}) js.Value {
	return vp.threejs.Get(name).New(args...)
}

func (vp *Viewport) calcShiftDxDy(frameProgress float64) (int, int) {
	rate := 1 - frameProgress
	// if rate < 0 {
	// 	rate = 0
	// }
	// if rate > 1 {
	// 	rate = 1
	// }
	dx := int(float64(DstCellSize) * rate)
	dy := int(float64(DstCellSize) * rate)
	return dx, dy
}
