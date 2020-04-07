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
	"syscall/js"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/wasmclient/clientfloor"
	"github.com/kasworld/goguelike/game/wasmclient/webtilegroup"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/wrapper"
)

func (vp *Viewport2d) Hide() {
	vp.Canvas.Get("style").Set("display", "none")
}
func (vp *Viewport2d) Show() {
	vp.Canvas.Get("style").Set("display", "initial")
	vp.calcViewCellValue()
}

// resize canvas
func (vp *Viewport2d) Resize() {
	vp.calcViewCellValue()
}

func (vp *Viewport2d) Focus() {
	vp.Canvas.Call("focus")
}

func (vp *Viewport2d) Zoom(state int) {
	vp.zoomState = state
	// vp.calcViewCellValue()
}

func (vp *Viewport2d) AddEventListener(evt string, fn func(this js.Value, args []js.Value) interface{}) {
	vp.Canvas.Call("addEventListener", evt, js.FuncOf(fn))
}

// resize with zoom state
func (vp *Viewport2d) calcViewCellValue() {
	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Int()
	winH := win.Get("innerHeight").Int()

	// make cell square rect , to smaller
	if winW/gameconst.ViewPortW > winH/gameconst.ViewPortH {
		vp.CellSize = winH / (gameconst.ViewPortH)
	} else {
		vp.CellSize = winW / (gameconst.ViewPortW)
	}
	vp.RefSize = vp.CellSize
	if vp.RefSize < 48 {
		vp.RefSize = 48
	}
	for i := 0; i < vp.zoomState; i++ {
		vp.CellSize = vp.CellSize * 2 / 3
	}
	if vp.CellSize < 8 {
		vp.CellSize = 8
	}
	vp.CellXCount = winW / vp.CellSize
	vp.CellYCount = winH / vp.CellSize
	vp.ViewWidth = vp.CellSize * vp.CellXCount
	vp.ViewHeight = vp.CellSize * vp.CellYCount
	vp.Canvas.Call("setAttribute", "width", vp.ViewWidth)
	vp.Canvas.Call("setAttribute", "height", vp.ViewHeight)

	vp.CellIndexPos = make([]CellVPCarryObjS, (vp.CellXCount+2)*(vp.CellYCount+2))
	i := 0
	for x := -1; x < vp.CellXCount+1; x++ {
		for y := -1; y < vp.CellYCount+1; y++ {
			dstX, dstY := vp.calcDstPos(x, y, 0, 0)
			cv := CellVPCarryObjS{x, y, dstX, dstY}
			vp.CellIndexPos[i] = cv
			i++
		}
	}
	for i, v := range tile.TileScrollAttrib {
		if v.Texture {
			wrapInfo[i].Xcount = vp.TileImgCnvList[i].W / vp.CellSize
			wrapInfo[i].Ycount = vp.TileImgCnvList[i].H / vp.CellSize
			wrapInfo[i].WrapX = wrapper.New(vp.TileImgCnvList[i].W - vp.CellSize).WrapSafe
			wrapInfo[i].WrapY = wrapper.New(vp.TileImgCnvList[i].H - vp.CellSize).WrapSafe
		}
	}
}

func IsInViewRange(floorLen, floorVal, viewportStart, viewportLen int) (int, bool) {
	viewportVal := floorVal - viewportStart
	if IsIn(viewportVal, 0, viewportLen) {
		return viewportVal, true
	}
	if viewportVal < 0 {
		if IsIn(viewportVal+floorLen, 0, viewportLen) {
			return viewportVal + floorLen, true
		}
	}
	if viewportVal >= viewportLen {
		if IsIn(viewportVal-floorLen, 0, viewportLen) {
			return viewportVal - floorLen, true
		}
	}
	return 0, false
}

func IsIn(x int, r1, r2 int) bool {
	return x >= r1 && x < r2
}

func (vp *Viewport2d) calcDstRect(vpx, vpy int, dx, dy int) webtilegroup.Rect {
	return webtilegroup.Rect{
		vpx*vp.CellSize + dx,
		vpy*vp.CellSize + dy,
		vp.CellSize,
		vp.CellSize}
}

