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

	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/lib/g2id"
)

type Viewport struct {
	rnd        *rand.Rand
	ViewWidth  int
	ViewHeight int

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
	floorMeshCache map[g2id.G2ID]js.Value

	// cache
	colorMaterialCache map[uint32]js.Value

	aoGeometryCache map[factiontype.FactionType]js.Value
	foGeometryCache map[fieldobjdisplaytype.FieldObjDisplayType]js.Value

	eqGeometryCache     map[equipslottype.EquipSlotType]js.Value
	potionGeometryCache map[potiontype.PotionType]js.Value
	scrollGeometryCache map[scrolltype.ScrollType]js.Value
	moneyGeometry       js.Value
}

func NewViewport() *Viewport {
	vp := &Viewport{
		rnd:                rand.New(rand.NewSource(time.Now().UnixNano())),
		jsSceneObjs:        make(map[g2id.G2ID]js.Value),
		aoGeometryCache:    make(map[factiontype.FactionType]js.Value),
		colorMaterialCache: make(map[uint32]js.Value),
	}

	vp.threejs = js.Global().Get("THREE")
	vp.renderer = vp.ThreeJsNew("WebGLRenderer")
	vp.CanvasGL = vp.renderer.Get("domElement")
	js.Global().Get("document").Call("getElementById", "canvasglholder").Call("appendChild", vp.CanvasGL)
	vp.CanvasGL.Set("tabindex", "1")

	vp.scene = vp.ThreeJsNew("Scene")
	vp.camera = vp.ThreeJsNew("PerspectiveCamera", 75, 1, 1, StageSize*2)
	vp.textureLoader = vp.ThreeJsNew("TextureLoader")
	vp.fontLoader = vp.ThreeJsNew("FontLoader")
	vp.initGrid()
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
	vp.CanvasGL.Call("setAttribute", "width", winW)
	vp.CanvasGL.Call("setAttribute", "height", winH)
	vp.ViewWidth = winW
	vp.ViewHeight = winH

	vp.camera.Set("aspect", float64(winW)/float64(winH))
	vp.camera.Call("updateProjectionMatrix")

	vp.CanvasGL.Call("setAttribute", "width", winW)
	vp.CanvasGL.Call("setAttribute", "height", winH)
	vp.renderer.Call("setSize", winW, winH)
}

func (vp *Viewport) Focus() {
	vp.CanvasGL.Call("focus")
}

func (vp *Viewport) Zoom(state int) {
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
