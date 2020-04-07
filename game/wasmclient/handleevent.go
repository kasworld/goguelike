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

package wasmclient

import (
	"syscall/js"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/clientcontroltype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

func (app *WasmClient) makePathToMouseClick(mouseX, mouseY int) {
	switch gameOptions.GetByIDBase("Viewport").State {
	case 0: // play viewpot mode
		cf := app.currentFloor()
		if cf == nil {
			jslog.Error("no current floor")
			return
		}
		tand := app.taNotiData
		if tand == nil || cf.FloorInfo.G2ID != tand.FloorG2ID {
			jslog.Error("invalid floor x,y")
			return
		}
		flX, flY := gVP2d.CanvasXY2FloorXY(
			cf.XWrapSafe, cf.YWrapSafe,
			tand.VPX, tand.VPY,
			mouseX/gVP2d.CellSize, mouseY/gVP2d.CellSize)

		autoPlayButton := autoActs.GetByIDBase("AutoPlay")
		if autoPlayButton.State == 0 {
			autoPlayButton.MakeJSFn(app)(js.Null(), nil)
		}
		playerX, playerY := app.GetPlayerXY()

		if playerX == flX && playerY == flY {
			app.tryEnterPortal(flX, flY)
		} else {
			newPath := cf.Tiles4PathFind.FindPath(
				flX, flY, playerX, playerY, gameconst.ViewPortWH)
			if newPath != nil {
				app.Path2dst = newPath
			}
			app.ClientColtrolMode = clientcontroltype.MoveToDest
			app.KeyDir = way9type.Center
			// fmt.Printf("move2dest [%v %v] to [%v %v] %v", playerX, playerY, floorX, flY, newPath)
		}
	case 1: // floor viewport mode
	}
}

func (app *WasmClient) tryEnterPortal(x, y int) {
	for _, iao := range app.olNotiData.FieldObjList {
		if iao.X == x && iao.Y == y &&
			(iao.ActType == fieldobjacttype.PortalIn || iao.ActType == fieldobjacttype.PortalInOut) {
			// fmt.Printf("enter portal %v %v %v", x, y, iao)
			go app.sendPacket(c2t_idcmd.EnterPortal,
				&c2t_obj.ReqEnterPortal_data{},
			)
		}
	}
}

var Key2Dir = map[string]way9type.Way9Type{
	"End":        way9type.SouthWest,
	"ArrowDown":  way9type.South,
	"PageDown":   way9type.SouthEast,
	"ArrowLeft":  way9type.West,
	"Clear":      way9type.Center,
	"ArrowRight": way9type.East,
	"Home":       way9type.NorthWest,
	"ArrowUp":    way9type.North,
	"PageUp":     way9type.NorthEast,
}

func (app *WasmClient) actByKeyPressMap(kcode string) bool {
	oldDir := app.KeyDir
	dx, dy := app.KeyboardPressedMap.SumMoveDxDy(Key2Dir)
	app.KeyDir = way9type.RemoteDxDy2Way9(dx, dy)

	switch gameOptions.GetByIDBase("Viewport").State {
	case 0: // play viewpot mode
		if app.KeyDir != way9type.Center {
			app.ClientColtrolMode = clientcontroltype.Keyboard
			app.Path2dst = nil
			autoPlayButton := autoActs.GetByIDBase("AutoPlay")
			if autoPlayButton.State == 0 {
				autoPlayButton.MakeJSFn(app)(js.Null(), nil)
			}
			if oldDir != app.KeyDir {
				app.sendMovePacketByInput(app.KeyDir)
			}
			return true
		}
	case 1: // floor viewport mode
		if app.KeyDir != way9type.Center {
			app.ClientColtrolMode = clientcontroltype.Keyboard
			app.Path2dst = nil
			dir := app.KeyDir
			app.floorVPPosX += dir.Dx()
			app.floorVPPosY += dir.Dy()
			return true
		}
	}
	return false
}

func (app *WasmClient) processKeyUpEvent(kcode string) bool {
	// jslog.Errorf("keyup %v", kcode)
	if kcode == "Escape" {
		// reset to default
		app.reset2Default()
		return true
	}

	dx, dy := app.KeyboardPressedMap.SumMoveDxDy(Key2Dir)
	app.KeyDir = way9type.RemoteDxDy2Way9(dx, dy)

	if btn := adminCommandButtons.GetByKeyCode(kcode); btn != nil {
		btn.MakeJSFn(app)(js.Null(), nil)
		return true
	}
	if btn := commandButtons.GetByKeyCode(kcode); btn != nil {
		btn.MakeJSFn(app)(js.Null(), nil)
		return true
	}
	if btn := gameOptions.GetByKeyCode(kcode); btn != nil {
		btn.MakeJSFn(app)(js.Null(), nil)
		return true
	}
	if btn := autoActs.GetByKeyCode(kcode); btn != nil {
		btn.MakeJSFn(app)(js.Null(), nil)
		return true
	}
	return false
}

/////////

func (app *WasmClient) actByMouseRightDown() {
	// follow mouse mode
	switch gameOptions.GetByIDBase("Viewport").State {
	case 0: // play viewpot mode
		autoPlayButton := autoActs.GetByIDBase("AutoPlay")
		if autoPlayButton.State == 0 {
			autoPlayButton.MakeJSFn(app)(js.Null(), nil)
		}
	case 1: // floor viewport mode
	}
	if app.ClientColtrolMode != clientcontroltype.FollowMouse {
		app.ClientColtrolMode = clientcontroltype.FollowMouse
		app.KeyDir = way9type.Center
		app.Path2dst = nil
	} else {
		app.ClientColtrolMode = clientcontroltype.Keyboard
	}
}

func (app *WasmClient) actByMouseMove(mouseX, mouseY int) {
	// update mouse pos
	if gVP2d.CellSize == 0 {
		return
	}
	oldDir := app.MouseDir
	gVP2d.MouseX = mouseX
	gVP2d.MouseY = mouseY
	app.MouseDir = way9type.RemoteDxDy2Way9(
		(mouseX-gVP2d.ViewWidth/2)/gVP2d.CellSize,
		(mouseY-gVP2d.ViewHeight/2)/gVP2d.CellSize,
	)
	switch gameOptions.GetByIDBase("Viewport").State {
	case 0: // play viewpot mode
		if app.ClientColtrolMode == clientcontroltype.FollowMouse {
			if oldDir != app.MouseDir {
				app.sendMovePacketByInput(app.MouseDir)
			}
		}
	case 1: // floor viewport mode
		dir := app.MouseDir
		app.floorVPPosX += dir.Dx()
		app.floorVPPosY += dir.Dy()
	}
}

/////////

// from noti obj list
func (app *WasmClient) actByControlMode() {
	if app.olNotiData == nil || app.olNotiData.ActiveObj.HP <= 0 {
		return
	}

	switch gameOptions.GetByIDBase("Viewport").State {
	case 0: // play viewpot mode
		if app.moveByUserInput() {
			return
		}
	case 1: // floor viewport mode
	}

	if fo := app.onFieldObj; fo == nil ||
		(fo != nil && fieldobjacttype.ClientData[fo.ActType].ActOn) {
		for i, v := range autoActs.ButtonList {
			if v.State == 0 {
				if tryAutoActFn[i](app, v) {
					return
				}
			}
		}
	}

}

func (app *WasmClient) moveByUserInput() bool {
	cf := app.currentFloor()
	playerX, playerY := app.GetPlayerXY()
	if !cf.IsValidPos(playerX, playerY) {
		jslog.Errorf("ao out of floor %v %v", app.olNotiData.ActiveObj, cf)
		return false
	}
	w, h := cf.Tiles.GetXYLen()
	switch app.ClientColtrolMode {
	default:
		jslog.Errorf("invalid ClientColtrolMode %v", app.ClientColtrolMode)
	case clientcontroltype.Keyboard:
		if app.sendMovePacketByInput(app.KeyDir) {
			return true
		}
	case clientcontroltype.FollowMouse:
		if app.sendMovePacketByInput(app.MouseDir) {
			return true
		}
	case clientcontroltype.MoveToDest:
		playerPos := [2]int{playerX, playerY}
		if app.Path2dst == nil || len(app.Path2dst) == 0 {
			app.ClientColtrolMode = clientcontroltype.Keyboard
			return false
		}
		for i := len(app.Path2dst) - 1; i >= 0; i-- {
			nextPos := app.Path2dst[i]
			isContact, dir := way9type.CalcContactDirWrapped(playerPos, nextPos, w, h)
			if isContact {
				if dir == way9type.Center {
					// arrived
					app.Path2dst = nil
					app.ClientColtrolMode = clientcontroltype.Keyboard
					return false
				} else {
					go app.sendPacket(c2t_idcmd.Move,
						&c2t_obj.ReqMove_data{Dir: dir},
					)
					return true
				}
			}
		}
		// fail to path2dest
		app.Path2dst = nil
		app.ClientColtrolMode = clientcontroltype.Keyboard
	}
	return false
}