func (vp *Viewport2d) calcDstPos(vpx, vpy int, dx, dy int) (int, int) {
	return vpx*vp.CellSize + dx, vpy*vp.CellSize + dy
}

func (vp *Viewport2d) CanvasXY2FloorXY(xwrap, ywrap func(int) int, vpCenterX, vpCenterY, x, y int) (int, int) {
	return xwrap(x + vpCenterX - vp.CellXCount/2), ywrap(y + vpCenterY - vp.CellYCount/2)
}

func (vp *Viewport2d) floorXY2canvasXY(
	w, h, vpCenterX, vpCenterY, x, y int) (int, int, bool) {
	xv, xok := IsInViewRange(w, x, vpCenterX-vp.CellXCount/2, vp.CellXCount)
	if !xok {
		return 0, 0, false
	}
	yv, yok := IsInViewRange(h, y, vpCenterY-vp.CellYCount/2, vp.CellYCount)
	return xv, yv, yok
}

func (vp *Viewport2d) calcShiftDxDy(frameProgress float64) (int, int) {
	rate := 1 - frameProgress
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	dx := int(float64(vp.CellSize) * rate)
	dy := int(float64(vp.CellSize) * rate)
	return dx, dy
}

func (vp *Viewport2d) drawNickOrChat(ao *c2t_obj.ActiveObjClient, dstRect webtilegroup.Rect) {
	vp.context2d.Set("font", fmt.Sprintf("%dpx sans-serif", vp.CellSize/3))
	if ao.Chat != "" {
		vp.context2d.Set("fillStyle", "gray")
		vp.context2d.Call("fillText", ao.Chat, dstRect.CenterX(), dstRect.Y2()+vp.CellSize/4)
	} else {
		vp.context2d.Set("fillStyle", ao.Faction.Color24().String())
		vp.context2d.Call("fillText", ao.NickName, dstRect.CenterX(), dstRect.Y2()+vp.CellSize/4)
	}
}

func calcBar_WH(dstRect webtilegroup.Rect, pt, ptmax int) (int, int) {
	if pt < 0 {
		pt = 0
	}
	ptLen := pt * dstRect.W / ptmax
	ptWidth := ptmax/gameconst.ActiveObjBaseBiasLen + 1
	if ptWidth > dstRect.H/3 {
		ptWidth = dstRect.H / 3
	}
	if ptWidth < 2 {
		ptWidth = 2
	}
	return ptWidth, ptLen
}

func (vp *Viewport2d) drawBars(
	dstRect webtilegroup.Rect,
	hp, hpmax int,
	sp, spmax int,
	remainturn2act, frameProgress float64,
) {
	tpW := dstRect.H / 20
	if tpW < 2 {
		tpW = 2
	}

	hpW, hpLen := calcBar_WH(dstRect, hp, hpmax)
	vp.context2d.Set("fillStyle", "red")
	vp.context2d.Call("fillRect", dstRect.X, dstRect.Y+1, hpLen, hpW)

	spW, spLen := calcBar_WH(dstRect, sp, spmax)
	vp.context2d.Set("fillStyle", "yellow")
	vp.context2d.Call("fillRect", dstRect.X, dstRect.Y-spW-tpW-1, spLen, spW)

	if remainturn2act > 0 {
		tpLen := int(float64(dstRect.W) * frameProgress)
		vp.context2d.Set("fillStyle", "Lime")
		vp.context2d.Call("fillRect", dstRect.X, dstRect.Y-tpW, tpLen, tpW)
	} else {
		tpLen := int(float64(dstRect.W) * -remainturn2act)
		vp.context2d.Set("fillStyle", "LimeGreen")
		vp.context2d.Call("fillRect", dstRect.X2()-tpLen, dstRect.Y-tpW, tpLen, tpW)
	}
}

