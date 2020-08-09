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
	"math"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/kasworld/findnear"
	"github.com/kasworld/go-abs"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/game/clientinitdata"
)

const (
	DisplayLineLimit = 3*gameconst.ViewPortH - gameconst.ViewPortH/2
	DstCellSize      = 64.0
	HelperSize       = DstCellSize * 32
	UnitTileZ        = DstCellSize / 16

	PoolSizeActiveObj3D = 1000
	PoolSizeCarryObj3D  = 1000
	PoolSizeFieldObj3D  = 1000

	// invalid value
	Tile3DHeightMin = -1000.0
)

var gRnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var gXYLenListView findnear.XYLenList = findnear.NewXYLenList(
	gameconst.ClientViewPortW, gameconst.ClientViewPortH)
var gInitData *clientinitdata.InitData = clientinitdata.New()

var gTextureLoader js.Value = ThreeJsNew("TextureLoader")
var gFontLoader js.Value = ThreeJsNew("FontLoader")
var gFont_helvetiker_regular js.Value
var gFont_droid_sans_mono_regular js.Value

var gColorMaterialCache map[string]js.Value = make(map[string]js.Value)

func GetColorMaterialByCache(co string) js.Value {
	mat, exist := gColorMaterialCache[co]
	if !exist {
		mat = ThreeJsNew("MeshStandardMaterial",
			map[string]interface{}{
				"color": co,
			},
		)
		mat.Set("transparent", true)
		gColorMaterialCache[co] = mat
	}
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

// 0~1 -> 1->1.5->1->0.5->1
func CalcScaleFrameProgress(frameProgress float64, damage int) float64 {

	amplitude := math.Log10(float64(damage)) / 2
	if amplitude < 0.1 {
		amplitude = 0.1
	}
	if amplitude > 2 {
		amplitude = 2
	}

	// 0~1 -> 1->2->1->0->1
	progress := math.Sin(frameProgress * math.Pi * 2)

	rtn := 1 + progress*amplitude
	if rtn < 0 {
		rtn = -rtn
	}
	return rtn
}

// 0~1 -> 0-> -pi -> 0 -> +pi
func CalcRotateFrameProgress(frameProgress float64) float64 {
	return math.Sin(frameProgress*math.Pi*2) * math.Pi / 4
}

func NextPowerOf2(n int) int {
	if IsPowerOfTwo(n) {
		return n
	}
	rtn := 1
	for rtn < n {
		rtn <<= 1
	}
	return rtn
}

func IsPowerOfTwo(i int) bool {
	return (i & (i - 1)) == 0
}
