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
	"fmt"
	"math"
	"time"

	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

// process noti vptiles
// viewport x,y changed == need scroll
func (vp *GameScene) ProcessNotiVPTiles(
	cf *clientfloor.ClientFloor,
	taNoti *c2t_obj.NotiVPTiles_data,
	olNoti *c2t_obj.NotiObjectList_data,
	path2dst [][2]int,
) error {

	if cf.FloorInfo.Name != taNoti.FloorName {
		return fmt.Errorf("vptile data floor not match %v %v",
			cf.FloorInfo.Name, taNoti.FloorName)

	}
	vp.makeClientTile4PlayView(cf, taNoti)
	vp.updateFieldObjInView(cf, taNoti.VPX, taNoti.VPY)
	vp.makeMovePathInView(cf, taNoti.VPX, taNoti.VPY, path2dst)
	vp.raycastPlane.MoveCenterTo(taNoti.VPX, taNoti.VPY)
	return nil
}

// place obj around vpx, vpy
func (vp *GameScene) processNotiObjectList(
	cf *clientfloor.ClientFloor, olNoti *c2t_obj.NotiObjectList_data, vpx, vpy int) {

	floorW := cf.XWrapper.GetWidth()
	floorH := cf.YWrapper.GetWidth()

	addAOuuid := make(map[string]bool)
	addCOuuid := make(map[string]bool)
	playerUUID := gInitData.AccountInfo.ActiveObjUUID

	// make activeobj
	for _, ao := range olNoti.ActiveObjList {
		ao3d, exist := vp.jsSceneAOs[ao.UUID]
		if !exist {
			ao3d = NewActiveObj3D(ao)
			vp.scene.Call("add", ao3d.Mesh)
			vp.jsSceneAOs[ao.UUID] = ao3d
			vp.scene.Call("add", ao3d.Name.Mesh)
			for _, v := range ao3d.Condition {
				vp.scene.Call("add", v.Mesh)
			}
		}
		toadds, todels := ao3d.UpdateAOC(ao)
		for _, v := range toadds {
			vp.scene.Call("add", v)
		}
		for _, v := range todels {
			vp.scene.Call("remove", v)
		}

		fx, fy := CalcAroundPos(floorW, floorH, vpx, vpy, ao.X, ao.Y)
		tl := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		shZ := CalcTile3DStepOn(tl)
		ao3d.SetFieldPosition(fx, fy, 0, 0, shZ)

		addAOuuid[ao.UUID] = true

		if ao.UUID == playerUUID { // player ao
			aop := olNoti.ActiveObj
			if aop.Conditions.TestByCondition(condition.Invisible) {
				ao3d.Visible(false)
			} else {
				ao3d.Visible(true)
			}
			vp.UpdatePlayerAO(cf, ao, aop)
		}
		if !ao.Alive {
			// ao3d.RotateX(-math.Pi / 2)
			ao3d.ScaleX(0.5)
			ao3d.ScaleY(0.5)
			ao3d.ScaleZ(0.5)
		}
	}

	for id, ao3d := range vp.jsSceneAOs {
		if !addAOuuid[id] {
			for _, v := range ao3d.Condition {
				vp.scene.Call("remove", v.Mesh)
			}
			vp.scene.Call("remove", ao3d.Name.Mesh)
			vp.scene.Call("remove", ao3d.Mesh)
			if ao3d.Chat != nil {
				vp.scene.Call("remove", ao3d.Chat.Mesh)
				ao3d.Chat.Dispose()
				ao3d.Chat = nil
			}
			delete(vp.jsSceneAOs, id)
			ao3d.Dispose()
		}
	}

	// make carryobj
	for _, cro := range olNoti.CarryObjList {
		cr3d, exist := vp.jsSceneCOs[cro.UUID]
		if !exist {
			str, color := CarryObj2StrColor(cro)
			cr3d = NewCarryObj3D(str, color)
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
			cr3d.Dispose()
		}
	}

	// make dangerobj
	addDOuuid := make(map[string]bool)
	for _, dao := range olNoti.DangerObjList {
		dao3d, exist := vp.jsSceneDOs[dao.UUID]
		if !exist {
			dao3d = gPoolDangerObj3D.Get(dao)
			vp.scene.Call("add", dao3d.Mesh)
			vp.jsSceneDOs[dao.UUID] = dao3d
		}

		fx, fy := CalcAroundPos(floorW, floorH, vpx, vpy, dao.X, dao.Y)
		dao3d.SetFieldPosition(fx, fy, 0, 0, 0)
		addDOuuid[dao.UUID] = true
	}

	for id, dao3d := range vp.jsSceneDOs {
		if !addDOuuid[id] {
			vp.scene.Call("remove", dao3d.Mesh)
			delete(vp.jsSceneDOs, id)
			gPoolDangerObj3D.Put(dao3d)
		}
	}

}

