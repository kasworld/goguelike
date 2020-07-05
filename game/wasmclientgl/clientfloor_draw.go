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
	"syscall/js"

	"github.com/kasworld/go-abs"

	"github.com/kasworld/goguelike/lib/webtilegroup"

	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

// fx,fy wrapped, no need wrap again
func (cf *ClientFloorGL) drawTileAt(fx, fy int, newTile tile_flag.TileFlag) {
	oldTile := cf.Tiles[fx][fy]

	for _, tl := range []tile.Tile{
		tile.Tree, tile.Grass,
		tile.Wall, tile.Door, tile.Window,
		tile.Swamp, tile.Soil, tile.Stone, tile.Sand, tile.Sea,
		tile.Magma, tile.Ice, tile.Road, tile.Room,
		tile.Fog, tile.Smoke,
	} {
		if oldTile.TestByTile(tl) && !newTile.TestByTile(tl) {
			// del from scene
			v := cf.jsScene9Tile3D[tl][[2]int{fx, fy}]
			cf.scene.Call("remove", v)
		}
	}

	for i := 0; i < tile.Tile_Count; i++ {
		tlt := tile.Tile(i)
		if !newTile.TestByTile(tlt) {
			continue
		}
		switch tlt {
		default:
			jslog.Errorf("unhandled tile %v", tlt)

		case tile.Tree:
			if oldTile.TestByTile(tlt) {
				continue // skip exist
			}
			mesh, exist := cf.jsScene9Tile3D[tlt][[2]int{fx, fy}]
			if !exist {
				mat := GetTextureTileMaterialByCache(tile.Grass)
				geo := GetConeGeometryByCache(DstCellSize/2-1, DstCellSize-1)
				mesh = cf.make9InstancedMeshAt(mat, geo, fx, fy)
				cf.jsScene9Tile3D[tlt][[2]int{fx, fy}] = mesh
			}
			cf.scene.Call("add", mesh)

		case tile.Grass:
			if oldTile.TestByTile(tlt) {
				continue // skip exist
			}
			mesh, exist := cf.jsScene9Tile3D[tlt][[2]int{fx, fy}]
			if !exist {
				mat := GetTextureTileMaterialByCache(tile.Grass)
				geo := GetBoxGeometryByCache(DstCellSize, DstCellSize, DstCellSize/8)
				mesh = cf.make9InstancedMeshAt(mat, geo, fx, fy)
				cf.jsScene9Tile3D[tlt][[2]int{fx, fy}] = mesh
			}
			cf.scene.Call("add", mesh)

		case tile.Wall:
			if oldTile.TestByTile(tlt) {
				continue // skip exist
			}
			mesh, exist := cf.jsScene9Tile3D[tlt][[2]int{fx, fy}]
			if !exist {
				mat := GetTextureTileMaterialByCache(tile.Stone)
				geo := GetBoxGeometryByCache(DstCellSize, DstCellSize, DstCellSize)
				mesh = cf.make9InstancedMeshAt(mat, geo, fx, fy)
				cf.jsScene9Tile3D[tlt][[2]int{fx, fy}] = mesh
			}
			cf.scene.Call("add", mesh)
		case tile.Window:
			if oldTile.TestByTile(tlt) {
				continue // skip exist
			}
			mesh, exist := cf.jsScene9Tile3D[tlt][[2]int{fx, fy}]
			if !exist {
				ti := gClientTile.CursorTiles[0]
				mat := GetTileMaterialByCache(ti)
				// mat := GetTextureTileMaterialByCache(tile.Fog)
				geo := GetBoxGeometryByCache(DstCellSize, DstCellSize, DstCellSize)
				mesh = cf.make9InstancedMeshAt(mat, geo, fx, fy)
				cf.jsScene9Tile3D[tlt][[2]int{fx, fy}] = mesh
			}
			cf.scene.Call("add", mesh)
		case tile.Door:
			if oldTile.TestByTile(tlt) {
				continue // skip exist
			}
			mesh, exist := cf.jsScene9Tile3D[tlt][[2]int{fx, fy}]
			if !exist {
				tlList := gClientTile.FloorTiles[tile.Door]
				ti := tlList[0]
				mat := GetTileMaterialByCache(ti)
				geo := GetBoxGeometryByCache(DstCellSize, DstCellSize, DstCellSize)
				mesh = cf.make9InstancedMeshAt(mat, geo, fx, fy)
				cf.jsScene9Tile3D[tlt][[2]int{fx, fy}] = mesh
			}
			cf.scene.Call("add", mesh)

		case
			tile.Swamp,
			tile.Soil,
			tile.Stone,
			tile.Sand,
			tile.Sea,
			tile.Magma,
			tile.Ice,
			tile.Road,
			tile.Room:
			if oldTile.TestByTile(tlt) {
				continue // skip exist
			}
			mesh, exist := cf.jsScene9Tile3D[tlt][[2]int{fx, fy}]
			if !exist {
				mat := GetTextureTileMaterialByCache(tlt)
				geo := GetPlaneGeometryByCache(DstCellSize, DstCellSize)
				mesh = cf.make9InstancedMeshAt(mat, geo, fx, fy)
				cf.jsScene9Tile3D[tlt][[2]int{fx, fy}] = mesh
			}
			cf.scene.Call("add", mesh)

		case
			tile.Fog,
			tile.Smoke:
			if oldTile.TestByTile(tlt) {
				continue // skip exist
			}
			mesh, exist := cf.jsScene9Tile3D[tlt][[2]int{fx, fy}]
			if !exist {
				mat := GetTextureTileMaterialByCache(tlt)
				geo := GetBoxGeometryByCache(DstCellSize, DstCellSize, DstCellSize)
				mesh = cf.make9InstancedMeshAt(mat, geo, fx, fy)
				cf.jsScene9Tile3D[tlt][[2]int{fx, fy}] = mesh
			}
			cf.scene.Call("add", mesh)

		}
	}
}

