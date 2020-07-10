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

	"github.com/kasworld/goguelike/config/viewportdata"
)

type SightPlane struct {
	W   int // canvas width
	H   int // canvas height
	Cnv js.Value
	Ctx js.Value
	Tex js.Value
	// Mat  js.Value
	// Geo  js.Value
	Mesh js.Value
}

func NewSightPlane() *SightPlane {
	w := ClientViewLen * DstCellSize
	h := ClientViewLen * DstCellSize
	Cnv := js.Global().Get("document").Call("createElement",
		"CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", w)
	Cnv.Set("height", h)

	Tex := ThreeJsNew("CanvasTexture", Cnv)
	Mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	Mat.Set("transparent", true)
	Geo := ThreeJsNew("PlaneGeometry", w, h)
	Mesh := ThreeJsNew("Mesh", Geo, Mat)

	return &SightPlane{
		W:    w,
		H:    h,
		Cnv:  Cnv,
		Ctx:  Ctx,
		Tex:  Tex,
		Mesh: Mesh,
	}
}

func (pl *SightPlane) MoveCenterTo(fx, fy int) {
	SetPosition(pl.Mesh,
		fx*DstCellSize,
		-fy*DstCellSize,
		DstCellSize/2)
}

func (pl *SightPlane) FillColor(co string) {
	pl.Ctx.Set("fillStyle", co)
	pl.Ctx.Call("fillRect", 0, 0, pl.W, pl.H)
}

func (pl *SightPlane) ClearRect() {
	pl.Ctx.Call("clearRect", 0, 0, pl.W, pl.H)
}

func (pl *SightPlane) ClearSight(vpTiles *viewportdata.ViewportTileArea2) {
	cx := ClientViewLen / 2
	cy := ClientViewLen / 2
	for i, v := range gInitData.ViewportXYLenList {
		if vpTiles[i] == 0 {
			continue
		}
		posx := (cx + v.X) * DstCellSize
		posy := (cy + v.Y) * DstCellSize
		pl.Ctx.Call("clearRect", posx, posy, DstCellSize, DstCellSize)
	}
	pl.Tex.Set("needsUpdate", true)
}

func (pl *SightPlane) Death() {
	pl.ClearRect()
	pl.FillColor("#000000a0")
	pl.Tex.Set("needsUpdate", true)

}

func (pl *SightPlane) Ready2Rebirth() {
	pl.ClearRect()
	pl.FillColor("#00000020")
	pl.Tex.Set("needsUpdate", true)

}
