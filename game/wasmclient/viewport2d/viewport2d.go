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

package viewport2d

import (
	"fmt"
	"syscall/js"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/lib/canvastext"
	"github.com/kasworld/goguelike/lib/clienttile"
	"github.com/kasworld/goguelike/lib/imagecanvas"
	"github.com/kasworld/gowasmlib/jslog"
)

type CellVPCarryObjS struct {
	X, Y       int // cell index
	dstX, dstY int // cell pos
}

type Viewport2d struct {
	clientTile       *clienttile.ClientTile
	TileImgCnvList   [tile.Tile_Count]*imagecanvas.ImageCanvas
	DarkerTileImgCnv *imagecanvas.ImageCanvas

	Canvas    js.Value
	context2d js.Value

	zoomState int

	ViewWidth         int
	ViewHeight        int
	RefSize           int // cellsize w/o zoom
	CellSize          int // cellsize with zoom
	CellXCount        int
	CellYCount        int
	CellIndexPos      []CellVPCarryObjS
	ViewportPos2Index map[[2]int]int

	Grid         bool
	DrawEquip    bool
	TileScroll   bool
	DamageEffect bool

	// for display cursor
	MouseX int
	MouseY int

	NotiMessage      canvastext.CanvasTextList
	ActiveObjMessage canvastext.CanvasTextList
}

func New(cnvid string, ct *clienttile.ClientTile) *Viewport2d {
	vp := &Viewport2d{
		clientTile:   ct,
		Grid:         true, // from option
		DrawEquip:    true, // from option
		TileScroll:   true, // from option
		DamageEffect: true, // from option
	}

	vp.Canvas = js.Global().Get("document").Call("getElementById", cnvid)
	if !vp.Canvas.Truthy() {
		jslog.Errorf("fail to get canvas %v", cnvid)
	}
	vp.context2d = vp.Canvas.Call("getContext", "2d")
	vp.context2d.Set("imageSmoothingEnabled", false)
	if !vp.context2d.Truthy() {
		jslog.Error("fail to get context2d context2d")
	}
	for i, v := range tile.TileScrollAttrib {
		if v.Texture {
			idstr := fmt.Sprintf("%vPng", tile.Tile(i))
			vp.TileImgCnvList[i] = imagecanvas.NewByID(idstr)
		}
	}
	vp.DarkerTileImgCnv = imagecanvas.NewByID("DarkerPng")
	return vp
}

var textureTileWrapInfo = [tile.Tile_Count]struct {
	Xcount int
	Ycount int
	WrapX  func(int) int
	WrapY  func(int) int
}{}
