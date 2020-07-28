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
	"sync"
	"syscall/js"

	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/lib/webtilegroup"
)

var gPoolActiveObj3D = NewPoolActiveObj3D(PoolSizeActiveObj3D)

type PoolActiveObj3D struct {
	mutex    sync.Mutex
	poolData []*ActiveObj3D
	newCount int
	getCount int
	putCount int
}

func NewPoolActiveObj3D(initCap int) *PoolActiveObj3D {
	return &PoolActiveObj3D{
		poolData: make([]*ActiveObj3D, 0, initCap),
	}
}

func (p *PoolActiveObj3D) String() string {
	return fmt.Sprintf("PoolActiveObj3D[%v/%v new:%v get:%v put:%v]",
		len(p.poolData), cap(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolActiveObj3D) Get() *ActiveObj3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *ActiveObj3D
	if l := len(p.poolData); l > 0 {
		rtn = p.poolData[l-1]
		p.poolData = p.poolData[:l-1]
		p.getCount++
	} else {
		rtn = NewActiveObj3D()
		p.newCount++
	}
	return rtn
}

func (p *PoolActiveObj3D) Put(pb *ActiveObj3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData = append(p.poolData, pb)
	p.putCount++
}

type ActiveObj3D struct {
	Name    *Label3D
	Chat    *Label3D
	Cnv     js.Value
	Ctx     js.Value
	Tex     js.Value
	GeoInfo GeoInfo
	Mesh    js.Value
}

func NewActiveObj3D() *ActiveObj3D {
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
	geo := ThreeJsNew("BoxGeometry", DstCellSize, DstCellSize, DstCellSize)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &ActiveObj3D{
		Cnv:     cnv,
		Ctx:     ctx,
		Tex:     tex,
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func (aog *ActiveObj3D) ChangeTile(ti webtilegroup.TileInfo) {
	aog.Ctx.Call("clearRect", 0, 0, DstCellSize, DstCellSize)
	aog.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)
	aog.Tex.Set("needsUpdate", true)
}

func (aog *ActiveObj3D) SetFieldPosition(fx, fy int, tl tile_flag.TileFlag) {
	height := GetTile3DHeightByCache(tl)
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+DstCellSize/2,
		-float64(fy)*DstCellSize-DstCellSize/2,
		aog.GeoInfo.Len[2]/2+1+height,
	)
}

func (aog *ActiveObj3D) ResetMatrix() {
	aog.ScaleX(1.0)
	aog.ScaleY(1.0)
	aog.ScaleZ(1.0)
	aog.RotateX(0.0)
	aog.RotateY(0.0)
	aog.RotateZ(0.0)
}

func (aog *ActiveObj3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *ActiveObj3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *ActiveObj3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *ActiveObj3D) ScaleX(x float64) {
	aog.Mesh.Get("scale").Set("x", x)
}
func (aog *ActiveObj3D) ScaleY(y float64) {
	aog.Mesh.Get("scale").Set("y", y)
}
func (aog *ActiveObj3D) ScaleZ(z float64) {
	aog.Mesh.Get("scale").Set("z", z)
}

func (aog *ActiveObj3D) Dispose() {
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