func (vp *Viewport2d) drawGrid(
	scrollDx, scrollDy int, towerBias, floorBias bias.Bias, frameProgress float64) {
	lineWidth := vp.CellSize / 32
	if lineWidth < 1 {
		lineWidth = 1
	}
	vp.context2d.Set("globalAlpha", 0.3)
	vp.context2d.Call("beginPath")
	vp.context2d.Set("lineWidth", lineWidth)

	vp.context2d.Set("strokeStyle", towerBias.ToHTMLColorString())
	for x := -vp.CellSize + scrollDx; x < vp.ViewWidth; x += vp.CellSize {
		vp.context2d.Call("moveTo", x, 0)
		vp.context2d.Call("lineTo", x, vp.ViewHeight)
	}
	vp.context2d.Call("stroke")

	vp.context2d.Call("beginPath")
	vp.context2d.Set("lineWidth", lineWidth)
	vp.context2d.Set("strokeStyle", floorBias.ToHTMLColorString())
	for y := -vp.CellSize + scrollDy; y < vp.ViewHeight; y += vp.CellSize {
		vp.context2d.Call("moveTo", 0, y)
		vp.context2d.Call("lineTo", vp.ViewWidth, y)
	}
	vp.context2d.Call("stroke")
	vp.context2d.Set("globalAlpha", 1)
}

func (vp *Viewport2d) drawActiveObjMessage(dstRect webtilegroup.Rect) {
	vp.ActiveObjMessage = vp.ActiveObjMessage.Compact()
	lineH := float64(vp.CellSize) * 2 / 8
	n := 0
	posy := dstRect.Y - dstRect.H/6
	for i := len(vp.ActiveObjMessage) - 1; i >= 0; i-- {
		v := vp.ActiveObjMessage[i]
		if v.IsEnded() {
			continue
		}
		fontH := int(lineH * v.SizeRate)
		vp.context2d.Set("font", fmt.Sprintf("%dpx sans-serif", fontH))
		posx := dstRect.X + int(float64(dstRect.W)*0.5)
		vp.context2d.Set("globalAlpha", 1-v.ProgressRate())
		vp.context2d.Set("fillStyle", v.Color)
		vp.context2d.Call("fillText", v.Text, posx, posy)
		posy += -fontH
		n++
	}
	vp.context2d.Set("globalAlpha", 1)
}

func (vp *Viewport2d) drawNoti() {
	vp.NotiMessage = vp.NotiMessage.Compact()
	lineH := float64(vp.RefSize) / 1.5
	n := 0
	posy := vp.ViewHeight/2 - 2*vp.RefSize
	for i := len(vp.NotiMessage) - 1; i >= 0; i-- {
		v := vp.NotiMessage[i]
		if v.IsEnded() {
			continue
		}
		fontH := int(lineH * v.SizeRate)
		vp.context2d.Set("font", fmt.Sprintf("%dpx sans-serif", fontH))
		posx := vp.ViewWidth / 2
		posy += -fontH
		vp.context2d.Set("globalAlpha", 1-v.ProgressRate())
		vp.context2d.Set("fillStyle", v.Color)
		vp.context2d.Call("fillText", v.Text, posx, posy)
		n++
	}
	vp.context2d.Set("globalAlpha", 1)
}

