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

	"github.com/kasworld/goguelike/enum/factiontype"
)

var gActiveObj3DGeo [factiontype.FactionType_Count]struct {
	Geo     js.Value
	GeoInfo GeoInfo
}

func preMakeActiveObj3DGeo() {
	ftList := [factiontype.FactionType_Count]string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}
	for i, str := range ftList {
		geo := ThreeJsNew("TextGeometry", str,
			map[string]interface{}{
				"font":           gFont_droid_sans_mono_regular,
				"size":           DstCellSize / 2,
				"height":         DstCellSize / 3,
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
	Faction factiontype.FactionType
	Name    *Label3D
	Chat    *Label3D
	Mesh    js.Value
}

func NewActiveObj3D(ft factiontype.FactionType) *ActiveObj3D {
	mat := GetColorMaterialByCache(ft.Color24().ToHTMLColorString())
	geo := gActiveObj3DGeo[ft].Geo
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &ActiveObj3D{
		Faction: ft,
		Mesh:    mesh,
	}
}

// return changed
func (aog *ActiveObj3D) ChangeFaction(ft factiontype.FactionType) (js.Value, bool) {
	if ft == aog.Faction {
		return aog.Mesh, false
	}
	oldmesh := aog.Mesh
	aog.Mesh.Get("geometry").Call("dispose")
	aog.Mesh.Get("material").Call("dispose")
	mat := GetColorMaterialByCache(ft.Color24().ToHTMLColorString())
	geo := gActiveObj3DGeo[ft].Geo
	mesh := ThreeJsNew("Mesh", geo, mat)
	aog.Faction = ft
	aog.Mesh = mesh
	return oldmesh, true
}

func (aog *ActiveObj3D) SetFieldPosition(fx, fy int, shZ float64) {
	geoinfo := gActiveObj3DGeo[aog.Faction].GeoInfo
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+DstCellSize/2,
		-float64(fy)*DstCellSize-DstCellSize/2,
		geoinfo.Len[2]/2+1+shZ,
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
	aog.Mesh = js.Undefined()
	// no need createElement canvas dom obj
}
