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
	"math/rand"
	"syscall/js"
	"time"

	"github.com/kasworld/findnear"
	"github.com/kasworld/go-abs"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/clientinitdata"
	"github.com/kasworld/goguelike/lib/clienttile"
	"github.com/kasworld/goguelike/lib/webtilegroup"
)

const (
	DisplayLineLimit = 3*gameconst.ViewPortH - gameconst.ViewPortH/2
	DstCellSize      = 32
	HelperSize       = DstCellSize * 32

	PoolSizeActiveObj3D = 1000
	PoolSizeCarryObj3D  = 1000
	PoolSizeFieldObj3D  = 1000
	PoolSizeArrow3D     = 100

	// invalid value
	Tile3DHeightMin = -1000.0
)

var gRnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var gXYLenListView findnear.XYLenList = findnear.NewXYLenList(
	gameconst.ClientViewPortW, gameconst.ClientViewPortH)
var gInitData *clientinitdata.InitData = clientinitdata.New()
var gClientTile *clienttile.ClientTile = clienttile.New()

var gTextureLoader js.Value = ThreeJsNew("TextureLoader")
var gFontLoader js.Value = ThreeJsNew("FontLoader")
var gFont_helvetiker_regular js.Value

func NewTileMaterial(ti webtilegroup.TileInfo) js.Value {
	Cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", DstCellSize)
	Cnv.Set("height", DstCellSize)
	Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)

	Tex := ThreeJsNew("CanvasTexture", Cnv)
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	mat.Set("transparent", true)
	// mat.Set("side", ThreeJs().Get("DoubleSide"))
	return mat
}

func NewTextureTileMaterial(ti tile.Tile) js.Value {
	Cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	Ctx := Cnv.Call("getContext", "2d")
	Ctx.Set("imageSmoothingEnabled", false)
	Cnv.Set("width", DstCellSize)
	Cnv.Set("height", DstCellSize)

	img := GetElementById(fmt.Sprintf("%vPng", ti))
	Ctx.Call("drawImage", img,
		0, 0, DstCellSize, DstCellSize,
		0, 0, DstCellSize, DstCellSize)

	Tex := ThreeJsNew("CanvasTexture", Cnv)
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"map": Tex,
		},
	)
	mat.Set("transparent", true)
	// mat.Set("side", ThreeJs().Get("DoubleSide"))
	return mat
}

func CalcCurrentFrame(difftick int64, fps float64) int {
	diffsec := float64(difftick) / float64(time.Second)
	frame := fps * diffsec
	return int(frame)
}

func CalcShiftDxDy(frameProgress float64) (int, int) {
	rate := 1 - frameProgress
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	dx := int(float64(DstCellSize) * rate)
	dy := int(float64(DstCellSize) * rate)
	return dx, dy
}

func SetPosition(jso js.Value, pos ...interface{}) {
	po := jso.Get("position")
	if len(pos) >= 1 {
		po.Set("x", pos[0])
	}
	if len(pos) >= 2 {
		po.Set("y", pos[1])
	}
	if len(pos) >= 3 {
		po.Set("z", pos[2])
	}
}

func ThreeJsNew(name string, args ...interface{}) js.Value {
	return js.Global().Get("THREE").Get(name).New(args...)
}

func ThreeJs() js.Value {
	return js.Global().Get("THREE")
}

func GetElementById(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

// make fx,fy around vpx, vpy
func CalcAroundPos(w, h, vpx, vpy, fx, fy int) (int, int) {
	if abs.Absi(fx-vpx) > w/2 {
		if fx > vpx {
			fx -= w
		} else {
			fx += w
		}
	}
	if abs.Absi(fy-vpy) > h/2 {
		if fy > vpy {
			fy -= h
		} else {
			fy += h
		}
	}
	return fx, fy
}