func (vp *Viewport2d) drawFloorTiles(
	envBias bias.Bias,
	cf *clientfloor.ClientFloor,
	scrollDx, scrollDy int, vpx, vpy int) {
	canvasStartx := vpx - vp.CellXCount/2
	canvasStarty := vpy - vp.CellYCount/2

	for _, v := range vp.CellIndexPos {
		fx := cf.XWrapSafe(v.X + canvasStartx)
		fy := cf.YWrapSafe(v.Y + canvasStarty)
		tl := cf.Tiles[fx][fy]
		diffbase := fx*5 + fy*3
		for i := 0; i < tile.Tile_Count; i++ {
			shX := 0.0
			shY := 0.0
			if vp.TileScroll {
				tilc := tile.TileScrollAttrib[i]
				shX = envBias[tilc.SrcXBiasAxis] * tilc.AniXSpeed
				shY = envBias[tilc.SrcYBiasAxis] * tilc.AniYSpeed
			}
			tlt := tile.Tile(i)
			if tl.TestByTile(tlt) {
				if vp.TileImgCnvList[i] != nil {
					// texture tile
					tlic := vp.TileImgCnvList[i]
					wrap := wrapInfo[i]
					tx := fx % wrap.Xcount
					ty := fy % wrap.Ycount
					srcx := wrap.WrapX(tx*vp.CellSize + int(shX))
					srcy := wrap.WrapY(ty*vp.CellSize + int(shY))
					vp.context2d.Call("drawImage", tlic.Cnv,
						srcx, srcy, vp.CellSize, vp.CellSize,
						v.dstX+scrollDx, v.dstY+scrollDy, vp.CellSize, vp.CellSize)

				} else if tlt == tile.Wall {
					// wall tile process
					tlList := vp.clientTile.FloorTiles[i]
					tilediff := vp.calcWallTileDiff(cf, fx, fy)
					ti := tlList[tilediff%len(tlList)]
					vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
						ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
						v.dstX+scrollDx, v.dstY+scrollDy, vp.CellSize, vp.CellSize)
				} else if tlt == tile.Window {
					// window tile process
					tlList := vp.clientTile.FloorTiles[i]
					tlindex := 0
					if vp.checkWallAt(cf, fx, fy-1) && vp.checkWallAt(cf, fx, fy+1) { // n-s window
						tlindex = 1
					}
					ti := tlList[tlindex]
					vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
						ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
						v.dstX+scrollDx, v.dstY+scrollDy, vp.CellSize, vp.CellSize)
				} else {
					// bitmap tile
					tlList := vp.clientTile.FloorTiles[i]
					tilediff := diffbase + int(shX)
					ti := tlList[tilediff%len(tlList)]
					vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
						ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
						v.dstX+scrollDx, v.dstY+scrollDy, vp.CellSize, vp.CellSize)
				}
			}
		}
	}
}

func (vp *Viewport2d) calcWallTileDiff(cf *clientfloor.ClientFloor, flx, fly int) int {
	rtn := 0
	if vp.checkWallAt(cf, flx, fly-1) {
		rtn |= 1
	}
	if vp.checkWallAt(cf, flx+1, fly) {
		rtn |= 1 << 1
	}
	if vp.checkWallAt(cf, flx, fly+1) {
		rtn |= 1 << 2
	}
	if vp.checkWallAt(cf, flx-1, fly) {
		rtn |= 1 << 3
	}
	return rtn
}

func (vp *Viewport2d) checkWallAt(cf *clientfloor.ClientFloor, flx, fly int) bool {
	flx = cf.XWrapSafe(flx)
	fly = cf.YWrapSafe(fly)
	tl := cf.Tiles[flx][fly]
	return tl.TestByTile(tile.Wall) ||
		tl.TestByTile(tile.Door) ||
		tl.TestByTile(tile.Window)
}

func (vp *Viewport2d) drawFloorFieldObjs(cf *clientfloor.ClientFloor,
	scrollDx, scrollDy int, vpx, vpy int) {
	vp.context2d.Set("font", fmt.Sprintf("%dpx sans-serif", vp.CellSize/3))
	allFieldObj := cf.FieldObjPosMan.GetAllList()
	for _, v := range allFieldObj {
		iao, ok := v.(*c2t_obj.FieldObjClient)
		if !ok {
			continue
		}
		x, y, ok := vp.floorXY2canvasXY(
			cf.FloorInfo.W, cf.FloorInfo.H,
			vpx, vpy, iao.X, iao.Y)
		if !ok {
			continue
		}
		dstRect := vp.calcDstRect(x, y, scrollDx, scrollDy)
		tilediff := iao.X*5 + iao.Y*3
		tlList := vp.clientTile.FieldObjTiles[iao.DisplayType]
		if len(tlList) == 0 {
			jslog.Errorf("len=0 %v", iao.DisplayType)
			continue
		}
		ti := tlList[tilediff%len(tlList)]
		vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
			ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
			dstRect.X, dstRect.Y, dstRect.W, dstRect.H)
		vp.context2d.Set("fillStyle", iao.ActType.Color24().ToHTMLColorString())
		vp.context2d.Call("fillText", iao.Message,
			dstRect.CenterX(), dstRect.CenterY()+vp.CellSize/2)
	}
}

