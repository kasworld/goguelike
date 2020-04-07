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

package viewport2d

import (
	"fmt"
	"math"

	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/way9type"

	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/wasmclient/clientfloor"
	"github.com/kasworld/goguelike/game/wasmclient/webtilegroup"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (vp *Viewport2d) DrawPlayVP(
	frameProgress float64, envBias, towerBias, floorBias bias.Bias,
	accountG2ID g2id.G2ID,
	cf *clientfloor.ClientFloor,
	scrollDir way9type.Way9Type,
	actMoveDir way9type.Way9Type,
	taNoti *c2t_obj.NotiVPTiles_data,
	olNoti *c2t_obj.NotiObjectList_data,
	lastOLNoti *c2t_obj.NotiObjectList_data,
	Path2dst [][2]int,
) error {

	defer vp.drawNoti()

	sx, sy := vp.calcShiftDxDy(frameProgress)
	scrollDx := scrollDir.Dx() * sx
	scrollDy := scrollDir.Dy() * sy

	vp.context2d.Set("fillStyle", "black")
	vp.context2d.Call("fillRect", 0, 0, vp.ViewWidth, vp.ViewHeight)

	vp.drawFloorTiles(envBias, cf, scrollDx, scrollDy, taNoti.VPX, taNoti.VPY)
	if vp.Grid {
		vp.drawGrid(scrollDx, scrollDy, towerBias, floorBias, frameProgress)
	}
	vp.drawFloorFieldObjs(cf, scrollDx, scrollDy, taNoti.VPX, taNoti.VPY)
	vp.drawShadowSight(cf, scrollDx, scrollDy, taNoti)

	err := vp.drawCarryObjList(cf, scrollDx, scrollDy, olNoti.CarryObjList, taNoti.VPX, taNoti.VPY)
	if err != nil {
		return err
	}

	vp.context2d.Set("textAlign", "center")
	for _, ao := range olNoti.ActiveObjList {
		vpx, vpy, ok := vp.floorXY2canvasXY(
			cf.FloorInfo.W, cf.FloorInfo.H,
			taNoti.VPX, taNoti.VPY, ao.X, ao.Y)
		if !ok {
			continue
		}
		if ao.G2ID == accountG2ID {
			vp.drawUserActiveObj(
				frameProgress, vpx, vpy, ao,
				olNoti.ActiveObj, lastOLNoti.ActiveObj,
				actMoveDir)
		} else {
			vp.drawActiveObj(frameProgress, vpx, vpy, sx, sy, ao, scrollDir)
		}
	}

	vp.drawMovePath(cf, scrollDx, scrollDy, Path2dst, taNoti)
	vp.drawMouseCursorViewport(cf, scrollDx, scrollDy, taNoti)
	return nil
}

func (vp *Viewport2d) drawActiveObj(
	frameProgress float64,
	vpx, vpy, sx, sy int,
	ao *c2t_obj.ActiveObjClient,
	scrollDir way9type.Way9Type) {

	aoDir := way9type.Center
	if ao.Act == c2t_idcmd.Move && ao.Result == c2t_error.None {
		aoDir = ao.Dir
	}
	scrollDx := (scrollDir.Dx() - aoDir.Dx()) * sx
	scrollDy := (scrollDir.Dy() - aoDir.Dy()) * sy
	dstRect := vp.calcDstRect(vpx, vpy, scrollDx, scrollDy)

	if vp.DrawEquip {
		vp.drawActiveObjAt(frameProgress, dstRect.ScaleCenterRate(0.7, 0.7), ao)
		for _, po := range ao.EquippedPo {
			vp.drawEquipItemsAt(dstRect, po.EquipType, po.Faction)
		}
	} else {
		vp.drawActiveObjAt(frameProgress, dstRect, ao)
	}

	vp.drawNickOrChat(ao, dstRect)

	if ao.Act == c2t_idcmd.Attack && ao.Dir != way9type.Center {
		vp.drawMovingArrowDir(frameProgress, dstRect, ao.Dir, ao.Faction)
	}
	if ao.DamageTake > 0 {
		fontSize := float64(vp.CellSize/4) * (math.Log10(float64(ao.DamageTake))/2 + 1)
		vp.context2d.Set("font", fmt.Sprintf("%dpx sans-serif", int(fontSize)))
		vp.context2d.Set("fillStyle", "Red")
		vp.context2d.Call("fillText", -ao.DamageTake,
			dstRect.CenterX(), dstRect.Y)
	}
}

func (vp *Viewport2d) drawUserActiveObj(
	frameProgress float64,
	vpx, vpy int,
	ao *c2t_obj.ActiveObjClient,
	plao *c2t_obj.PlayerActiveObjInfo,
	lastPlao *c2t_obj.PlayerActiveObjInfo,
	actMoveDir way9type.Way9Type,
) {

	dstRect := vp.calcDstRect(vpx, vpy, 0, 0)
	if actMoveDir != way9type.Center {
		vp.drawCenterArrowDir(dstRect, actMoveDir)
	}

	vp.drawBars(dstRect, plao.HP, plao.HPMax, plao.SP, plao.SPMax,
		lastPlao.RemainTurn2Act, frameProgress)

	// cannot see self if invisible
	if !plao.Conditions.TestByCondition(condition.Invisible) {
		// drawScale := 1.5
		if vp.DrawEquip {
			vp.drawActiveObjAt(frameProgress, dstRect.ScaleCenterRate(1.0, 1.0), ao)
			for _, po := range ao.EquippedPo {
				vp.drawEquipItemsAt(dstRect.ScaleCenterRate(1.5, 1.5), po.EquipType, po.Faction)
			}
		} else {
			vp.drawActiveObjAt(frameProgress, dstRect.ScaleCenterRate(1.0, 1.0), ao)
		}
	}

	vp.drawNickOrChat(ao, dstRect)

	if ao.Dir != way9type.Center {
		vp.drawMovingArrowDir(frameProgress, dstRect, ao.Dir, ao.Faction)
	}
	vp.drawActiveObjMessage(dstRect)
}

func (vp *Viewport2d) drawCarryObjList(cf *clientfloor.ClientFloor, scrollDx, scrollDy int,
	poList []*c2t_obj.CarryObjClientOnFloor,
	vpx, vpy int) error {

	for _, po := range poList {
		cnvX, cnvY, in := vp.floorXY2canvasXY(
			cf.FloorInfo.W, cf.FloorInfo.H,
			vpx, vpy, po.X, po.Y)
		if !in {
			continue
		}
		var ti webtilegroup.TileInfo
		dstRect := vp.calcDstRect(cnvX, cnvY, scrollDx, scrollDy)
		switch po.CarryingObjectType {
		default:
			return fmt.Errorf("unknonw po %v", po)
		case carryingobjecttype.Equip:
			vp.drawEquipItemsAt(dstRect, po.EquipType, po.Faction)
			continue
		case carryingobjecttype.Money:
			var find bool
			for i, v := range moneycolor.Attrib {
				if po.Value < v.UpLimit {
					ti = vp.clientTile.GoldTiles[i]
					find = true
					break
				}
			}
			if !find {
				ti = vp.clientTile.GoldTiles[len(vp.clientTile.GoldTiles)-1]
			}
		case carryingobjecttype.Potion:
			ti = vp.clientTile.PotionTiles[po.PotionType]
		case carryingobjecttype.Scroll:
			ti = vp.clientTile.ScrollTiles[po.ScrollType]
		}
		srcRect := ti.Rect
		coshift := coShiftOther[po.CarryingObjectType]
		dstRect2 := dstRect
		dstRect2.X += int(float64(dstRect.W) * coshift.X)
		dstRect2.Y += int(float64(dstRect.H) * coshift.Y)
		dstRect2.W = int(float64(dstRect.W) * coshift.W)
		dstRect2.H = int(float64(dstRect.H) * coshift.H)

		vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
			srcRect.X, srcRect.Y, srcRect.W, srcRect.H,
			dstRect2.X, dstRect2.Y, dstRect2.W, dstRect2.H)
	}
	return nil
}

