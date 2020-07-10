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
	"math"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

func (cf *ClientFloorGL) UpdateFrame(
	frameProgress float64,
	scrollDir way9type.Way9Type,
	taNoti *c2t_obj.NotiVPTiles_data,
	envBias bias.Bias) {

	zoom := gameOptions.GetByIDBase("Zoom").State
	sx, sy := CalcShiftDxDy(frameProgress)
	scrollDx := -scrollDir.Dx() * sx
	scrollDy := scrollDir.Dy() * sy

	// move camera, light
	cameraX := taNoti.VPX*DstCellSize + scrollDx
	cameraY := -taNoti.VPY*DstCellSize + scrollDy
	cameraZ := HelperSize

	envBias = envBias.MakeAbsSumTo(1)
	cx := float64(cf.XWrapper.GetWidth()) * DstCellSize / 2
	cy := float64(cf.YWrapper.GetWidth()) * DstCellSize / 2
	r := cx
	if cy < cx {
		r = cy
	}
	r *= 0.7
	for i := range cf.light {
		rad := envBias[i] * 2 * math.Pi
		x := cx + r*math.Sin(rad)
		y := cy + r*math.Cos(rad)
		SetPosition(cf.light[i],
			x, -y, DstCellSize*8,
		)
	}
	SetPosition(cf.lightW,
		cameraX, cameraY-cameraZ/2, cameraZ,
	)

	SetPosition(cf.camera,
		cameraX, cameraY-cameraZ/2, cameraZ,
	)
	cf.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			cameraX, cameraY, 0,
		),
	)
	cf.camera.Set("zoom", 1.0+float64(zoom)/2)
	cf.camera.Call("updateProjectionMatrix")
}

// add tiles in gXYLenListView
func (cf *ClientFloorGL) makeClientTileView(vpx, vpy int) {
	for i := 0; i < tile.Tile_Count; i++ {
		cf.jsInstacedCount[i] = 0 // clear use count
	}
	matrix := ThreeJsNew("Matrix4")
	for _, v := range gXYLenListView {
		fx := v.X + vpx
		fy := v.Y + vpy
		newTile := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		for i := 0; i < tile.Tile_Count; i++ {
			if newTile.TestByTile(tile.Tile(i)) {
				geolen := gTile3D[i].GeoInfo.Len
				sh := gTile3D[i].Shift
				matrix.Call("setPosition",
					ThreeJsNew("Vector3",
						sh[0]+float64(fx)*DstCellSize+geolen[0]/2,
						-sh[1]+-float64(fy)*DstCellSize-geolen[1]/2,
						sh[2]+geolen[2]/2,
					),
				)
				cf.jsInstacedMesh[i].Call("setMatrixAt", cf.jsInstacedCount[i], matrix)
				cf.jsInstacedCount[i]++
			}
		}
	}
	for i := 0; i < tile.Tile_Count; i++ {
		cf.jsInstacedMesh[i].Set("count", cf.jsInstacedCount[i])
		cf.jsInstacedMesh[i].Get("instanceMatrix").Set("needsUpdate", true)
	}
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
	mat, geo := MakeFieldObjMatGeo(o, o.X, o.Y)
	geoInfo := GetGeoInfo(geo)
	mesh := ThreeJsNew("InstancedMesh", geo, mat, 9)
	cf.scene.Call("add", mesh)
	matrix := ThreeJsNew("Matrix4")
	for i := 0; i < way9type.Way9Type_Count; i++ {
		dx, dy := way9type.Way9Type(i).DxDy()
		fx := o.X + dx*cf.XWrapper.GetWidth()
		fy := o.Y + dy*cf.YWrapper.GetWidth()
		matrix.Call("setPosition",
			ThreeJsNew("Vector3",
				float64(fx)*DstCellSize+geoInfo.Len[0]/2,
				-float64(fy)*DstCellSize-geoInfo.Len[1]/2,
				geoInfo.Len[2]/2,
			),
		)
		mesh.Call("setMatrixAt", i, matrix)
	}
}

func (cf *ClientFloorGL) processNotiObjectList(
	olNoti *c2t_obj.NotiObjectList_data,
	vpx, vpy int, // place obj around
) {
	floorW := cf.XWrapper.GetWidth()
	floorH := cf.YWrapper.GetWidth()

	addAOUUID := make(map[string]bool)
	addCOUUID := make(map[string]bool)

	// make activeobj
	for _, ao := range olNoti.ActiveObjList {
		ao3d, exist := cf.jsSceneAOs[ao.UUID]
		if !exist {
			ao3d = NewActiveObj3D()
			cf.scene.Call("add", ao3d.Mesh)
			cf.jsSceneAOs[ao.UUID] = ao3d
		}
		tlList := gClientTile.CharTiles[ao.Faction]
		if ao.Alive {
			ao3d.ChangeTile(tlList[0])
		} else {
			ao3d.ChangeTile(tlList[1])
		}
		geoInfo := ao3d.GeoInfo
		fx, fy := CalcAroundPos(floorW, floorH, vpx, vpy, ao.X, ao.Y)
		SetPosition(
			ao3d.Mesh,
			float64(fx)*DstCellSize+geoInfo.Len[0]/2,
			-float64(fy)*DstCellSize-geoInfo.Len[1]/2,
			geoInfo.Len[2]/2)
		addAOUUID[ao.UUID] = true

		for _, eqo := range ao.EquippedPo {
			mesh, exist := cf.jsSceneCOs[eqo.UUID]
			if !exist {
				mesh = MakeEquipedMesh(eqo)
				cf.scene.Call("add", mesh)
				cf.jsSceneCOs[eqo.UUID] = mesh
			}
			shInfo := aoEqPosShift[eqo.EquipType]
			geo := mesh.Get("geometry")
			geoInfo := GetGeoInfo(geo)
			SetPosition(
				mesh,
				float64(fx)*DstCellSize+geoInfo.Len[0]/2+DstCellSize*shInfo.X,
				-float64(fy)*DstCellSize-geoInfo.Len[1]/2-DstCellSize*shInfo.Y,
				geoInfo.Len[2]/2+DstCellSize*shInfo.Z,
			)
			addCOUUID[eqo.UUID] = true
		}
	}

	for id, ao3d := range cf.jsSceneAOs {
		if !addAOUUID[id] {
			cf.scene.Call("remove", ao3d.Mesh)
			delete(cf.jsSceneAOs, id)
			ao3d.Dispose()
		}
	}

	// make carryobj
	for _, cro := range olNoti.CarryObjList {
		mesh, exist := cf.jsSceneCOs[cro.UUID]
		if !exist {
			mesh = MakeCarryObjMesh(cro)
			cf.scene.Call("add", mesh)
			cf.jsSceneCOs[cro.UUID] = mesh
		}

		fx, fy := CalcAroundPos(floorW, floorH, vpx, vpy, cro.X, cro.Y)
		shInfo := CarryObjClientOnFloor2DrawInfo(cro)
		geo := mesh.Get("geometry")
		geoInfo := GetGeoInfo(geo)
		SetPosition(
			mesh,
			float64(fx)*DstCellSize+geoInfo.Len[0]/2+DstCellSize*shInfo.X,
			-float64(fy)*DstCellSize-geoInfo.Len[1]/2-DstCellSize*shInfo.Y,
			geoInfo.Len[2]/2+DstCellSize*shInfo.Z,
		)

		addCOUUID[cro.UUID] = true
	}

	for id, mesh := range cf.jsSceneCOs {
		if !addCOUUID[id] {
			cf.scene.Call("remove", mesh)
			delete(cf.jsSceneCOs, id)
		}
	}
}
