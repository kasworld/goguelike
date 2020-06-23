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
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

// for tile texture wrap
type textureTileWrapInfo struct {
	CellSize int
	Xcount   int
	Ycount   int
	WrapX    func(int) int
	WrapY    func(int) int
}

func (ttwi textureTileWrapInfo) CalcSrc(fx, fy int, shiftx, shifty float64) (int, int, int) {
	tx := fx % ttwi.Xcount
	ty := fy % ttwi.Ycount
	srcx := ttwi.WrapX(tx*ttwi.CellSize + int(shiftx))
	srcy := ttwi.WrapY(ty*ttwi.CellSize + int(shifty))
	return srcx, srcy, ttwi.CellSize
}

func (vp *Viewport) ReplaceFloorTiles(cf *ClientFloorGL) {
	for fx, xv := range cf.Tiles {
		for fy, yv := range xv {
			vp.drawTileAt(cf, fx, fy, yv)
		}
	}
}

func (vp *Viewport) UpdateClientField(
	cf *ClientFloorGL,
	taNoti *c2t_obj.NotiVPTiles_data,
	ViewportXYLenList findnear.XYLenList,
) {
	for i, v := range ViewportXYLenList {
		fx := cf.XWrapSafe(v.X + taNoti.VPX)
		fy := cf.YWrapSafe(v.Y + taNoti.VPY)
		tl := taNoti.VPTiles[i]
		vp.drawTileAt(cf, fx, fy, tl)
	}
	cf.Tex.Set("needsUpdate", true)

}

func (vp *Viewport) drawTileAt(
	cf *ClientFloorGL, fx, fy int, tl tile_flag.TileFlag) {
	dstX := fx * DstCellSize
	dstY := fy * DstCellSize
	diffbase := fx*5 + fy*3

	for i := 0; i < tile.Tile_Count; i++ {
		shX := 0.0
		shY := 0.0
		tlt := tile.Tile(i)
		if tl.TestByTile(tlt) {
			if vp.textureTileList[i] != nil {
				// texture tile
				tlic := vp.textureTileList[i]
				srcx, srcy, srcCellSize := vp.textureTileWrapInfoList[i].CalcSrc(fx, fy, shX, shY)
				cf.Ctx.Call("drawImage", tlic.Cnv,
					srcx, srcy, srcCellSize, srcCellSize,
					dstX, dstY, DstCellSize, DstCellSize)

			} else if tlt == tile.Wall {
				// wall tile process
				tlList := gClientTile.FloorTiles[i]
				tilediff := cf.calcWallTileDiff(fx, fy)
				ti := tlList[tilediff%len(tlList)]
				cf.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, DstCellSize, DstCellSize)
			} else if tlt == tile.Window {
				// window tile process
				tlList := gClientTile.FloorTiles[i]
				tlindex := 0
				if cf.checkWallAt(fx, fy-1) && cf.checkWallAt(fx, fy+1) { // n-s window
					tlindex = 1
				}
				ti := tlList[tlindex]
				cf.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, DstCellSize, DstCellSize)
			} else {
				// bitmap tile
				tlList := gClientTile.FloorTiles[i]
				tilediff := diffbase + int(shX)
				ti := tlList[tilediff%len(tlList)]
				cf.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, DstCellSize, DstCellSize)
			}
		}
	}
}