func (cf *ClientFloorGL) make9InstancedMeshAt(
	mat, geo js.Value, fx, fy int) js.Value {
	w := cf.XWrapper.GetWidth()
	h := cf.YWrapper.GetWidth()
	geoXmin, geoXmax := CalcGeoMinMaxX(geo)
	geoYmin, geoYmax := CalcGeoMinMaxY(geo)
	geoZmin, geoZmax := CalcGeoMinMaxZ(geo)
	mesh := ThreeJsNew("InstancedMesh", geo, mat, 9)
	matrix := ThreeJsNew("Matrix4")
	for i := 0; i < way9type.Way9Type_Count; i++ {
		dx, dy := way9type.Way9Type(i).DxDy()
		x := fx + dx*w
		y := fy + dy*h
		matrix.Call("setPosition",
			ThreeJsNew("Vector3",
				float64(x)*DstCellSize+(geoXmax-geoXmin)/2,
				-float64(y)*DstCellSize-(geoYmax-geoYmin)/2,
				(geoZmax-geoZmin)/2,
			),
		)
		mesh.Call("setMatrixAt", i, matrix)
	}
	return mesh
}

func (cf *ClientFloorGL) add9TileAt(mat, geo js.Value, fx, fy int) []js.Value {
	w := cf.XWrapper.GetWidth()
	h := cf.YWrapper.GetWidth()
	rtn := make([]js.Value, 0, 9)
	geoXmin, geoXmax := CalcGeoMinMaxX(geo)
	geoYmin, geoYmax := CalcGeoMinMaxY(geo)
	geoZmin, geoZmax := CalcGeoMinMaxZ(geo)
	for i := 0; i < way9type.Way9Type_Count; i++ {
		dx, dy := way9type.Way9Type(i).DxDy()
		mesh := ThreeJsNew("Mesh", geo, mat)
		x := fx + dx*w
		y := fy + dy*h
		SetPosition(
			mesh,
			float64(x)*DstCellSize+(geoXmax-geoXmin)/2,
			-float64(y)*DstCellSize-(geoYmax-geoYmin)/2,
			(geoZmax-geoZmin)/2)
		cf.scene.Call("add", mesh)
		rtn = append(rtn, mesh)
	}
	return rtn
}

func (cf *ClientFloorGL) UpdateFrame(
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
	cameraZ := HelperSize - HelperSize*zoom/4
	SetPosition(cf.light,
		cameraX, cameraY, DstCellSize*10,
	)
	SetPosition(cf.camera,
		cameraX, cameraY, cameraZ,
	)
	cf.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			cameraX, cameraY, 0,
		),
	)
}

