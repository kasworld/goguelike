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
)

type RaycastPlane struct {
	W    int // canvas width
	H    int // canvas height
	Cnv  js.Value
	Ctx  js.Value
	Tex  js.Value
	Mesh js.Value
}

func NewRaycastPlane() *RaycastPlane {
	w := ClientViewLen * DstCellSize
	h := ClientViewLen * DstCellSize
	Cnv := js.Global().Get("document").Call("createElement",
		"CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", w)
	Cnv.Set("height", h)
	Ctx.Call("clearRect", 0, 0, w, h)
	Tex := ThreeJsNew("CanvasTexture", Cnv)
	Mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	Mat.Set("transparent", true)
	Geo := ThreeJsNew("PlaneGeometry", w, h)
	Mesh := ThreeJsNew("Mesh", Geo, Mat)

	return &RaycastPlane{
		W:    w,
		H:    h,
		Cnv:  Cnv,
		Ctx:  Ctx,
		Tex:  Tex,
		Mesh: Mesh,
	}
}

func (pl *RaycastPlane) MoveCenterTo(fx, fy int) {
	SetPosition(pl.Mesh,
		fx*DstCellSize,
		-fy*DstCellSize,
		0,
	)
}
