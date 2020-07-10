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

	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"

	"github.com/kasworld/goguelike/lib/webtilegroup"
)

type FieldObj3D struct {
	Cnv     js.Value
	Ctx     js.Value
	Tex     js.Value
	GeoInfo GeoInfo
	Mesh    js.Value
}

func NewFieldObj3D() *FieldObj3D {
	cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	ctx := cnv.Call("getContext", "2d")
	ctx.Set("imageSmoothingEnabled", false)
	cnv.Set("width", DstCellSize)
	cnv.Set("height", DstCellSize)
	tex := ThreeJsNew("CanvasTexture", cnv)
	// mat := ThreeJsNew("MeshPhongMaterial",
	mat := ThreeJsNew("MeshBasicMaterial",
		map[string]interface{}{
			"map": tex,
		},
	)
	mat.Set("transparent", true)
	geo := ThreeJsNew("BoxGeometry", DstCellSize, DstCellSize, DstCellSize)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &FieldObj3D{
		Cnv:     cnv,
		Ctx:     ctx,
		Tex:     tex,
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func (aog *FieldObj3D) ChangeTile(ti webtilegroup.TileInfo) {
	aog.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)
	aog.Tex.Set("needsUpdate", true)
}

func (aog *FieldObj3D) SetFieldPosition(fx, fy int) {
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+aog.GeoInfo.Len[0]/2,
		-float64(fy)*DstCellSize-aog.GeoInfo.Len[1]/2,
		aog.GeoInfo.Len[2]/2+2,
	)
}

func (aog *FieldObj3D) Dispose() {
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

func FieldObj2TileInfo(
	fotype fieldobjdisplaytype.FieldObjDisplayType,
	fx, fy int) webtilegroup.TileInfo {
	tlList := gClientTile.FieldObjTiles[fotype]
	tilediff := fx*5 + fy*3
	ti := tlList[tilediff%len(tlList)]
	return ti
}