func calcShiftDxDy(frameProgress float64) (int, int) {
	rate := 1 - frameProgress
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	dx := int(float64(DstCellSize) * rate)
	dy := int(float64(DstCellSize) * rate)
	return dx, dy
}

func (cf *ClientFloorGL) addFieldObj(o *c2t_obj.FieldObjClient) {
	oldx, oldy, exist := cf.FieldObjPosMan.GetXYByUUID(o.ID)
	if exist && o.X == oldx && o.Y == oldy {
		return // no need to add
	}
	if exist { // handle obj move
		// something wrong, field obj do not move
		jslog.Errorf("fieldobj move? %v %v %v", o, oldx, oldy)
		return
	}
	// add new obj
	cf.FieldObjPosMan.AddToXY(o, o.X, o.Y)
	for dx := -1; dx < 2; dx++ {
		for dy := -1; dy < 2; dy++ {
			cf.addFieldObjAt(o,
				o.X+dx*cf.XWrapper.GetWidth(),
				o.Y+dy*cf.YWrapper.GetWidth(),
			)
		}
	}
}

func (cf *ClientFloorGL) addFieldObjAt(
	o *c2t_obj.FieldObjClient,
	fx, fy int,
) {
	tlList := gClientTile.FieldObjTiles[o.DisplayType]
	tilediff := fx*5 + fy*3
	if tilediff < 0 {
		tilediff = -tilediff
	}
	ti := tlList[tilediff%len(tlList)]

	mat := GetTileMaterialByCache(ti)
	// geo := GetSphereGeometryByCache(DstCellSize/2, DstCellSize, DstCellSize)
	geo := GetBoxGeometryByCache(DstCellSize-1, DstCellSize-1, DstCellSize-1)
	geoXmin, geoXmax := CalcGeoMinMaxX(geo)
	geoYmin, geoYmax := CalcGeoMinMaxY(geo)
	geoZmin, geoZmax := CalcGeoMinMaxZ(geo)
	mesh := ThreeJsNew("Mesh", geo, mat)
	SetPosition(
		mesh,
		float64(fx)*DstCellSize+(geoXmax-geoXmin)/2,
		-float64(fy)*DstCellSize-(geoYmax-geoYmin)/2,
		(geoZmax-geoZmin)/2)
	cf.scene.Call("add", mesh)
}

func (cf *ClientFloorGL) processNotiObjectList(
	olNoti *c2t_obj.NotiObjectList_data,
	vpx, vpy int, // place obj around
) {
	floorW := cf.XWrapper.GetWidth()
	floorH := cf.YWrapper.GetWidth()

	addUUID := make(map[string]bool)

	// make activeobj
	for _, ao := range olNoti.ActiveObjList {
		mesh, exist := cf.jsSceneObjs[ao.UUID]
		tlList := gClientTile.CharTiles[ao.Faction]
		ti := tlList[0]
		if !ao.Alive {
			ti = tlList[1]
		}
		mat := GetTileMaterialByCache(ti)
		if !exist {
			geo := GetBoxGeometryByCache(DstCellSize, DstCellSize, DstCellSize)
			mesh = ThreeJsNew("Mesh", geo, mat)
			cf.scene.Call("add", mesh)
			cf.jsSceneObjs[ao.UUID] = mesh
		}
		mesh.Set("material", mat)

		fx, fy := cf.calcAroundPos(floorW, floorH, vpx, vpy, ao.X, ao.Y)
		geo := mesh.Get("geometry")
		geoXmin, geoXmax := CalcGeoMinMaxX(geo)
		geoYmin, geoYmax := CalcGeoMinMaxY(geo)
		geoZmin, geoZmax := CalcGeoMinMaxZ(geo)
		SetPosition(
			mesh,
			float64(fx)*DstCellSize+(geoXmax-geoXmin)/2,
			-float64(fy)*DstCellSize-(geoYmax-geoYmin)/2,
			(geoZmax-geoZmin)/2)

		addUUID[ao.UUID] = true

		for _, eqo := range ao.EquippedPo {
			mesh, exist := cf.jsSceneObjs[eqo.UUID]
			if !exist {
				mesh = cf.makeEquipedMesh(eqo)
				cf.scene.Call("add", mesh)
				cf.jsSceneObjs[eqo.UUID] = mesh
			}

			fx, fy := cf.calcAroundPos(floorW, floorH, vpx, vpy, ao.X, ao.Y)
			shInfo := aoeqposShift[eqo.EquipType]
			geo := mesh.Get("geometry")
			geoXmin, geoXmax := CalcGeoMinMaxX(geo)
			geoYmin, geoYmax := CalcGeoMinMaxY(geo)
			geoZmin, geoZmax := CalcGeoMinMaxZ(geo)
			SetPosition(
				mesh,
				float64(fx)*DstCellSize+(geoXmax-geoXmin)/2+DstCellSize*shInfo.X,
				-float64(fy)*DstCellSize-(geoYmax-geoYmin)/2-DstCellSize*shInfo.Y,
				(geoZmax-geoZmin)/2+DstCellSize*shInfo.Z,
			)
			addUUID[eqo.UUID] = true
		}
	}

	// make carryobj
	for _, cro := range olNoti.CarryObjList {
		mesh, exist := cf.jsSceneObjs[cro.UUID]
		if !exist {
			mesh = cf.makeCarryObjMesh(cro)
			cf.scene.Call("add", mesh)
			cf.jsSceneObjs[cro.UUID] = mesh
		}

		fx, fy := cf.calcAroundPos(floorW, floorH, vpx, vpy, cro.X, cro.Y)
		shInfo := carryObjClientOnFloor2DrawInfo(cro)
		geo := mesh.Get("geometry")
		geoXmin, geoXmax := CalcGeoMinMaxX(geo)
		geoYmin, geoYmax := CalcGeoMinMaxY(geo)
		geoZmin, geoZmax := CalcGeoMinMaxZ(geo)
		SetPosition(
			mesh,
			float64(fx)*DstCellSize+(geoXmax-geoXmin)/2+DstCellSize*shInfo.X,
			-float64(fy)*DstCellSize-(geoYmax-geoYmin)/2-DstCellSize*shInfo.Y,
			(geoZmax-geoZmin)/2+DstCellSize*shInfo.Z,
		)

		addUUID[cro.UUID] = true
	}

	for id, mesh := range cf.jsSceneObjs {
		if !addUUID[id] {
			cf.scene.Call("remove", mesh)
			delete(cf.jsSceneObjs, id)
		}
	}
}

