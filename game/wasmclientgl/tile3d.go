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
	"github.com/kasworld/goguelike/lib/webtilegroup"
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
	Shift   [3]float64
	GeoInfo GeoInfo
}

func newTile3D() *Tile3D {
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
	return &Tile3D{
		Cnv: cnv,
		Ctx: ctx,
		Tex: tex,
		Mat: mat,
	}
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

func NewTile3D_PlaneGeo(tl tile.Tile, shiftZ float64) *Tile3D {
	t3d := newTile3D().initSrc(tl)
	t3d.Geo = ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	t3d.Shift = [3]float64{0, 0, shiftZ}
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}

func NewTile3D_Grass(shiftZ float64) *Tile3D {
	tl := tile.Grass
	t3d := newTile3D().initSrc(tl)
	t3d.Geo = ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize/8)
	t3d.Shift = [3]float64{0, 0, shiftZ}
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}

func NewTile3D_Tree(shiftZ float64) *Tile3D {
	tl := tile.Grass
	t3d := newTile3D().initSrc(tl)
	t3d.Geo = MakeTreeGeo()
	t3d.Shift = [3]float64{0, 0, shiftZ}
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}
func MakeTreeGeo() js.Value {
	matrix := ThreeJsNew("Matrix4")
	geo := ThreeJsNew("CylinderGeometry", 1, 3, DstCellSize-1)
	geo.Call("rotateX", math.Pi/2)
	geo1 := ThreeJsNew("ConeGeometry", DstCellSize/3, DstCellSize/2-2)
	geo1.Call("rotateX", math.Pi/2)
	geo2 := ThreeJsNew("ConeGeometry", DstCellSize/2-2, DstCellSize/2-2)
	geo2.Call("rotateX", math.Pi/2)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, 0, DstCellSize/2-2,
	))
	geo.Call("merge", geo1, matrix)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, 0, 2,
	))
	geo.Call("merge", geo2, matrix)
	geo1.Call("dispose")
	geo2.Call("dispose")
	return geo
}

func NewTile3D_Wall(shiftZ float64) *Tile3D {
	tl := tile.Stone
	t3d := newTile3D().initSrc(tl)
	t3d.Geo = ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize)
	t3d.Shift = [3]float64{0, 0, shiftZ}
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}

func NewTile3D_Window(shiftZ float64) *Tile3D {
	tl := tile.Window
	t3d := newTile3D().initSrc(tl)
	t3d.Geo = ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize)
	t3d.Shift = [3]float64{0, 0, shiftZ}
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.DrawTexture(0, 0)
	return t3d
}

func NewTile3D_BoxTile(ti webtilegroup.TileInfo, shiftZ float64) *Tile3D {
	t3d := newTile3D()
	t3d.Geo = ThreeJsNew("BoxGeometry", DstCellSize-1, DstCellSize-1, DstCellSize-1)
	t3d.Shift = [3]float64{0, 0, shiftZ}
	t3d.GeoInfo = GetGeoInfo(t3d.Geo)
	t3d.ChangeTile(ti)
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

func (t3d *Tile3D) ChangeTile(ti webtilegroup.TileInfo) {
	t3d.Ctx.Call("clearRect", 0, 0, DstCellSize, DstCellSize)
	t3d.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)
	t3d.Tex.Set("needsUpdate", true)
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
