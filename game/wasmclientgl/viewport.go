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
	"syscall/js"
)

type Viewport struct {
	CanvasGL js.Value
	renderer js.Value
	jsMouse  js.Value
}

func NewViewport() *Viewport {
	vp := &Viewport{}

	vp.jsMouse = ThreeJsNew("Vector2")
	vp.renderer = ThreeJsNew("WebGLRenderer")
	vp.CanvasGL = vp.renderer.Get("domElement")
	GetElementById("canvasglholder").Call("appendChild", vp.CanvasGL)
	vp.CanvasGL.Set("tabindex", "1")

	return vp
}

func (vp *Viewport) Hide() {
	vp.CanvasGL.Get("style").Set("display", "none")
}
func (vp *Viewport) Show() {
	vp.CanvasGL.Get("style").Set("display", "initial")
}

func (vp *Viewport) Resize(w, h float64) {
	vp.CanvasGL.Call("setAttribute", "width", w)
	vp.CanvasGL.Call("setAttribute", "height", h)
	vp.renderer.Call("setSize", w, h)
}
