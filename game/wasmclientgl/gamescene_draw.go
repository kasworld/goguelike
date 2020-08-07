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

	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

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
			if lastOLNoti.ActiveObj.RemainTurn2Act > 0 {
				aod.RotateY(CalcRotateFrameProgress(frameProgress))
				vp.AP.ScaleX(frameProgress)
			}
		}
		if ao.DamageTake > 0 {
			if i%2 == 0 {
				aod.ScaleX(CalcScaleFrameProgress(frameProgress, ao.DamageTake))
			} else {
				aod.ScaleY(CalcScaleFrameProgress(frameProgress, ao.DamageTake))
			}
		}
		vp.animateMoveArrow(cf, aod, ao.X, ao.Y, ao.Dir, frameProgress)
	}

	vp.animateFieldObj()
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

// animate move arrow
func (vp *GameScene) animateMoveArrow(
	cf *clientfloor.ClientFloor, ao3d *ActiveObj3D,
	fx, fy int, dir way9type.Way9Type, frameProgress float64) {

	if dir != way9type.Center {
		dx, dy := dir.DxDy()
		tl := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		shX := DstCellSize * frameProgress * float64(dx)
		shY := DstCellSize * frameProgress * float64(dy)
		shZ := GetTile3DOnByCache(tl)
		ao3d.MoveArrow.SetFieldPosition(fx, fy, shX, shY, shZ)
	}
}

// floorview frame update
func (vp *GameScene) UpdateFloorViewFrame(
	cf *clientfloor.ClientFloor, vpx, vpy int, envBias bias.Bias) {

	vp.makeClientTile4FloorView(cf, vpx, vpy)
	vp.updateFieldObjInView(cf, vpx, vpy)
	vp.animateFieldObj()
	vp.animateTile(envBias)
	vp.moveCameraLight(
		cf, vpx, vpy,
		0, 0,
		envBias,
	)
	vp.renderer.Call("render", vp.scene, vp.camera)
}

// fieldobj animate
// common to playview, floorview
func (vp *GameScene) animateFieldObj() {
	rad := time.Now().Sub(gInitData.TowerInfo.StartTime).Seconds()
	for _, fo := range vp.jsSceneFOs {
		fo.RotateZ(rad)
	}
}

// tile scroll list of animate
// common to playview, floorview
func (vp *GameScene) animateTile(envBias bias.Bias) {
	for _, i := range []tile.Tile{
		tile.Swamp,
		tile.Soil,
		tile.Stone,
		tile.Sand,
		tile.Sea,
		tile.Magma,
		tile.Ice,
		tile.Grass,
		tile.Tree,
		// tile.Road,
		// tile.Room,
		// tile.Wall,
		tile.Window,
		tile.Door,
		tile.Fog,
		tile.Smoke,
	} {
		tilc := tile.TileScrollAttrib[i]
		shX := int(envBias[tilc.SrcXBiasAxis] * tilc.AniXSpeed)
		shY := int(envBias[tilc.SrcYBiasAxis] * tilc.AniYSpeed)
		gTile3D[i].DrawTexture(shX, shY)
		gTile3DDark[i].DrawTexture(shX, shY)
	}
}

// common to playview, floorview
func (vp *GameScene) moveCameraLight(
	cf *clientfloor.ClientFloor, vpx, vpy int, frameProgress float64,
	scrollDir way9type.Way9Type, envBias bias.Bias) {

	sx, sy := CalcShiftDxDy(frameProgress)
	scrollDx := -scrollDir.Dx() * sx
	scrollDy := scrollDir.Dy() * sy

	// move camera, light
	cameraX := float64(vpx*DstCellSize + scrollDx)
	cameraY := float64(-vpy*DstCellSize + scrollDy)
	cameraR := float64(HelperSize)

	carmeraDy := float64(
		(2-gameOptions.GetByIDBase("Zoom").State)*gameOptions.GetByIDBase("Angle").State) * DstCellSize

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
		cameraX,
		cameraY-cameraR*math.Cos(cameraRad)-carmeraDy,
		cameraR*math.Sin(cameraRad),
	)
	vp.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			cameraX,
			cameraY-carmeraDy,
			0,
		),
	)
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

