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
	"fmt"
	"syscall/js"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/wrapper"
)

var gTextureTileList [tile.Tile_Count]*TextureTile

func LoadTextureTileList() [tile.Tile_Count]*TextureTile {
	var rtn [tile.Tile_Count]*TextureTile
	for i, v := range tile.TileScrollAttrib {
		if v.Texture {
			idstr := fmt.Sprintf("%vPng", tile.Tile(i))
			rtn[i] = NewTextureTile(idstr)
		}
	}
	return rtn
}

type TextureTile struct {
	// Img js.Value
	W   int
	H   int
	Cnv js.Value
	// Ctx js.Value

	// wrap info
	CellSize int
	Xcount   int
	Ycount   int
	WrapX    func(int) int
	WrapY    func(int) int
}

func NewTextureTile(srcImageID string) *TextureTile {
	img := js.Global().Get("document").Call("getElementById", srcImageID)
	if !img.Truthy() {
		jslog.Errorf("fail to get %v", srcImageID)
		return nil
	}
	srcw := img.Get("naturalWidth").Int()
	srch := img.Get("naturalHeight").Int()

	cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	ctx := cnv.Call("getContext", "2d")
	if !ctx.Truthy() {
		jslog.Errorf("fail to get context", srcImageID)
		return nil
	}
	ctx.Set("imageSmoothingEnabled", false)
	cnv.Set("width", srcw)
	cnv.Set("height", srch)
	ctx.Call("clearRect", 0, 0, srcw, srch)
	ctx.Call("drawImage", img, 0, 0)

	srcCellSize := 64
	return &TextureTile{
		// Img: img,
		W:   srcw,
		H:   srch,
		Cnv: cnv,
		// Ctx: ctx,

		CellSize: srcCellSize,
		Xcount:   srcw / srcCellSize,
		Ycount:   srch / srcCellSize,
		WrapX:    wrapper.New(srcw - srcCellSize).WrapSafe,
		WrapY:    wrapper.New(srch - srcCellSize).WrapSafe,
	}
}

func (tt *TextureTile) CalcSrc(fx, fy int, shiftx, shifty float64) (int, int, int) {
	tx := fx % tt.Xcount
	ty := fy % tt.Ycount
	srcx := tt.WrapX(tx*tt.CellSize + int(shiftx))
	srcy := tt.WrapY(ty*tt.CellSize + int(shifty))
	return srcx, srcy, tt.CellSize
}