func (vp *Viewport2d) drawMovePath(
	cf *clientfloor.ClientFloor, scrollDx, scrollDy int,
	Path2dst [][2]int,
	taNoti *c2t_obj.NotiVPTiles_data) {

	if len(Path2dst) == 0 {
		return
	}
	w, h := cf.XWrapper.GetWidth(), cf.YWrapper.GetWidth()
	for i, v := range Path2dst[:len(Path2dst)-1] {
		if x, y, isIn := vp.floorXY2canvasXY(
			cf.FloorInfo.W, cf.FloorInfo.H,
			taNoti.VPX, taNoti.VPY, v[0], v[1]); isIn {

			v2 := Path2dst[i+1]
			dx, dy := way9type.CalcDxDyWrapped(
				v2[0]-v[0], v2[1]-v[1],
				w, h,
			)
			diri := way9type.RemoteDxDy2Way9(dx, dy)
			ti := vp.clientTile.Dir2Tiles[diri]

			dstRect := vp.calcDstRect(x, y, scrollDx, scrollDy)
			srcRect := ti.Rect
			vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
				srcRect.X, srcRect.Y, srcRect.W, srcRect.H,
				dstRect.X, dstRect.Y, dstRect.W, dstRect.H)

		}
	}
	if len(Path2dst) >= 1 {
		v := Path2dst[len(Path2dst)-1]
		if x, y, isIn := vp.floorXY2canvasXY(
			cf.FloorInfo.W, cf.FloorInfo.H,
			taNoti.VPX, taNoti.VPY, v[0], v[1]); isIn {
			ti := vp.clientTile.DirTiles[0]
			dstRect := vp.calcDstRect(x, y, scrollDx, scrollDy)
			srcRect := ti.Rect
			vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
				srcRect.X, srcRect.Y, srcRect.W, srcRect.H,
				dstRect.X, dstRect.Y, dstRect.W, dstRect.H)
		}
	}
}

