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

	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"

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

func (p *PoolDangerObj3D) Get(Dao *c2t_obj.DangerObjClient) *DangerObj3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *DangerObj3D
	if l := len(p.poolData[Dao.DangerType]); l > 0 {
		rtn = p.poolData[Dao.DangerType][l-1]
		p.poolData[Dao.DangerType] = p.poolData[Dao.DangerType][:l-1]
		p.getCount++
	} else {
		rtn = NewDangerObj3D(Dao)
		p.newCount++
	}
	return rtn
}

func (p *PoolDangerObj3D) Put(pb *DangerObj3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[pb.Dao.DangerType] = append(p.poolData[pb.Dao.DangerType], pb)
	p.putCount++
}

type DangerObj3D struct {
	Dao      *c2t_obj.DangerObjClient
	ColorStr string
	GeoInfo  GeoInfo
	Mesh     js.Value
}

func NewDangerObj3D(Dao *c2t_obj.DangerObjClient) *DangerObj3D {
	colorstr := Dao.DangerType.Color24().ToHTMLColorString()
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"color": colorstr,
		},
	)
	mat.Set("transparent", true)

	geo := MakeDanger3DGeo()
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &DangerObj3D{
		Dao:     Dao,
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func MakeDanger3DGeo() js.Value {
	geo := ThreeJsNew("TorusGeometry", DstCellSize/1.4, DstCellSize/16, 8, 4)
	geo.Call("center")
	return geo
}

func (dao3d *DangerObj3D) Visible(b bool) {
	dao3d.Mesh.Set("visible", b)
}

func (dao3d *DangerObj3D) SetFieldPosition(fx, fy int, shX, shY, shZ float64) {
	SetPosition(
		dao3d.Mesh,
		shX+float64(fx)*DstCellSize+DstCellSize/2,
		-shY-float64(fy)*DstCellSize-DstCellSize/2,
		shZ+DstCellSize/2,
	)
}

func (dao3d *DangerObj3D) RotateX(rad float64) {
	dao3d.Mesh.Get("rotation").Set("x", rad)
}
func (dao3d *DangerObj3D) RotateY(rad float64) {
	dao3d.Mesh.Get("rotation").Set("y", rad)
}
func (dao3d *DangerObj3D) RotateZ(rad float64) {
	dao3d.Mesh.Get("rotation").Set("z", rad)
}

func (dao3d *DangerObj3D) ScaleX(x float64) {
	dao3d.Mesh.Get("scale").Set("x", x)
}
func (dao3d *DangerObj3D) ScaleY(y float64) {
	dao3d.Mesh.Get("scale").Set("y", y)
}
func (dao3d *DangerObj3D) ScaleZ(z float64) {
	dao3d.Mesh.Get("scale").Set("z", z)
}

func (dao3d *DangerObj3D) Dispose() {
	// mesh do not need dispose
	dao3d.Mesh.Get("geometry").Call("dispose")
	dao3d.Mesh.Get("material").Call("dispose")
	dao3d.Mesh = js.Undefined()
	dao3d.Dao = nil
	// no need createElement canvas dom obj
}
