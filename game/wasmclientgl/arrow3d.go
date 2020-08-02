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
	"sync"
	"syscall/js"

	"github.com/kasworld/goguelike/enum/way9type"
)

var gPoolColorArrow3D = NewPoolColorArrow3D()

type PoolColorArrow3D struct {
	mutex    sync.Mutex
	poolData map[string][]*ColorArrow3D
	newCount int
	getCount int
	putCount int
}

func NewPoolColorArrow3D() *PoolColorArrow3D {
	return &PoolColorArrow3D{
		poolData: make(map[string][]*ColorArrow3D),
	}
}

func (p *PoolColorArrow3D) String() string {
	return fmt.Sprintf("PoolColorArrow3D[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolColorArrow3D) Get(colorstr string) *ColorArrow3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *ColorArrow3D
	if l := len(p.poolData[colorstr]); l > 0 {
		rtn = p.poolData[colorstr][l-1]
		p.poolData[colorstr] = p.poolData[colorstr][:l-1]
		p.getCount++
	} else {
		rtn = NewColorArrow3D(colorstr)
		p.newCount++
	}
	return rtn
}

func (p *PoolColorArrow3D) Put(pb *ColorArrow3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[pb.ColorStr] = append(p.poolData[pb.ColorStr], pb)
	p.putCount++
}

type ColorArrow3D struct {
	Dir      way9type.Way9Type
	ColorStr string
	GeoInfo  GeoInfo
	Mesh     js.Value
}

func NewColorArrow3D(colorstr string) *ColorArrow3D {
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"color": colorstr,
		},
	)
	mat.Set("transparent", true)

	geo := MakeArrowGeo()
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &ColorArrow3D{
		ColorStr: colorstr,
		GeoInfo:  GetGeoInfo(geo),
		Mesh:     mesh,
	}
}

func MakeArrowGeo() js.Value {
	matrix := ThreeJsNew("Matrix4")
	geoLine := ThreeJsNew("CylinderGeometry", 1, 3, DstCellSize-2)
	geoCone := ThreeJsNew("ConeGeometry", DstCellSize/6, DstCellSize/3)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, DstCellSize/4, 0,
	))
	geoLine.Call("merge", geoCone, matrix)
	geoCone.Call("dispose")
	geoLine.Call("center")
	return geoLine
}

func (aog *ColorArrow3D) Visible(b bool) {
	aog.Mesh.Set("visible", b)
}

func (aog *ColorArrow3D) SetFieldPosition(fx, fy int, shX, shY, shZ float64) {
	SetPosition(
		aog.Mesh,
		shX+float64(fx)*DstCellSize+DstCellSize/2,
		-shY-float64(fy)*DstCellSize-DstCellSize/2,
		shZ+DstCellSize/2,
	)
}

func (aog *ColorArrow3D) SetDir(dir way9type.Way9Type) {
	aog.Dir = dir
	aog.RotateZ(-float64(dir-1) * math.Pi / 4)
}

func (aog *ColorArrow3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *ColorArrow3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *ColorArrow3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *ColorArrow3D) Dispose() {
	// mesh do not need dispose
	aog.Mesh.Get("geometry").Call("dispose")
	aog.Mesh.Get("material").Call("dispose")
	aog.Mesh = js.Undefined()
	// no need createElement canvas dom obj
}
