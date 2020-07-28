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

	"github.com/kasworld/goguelike/enum/tile_flag"

	"github.com/kasworld/goguelike/lib/webtilegroup"
)

type Cursor3D struct {
	Cnv     js.Value
	Ctx     js.Value
	Tex     js.Value
	GeoInfo GeoInfo
	Mesh    js.Value
}

func NewCursor3D() *Cursor3D {
	cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	ctx := cnv.Call("getContext", "2d")
	ctx.Set("imageSmoothingEnabled", false)
	cnv.Set("width", DstCellSize)
	cnv.Set("height", DstCellSize)
	tex := ThreeJsNew("CanvasTexture", cnv)
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"map": tex,
		},
	)
	mat.Set("transparent", true)

	geo := ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &Cursor3D{
		Cnv:     cnv,
		Ctx:     ctx,
		Tex:     tex,
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func (aog *Cursor3D) ChangeTile(ti webtilegroup.TileInfo) {
	aog.Ctx.Call("clearRect", 0, 0, DstCellSize, DstCellSize)
	aog.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)
	aog.Tex.Set("needsUpdate", true)
}

func (aog *Cursor3D) SetFieldPosition(fx, fy int, tl tile_flag.TileFlag) {
	height := GetTile3DHeightByCache(tl)
	if !tl.CharPlaceable() {
		aog.ChangeTile(gClientTile.CursorTiles[2])
	} else {
		if tl.NoBattle() {
			aog.ChangeTile(gClientTile.CursorTiles[0])

		} else {
			aog.ChangeTile(gClientTile.CursorTiles[1])
		}
	}
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+DstCellSize/2,
		-float64(fy)*DstCellSize-DstCellSize/2,
		aog.GeoInfo.Len[2]/2+1+height,
	)
}

func (aog *Cursor3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *Cursor3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *Cursor3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *Cursor3D) Dispose() {
	// mesh do not need dispose
	aog.Mesh.Get("geometry").Call("dispose")
	aog.Mesh.Get("material").Call("dispose")
	aog.Tex.Call("dispose")

	aog.Cnv = js.Undefined()
	aog.Ctx = js.Undefined()
	aog.Mesh = js.Undefined()
	aog.Tex = js.Undefined()
	// no need createElement canvas dom obj
}
