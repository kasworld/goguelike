// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

func (cf *ClientFloorGL) drawTileAt(fx, fy int, tl tile_flag.TileFlag) {
	dstX := fx * DstCellSize
	dstY := fy * DstCellSize
	diffbase := fx*5 + fy*3

	for i := 0; i < tile.Tile_Count; i++ {
		shX := 0.0
		shY := 0.0
		tlt := tile.Tile(i)
		if tl.TestByTile(tlt) {
			if gTextureTileList[i] != nil {
				// texture tile
				tlic := gTextureTileList[i]
				srcx, srcy, srcCellSize := gTextureTileList[i].CalcSrc(fx, fy, shX, shY)
				cf.PlaneTile.Ctx.Call("drawImage", tlic.Cnv,
					srcx, srcy, srcCellSize, srcCellSize,
					dstX, dstY, DstCellSize, DstCellSize)

			} else if tlt == tile.Wall {
				// wall tile process
				tlList := gClientTile.FloorTiles[i]
				tilediff := cf.calcWallTileDiff(fx, fy)
				ti := tlList[tilediff%len(tlList)]
				cf.PlaneTile.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, DstCellSize, DstCellSize)
			} else if tlt == tile.Window {
				// window tile process
				tlList := gClientTile.FloorTiles[i]
				tlindex := 0
				if cf.checkWallAt(fx, fy-1) && cf.checkWallAt(fx, fy+1) { // n-s window
					tlindex = 1
				}
				ti := tlList[tlindex]
				cf.PlaneTile.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, DstCellSize, DstCellSize)
			} else {
				// bitmap tile
				tlList := gClientTile.FloorTiles[i]
				tilediff := diffbase + int(shX)
				ti := tlList[tilediff%len(tlList)]
				cf.PlaneTile.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, DstCellSize, DstCellSize)
			}
		}
	}
}

func (cf *ClientFloorGL) calcWallTileDiff(flx, fly int) int {
	rtn := 0
	if cf.checkWallAt(flx, fly-1) {
		rtn |= 1
	}
	if cf.checkWallAt(flx+1, fly) {
		rtn |= 1 << 1
	}
	if cf.checkWallAt(flx, fly+1) {
		rtn |= 1 << 2
	}
	if cf.checkWallAt(flx-1, fly) {
		rtn |= 1 << 3
	}
	return rtn
}

func (cf *ClientFloorGL) checkWallAt(flx, fly int) bool {
	flx = cf.XWrapSafe(flx)
	fly = cf.YWrapSafe(fly)
	tl := cf.Tiles[flx][fly]
	return tl.TestByTile(tile.Wall) ||
		tl.TestByTile(tile.Door) ||
		tl.TestByTile(tile.Window)
}

func (cf *ClientFloorGL) Draw(
	frameProgress float64,
	scrollDir way9type.Way9Type,
	taNoti *c2t_obj.NotiVPTiles_data,
) {
	zoom := gameOptions.GetByIDBase("Zoom").State
	sx, sy := calcShiftDxDy(frameProgress)
	scrollDx := -scrollDir.Dx() * sx
	scrollDy := scrollDir.Dy() * sy

	// move camera, light
	cameraX := taNoti.VPX*DstCellSize + scrollDx
	cameraY := -taNoti.VPY*DstCellSize + scrollDy
	SetPosition(cf.light,
		cameraX, cameraY, DstCellSize*2,
	)
	cameraZ := HelperSize / (zoom + 1)
	SetPosition(cf.camera,
		cameraX, cameraY, cameraZ,
	)
	cf.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			cameraX, cameraY, 0,
		),
	)

	cf.PlaneSight.fillColor("#00000010")
	cf.PlaneSight.clearSight(taNoti.VPX, taNoti.VPY, taNoti.VPTiles)
}

func (cf *ClientFloorGL) Resize(w, h float64) {
	cf.camera.Set("aspect", w/h)
	cf.camera.Call("updateProjectionMatrix")
}

