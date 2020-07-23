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
)

var gPoolColorBar3D = NewPoolColorBar3D()

type PoolColorBar3D struct {
	mutex    sync.Mutex
	poolData map[string][]*ColorBar3D
	newCount int
	getCount int
	putCount int
}

func NewPoolColorBar3D() *PoolColorBar3D {
	return &PoolColorBar3D{
		poolData: make(map[string][]*ColorBar3D),
	}
}

func (p *PoolColorBar3D) String() string {
	return fmt.Sprintf("PoolColorBar3D[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolColorBar3D) Get(str string) *ColorBar3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *ColorBar3D
	if l := len(p.poolData[str]); l > 0 {
		rtn = p.poolData[str][l-1]
		p.poolData[str] = p.poolData[str][:l-1]
		p.getCount++
	} else {
		rtn = NewColorBar3D(str)
		p.newCount++
	}
	return rtn
}

func (p *PoolColorBar3D) Put(pb *ColorBar3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[pb.ColorStr] = append(p.poolData[pb.ColorStr], pb)
	p.putCount++
}

type ColorBar3D struct {
	ColorStr string
	Tex      js.Value
	GeoInfo  GeoInfo
	Mesh     js.Value
}

func NewColorBar3D(colorstr string) *ColorBar3D {
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"color": colorstr,
		},
	)
	mat.Set("transparent", true)

	geo := ThreeJsNew("CylinderGeometry", 1, 1, DstCellSize)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &ColorBar3D{
		ColorStr: colorstr,
		GeoInfo:  GetGeoInfo(geo),
		Mesh:     mesh,
	}
}

func (aog *ColorBar3D) SetFieldPositionUp(fx, fy int, shZ float64) {
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+DstCellSize/2, //+aog.GeoInfo.Len[0]/2,
		-float64(fy)*DstCellSize+aog.GeoInfo.Len[1]/2,
		aog.GeoInfo.Len[2]/2+1+shZ,
	)
}

func (aog *ColorBar3D) SetFieldPositionDown(fx, fy int, shZ float64) {
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+DstCellSize/2, //+aog.GeoInfo.Len[0]/2,
		-float64(fy)*DstCellSize-aog.GeoInfo.Len[1]/2-DstCellSize,
		aog.GeoInfo.Len[2]/2+1+shZ,
	)
}

func (aog *ColorBar3D) ResetMatrix() {
	aog.ScaleX(1.0)
	aog.ScaleY(1.0)
	aog.ScaleZ(1.0)
	aog.RotateX(0.0)
	aog.RotateY(0.0)
	aog.RotateZ(0.0)
}

func (aog *ColorBar3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *ColorBar3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *ColorBar3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *ColorBar3D) ScaleX(x float64) {
	aog.Mesh.Get("scale").Set("x", x)
}
func (aog *ColorBar3D) ScaleY(y float64) {
	aog.Mesh.Get("scale").Set("y", y)
}
func (aog *ColorBar3D) ScaleZ(z float64) {
	aog.Mesh.Get("scale").Set("z", z)
}

func (aog *ColorBar3D) Dispose() {
	// mesh do not need dispose
	aog.Mesh.Get("geometry").Call("dispose")
	aog.Mesh.Get("material").Call("dispose")

	aog.Mesh = js.Undefined()
	// no need createElement canvas dom obj
}