// update hp,ap,sp bar movearrow for player ao
func (vp *GameScene) UpdatePlayerAO(
	cf *clientfloor.ClientFloor, ao *c2t_obj.ActiveObjClient, aop *c2t_obj.PlayerActiveObjInfo) {

	fx, fy := ao.X, ao.Y
	shZ := CalcTile3DStepOn(cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)])
	var spw, hpw float64
	apw := math.Sqrt(leveldata.CalcLevelFromExp(float64(aop.Exp))) + 1
	vp.AP.ScaleY(apw)
	vp.AP.ScaleZ(apw)
	if ao.Alive {
		_, hpw = vp.HP.SetWH(aop.HP, aop.HPMax)
		_, spw = vp.SP.SetWH(aop.SP, aop.SPMax)
		if aop.AP < 0 {
			vp.AP.ScaleX(-aop.AP)
		}
	} else {
		_, hpw = vp.HP.SetWH(0, aop.HPMax)
		_, spw = vp.SP.SetWH(0, aop.SPMax)
		vp.AP.ScaleX(0)
	}
	vp.HP.SetFieldPosition(fx, fy, 0, -(spw+apw+hpw)*2, DstCellSize+6+shZ)
	vp.AP.SetFieldPosition(fx, fy, 0, -(spw+apw)*2, DstCellSize+4+shZ)
	vp.SP.SetFieldPosition(fx, fy, 0, -(spw)*2, DstCellSize+2+shZ)
}

// playview frame update
func (vp *GameScene) UpdatePlayViewFrame(
	cf *clientfloor.ClientFloor,
	frameProgress float64,
	scrollDir way9type.Way9Type,
	taNoti *c2t_obj.NotiVPTiles_data,
	olNoti *c2t_obj.NotiObjectList_data,
	lastOLNoti *c2t_obj.NotiObjectList_data,
	envBias bias.Bias,
) {
	playerUUID := gInitData.AccountInfo.ActiveObjUUID

	// activeobj animate
	for i, ao := range olNoti.ActiveObjList {
		aod, exist := vp.jsSceneAOs[ao.UUID]
		if !exist {
			continue // ??
		}
		if !ao.Alive {
			continue
		}
		aod.ResetMatrix()
		if ao.UUID == playerUUID {
			// player
			if lastOLNoti.ActiveObj.AP < 0 {
				aod.RotateY(CalcRotateFrameProgress(frameProgress))
				vp.AP.ScaleX(-frameProgress)
			}
		}
		if ao.DamageTake > 0 {
			if i%2 == 0 {
				aod.ScaleX(CalcScaleFrameProgress(frameProgress, ao.DamageTake))
			} else {
				aod.ScaleY(CalcScaleFrameProgress(frameProgress, ao.DamageTake))
			}
		}
		// if ao.Act == c2t_idcmd.Move && ao.Result == c2t_error.None {
		// 	vp.animateActiveObj(cf, aod, ao.X, ao.Y, ao.Dir, frameProgress)
		// }
	}

	vp.animateFieldObj()
	vp.animateDangerObj(frameProgress)
	vp.animateTile(envBias)
	vp.moveCameraLight(
		cf, taNoti.VPX, taNoti.VPY,
		frameProgress, scrollDir,
		envBias,
	)

	// move cursor
	fx, fy := vp.mouseCursorFx, vp.mouseCursorFy
	tl := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
	vp.cursor.SetFieldPosition(fx, fy, tl)

	vp.renderer.Call("render", vp.scene, vp.camera)
}

// add tiles in gXYLenListView for playview
func (vp *GameScene) makeClientTile4PlayView(
	cf *clientfloor.ClientFloor,
	taNoti *c2t_obj.NotiVPTiles_data) {
	vpx, vpy := taNoti.VPX, taNoti.VPY
	for ti := 0; ti < tile.Tile_Count; ti++ {
		vp.jsTile3DCount[ti] = 0     // clear use count
		vp.jsTile3DDarkCount[ti] = 0 // clear use count
	}
	// matrix := ThreeJsNew("Matrix4")
	rad := time.Now().Sub(gInitData.TowerInfo.StartTime).Seconds()
	for vpi, v := range gXYLenListView {
		fx := v.X + vpx
		fy := v.Y + vpy
		newTile := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		dark := false
		if vpi >= len(taNoti.VPTiles) || taNoti.VPTiles[vpi] == 0 {
			dark = true
		}
		for ti := 0; ti < tile.Tile_Count; ti++ {
			if !newTile.TestByTile(tile.Tile(ti)) {
				continue
			}
			matrix := ThreeJsNew("Matrix4")
			if dark {
				if tile.Tile(ti) == tile.Door {
					matrix.Call("makeRotationZ", rad)
				}
				matrix.Call("setPosition",
					gTile3DDark[ti].MakePosVector3(fx, fy),
				)
				vp.jsTile3DDarkMesh[ti].Call("setMatrixAt",
					vp.jsTile3DDarkCount[ti], matrix)
				vp.jsTile3DDarkCount[ti]++
			} else {
				if tile.Tile(ti) == tile.Door {
					matrix.Call("makeRotationZ", rad)
				}
				matrix.Call("setPosition",
					gTile3D[ti].MakePosVector3(fx, fy),
				)
				vp.jsTile3DMesh[ti].Call("setMatrixAt",
					vp.jsTile3DCount[ti], matrix)
				vp.jsTile3DCount[ti]++
			}
		}
	}
	for ti := 0; ti < tile.Tile_Count; ti++ {
		vp.jsTile3DMesh[ti].Set("count", vp.jsTile3DCount[ti])
		vp.jsTile3DMesh[ti].Get("instanceMatrix").Set("needsUpdate", true)
		vp.jsTile3DDarkMesh[ti].Set("count", vp.jsTile3DDarkCount[ti])
		vp.jsTile3DDarkMesh[ti].Get("instanceMatrix").Set("needsUpdate", true)
	}
}
