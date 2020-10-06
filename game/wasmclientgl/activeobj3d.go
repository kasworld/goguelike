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
	"syscall/js"

	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

var gActiveObj3DGeo [factiontype.FactionType_Count]struct {
	Geo     js.Value
	GeoInfo GeoInfo
}

func preMakeActiveObj3DGeo() {
	for i := 0; i < factiontype.FactionType_Count; i++ {
		// ftstr := factiontype.FactionType(i).Rune()
		// geo := MakeAO3DGeoByRune(ftstr)
		geo := MakeAO3DGeo()
		gActiveObj3DGeo[i].Geo = geo
		gActiveObj3DGeo[i].GeoInfo = GetGeoInfo(geo)
	}
}

func MakeAO3DGeoByRune(ftstr string) js.Value {
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
	return geo
}

func MakeAO3DGeo() js.Value {
	matrix := ThreeJsNew("Matrix4")
	geoHead := ThreeJsNew("SphereGeometry", DstCellSize/6, DstCellSize, DstCellSize)

	geoArm := ThreeJsNew("CylinderGeometry", DstCellSize/16, DstCellSize/16, DstCellSize)
	geoArm.Call("rotateZ", math.Pi/2)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, -DstCellSize/4, 0,
	))
	geoHead.Call("merge", geoArm, matrix)
	geoArm.Call("dispose")

	geoBody := ThreeJsNew("CylinderGeometry", DstCellSize/16, DstCellSize/16, DstCellSize/2)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		0, -DstCellSize/2.5, 0,
	))
	geoHead.Call("merge", geoBody, matrix)
	geoBody.Call("dispose")

	geoLegL := ThreeJsNew("CylinderGeometry", DstCellSize/16, DstCellSize/16, DstCellSize/2)
	geoLegL.Call("rotateZ", -math.Pi/4)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		-DstCellSize/4, -DstCellSize/2-DstCellSize/4, 0,
	))
	geoHead.Call("merge", geoLegL, matrix)
	geoLegL.Call("dispose")

	geoLegR := ThreeJsNew("CylinderGeometry", DstCellSize/16, DstCellSize/16, DstCellSize/2)
	geoLegR.Call("rotateZ", math.Pi/4)
	matrix.Call("setPosition", ThreeJsNew("Vector3",
		DstCellSize/4, -DstCellSize/2-DstCellSize/4, 0,
	))
	geoHead.Call("merge", geoLegR, matrix)
	geoLegR.Call("dispose")

	geoHead.Call("rotateX", math.Pi/2)

	geoHead.Call("center")
	return geoHead
}

type ActiveObj3D struct {
	AOC       *c2t_obj.ActiveObjClient
	Co3d      [equipslottype.EquipSlotType_Count]*CarryObj3D
	Condition [condition.Condition_Count]*Condition3D
	Name      *Label3D
	Chat      *Label3D
	Mesh      js.Value
}

func NewActiveObj3D(aoc *c2t_obj.ActiveObjClient) *ActiveObj3D {
	mat := gPoolColorMaterial.Get(aoc.Faction.Color24().ToHTMLColorString())
	// mat.Set("transparent", true)
	mat.Set("opacity", 1)

	geo := gActiveObj3DGeo[aoc.Faction].Geo
	mesh := ThreeJsNew("Mesh", geo, mat)
	mesh.Set("castShadow", true)
	mesh.Set("receiveShadow", false)
	ao3d := &ActiveObj3D{
		AOC:  aoc,
		Name: gPoolLabel3D.Get(aoc.NickName),
		Mesh: mesh,
	}
	for i := range ao3d.Condition {
		ao3d.Condition[i] = gPoolCondition3D.Get(condition.Condition(i))
	}
	return ao3d
}