// make floor tiles for floorview
func (vp *GameScene) makeClientTile4FloorView(
	cf *clientfloor.ClientFloor, vpx, vpy int) {
	for i := 0; i < tile.Tile_Count; i++ {
		vp.jsTile3DCount[i] = 0     // clear use count
		vp.jsTile3DDarkCount[i] = 0 // clear use count
	}
	matrix := ThreeJsNew("Matrix4")
	for _, v := range gXYLenListView {
		fx := v.X + vpx
		fy := v.Y + vpy
		newTile := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		for ti := 0; ti < tile.Tile_Count; ti++ {
			if !newTile.TestByTile(tile.Tile(ti)) {
				continue
			}
			matrix.Call("setPosition",
				gTile3D[ti].MakePosVector3(fx, fy),
			)
			vp.jsTile3DMesh[ti].Call("setMatrixAt",
				vp.jsTile3DCount[ti], matrix)
			vp.jsTile3DCount[ti]++
		}
	}
	for ti := 0; ti < tile.Tile_Count; ti++ {
		vp.jsTile3DMesh[ti].Set("count", vp.jsTile3DCount[ti])
		vp.jsTile3DMesh[ti].Get("instanceMatrix").Set("needsUpdate", true)
		vp.jsTile3DDarkMesh[ti].Set("count", vp.jsTile3DDarkCount[ti])
		vp.jsTile3DDarkMesh[ti].Get("instanceMatrix").Set("needsUpdate", true)
	}
}

// add fo to clientview by vp.FieldObjPosMan
func (vp *GameScene) updateFieldObjInView(
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

			lb3d := gPoolLabel3D.Get(fo.Message)
			fo3d.Label = lb3d
			vp.scene.Call("add", lb3d.Mesh)
		}
		if !addFOuuid[obj.GetUUID()] {
			tl := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
			shZ := GetTile3DOnByCache(tl)
			fo3d.SetFieldPosition(fx, fy, shZ)
			addFOuuid[obj.GetUUID()] = true
			fo3d.Label.SetFieldPosition(fx, fy,
				0, 0, DstCellSize+2+shZ)
		} else {
			// same fo in gXYLenListView
			// field too small
			// skip farther fo
		}
	}

	// del removed obj
	for id, fo3d := range vp.jsSceneFOs {
		if !addFOuuid[id] {
			vp.scene.Call("remove", fo3d.Mesh)
			gPoolFieldObj3D.Put(fo3d)
			delete(vp.jsSceneFOs, id)
			lb3d := fo3d.Label
			fo3d.Label = nil
			vp.scene.Call("remove", lb3d.Mesh)
			gPoolLabel3D.Put(lb3d)
		}
	}
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
			ao3d = NewActiveObj3D(ao.Faction, ao.NickName)
			vp.scene.Call("add", ao3d.Mesh)
			vp.jsSceneAOs[ao.UUID] = ao3d
			vp.scene.Call("add", ao3d.Name.Mesh)
			vp.scene.Call("add", ao3d.MoveArrow.Mesh)
		}
		if oldmesh, changed := ao3d.ChangeFaction(ao.Faction); changed {
			vp.scene.Call("remove", oldmesh)
			vp.scene.Call("add", ao3d.Mesh)
		}

		fx, fy := CalcAroundPos(floorW, floorH, vpx, vpy, ao.X, ao.Y)
		tl := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		shZ := GetTile3DOnByCache(tl)
		if ao.Conditions.TestByCondition(condition.Float) {
			ao3d.SetFieldPosition(fx, fy, shZ+DstCellSize)
		} else {
			ao3d.SetFieldPosition(fx, fy, shZ)
		}
		if ao.Dir != way9type.Center {
			ao3d.MoveArrow.Visible(true)
			ao3d.MoveArrow.SetDir(ao.Dir)
		} else {
			ao3d.MoveArrow.Visible(false)
		}

		addAOuuid[ao.UUID] = true
		if len(ao.Chat) == 0 {
			if ao3d.Chat != nil {
				vp.scene.Call("remove", ao3d.Chat.Mesh)
				ao3d.Chat.Dispose()
				ao3d.Chat = nil
			}
		} else {
			if ao3d.Chat == nil {
				// add new chat
				ao3d.Chat = NewLabel3D(ao.Chat)
				vp.scene.Call("add", ao3d.Chat.Mesh)
			} else {
				if ao.Chat != ao3d.Chat.Str {
					// remove old chat , add new chat
					vp.scene.Call("remove", ao3d.Chat.Mesh)
					ao3d.Chat.Dispose()
					ao3d.Chat = NewLabel3D(ao.Chat)
					vp.scene.Call("add", ao3d.Chat.Mesh)
				}
			}
		}
		if ao3d.Chat != nil {
			ao3d.Chat.SetFieldPosition(fx, fy,
				0, -DstCellSize/2, DstCellSize+2+shZ)
		}

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
			vp.scene.Call("remove", ao3d.Name.Mesh)
			vp.scene.Call("remove", ao3d.MoveArrow.Mesh)
			vp.scene.Call("remove", ao3d.Mesh)
			delete(vp.jsSceneAOs, id)
			ao3d.Dispose()
			if ao3d.Chat != nil {
				vp.scene.Call("remove", ao3d.Chat.Mesh)
				ao3d.Chat.Dispose()
				ao3d.Chat = nil
			}
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

