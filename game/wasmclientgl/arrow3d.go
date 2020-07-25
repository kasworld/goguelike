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

	"github.com/kasworld/goguelike/enum/tile_flag"
)

var gPoolArrow3D = NewPoolArrow3D(PoolSizeArrow3D)

type PoolArrow3D struct {
	mutex    sync.Mutex
	poolData []*Arrow3D
	newCount int
	getCount int
	putCount int
}

func NewPoolArrow3D(initCap int) *PoolArrow3D {
	return &PoolArrow3D{
		poolData: make([]*Arrow3D, 0, initCap),
	}
}

func (p *PoolArrow3D) String() string {
	return fmt.Sprintf("PoolArrow3D[%v/%v new:%v get:%v put:%v]",
		len(p.poolData), cap(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolArrow3D) Get() *Arrow3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *Arrow3D
	if l := len(p.poolData); l > 0 {
		rtn = p.poolData[l-1]
		p.poolData = p.poolData[:l-1]
		p.getCount++
	} else {
		rtn = NewArrow3D()
		p.newCount++
	}
	return rtn
}

func (p *PoolArrow3D) Put(pb *Arrow3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData = append(p.poolData, pb)
	p.putCount++
}

type Arrow3D struct {
	GeoInfo GeoInfo
	Mesh    js.Value
}

func NewArrow3D() *Arrow3D {
	colorstr := "#ffffff"
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"color": colorstr,
		},
	)
	mat.Set("transparent", true)

	geo := MakeArrowGeo()
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &Arrow3D{
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func MakeArrowGeo() js.Value {
	matrix := ThreeJsNew("Matrix4")
	geoLine := ThreeJsNew("CylinderGeometry", 1, 3, DstCellSize)
	geoCone := ThreeJsNew("ConeGeometry", DstCellSize/4, DstCellSize/2)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, DstCellSize/2, 0,
	))
	geoLine.Call("merge", geoCone, matrix)
	geoCone.Call("dispose")
	return geoLine
}

func (aog *Arrow3D) Visible(b bool) {
	aog.Mesh.Set("visible", b)
}

func (aog *Arrow3D) SetFieldPosition(fx, fy int, tl tile_flag.TileFlag) {
	height := GetTile3DHeightByCache(tl)
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+aog.GeoInfo.Len[0]/2,
		-float64(fy)*DstCellSize-aog.GeoInfo.Len[1]/2,
		aog.GeoInfo.Len[2]/2+1+height,
	)
}

func (aog *Arrow3D) SetDir(dir way9type.Way9Type) {
	aog.RotateZ(-float64(dir-1) * math.Pi / 4)
}

func (aog *Arrow3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *Arrow3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *Arrow3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *Arrow3D) Dispose() {
	// mesh do not need dispose
	aog.Mesh.Get("geometry").Call("dispose")
	aog.Mesh.Get("material").Call("dispose")
	aog.Mesh = js.Undefined()
	// no need createElement canvas dom obj
}