func (vp *Viewport2d) drawMouseCursorViewport(
	cf *clientfloor.ClientFloor, scrollDx, scrollDy int,
	taNoti *c2t_obj.NotiVPTiles_data) {

	flX, flY := vp.CanvasXY2FloorXY(cf.XWrapSafe, cf.YWrapSafe,
		taNoti.VPX, taNoti.VPY, vp.MouseX/vp.CellSize, vp.MouseY/vp.CellSize)
	canvasX, canvasY, isIn := vp.floorXY2canvasXY(
		cf.FloorInfo.W, cf.FloorInfo.H, taNoti.VPX, taNoti.VPY, flX, flY)
	if !isIn {
		return
	}
	tl := cf.Tiles[flX][flY]
	ti := vp.clientTile.CursorTiles[1]
	if !tl.CharPlaceable() {
		ti = vp.clientTile.CursorTiles[2]
	} else if tl.Safe() {
		ti = vp.clientTile.CursorTiles[0]
	}
	dstRect := vp.calcDstRect(canvasX, canvasY, scrollDx, scrollDy)
	srcRect := ti.Rect
	vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
		srcRect.X, srcRect.Y, srcRect.W, srcRect.H,
		dstRect.X, dstRect.Y, dstRect.W, dstRect.H)
}

func (vp *Viewport2d) drawShadowSight(cf *clientfloor.ClientFloor,
	scrollDx, scrollDy int,
	taNoti *c2t_obj.NotiVPTiles_data) {

	tlic := vp.DarkerTileImgCnv
	for _, v := range vp.CellIndexPos {
		pos := [2]int{v.X - vp.CellXCount/2, v.Y - vp.CellYCount/2}
		index, exist := vp.ViewportPos2Index[pos]
		if !exist || taNoti.VPTiles[index] == 0 {
			vp.context2d.Call("drawImage", tlic.Cnv,
				0, 0, vp.CellSize, vp.CellSize,
				v.dstX+scrollDx, v.dstY+scrollDy, vp.CellSize, vp.CellSize)
		}
	}
}
