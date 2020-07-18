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

// 0~1 -> 1->1.5->1->0.5->1
func CalcScaleFrameProgress(frameProgress float64, damage int) float64 {

	amplitude := math.Log10(float64(damage)) / 2
	if amplitude < 0.1 {
		amplitude = 0.1
	}
	if amplitude > 2 {
		amplitude = 2
	}

	// 0~1 -> 1->2->1->0->1
	progress := math.Sin(frameProgress * math.Pi * 2)

	rtn := 1 + progress*amplitude
	if rtn < 0 {
		rtn = -rtn
	}
	return rtn
}

// 0~1 -> 0-> -pi -> 0 -> +pi
func CalcRotateFrameProgress(frameProgress float64) float64 {
	return math.Sin(frameProgress*math.Pi*2) * math.Pi / 4
}

func (vp *Viewport) UpdateFrame(
	cf *clientfloor.ClientFloor,
	frameProgress float64,
	scrollDir way9type.Way9Type,
	taNoti *c2t_obj.NotiVPTiles_data,
	olNoti *c2t_obj.NotiObjectList_data,
	lastOLNoti *c2t_obj.NotiObjectList_data,
	envBias bias.Bias,
) {
	playerUUID := gInitData.AccountInfo.ActiveObjUUID
	sx, sy := CalcShiftDxDy(frameProgress)
	scrollDx := -scrollDir.Dx() * sx
	scrollDy := scrollDir.Dy() * sy

	// fieldobj animate
	rad := time.Now().Sub(gInitData.TowerInfo.StartTime).Seconds()
	for _, fo := range vp.jsSceneFOs {
		fo.RotateZ(rad)
	}

	// activeobj animate
	for i, ao := range olNoti.ActiveObjList {
		aod, exist := vp.jsSceneAOs[ao.UUID]
		if !exist {
			continue // ??
		}
		aod.ResetMatrix()
		if ao.UUID == playerUUID {
			// player
			if lastOLNoti.ActiveObj.RemainTurn2Act > 0 {
				aod.RotateY(CalcRotateFrameProgress(frameProgress))
			}
		}
		if ao.DamageTake > 0 {
			if i%2 == 0 {
				aod.ScaleX(CalcScaleFrameProgress(frameProgress, ao.DamageTake))
			} else {
				aod.ScaleY(CalcScaleFrameProgress(frameProgress, ao.DamageTake))
			}
		}
	}

	// tile scroll list of animate
	for _, i := range []tile.Tile{
		tile.Swamp,
		tile.Soil,
		tile.Stone,
		tile.Sand,
		tile.Sea,
		tile.Magma,
		tile.Ice,
		tile.Grass,
		// tile.Tree,
		// tile.Road,
		// tile.Room,
		// tile.Wall,
		// tile.Window,
		// tile.Door,
		tile.Fog,
		tile.Smoke,
	} {
		tilc := tile.TileScrollAttrib[i]
		shX := int(envBias[tilc.SrcXBiasAxis] * tilc.AniXSpeed)
		shY := int(envBias[tilc.SrcYBiasAxis] * tilc.AniYSpeed)
		gTile3D[i].DrawTexture(tile.Tile(i), shX, shY)
		gTile3DDark[i].DrawTexture(tile.Tile(i), shX, shY)
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
	cf *clientfloor.ClientFloor,
	taNoti *c2t_obj.NotiVPTiles_data) {
	vpx, vpy := taNoti.VPX, taNoti.VPY
	for i := 0; i < tile.Tile_Count; i++ {
		vp.jsTile3DCount[i] = 0     // clear use count
		vp.jsTile3DDarkCount[i] = 0 // clear use count
	}
	matrix := ThreeJsNew("Matrix4")
	for i, v := range gXYLenListView {
		fx := v.X + vpx
		fy := v.Y + vpy
		newTile := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		dark := false
		if i >= len(taNoti.VPTiles) || taNoti.VPTiles[i] == 0 {
			dark = true
		}
		for i := 0; i < tile.Tile_Count; i++ {
			if !newTile.TestByTile(tile.Tile(i)) {
				continue
			}
			if dark {
				geolen := gTile3DDark[i].GeoInfo.Len
				sh := gTile3DDark[i].Shift
				matrix.Call("setPosition",
					ThreeJsNew("Vector3",
						sh[0]+float64(fx)*DstCellSize+geolen[0]/2,
						-sh[1]+-float64(fy)*DstCellSize-geolen[1]/2,
						sh[2]+geolen[2]/2,
					),
				)
				vp.jsTile3DDarkMesh[i].Call("setMatrixAt",
					vp.jsTile3DDarkCount[i], matrix)
				vp.jsTile3DDarkCount[i]++
			} else {
				geolen := gTile3D[i].GeoInfo.Len
				sh := gTile3D[i].Shift
				matrix.Call("setPosition",
					ThreeJsNew("Vector3",
						sh[0]+float64(fx)*DstCellSize+geolen[0]/2,
						-sh[1]+-float64(fy)*DstCellSize-geolen[1]/2,
						sh[2]+geolen[2]/2,
					),
				)
				vp.jsTile3DMesh[i].Call("setMatrixAt",
					vp.jsTile3DCount[i], matrix)
				vp.jsTile3DCount[i]++
			}
		}
	}
	for i := 0; i < tile.Tile_Count; i++ {
		vp.jsTile3DMesh[i].Set("count", vp.jsTile3DCount[i])
		vp.jsTile3DMesh[i].Get("instanceMatrix").Set("needsUpdate", true)
		vp.jsTile3DDarkMesh[i].Set("count", vp.jsTile3DDarkCount[i])
		vp.jsTile3DDarkMesh[i].Get("instanceMatrix").Set("needsUpdate", true)
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
		if !addFOuuid[obj.GetUUID()] {
			fo3d.SetFieldPosition(fx, fy)
			addFOuuid[obj.GetUUID()] = true
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
