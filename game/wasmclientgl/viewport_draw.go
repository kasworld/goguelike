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
	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

func (vp *Viewport) processNotiObjectList(
	cf *ClientFloorGL,
	olNoti *c2t_obj.NotiObjectList_data) {
	// shY := int(-float64(DstCellSize) * 0.8)
	addUUID := make(map[g2id.G2ID]bool)

	// make activeobj
	for _, o := range olNoti.ActiveObjList {
		jso, exist := vp.jsSceneObjs[o.G2ID]
		if !exist {
			geo := vp.getTextGeometry(
				o.Faction.String()[:2],
				DstCellSize/2.0,
			)
			mat := vp.getColorMaterial(uint32(o.Faction.Color24()))
			jso = vp.ThreeJsNew("Mesh", geo, mat)
			vp.scene.Call("add", jso)
			vp.jsSceneObjs[o.G2ID] = jso
		}
		miny, maxy := vp.calcGeoMinMaxY(jso.Get("geometry"))
		SetPosition(
			jso,
			float64(o.X)*DstCellSize,
			-float64(o.Y)*DstCellSize-(maxy-miny)/2-DstCellSize/2,
			0)
		addUUID[o.G2ID] = true
	}

	// make carryobj
	for _, o := range olNoti.CarryObjList {
		jso, exist := vp.jsSceneObjs[o.G2ID]
		str, co, posinfo := carryObjClientOnFloor2DrawInfo(o)
		if !exist {
			geo := vp.getTextGeometry(
				str,
				DstCellSize/2*posinfo.W,
			)
			mat := vp.getColorMaterial(co)
			jso = vp.ThreeJsNew("Mesh", geo, mat)
			vp.scene.Call("add", jso)
			vp.jsSceneObjs[o.G2ID] = jso
		}
		miny, maxy := vp.calcGeoMinMaxY(jso.Get("geometry"))
		SetPosition(
			jso,
			float64(o.X)*DstCellSize+DstCellSize*posinfo.X,
			-float64(o.Y)*DstCellSize-DstCellSize*posinfo.Y-(maxy-miny)/2,
			0)
		addUUID[o.G2ID] = true
	}

	for id, jso := range vp.jsSceneObjs {
		if !addUUID[id] {
			vp.scene.Call("remove", jso)
			delete(vp.jsSceneObjs, id)
		}
	}

	// draw fieldobj
	for _, o := range olNoti.FieldObjList {
		tlList := gClientTile.FieldObjTiles[o.DisplayType]
		if len(tlList) == 0 {
			jslog.Errorf("len=0 %v", o.DisplayType)
			continue
		}
		fx := o.X
		fy := o.Y
		dstX := fx * DstCellSize
		dstY := fy * DstCellSize

		diffbase := fx*5 + fy*3
		tilediff := diffbase
		ti := tlList[tilediff%len(tlList)]
		cf.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
			ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
			dstX, dstY, DstCellSize, DstCellSize)

	}
}

func carryObjClientOnFloor2DrawInfo(
	co *c2t_obj.CarryObjClientOnFloor) (string, uint32, coShift) {
	switch co.CarryingObjectType {
	case carryingobjecttype.Money:
		for _, v := range moneycolor.Attrib {
			if co.Value < v.UpLimit {
				return "$", uint32(v.Color),
					coShiftOther[co.CarryingObjectType]
			}
		}
		v := moneycolor.Attrib[len(moneycolor.Attrib)-1]
		return "$", uint32(v.Color),
			coShiftOther[co.CarryingObjectType]

	case carryingobjecttype.Potion:
		return "!", uint32(co.PotionType.Color24()),
			coShiftOther[co.CarryingObjectType]
	case carryingobjecttype.Scroll:
		return "~", uint32(co.ScrollType.Color24()),
			coShiftOther[co.CarryingObjectType]
	case carryingobjecttype.Equip:
		return eqtype2string[co.EquipType], uint32(co.Faction),
			eqposShift[co.EquipType]
	}
	return "?", 0x000000, coShiftOther[0]
}

var eqtype2string = [equipslottype.EquipSlotType_Count]string{
	equipslottype.Weapon:   "/",
	equipslottype.Shield:   "#",
	equipslottype.Helmet:   "^",
	equipslottype.Armor:    "%",
	equipslottype.Gauntlet: "=",
	equipslottype.Footwear: "_",
	equipslottype.Ring:     "o",
	equipslottype.Amulet:   "*",
}

type coShift struct {
	X float64
	Y float64
	W float64
}

var eqposShift = [equipslottype.EquipSlotType_Count]coShift{
	equipslottype.Helmet: {0.0, 0.0, 0.33},
	equipslottype.Amulet: {0.75, 0.0, 0.33},

	equipslottype.Weapon: {0.0, 0.25, 0.33},
	equipslottype.Shield: {0.75, 0.25, 0.33},

	equipslottype.Ring:     {0.0, 0.50, 0.33},
	equipslottype.Gauntlet: {0.75, 0.50, 0.33},

	equipslottype.Armor:    {0.0, 0.75, 0.33},
	equipslottype.Footwear: {0.75, 0.75, 0.33},
}

var coShiftOther = [carryingobjecttype.CarryingObjectType_Count]coShift{
	carryingobjecttype.Money:  {0.33, 0.0, 0.33},
	carryingobjecttype.Potion: {0.33, 0.33, 0.33},
	carryingobjecttype.Scroll: {0.33, 0.66, 0.33},
}