// make fx,fy around vpx, vpy
func (cf *ClientFloorGL) calcAroundPos(w, h, vpx, vpy, fx, fy int) (int, int) {
	if abs.Absi(fx-vpx) > w/2 {
		if fx > vpx {
			fx -= w
		} else {
			fx += w
		}
	}
	if abs.Absi(fy-vpy) > h/2 {
		if fy > vpy {
			fy -= h
		} else {
			fy += h
		}
	}
	return fx, fy
}

func (cf *ClientFloorGL) makeEquipedMesh(o *c2t_obj.EquipClient) js.Value {
	ti := gClientTile.EquipTiles[o.EquipType][o.Faction]
	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(
		DstCellSize/3, DstCellSize/3, DstCellSize/3,
	)
	return ThreeJsNew("Mesh", geo, mat)
}

func (cf *ClientFloorGL) makeCarryObjMesh(o *c2t_obj.CarryObjClientOnFloor) js.Value {
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
	mat := GetTileMaterialByCache(ti)
	geo := GetBoxGeometryByCache(
		DstCellSize/3, DstCellSize/3, DstCellSize/3,
	)
	return ThreeJsNew("Mesh", geo, mat)
}

func carryObjClientOnFloor2DrawInfo(
	co *c2t_obj.CarryObjClientOnFloor) coShift {
	switch co.CarryingObjectType {
	default:
		return coShiftOther[co.CarryingObjectType]
	case carryingobjecttype.Equip:
		return eqposShift[co.EquipType]
	}
}

type coShift struct {
	X float64
	Y float64
	Z float64
}

var aoeqposShift = [equipslottype.EquipSlotType_Count]coShift{
	equipslottype.Helmet: {-0.33, 0.0, 0.66},
	equipslottype.Amulet: {1.00, 0.0, 0.66},

	equipslottype.Weapon: {-0.33, 0.25, 0.66},
	equipslottype.Shield: {1.00, 0.25, 0.66},

	equipslottype.Ring:     {-0.33, 0.50, 0.66},
	equipslottype.Gauntlet: {1.00, 0.50, 0.66},

	equipslottype.Armor:    {-0.33, 0.75, 0.66},
	equipslottype.Footwear: {1.00, 0.75, 0.66},
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
