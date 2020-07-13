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
	"time"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (vp *Viewport) UpdateFrame(
	cf *clientfloor.ClientFloor,
	frameProgress float64,
	scrollDir way9type.Way9Type,
	taNoti *c2t_obj.NotiVPTiles_data,
	envBias bias.Bias,
) {

	sx, sy := CalcShiftDxDy(frameProgress)
	scrollDx := -scrollDir.Dx() * sx
	scrollDy := scrollDir.Dy() * sy

	rad := time.Now().Sub(gInitData.TowerInfo.StartTime).Seconds()
	for _, fo := range vp.jsSceneFOs {
		fo.RotateZ(rad)
	}

	// move camera, light
	cameraX := float64(taNoti.VPX*DstCellSize + scrollDx)
	cameraY := float64(-taNoti.VPY*DstCellSize + scrollDy)
	cameraR := float64(HelperSize)

	cameraRad := []float64{
		math.Pi / 2,
		math.Pi / 3,
		math.Pi / 4,
	}[gameOptions.GetByIDBase("Angle").State]

	envBias = envBias.MakeAbsSumTo(1)
	cx := float64(cf.XWrapper.GetWidth()) * DstCellSize / 2
	cy := float64(cf.YWrapper.GetWidth()) * DstCellSize / 2
	r := cx
	if cy < cx {
		r = cy
	}
	r *= 0.7
	for i := range vp.light {
		rad := envBias[i] * 2 * math.Pi
		x := cx + r*math.Sin(rad)
		y := cy + r*math.Cos(rad)
		SetPosition(vp.light[i],
			x, -y, DstCellSize*8,
		)
	}
	SetPosition(vp.lightW,
		cameraX, cameraY, cameraR,
	)

	SetPosition(vp.camera,
		cameraX, cameraY-cameraR*math.Cos(cameraRad), cameraR*math.Sin(cameraRad),
	)
	vp.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			cameraX, cameraY, 0,
		),
	)
}

// add tiles in gXYLenListView
func (vp *Viewport) makeClientTileView(
	cf *clientfloor.ClientFloor, vpx, vpy int) {
	for i := 0; i < tile.Tile_Count; i++ {
		vp.jsInstacedCount[i] = 0 // clear use count
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
				vp.jsInstacedMesh[i].Call("setMatrixAt", vp.jsInstacedCount[i], matrix)
				vp.jsInstacedCount[i]++
			}
		}
	}
	for i := 0; i < tile.Tile_Count; i++ {
		vp.jsInstacedMesh[i].Set("count", vp.jsInstacedCount[i])
		vp.jsInstacedMesh[i].Get("instanceMatrix").Set("needsUpdate", true)
	}
}

// add fo to clientview by vp.FieldObjPosMan
func (vp *Viewport) updateFieldObjInView(
	cf *clientfloor.ClientFloor, vpx, vpy int) {
	addFOuuid := make(map[string]bool)
	for _, v := range gXYLenListView {
		fx := v.X + vpx
		fy := v.Y + vpy
		obj := cf.FieldObjPosMan.Get1stObjAt(fx, fy) // wrap in fn
		if obj == nil {
			continue
		}
		fo := obj.(*c2t_obj.FieldObjClient)
		fo3d, exist := vp.jsSceneFOs[obj.GetUUID()]
		if !exist {
			// add new fieldobj
			fo3d = gPoolFieldObj3D.Get()
			ti := FieldObj2TileInfo(fo.DisplayType, fo.X, fo.Y)
			fo3d.ChangeTile(ti)
			vp.jsSceneFOs[obj.GetUUID()] = fo3d
			vp.scene.Call("add", fo3d.Mesh)
		}
		addFOuuid[obj.GetUUID()] = true
		fo3d.SetFieldPosition(fx, fy)
	}

	// del removed obj
	for id, fo3d := range vp.jsSceneFOs {
		if !addFOuuid[id] {
			vp.scene.Call("remove", fo3d.Mesh)
			gPoolFieldObj3D.Put(fo3d)
			delete(vp.jsSceneFOs, id)
		}
	}
}

func (vp *Viewport) processNotiObjectList(
	cf *clientfloor.ClientFloor,
	olNoti *c2t_obj.NotiObjectList_data,
	vpx, vpy int, // place obj around
) {
	floorW := cf.XWrapper.GetWidth()
	floorH := cf.YWrapper.GetWidth()

	addAOuuid := make(map[string]bool)
	addCOuuid := make(map[string]bool)

	// make activeobj
	for _, ao := range olNoti.ActiveObjList {
		ao3d, exist := vp.jsSceneAOs[ao.UUID]
		if !exist {
			ao3d = gPoolActiveObj3D.Get()
			vp.scene.Call("add", ao3d.Mesh)
			vp.jsSceneAOs[ao.UUID] = ao3d
		}
		tlList := gClientTile.CharTiles[ao.Faction]
		if ao.Alive {
			ao3d.ChangeTile(tlList[0])
		} else {
			ao3d.ChangeTile(tlList[1])
		}
		fx, fy := CalcAroundPos(floorW, floorH, vpx, vpy, ao.X, ao.Y)
		ao3d.SetFieldPosition(fx, fy)
		addAOuuid[ao.UUID] = true

		for _, eqo := range ao.EquippedPo {
			cr3d, exist := vp.jsSceneCOs[eqo.UUID]
			if !exist {
				cr3d = gPoolCarryObj3D.Get()
				ti := Equiped2TileInfo(eqo)
				cr3d.ChangeTile(ti)
				vp.scene.Call("add", cr3d.Mesh)
				vp.jsSceneCOs[eqo.UUID] = cr3d
			}
			shInfo := aoEqPosShift[eqo.EquipType]
			cr3d.SetFieldPosition(fx, fy, shInfo)
			addCOuuid[eqo.UUID] = true
		}
	}

	for id, ao3d := range vp.jsSceneAOs {
		if !addAOuuid[id] {
			vp.scene.Call("remove", ao3d.Mesh)
			delete(vp.jsSceneAOs, id)
			gPoolActiveObj3D.Put(ao3d)
		}
	}

	// make carryobj
	for _, cro := range olNoti.CarryObjList {
		cr3d, exist := vp.jsSceneCOs[cro.UUID]
		if !exist {
			cr3d = gPoolCarryObj3D.Get()
			ti := CarryObj2TileInfo(cro)
			cr3d.ChangeTile(ti)
			vp.scene.Call("add", cr3d.Mesh)
			vp.jsSceneCOs[cro.UUID] = cr3d
		}

		fx, fy := CalcAroundPos(floorW, floorH, vpx, vpy, cro.X, cro.Y)
		shInfo := CarryObjClientOnFloor2DrawInfo(cro)
		cr3d.SetFieldPosition(fx, fy, shInfo)
		addCOuuid[cro.UUID] = true
	}

	for id, cr3d := range vp.jsSceneCOs {
		if !addCOuuid[id] {
			vp.scene.Call("remove", cr3d.Mesh)
			delete(vp.jsSceneCOs, id)
			gPoolCarryObj3D.Put(cr3d)
		}
	}
}