func (vp *Viewport2d) drawEquipItemsAt(
	dstRect webtilegroup.Rect,
	EquipType equipslottype.EquipSlotType,
	Faction factiontype.FactionType,
) {
	ti := vp.clientTile.EquipTiles[EquipType][Faction]
	srcRect := ti.Rect
	dstRect2 := dstRect
	dstRect2.X += int(float64(dstRect.W) * eqposShift[EquipType].X)
	dstRect2.Y += int(float64(dstRect.H) * eqposShift[EquipType].Y)
	dstRect2.W = int(float64(dstRect.W) * eqposShift[EquipType].W)
	dstRect2.H = int(float64(dstRect.H) * eqposShift[EquipType].H)

	vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
		srcRect.X, srcRect.Y, srcRect.W, srcRect.H,
		dstRect2.X, dstRect2.Y, dstRect2.W, dstRect2.H)
}

type coShift struct {
	X float64
	Y float64
	W float64
	H float64
}

var eqposShift = [equipslottype.EquipSlotType_Count]coShift{
	equipslottype.Helmet: {0.0, 0.0, 0.25, 0.25},
	equipslottype.Amulet: {0.75, 0.0, 0.25, 0.25},

	equipslottype.Weapon: {0.0, 0.25, 0.25, 0.25},
	equipslottype.Shield: {0.75, 0.25, 0.25, 0.25},

	equipslottype.Ring:     {0.0, 0.50, 0.25, 0.25},
	equipslottype.Gauntlet: {0.75, 0.50, 0.25, 0.25},

	equipslottype.Armor:    {0.0, 0.75, 0.25, 0.25},
	equipslottype.Footwear: {0.75, 0.75, 0.25, 0.25},
}

var coShiftOther = [carryingobjecttype.CarryingObjectType_Count]coShift{
	carryingobjecttype.Money:  {0.33, 0.0, 0.33, 0.33},
	carryingobjecttype.Potion: {0.33, 0.33, 0.33, 0.33},
	carryingobjecttype.Scroll: {0.33, 0.66, 0.33, 0.33},
}

func (vp *Viewport2d) drawActiveObjAt(
	frameProgress float64,
	dstRect webtilegroup.Rect,
	ao *c2t_obj.ActiveObjClient,
) {
	tlList := vp.clientTile.CharTiles[ao.Faction]
	ti := tlList[0]
	if !ao.Alive {
		ti = tlList[1]
	}
	if vp.DamageEffect && ao.DamageTake > 1 {
		ti = tlList[int(frameProgress*2*float64(len(tlList)))%len(tlList)]
	}
	srcRect := ti.Rect
	vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
		srcRect.X, srcRect.Y, srcRect.W, srcRect.H,
		dstRect.X, dstRect.Y, dstRect.W, dstRect.H)
}

func (vp *Viewport2d) drawCenterArrowDir(dstRect webtilegroup.Rect, dir way9type.Way9Type) {
	ti := vp.clientTile.DirTiles[dir]
	dstRect = webtilegroup.Rect{
		X: dstRect.X + dstRect.W/2*dir.Dx(),
		Y: dstRect.Y + dstRect.H/2*dir.Dy(),
		W: dstRect.W,
		H: dstRect.H,
	}
	srcRect := ti.Rect
	vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
		srcRect.X, srcRect.Y, srcRect.W, srcRect.H,
		dstRect.X, dstRect.Y, dstRect.W, dstRect.H)
}

func (vp *Viewport2d) drawMovingArrowDir(
	frameProgress float64,
	dstRect webtilegroup.Rect,
	atkDir way9type.Way9Type,
	ft factiontype.FactionType,
) {
	ti := vp.clientTile.DirTiles[atkDir]
	srcRect := ti.Rect
	dx, dy := atkDir.TurnDir(2).DxDy()
	dx *= vp.CellSize / 12
	dy *= vp.CellSize / 12
	dx += atkDir.Dx() * int(frameProgress*float64(vp.CellSize))
	dy += atkDir.Dy() * int(frameProgress*float64(vp.CellSize))
	vp.context2d.Call("drawImage", vp.clientTile.TilePNG.Cnv,
		srcRect.X, srcRect.Y, srcRect.W, srcRect.H,
		dstRect.X+dx, dstRect.Y+dy, dstRect.W, dstRect.H)
}
