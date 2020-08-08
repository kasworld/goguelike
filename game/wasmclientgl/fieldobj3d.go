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

	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
)

var gPoolFieldObj3D = NewPoolFieldObj3D()

type FOKey struct {
	ActType     fieldobjacttype.FieldObjActType
	DisplayType fieldobjdisplaytype.FieldObjDisplayType
}
type PoolFieldObj3D struct {
	mutex    sync.Mutex
	poolData map[FOKey][]*FieldObj3D
	newCount int
	getCount int
	putCount int
}

func NewPoolFieldObj3D() *PoolFieldObj3D {
	return &PoolFieldObj3D{
		poolData: make(map[FOKey][]*FieldObj3D),
	}
}

func (p *PoolFieldObj3D) String() string {
	return fmt.Sprintf("PoolFieldObj3D[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolFieldObj3D) Get(
	ActType fieldobjacttype.FieldObjActType,
	DisplayType fieldobjdisplaytype.FieldObjDisplayType) *FieldObj3D {

	fokey := FOKey{
		ActType, DisplayType,
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *FieldObj3D
	if l := len(p.poolData[fokey]); l > 0 {
		rtn = p.poolData[fokey][l-1]
		p.poolData[fokey] = p.poolData[fokey][:l-1]
		p.getCount++
	} else {
		rtn = NewFieldObj3D(ActType, DisplayType)
		p.newCount++
	}
	return rtn
}

func (p *PoolFieldObj3D) Put(pb *FieldObj3D) {
	fokey := FOKey{
		pb.ActType, pb.DisplayType,
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[fokey] = append(p.poolData[fokey], pb)
	p.putCount++
}

func FieldObj2StrColor(
	ActType fieldobjacttype.FieldObjActType,
	DisplayType fieldobjdisplaytype.FieldObjDisplayType) (string, string) {

	str := DisplayType.String()[:2]
	co := ActType.Color24().ToHTMLColorString()
	return str, co
}

func NewFieldObjGeo(str string) js.Value {
	geo := ThreeJsNew("TextGeometry", str,
		map[string]interface{}{
			"font":           gFont_droid_sans_mono_regular,
			"size":           DstCellSize * 0.7,
			"height":         DstCellSize * 0.3,
			"curveSegments":  DstCellSize / 3,
			"bevelEnabled":   true,
			"bevelThickness": DstCellSize / 8,
			"bevelSize":      DstCellSize / 16,
			"bevelOffset":    0,
			"bevelSegments":  DstCellSize / 8,
		})
	geo.Call("rotateX", math.Pi/2)
	geo.Call("center")
	return geo
}

type FieldObj3D struct {
	ActType     fieldobjacttype.FieldObjActType
	DisplayType fieldobjdisplaytype.FieldObjDisplayType
	Label       *Label3D
	GeoInfo     GeoInfo
	Mesh        js.Value
}

func NewFieldObj3D(
	ActType fieldobjacttype.FieldObjActType,
	DisplayType fieldobjdisplaytype.FieldObjDisplayType) *FieldObj3D {

	str, color := FieldObj2StrColor(ActType, DisplayType)
	mat := GetColorMaterialByCache(color)
	mat.Set("transparent", true)
	geo := NewFieldObjGeo(str)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &FieldObj3D{
		ActType:     ActType,
		DisplayType: DisplayType,
		GeoInfo:     GetGeoInfo(geo),
		Mesh:        mesh,
	}
}

func (aog *FieldObj3D) SetFieldPosition(fx, fy int, shZ float64) {
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+DstCellSize/2,
		-float64(fy)*DstCellSize-DstCellSize/2,
		aog.GeoInfo.Len[2]/2+1+shZ,
	)
}

func (aog *FieldObj3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *FieldObj3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *FieldObj3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *FieldObj3D) Dispose() {
	// mesh do not need dispose
	// aog.Mesh.Get("geometry").Call("dispose")
	// aog.Mesh.Get("material").Call("dispose")
	aog.Mesh = js.Undefined()
	// no need createElement canvas dom obj
}
