// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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

	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

type PlaneLayer struct {
	W   int // canvas width
	H   int // canvas height
	Cnv js.Value
	Ctx js.Value
	Tex js.Value
	// Mat  js.Value
	// Geo  js.Value
	Mesh js.Value
}

func NewPlaneLayer(fi *c2t_obj.FloorInfo, zpos int) *PlaneLayer {
	w := fi.W * DstCellSize
	h := fi.H * DstCellSize
	xRepeat := 3
	yRepeat := 3
	Cnv := js.Global().Get("document").Call("createElement",
		"CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	if !Ctx.Truthy() {
		jslog.Errorf("fail to get context %v", fi)
		return nil
	}
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", w)
	Cnv.Set("height", h)
	// Ctx.Call("clearRect", 0, 0, w, h)

	Tex := ThreeJsNew("CanvasTexture", Cnv)
	Tex.Set("wrapS", ThreeJs().Get("RepeatWrapping"))
	Tex.Set("wrapT", ThreeJs().Get("RepeatWrapping"))
	Tex.Get("repeat").Set("x", xRepeat)
	Tex.Get("repeat").Set("y", yRepeat)
	Mat := ThreeJsNew("MeshBasicMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	Geo := ThreeJsNew("PlaneBufferGeometry",
		w*xRepeat, h*yRepeat)
	Mesh := ThreeJsNew("Mesh", Geo, Mat)

	SetPosition(Mesh, w/2, -h/2, zpos)

	return &PlaneLayer{
		W:    w,
		H:    h,
		Cnv:  Cnv,
		Ctx:  Ctx,
		Tex:  Tex,
		Mesh: Mesh,
	}
}

func (pl *PlaneLayer) drawBG(co string) {
	pl.Ctx.Set("fillStyle", co)
	pl.Ctx.Call("fillRect", 0, 0, pl.W, pl.H)
	pl.Tex.Set("needsUpdate", true)
}
