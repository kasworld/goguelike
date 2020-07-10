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
)

func (cf *ClientFloorGL) Zoom(zoom int) {
	cf.camera.Set("zoom", 1.0+float64(zoom)/2)
	cf.camera.Call("updateProjectionMatrix")
}

func (cf *ClientFloorGL) UpdateFrame(
	frameProgress float64,
	scrollDir way9type.Way9Type,
	taNoti *c2t_obj.NotiVPTiles_data,
	envBias bias.Bias) {

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

// add fo to clientview by cf.FieldObjPosMan
func (cf *ClientFloorGL) updateFieldObjInView(vpx, vpy int) {
	addFOuuid := make(map[string]bool)
	for _, v := range gXYLenListView {
		fx := v.X + vpx
		fy := v.Y + vpy
		obj := cf.FieldObjPosMan.Get1stObjAt(fx, fy) // wrap in fn
		if obj == nil {
			continue
		}
		fo := obj.(*c2t_obj.FieldObjClient)
		fo3d, exist := cf.jsSceneFOs[obj.GetUUID()]
		if !exist {
			// add new fieldobj
			fo3d = NewFieldObj3D()
			ti := FieldObj2TileInfo(fo.DisplayType, fo.X, fo.Y)
			fo3d.ChangeTile(ti)
			cf.jsSceneFOs[obj.GetUUID()] = fo3d
			cf.scene.Call("add", fo3d.Mesh)
		}
		addFOuuid[obj.GetUUID()] = true
		geoInfo := fo3d.GeoInfo
		SetPosition(
			fo3d.Mesh,
			float64(fx)*DstCellSize+geoInfo.Len[0]/2,
			-float64(fy)*DstCellSize-geoInfo.Len[1]/2,
			geoInfo.Len[2]/2,
		)
	}

	// del removed obj
	for id, fo3d := range cf.jsSceneFOs {
		if !addFOuuid[id] {
			cf.scene.Call("remove", fo3d.Mesh)
			fo3d.Dispose()
			delete(cf.jsSceneFOs, id)
		}
	}
}

func (cf *ClientFloorGL) processNotiObjectList(
	olNoti *c2t_obj.NotiObjectList_data,
	vpx, vpy int, // place obj around
) {
	floorW := cf.XWrapper.GetWidth()
	floorH := cf.YWrapper.GetWidth()

	addAOuuid := make(map[string]bool)
	addCOuuid := make(map[string]bool)

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
		addAOuuid[ao.UUID] = true

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
			addCOuuid[eqo.UUID] = true
		}
	}

	for id, ao3d := range cf.jsSceneAOs {
		if !addAOuuid[id] {
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

		addCOuuid[cro.UUID] = true
	}

	for id, mesh := range cf.jsSceneCOs {
		if !addCOuuid[id] {
			cf.scene.Call("remove", mesh)
			delete(cf.jsSceneCOs, id)
		}
	}
}