func calcShiftDxDy(frameProgress float64) (int, int) {
	rate := 1 - frameProgress
	// if rate < 0 {
	// 	rate = 0
	// }
	// if rate > 1 {
	// 	rate = 1
	// }
	dx := int(float64(DstCellSize) * rate)
	dy := int(float64(DstCellSize) * rate)
	return dx, dy
}

func (cf *ClientFloorGL) processNotiObjectList(
	olNoti *c2t_obj.NotiObjectList_data) {
	// shY := int(-float64(DstCellSize) * 0.8)
	addUUID := make(map[string]bool)

	// make activeobj
	for _, o := range olNoti.ActiveObjList {
		jso, exist := cf.jsSceneObjs[o.UUID]
		if !exist {
			geo := GetTextGeometryByCache(
				o.Faction.String()[:2],
				DstCellSize/2.0,
			)
			mat := GetColorMaterialByCache(uint32(o.Faction.Color24()))
			jso = ThreeJsNew("Mesh", geo, mat)
			cf.scene.Call("add", jso)
			cf.jsSceneObjs[o.UUID] = jso
		}
		miny, maxy := CalcGeoMinMaxY(jso.Get("geometry"))
		SetPosition(
			jso,
			float64(o.X)*DstCellSize,
			-float64(o.Y)*DstCellSize-(maxy-miny)/2-DstCellSize/2,
			0)
		addUUID[o.UUID] = true
	}

	// make carryobj
	for _, o := range olNoti.CarryObjList {
		jso, exist := cf.jsSceneObjs[o.UUID]
		str, co, posinfo := carryObjClientOnFloor2DrawInfo(o)
		if !exist {
			geo := GetTextGeometryByCache(
				str,
				DstCellSize/2*posinfo.W,
			)
			mat := GetColorMaterialByCache(co)
			jso = ThreeJsNew("Mesh", geo, mat)
			cf.scene.Call("add", jso)
			cf.jsSceneObjs[o.UUID] = jso
		}
		miny, maxy := CalcGeoMinMaxY(jso.Get("geometry"))
		SetPosition(
			jso,
			float64(o.X)*DstCellSize+DstCellSize*posinfo.X,
			-float64(o.Y)*DstCellSize-DstCellSize*posinfo.Y-(maxy-miny)/2,
			0)
		addUUID[o.UUID] = true
	}

	for id, jso := range cf.jsSceneObjs {
		if !addUUID[id] {
			cf.scene.Call("remove", jso)
			delete(cf.jsSceneObjs, id)
		}
	}
}

func (cf *ClientFloorGL) addFieldObj(o *c2t_obj.FieldObjClient) {
	oldx, oldy, exist := cf.FieldObjPosMan.GetXYByUUID(o.ID)
	if exist && o.X == oldx && o.Y == oldy {
		return // no need to add
	}
	if exist {
		cf.FieldObjPosMan.UpdateToXY(o, o.X, o.Y)
		// clear old rect

		cf.drawFieldObj(o) // draw at new pos
		return             // move exist obj
	}
	// add new obj
	cf.FieldObjPosMan.AddToXY(o, o.X, o.Y)
	cf.drawFieldObj(o)
}

func (cf *ClientFloorGL) drawFieldObj(o *c2t_obj.FieldObjClient) {
	tlList := gClientTile.FieldObjTiles[o.DisplayType]
	if len(tlList) == 0 {
		jslog.Errorf("len=0 %v", o.DisplayType)
		return
	}
	fx := o.X
	fy := o.Y
	diffbase := fx*5 + fy*3
	tilediff := diffbase
	ti := tlList[tilediff%len(tlList)]

	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(DstCellSize, DstCellSize, DstCellSize)
	mesh := ThreeJsNew("Mesh", geo, mat)
	cf.scene.Call("add", mesh)
	SetPosition(
		mesh,
		float64(o.X)*DstCellSize+DstCellSize/2,
		-float64(o.Y)*DstCellSize-DstCellSize/2,
		DstCellSize/2)
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
