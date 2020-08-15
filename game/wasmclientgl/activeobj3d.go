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
	"syscall/js"

	"github.com/kasworld/goguelike/enum/condition_flag"

	"github.com/kasworld/goguelike/enum/condition"

	"github.com/kasworld/goguelike/enum/factiontype"
)

var gActiveObj3DGeo [factiontype.FactionType_Count]struct {
	Geo     js.Value
	GeoInfo GeoInfo
}

func preMakeActiveObj3DGeo() {
	for i := 0; i < factiontype.FactionType_Count; i++ {
		ftstr := factiontype.FactionType(i).Rune()
		geo := ThreeJsNew("TextGeometry", ftstr,
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
		geo.Call("center")
		gActiveObj3DGeo[i].Geo = geo
		gActiveObj3DGeo[i].GeoInfo = GetGeoInfo(geo)
	}
}

type ActiveObj3D struct {
	Faction   factiontype.FactionType
	Condition [condition.Condition_Count]*Condition3D
	MoveArrow *ColorArrow3D
	Name      *Label3D
	Chat      *Label3D
	Mesh      js.Value
}

func NewActiveObj3D(ft factiontype.FactionType, name string) *ActiveObj3D {
	mat := gPoolColorMaterial.Get(ft.Color24().ToHTMLColorString())
	// mat.Set("transparent", true)
	mat.Set("opacity", 1)

	geo := gActiveObj3DGeo[ft].Geo
	mesh := ThreeJsNew("Mesh", geo, mat)
	ao3d := &ActiveObj3D{
		MoveArrow: gPoolColorArrow3D.Get("#ffffff"),
		Name:      gPoolLabel3D.Get(name),
		Faction:   ft,
		Mesh:      mesh,
	}
	for i := range ao3d.Condition {
		ao3d.Condition[i] = gPoolCondition3D.Get(condition.Condition(i))
	}

	return ao3d
}

// return changed
func (ao3d *ActiveObj3D) ChangeFaction(ft factiontype.FactionType) (js.Value, bool) {
	if ft == ao3d.Faction {
		return ao3d.Mesh, false
	}
	oldmesh := ao3d.Mesh
	gPoolColorMaterial.Put(ao3d.Mesh.Get("material"))
	mat := gPoolColorMaterial.Get(ft.Color24().ToHTMLColorString())
	mat.Set("opacity", 1)
	geo := gActiveObj3DGeo[ft].Geo
	mesh := ThreeJsNew("Mesh", geo, mat)
	ao3d.Faction = ft
	ao3d.Mesh = mesh
	return oldmesh, true
}

func (ao3d *ActiveObj3D) SetFieldPosition(fx, fy int, shZ float64, cnf condition_flag.ConditionFlag) {
	geoinfo := gActiveObj3DGeo[ao3d.Faction].GeoInfo
	SetPosition(
		ao3d.Mesh,
		float64(fx)*DstCellSize+DstCellSize/2,
		-float64(fy)*DstCellSize-DstCellSize/2,
		geoinfo.Len[2]/2+1+shZ,
	)
	ao3d.Name.SetFieldPosition(fx, fy, 0, DstCellSize, DstCellSize+2+shZ)
	for i, v := range ao3d.Condition {
		v.SetFieldPosition(fx, fy, float64(i)*DstCellSize/8, DstCellSize, DstCellSize+2+shZ)
		v.Visible(cnf.TestByCondition(condition.Condition(i)))
	}
}

func (ao3d *ActiveObj3D) Visible(b bool) {
	ao3d.Mesh.Set("visible", b)
}

func (ao3d *ActiveObj3D) ResetMatrix() {
	ao3d.ScaleX(1.0)
	ao3d.ScaleY(1.0)
	ao3d.ScaleZ(1.0)
	ao3d.RotateX(0.0)
	ao3d.RotateY(0.0)
	ao3d.RotateZ(0.0)
}

func (ao3d *ActiveObj3D) RotateX(rad float64) {
	ao3d.Mesh.Get("rotation").Set("x", rad)
}
func (ao3d *ActiveObj3D) RotateY(rad float64) {
	ao3d.Mesh.Get("rotation").Set("y", rad)
}
func (ao3d *ActiveObj3D) RotateZ(rad float64) {
	ao3d.Mesh.Get("rotation").Set("z", rad)
}

func (ao3d *ActiveObj3D) ScaleX(x float64) {
	ao3d.Mesh.Get("scale").Set("x", x)
}
func (ao3d *ActiveObj3D) ScaleY(y float64) {
	ao3d.Mesh.Get("scale").Set("y", y)
}
func (ao3d *ActiveObj3D) ScaleZ(z float64) {
	ao3d.Mesh.Get("scale").Set("z", z)
}

func (ao3d *ActiveObj3D) Dispose() {
	for i := range ao3d.Condition {
		gPoolCondition3D.Put(ao3d.Condition[i])
		ao3d.Condition[i] = nil
	}

	// mesh do not need dispose
	// ao3d.Mesh.Get("geometry").Call("dispose")
	// ao3d.Mesh.Get("material").Call("dispose")
	gPoolColorMaterial.Put(ao3d.Mesh.Get("material"))
	ao3d.Mesh = js.Undefined()

	gPoolColorArrow3D.Put(ao3d.MoveArrow)
	ao3d.MoveArrow = nil
	gPoolLabel3D.Put(ao3d.Name)
	ao3d.Name = nil
	// no need createElement canvas dom obj
}
