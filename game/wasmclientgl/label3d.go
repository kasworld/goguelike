// Copyright 2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
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
	"sync"
	"syscall/js"
)

var gPoolLabel3D = NewPoolLabel3D()

type PoolLabel3D struct {
	mutex    sync.Mutex
	poolData map[string][]*Label3D
	newCount int
	getCount int
	putCount int
}

func NewPoolLabel3D() *PoolLabel3D {
	return &PoolLabel3D{
		poolData: make(map[string][]*Label3D),
	}
}

func (p *PoolLabel3D) String() string {
	return fmt.Sprintf("PoolLabel3D[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolLabel3D) Get(str string) *Label3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *Label3D
	if l := len(p.poolData[str]); l > 0 {
		rtn = p.poolData[str][l-1]
		p.poolData[str] = p.poolData[str][:l-1]
		p.getCount++
	} else {
		rtn = NewLabel3D(str)
		p.newCount++
	}
	return rtn
}

func (p *PoolLabel3D) Put(pb *Label3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[pb.Str] = append(p.poolData[pb.Str], pb)
	p.putCount++
}

type Label3D struct {
	Str     string
	Cnv     js.Value
	Ctx     js.Value
	Tex     js.Value
	GeoInfo GeoInfo
	Mesh    js.Value
}

func NewLabel3D(str string) *Label3D {
	cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	ctx := cnv.Call("getContext", "2d")
	ctx.Set("imageSmoothingEnabled", false)

	cnv.Set("width", len(str)*DstCellSize)
	cnv.Set("height", DstCellSize)
	height := DstCellSize / 3
	font := fmt.Sprintf("%vpx sans-serif", height)

	ctx.Set("font", font)
	textWidth := ctx.Call("measureText", str).Get("width").Int()
	canvasWidth := NextPowerOf2(textWidth)
	cnv.Set("width", canvasWidth)
	canvasHeight := NextPowerOf2(int(height))
	cnv.Set("height", canvasHeight)

	ctx.Set("fillStyle", "#00000080")
	ctx.Call("fillRect", 0, 0, canvasWidth, canvasHeight)
	ctx.Set("font", font)
	ctx.Set("fillStyle", "white")
	ctx.Call("fillText", str, (canvasWidth-textWidth)/2, canvasHeight-canvasHeight/3)

	tex := ThreeJsNew("CanvasTexture", cnv)
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"map": tex,
		},
	)
	mat.Set("transparent", true)

	geo := ThreeJsNew("PlaneGeometry", canvasWidth, canvasHeight)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &Label3D{
		Str:     str,
		Cnv:     cnv,
		Ctx:     ctx,
		Tex:     tex,
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func (aog *Label3D) SetFieldPosition(fx, fy int, shX, shY, shZ float64) {
	SetPosition(
		aog.Mesh,
		shX+float64(fx)*DstCellSize+DstCellSize/2, //+aog.GeoInfo.Len[0]/2,
		-shY-float64(fy)*DstCellSize-aog.GeoInfo.Len[1]/2,
		shZ+aog.GeoInfo.Len[2]/2,
	)
}

func (aog *Label3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *Label3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *Label3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *Label3D) Dispose() {
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
