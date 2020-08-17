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
	"time"

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

	if cf.FloorInfo.UUID != taNoti.FloorUUID {
		return fmt.Errorf("vptile data floor not match %v %v",
			cf.FloorInfo.UUID, taNoti.FloorUUID)

	}
	vp.makeClientTile4PlayView(cf, taNoti)
	vp.updateFieldObjInView(cf, taNoti.VPX, taNoti.VPY)
	vp.makeMovePathInView(cf, taNoti.VPX, taNoti.VPY, path2dst)
	vp.raycastPlane.MoveCenterTo(taNoti.VPX, taNoti.VPY)
	return nil
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