// update hp,ap,sp bar movearrow for player ao
func (vp *GameScene) UpdatePlayerAO(
	cf *clientfloor.ClientFloor, ao *c2t_obj.ActiveObjClient, aop *c2t_obj.PlayerActiveObjInfo) {

	fx, fy := ao.X, ao.Y
	shZ := GetTile3DOnByCache(cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)])
	var spw, hpw float64
	apw := math.Sqrt(leveldata.CalcLevelFromExp(float64(aop.Exp))) + 1
	vp.AP.ScaleY(apw)
	vp.AP.ScaleZ(apw)
	if ao.Alive {
		_, hpw = vp.HP.SetWH(aop.HP, aop.HPMax)
		_, spw = vp.SP.SetWH(aop.SP, aop.SPMax)
		if aop.RemainTurn2Act > 0 {
		} else {
			vp.AP.ScaleX(-aop.RemainTurn2Act)
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

func (vp *GameScene) ClearMovePath() {
	for pos, ar3d := range vp.jsSceneMovePathArrows {
		vp.scene.Call("remove", ar3d.Mesh)
		gPoolColorArrow3D.Put(ar3d)
		delete(vp.jsSceneMovePathArrows, pos)
	}
}

func (vp *GameScene) makeMovePathInView(
	cf *clientfloor.ClientFloor,
	vpx, vpy int,
	path2dst [][2]int) {

	addAr3Duuid := make(map[[2]int]bool)
	if len(path2dst) > 0 {
		w, h := cf.XWrapper.GetWidth(), cf.YWrapper.GetWidth()
		for i, pos := range path2dst[:len(path2dst)-1] {
			ar3d, exist := vp.jsSceneMovePathArrows[pos]
			if !exist {
				ar3d = gPoolColorArrow3D.Get("#808080")
				vp.jsSceneMovePathArrows[pos] = ar3d
				vp.scene.Call("add", ar3d.Mesh)
			}
			addAr3Duuid[pos] = true

			// change tile for dir
			pos2 := path2dst[i+1]
			dx, dy := way9type.CalcDxDyWrapped(
				pos2[0]-pos[0], pos2[1]-pos[1],
				w, h,
			)
			diri := way9type.RemoteDxDy2Way9(dx, dy)
			ar3d.SetDir(diri)
			x, y := CalcAroundPos(w, h, vpx, vpy, pos[0], pos[1])
			tl := cf.Tiles[cf.XWrapSafe(x)][cf.YWrapSafe(y)]
			shZ := GetTile3DOnByCache(tl)
			ar3d.SetFieldPosition(x, y, 0, 0, shZ)
		}
		// add last
		// pos := path2dst[len(path2dst)-1]
		// ar3d, exist := vp.jsSceneMovePathArrows[pos]
		// if !exist {
		// 	ar3d = gPoolColorArrow3D.Get()
		// 	vp.jsSceneMovePathArrows[pos] = ar3d
		// 	vp.scene.Call("add", ar3d.Mesh)
		// }
		// addAr3Duuid[pos] = true
		// ti := gClientTile.Dir2Tiles[way9type.Center]
		// ar3d.SetDir(diri)
		// tl := cf.Tiles[cf.XWrapSafe(pos[0])][cf.YWrapSafe(pos[1])]
		// x, y := CalcAroundPos(w, h, vpx, vpy, pos[0], pos[1])
		// ar3d.SetFieldPosition(x, y, tl)
	}

	for pos, ar3d := range vp.jsSceneMovePathArrows {
		if !addAr3Duuid[pos] {
			vp.scene.Call("remove", ar3d.Mesh)
			gPoolColorArrow3D.Put(ar3d)
			delete(vp.jsSceneMovePathArrows, pos)
		}
	}
}
