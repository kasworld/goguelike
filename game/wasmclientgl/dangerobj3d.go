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

	"github.com/kasworld/goguelike/enum/dangertype"
)

var gPoolDangerObj3D = NewPoolDangerObj3D()

type PoolDangerObj3D struct {
	mutex    sync.Mutex
	poolData [dangertype.DangerType_Count][]*DangerObj3D
	newCount int
	getCount int
	putCount int
}

func NewPoolDangerObj3D() *PoolDangerObj3D {
	return &PoolDangerObj3D{
		// poolData: make([dangertype.DangerType_Count][]*DangerObj3D),
	}
}

func (p *PoolDangerObj3D) String() string {
	return fmt.Sprintf("PoolDangerObj3D[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolDangerObj3D) Get(dangerType dangertype.DangerType) *DangerObj3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *DangerObj3D
	if l := len(p.poolData[dangerType]); l > 0 {
		rtn = p.poolData[dangerType][l-1]
		p.poolData[dangerType] = p.poolData[dangerType][:l-1]
		p.getCount++
	} else {
		rtn = NewDangerObj3D(dangerType)
		p.newCount++
	}
	return rtn
}

func (p *PoolDangerObj3D) Put(pb *DangerObj3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[pb.DangerType] = append(p.poolData[pb.DangerType], pb)
	p.putCount++
}

type DangerObj3D struct {
	DangerType dangertype.DangerType
	ColorStr   string
	GeoInfo    GeoInfo
	Mesh       js.Value
}

func NewDangerObj3D(dangerType dangertype.DangerType) *DangerObj3D {
	colorstr := dangerType.Color24().ToHTMLColorString()
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"color": colorstr,
		},
	)
	mat.Set("transparent", true)

	geo := MakeDanger3DGeo()
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &DangerObj3D{
		DangerType: dangerType,
		GeoInfo:    GetGeoInfo(geo),
		Mesh:       mesh,
	}
}

func MakeDanger3DGeo() js.Value {
	matrix := ThreeJsNew("Matrix4")
	geoLine := ThreeJsNew("CylinderGeometry", 1, 3, DstCellSize-2)
	geoCone := ThreeJsNew("ConeGeometry", DstCellSize/6, DstCellSize/3)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, DstCellSize/4, 0,
	))
	geoLine.Call("merge", geoCone, matrix)
	geoCone.Call("dispose")
	geoLine.Call("rotateX", math.Pi/2)
	geoLine.Call("center")
	return geoLine
}

func (aog *DangerObj3D) Visible(b bool) {
	aog.Mesh.Set("visible", b)
}

func (aog *DangerObj3D) SetFieldPosition(fx, fy int, shX, shY, shZ float64) {
	SetPosition(
		aog.Mesh,
		shX+float64(fx)*DstCellSize+DstCellSize/2,
		-shY-float64(fy)*DstCellSize-DstCellSize/2,
		shZ+DstCellSize/2,
	)
}

func (aog *DangerObj3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *DangerObj3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *DangerObj3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *DangerObj3D) Dispose() {
	// mesh do not need dispose
	aog.Mesh.Get("geometry").Call("dispose")
	aog.Mesh.Get("material").Call("dispose")
	aog.Mesh = js.Undefined()
	// no need createElement canvas dom obj
}
