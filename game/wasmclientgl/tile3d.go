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

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/lib/webtilegroup"
)

type Tile3D struct {
	Cnv     js.Value
	Ctx     js.Value
	Tex     js.Value
	Mat     js.Value
	Geo     js.Value
	Shift   [3]float64
	GeoInfo GeoInfo
}

func NewTile3DCanvas(tl tile.Tile) Tile3D {
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
	return Tile3D{
		Cnv: cnv,
		Ctx: ctx,
		Tex: tex,
		Mat: mat,
	}
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

func (t3d Tile3D) ChangeTile(ti webtilegroup.TileInfo) {
	t3d.Ctx.Call("clearRect", 0, 0, DstCellSize, DstCellSize)
	t3d.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)
	t3d.Tex.Set("needsUpdate", true)
}

func (t3d Tile3D) DrawTexture(tl tile.Tile, srcx, srcy int) {
	t3d.Ctx.Call("clearRect", 0, 0, DstCellSize, DstCellSize)
	src := gTextureTileList[tl].Cnv
	t3d.Ctx.Call("drawImage", src,
		srcx, srcy, DstCellSize, DstCellSize,
		0, 0, DstCellSize, DstCellSize)
	t3d.Tex.Set("needsUpdate", true)
}