// return toaddmesh, toremove mesh
func (ao3d *ActiveObj3D) UpdateAOC(newaoc *c2t_obj.ActiveObjClient) ([]js.Value, []js.Value) {
	var toaddMeshs []js.Value
	var todelMeshs []js.Value

	oldaoc := ao3d.AOC
	ao3d.AOC = newaoc
	if newaoc.Faction != oldaoc.Faction {
		todelMeshs = append(todelMeshs, ao3d.Mesh)
		gPoolColorMaterial.Put(ao3d.Mesh.Get("material"))
		mat := gPoolColorMaterial.Get(newaoc.Faction.Color24().ToHTMLColorString())
		mat.Set("opacity", 1)
		geo := gActiveObj3DGeo[newaoc.Faction].Geo
		mesh := ThreeJsNew("Mesh", geo, mat)
		ao3d.Mesh = mesh
		toaddMeshs = append(toaddMeshs, ao3d.Mesh)
	}

	if oldaoc == newaoc { // just made
		for _, v := range newaoc.EquippedPo {
			str, color := Equiped2StrColor(v)
			toaddco3d := NewCarryObj3D(str, color)
			toaddco3d.RotateX(math.Pi / 2)
			toaddMeshs = append(toaddMeshs, toaddco3d.Mesh)
			ao3d.Co3d[v.EquipType] = toaddco3d
		}

	} else { // handle changed equipped carryobj
		coadded, codel := findEqChanged(oldaoc.EquippedPo, newaoc.EquippedPo)
		for _, v := range codel {
			todelco3d := ao3d.Co3d[v.EquipType]
			ao3d.Co3d[v.EquipType] = nil
			todelMeshs = append(todelMeshs, todelco3d.Mesh)
			todelco3d.Dispose()
		}
		for _, v := range coadded {
			str, color := Equiped2StrColor(v)
			toaddco3d := NewCarryObj3D(str, color)
			toaddco3d.RotateX(math.Pi / 2)
			toaddMeshs = append(toaddMeshs, toaddco3d.Mesh)
			ao3d.Co3d[v.EquipType] = toaddco3d
		}

	}

	if len(ao3d.AOC.Chat) == 0 {
		if ao3d.Chat != nil {
			todelMeshs = append(todelMeshs, ao3d.Chat.Mesh)
			ao3d.Chat.Dispose()
			ao3d.Chat = nil
		}
	} else {
		if ao3d.Chat == nil {
			// add new chat
			ao3d.Chat = NewLabel3D(ao3d.AOC.Chat)
			toaddMeshs = append(toaddMeshs, ao3d.Chat.Mesh)
		} else {
			if ao3d.AOC.Chat != ao3d.Chat.Str {
				// remove old chat , add new chat
				todelMeshs = append(todelMeshs, ao3d.Chat.Mesh)
				ao3d.Chat.Dispose()
				ao3d.Chat = NewLabel3D(ao3d.AOC.Chat)
				toaddMeshs = append(toaddMeshs, ao3d.Chat.Mesh)
			}
		}
	}

	return toaddMeshs, todelMeshs
}

func (ao3d *ActiveObj3D) SetFieldPosition(fx, fy int, shX, shY, shZ float64) {
	geoinfo := gActiveObj3DGeo[ao3d.AOC.Faction].GeoInfo
	if ao3d.AOC.Conditions.TestByCondition(condition.Float) {
		shZ += DstCellSize
	}
	SetPosition(
		ao3d.Mesh,
		shX+float64(fx)*DstCellSize+DstCellSize/2,
		-shY+-float64(fy)*DstCellSize-DstCellSize/2,
		shZ+geoinfo.Len[2]/2+1,
	)

	nameDy := float64(
		gameOptions.GetByIDBase("Angle").State) * DstCellSize / 2
	ao3d.Name.SetFieldPosition(fx, fy, shX, shY+DstCellSize+nameDy, shZ+DstCellSize+2)
	for i, v := range ao3d.Condition {
		v.SetFieldPosition(fx, fy,
			shX+float64(i)*DstCellSize/8,
			shY+DstCellSize+DstCellSize/2,
			shZ+DstCellSize+2)
		v.Visible(ao3d.AOC.Conditions.TestByCondition(condition.Condition(i)))
	}
	for _, v := range ao3d.AOC.EquippedPo {
		if ao3d.Co3d[v.EquipType] == nil {
			jslog.Errorf("%v", v)
			continue
		}
		shInfo := aoEqPosShift[v.EquipType]
		ao3d.Co3d[v.EquipType].SetFieldPosition(fx, fy, shInfo.X+shX, shInfo.Y+shY, shInfo.Z+shZ)
	}
	if ao3d.Chat != nil {
		ao3d.Chat.SetFieldPosition(fx, fy,
			shX+0, shY-DstCellSize/2, shZ+DstCellSize+2)
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
	for _, v := range ao3d.Co3d {
		if v == nil {
			continue
		}
		v.Dispose()
	}
	// mesh do not need dispose
	gPoolColorMaterial.Put(ao3d.Mesh.Get("material"))
	ao3d.Mesh = js.Undefined()
	gPoolLabel3D.Put(ao3d.Name)
	ao3d.Name = nil
	ao3d.AOC = nil
	// no need createElement canvas dom obj
}

// return added, deleted
func findEqChanged(old, new []*c2t_obj.EquipClient) (added, deleted []*c2t_obj.EquipClient) {
	var oldmap [equipslottype.EquipSlotType_Count]*c2t_obj.EquipClient
	var newmap [equipslottype.EquipSlotType_Count]*c2t_obj.EquipClient
	for _, v := range old {
		oldmap[v.EquipType] = v
	}
	for _, v := range new {
		newmap[v.EquipType] = v
	}
	for eqpos := 0; eqpos < equipslottype.EquipSlotType_Count; eqpos++ {
		if oldmap[eqpos] == nil && newmap[eqpos] != nil {
			added = append(added, newmap[eqpos])
		}
		if oldmap[eqpos] != nil && newmap[eqpos] == nil {
			deleted = append(deleted, oldmap[eqpos])
		}
		if oldmap[eqpos] != nil && newmap[eqpos] != nil && oldmap[eqpos].UUID != newmap[eqpos].UUID {
			deleted = append(deleted, oldmap[eqpos])
			added = append(added, newmap[eqpos])
		}
	}
	return
}
