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

	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/lib/webtilegroup"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type CarryObj3D struct {
	Cnv     js.Value
	Ctx     js.Value
	Tex     js.Value
	GeoInfo GeoInfo
	Mesh    js.Value
}

func NewCarryObj3D() *CarryObj3D {
	cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	ctx := cnv.Call("getContext", "2d")
	ctx.Set("imageSmoothingEnabled", false)
	cnv.Set("width", DstCellSize)
	cnv.Set("height", DstCellSize)
	tex := ThreeJsNew("CanvasTexture", cnv)
	mat := ThreeJsNew("MeshPhongMaterial",
		map[string]interface{}{
			"map": tex,
		},
	)
	mat.Set("transparent", true)
	geo := ThreeJsNew("BoxGeometry", DstCellSize, DstCellSize, DstCellSize)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &CarryObj3D{
		Cnv:     cnv,
		Ctx:     ctx,
		Tex:     tex,
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func (aog *CarryObj3D) ChangeTile(ti webtilegroup.TileInfo) {
	aog.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)
	aog.Tex.Set("needsUpdate", true)
}

func (aog *CarryObj3D) Dispose() {
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

func Equiped2TileInfo(o *c2t_obj.EquipClient) webtilegroup.TileInfo {
	return gClientTile.EquipTiles[o.EquipType][o.Faction]
}

func CarryObj2TileInfo(o *c2t_obj.CarryObjClientOnFloor) webtilegroup.TileInfo {
	var ti webtilegroup.TileInfo
	switch o.CarryingObjectType {
	case carryingobjecttype.Equip:
		ti = gClientTile.EquipTiles[o.EquipType][o.Faction]
	case carryingobjecttype.Money:
		var find bool
		for i, v := range moneycolor.Attrib {
			if o.Value < v.UpLimit {
				ti = gClientTile.GoldTiles[i]
				find = true
				break
			}
		}
		if !find {
			ti = gClientTile.GoldTiles[len(gClientTile.GoldTiles)-1]
		}
	case carryingobjecttype.Potion:
		ti = gClientTile.PotionTiles[o.PotionType]
	case carryingobjecttype.Scroll:
		ti = gClientTile.ScrollTiles[o.ScrollType]
	}
	return ti
}

func MakeEquipedMesh(o *c2t_obj.EquipClient) js.Value {
	ti := Equiped2TileInfo(o)
	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(
		DstCellSize/3, DstCellSize/3, DstCellSize/3,
	)
	return ThreeJsNew("Mesh", geo, mat)
}

func MakeCarryObjMesh(o *c2t_obj.CarryObjClientOnFloor) js.Value {
	ti := CarryObj2TileInfo(o)
	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(
		DstCellSize/3, DstCellSize/3, DstCellSize/3,
	)
	return ThreeJsNew("Mesh", geo, mat)
}

type ShiftInfo struct {
	X float64
	Y float64
	Z float64
}

// equipped shift, around ao
var aoEqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	equipslottype.Helmet: {-0.33, 0.0, 0.66},
	equipslottype.Amulet: {1.00, 0.0, 0.66},

	equipslottype.Weapon: {-0.33, 0.25, 0.66},
	equipslottype.Shield: {1.00, 0.25, 0.66},

	equipslottype.Ring:     {-0.33, 0.50, 0.66},
	equipslottype.Gauntlet: {1.00, 0.50, 0.66},

	equipslottype.Armor:    {-0.33, 0.75, 0.66},
	equipslottype.Footwear: {1.00, 0.75, 0.66},
}

func CarryObjClientOnFloor2DrawInfo(
	co *c2t_obj.CarryObjClientOnFloor) ShiftInfo {
	switch co.CarryingObjectType {
	default:
		return otherCarryObjShift[co.CarryingObjectType]
	case carryingobjecttype.Equip:
		return eqPosShift[co.EquipType]
	}
}

// on floor in tile
var eqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	equipslottype.Helmet: {0.0, 0.0, 0.33},
	equipslottype.Amulet: {0.75, 0.0, 0.33},

	equipslottype.Weapon: {0.0, 0.25, 0.33},
	equipslottype.Shield: {0.75, 0.25, 0.33},

	equipslottype.Ring:     {0.0, 0.50, 0.33},
	equipslottype.Gauntlet: {0.75, 0.50, 0.33},

	equipslottype.Armor:    {0.0, 0.75, 0.33},
	equipslottype.Footwear: {0.75, 0.75, 0.33},
}

var otherCarryObjShift = [carryingobjecttype.CarryingObjectType_Count]ShiftInfo{
	carryingobjecttype.Money:  {0.33, 0.0, 0.33},
	carryingobjecttype.Potion: {0.33, 0.33, 0.33},
	carryingobjecttype.Scroll: {0.33, 0.66, 0.33},
}
