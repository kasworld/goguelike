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
	"math"
	"syscall/js"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/wrapper"
)

type Tile3D struct {
	// for texture tile
	SrcW     int
	SrcH     int
	SrcCnv   js.Value
	SrcWrapX func(int) int
	SrcWrapY func(int) int

	Cnv     js.Value
	Ctx     js.Value
	Tex     js.Value
	Mat     js.Value
	Geo     js.Value
	GeoInfo GeoInfo
	Tile    tile.Tile
}

func (aog *Tile3D) MakePosVector3(fx, fy int) js.Value {
	return ThreeJsNew("Vector3",
		float64(fx)*DstCellSize+DstCellSize/2,
		-float64(fy)*DstCellSize-DstCellSize/2,
		aog.GeoInfo.Len[2]/2+gTileZInfo[aog.Tile].Shift,
	)
}

func newTile3D(tl tile.Tile, opacity float64) *Tile3D {
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
	if opacity < 1.0 {
		mat.Set("transparent", true)
		mat.Set("side", ThreeJs().Get("DoubleSide"))
		mat.Set("opacity", opacity)
	}

	aog := &Tile3D{
		Tile: tl,
		Cnv:  cnv,
		Ctx:  ctx,
		Tex:  tex,
		Mat:  mat,
	}
	aog.initSrc(tl)
	return aog
}

// for texture tile
func (aog *Tile3D) initSrc(tl tile.Tile) *Tile3D {
	srcImageID := fmt.Sprintf("%vPng", tile.Tile(tl))
	img := GetElementById(srcImageID)
	if !img.Truthy() {
		jslog.Errorf("fail to get %v", srcImageID)
		return aog
	}
	srcw := img.Get("naturalWidth").Int()
	srch := img.Get("naturalHeight").Int()

	cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	ctx := cnv.Call("getContext", "2d")
	if !ctx.Truthy() {
		jslog.Errorf("fail to get context", srcImageID)
		return aog
	}
	ctx.Set("imageSmoothingEnabled", false)
	cnv.Set("width", srcw)
	cnv.Set("height", srch)
	ctx.Call("clearRect", 0, 0, srcw, srch)
	ctx.Call("drawImage", img, 0, 0)
	aog.SrcW = srcw
	aog.SrcH = srch
	aog.SrcCnv = cnv
	aog.SrcWrapX = wrapper.New(srcw - DstCellSize).WrapSafe
	aog.SrcWrapY = wrapper.New(srch - DstCellSize).WrapSafe
	return aog
}

func (aog *Tile3D) MakeSrcDark() *Tile3D {
	ctx := aog.SrcCnv.Call("getContext", "2d")
	ctx.Set("fillStyle", "#000000b0")
	ctx.Call("fillRect", 0, 0, aog.SrcW, aog.SrcH)
	aog.DrawTexture(0, 0)
	return aog
}

// func NewTile3D_PlaneGeo(tl tile.Tile) *Tile3D {
// 	t3d := newTile3D(tl, false)
// 	// t3d.Geo = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
// 	t3d.Geo = ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize/16)
// 	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
// 	t3d.DrawTexture(0, 0)
// 	return t3d
// }

func NewTile3D_BoxTexture(tl tile.Tile, opacity float64) *Tile3D {
	t3d := newTile3D(tl, opacity)
	t3d.Geo = ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, gTileZInfo[tl].Size)
	t3d.Geo.Call("center")
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}

func NewTile3D_OctCylinderTexture(tl tile.Tile, opacity float64) *Tile3D {
	t3d := newTile3D(tl, opacity)
	t3d.Geo = ThreeJsNew("CylinderGeometry",
		DstCellSize/2, DstCellSize/2, gTileZInfo[tl].Size, 8)
	t3d.Geo.Call("center")

	t3d.Geo.Call("rotateX", math.Pi/2)
	t3d.Geo.Call("rotateZ", math.Pi*2/8/2)
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}

func NewTile3D_Tree(opacity float64) *Tile3D {
	tl := tile.Tree
	t3d := newTile3D(tl, opacity)
	t3d.Geo = makeTreeGeo()
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}
func makeTreeGeo() js.Value {
	matrix := ThreeJsNew("Matrix4")

	geo := ThreeJsNew("CylinderGeometry", 1, 3, DstCellSize-1)
	geo.Call("rotateX", math.Pi/2)

	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, 0, DstCellSize/3-2,
	))
	geo1 := ThreeJsNew("ConeGeometry", DstCellSize/4, DstCellSize/3)
	geo1.Call("rotateX", math.Pi/2)
	geo.Call("merge", geo1, matrix)
	geo1.Call("dispose")

	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, 0, 2,
	))
	geo2 := ThreeJsNew("ConeGeometry", DstCellSize/3, DstCellSize/3)
	geo2.Call("rotateX", math.Pi/2)
	geo.Call("merge", geo2, matrix)
	geo2.Call("dispose")

	geo.Call("center")
	return geo
}

func NewTile3D_Door(opacity float64) *Tile3D {
	t3d := newTile3D(tile.Door, opacity)

	t3d.Geo = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	t3d.Geo.Call("rotateX", math.Pi/2)

	geo2 := ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	geo2.Call("rotateY", math.Pi/2)
	matrix := ThreeJsNew("Matrix4")
	t3d.Geo.Call("merge", geo2, matrix)
	geo2.Call("dispose")

	// t3d.Geo.Call("rotateZ", math.Pi/4)
	t3d.Geo.Call("center")
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}

func (aog *Tile3D) Dispose() {
	// mesh do not need dispose
	aog.Geo.Call("dispose")
	aog.Mat.Call("dispose")
	aog.Tex.Call("dispose")

	aog.Cnv = js.Undefined()
	aog.Ctx = js.Undefined()
	aog.Tex = js.Undefined()
	// no need createElement canvas dom obj
}

func (t3d *Tile3D) DrawTexture(srcx, srcy int) {
	srcx = t3d.SrcWrapX(srcx)
	srcy = t3d.SrcWrapY(srcy)
	t3d.Ctx.Call("clearRect", 0, 0, DstCellSize, DstCellSize)
	t3d.Ctx.Call("drawImage", t3d.SrcCnv,
		srcx, srcy, DstCellSize, DstCellSize,
		0, 0, DstCellSize, DstCellSize)
	t3d.Tex.Set("needsUpdate", true)
}
