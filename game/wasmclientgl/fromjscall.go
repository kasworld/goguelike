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
	"strings"
	"syscall/js"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/clientcontroltype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/lib/jsobj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

func (app *WasmClient) registerJSButton() {
	js.Global().Set("moveFloor", js.FuncOf(app.jsMove2Floor))
	js.Global().Set("sendChat", js.FuncOf(app.jsSendChat))
	js.Global().Set("unequip", js.FuncOf(app.jsUnequipCarryObj))
	js.Global().Set("equip", js.FuncOf(app.jsEquipCarryObj))
	js.Global().Set("drop", js.FuncOf(app.jsDropCarryObj))
	js.Global().Set("drinkpotion", js.FuncOf(app.jsDrinkPotion))
	js.Global().Set("readscroll", js.FuncOf(app.jsReadScroll))
	js.Global().Set("recycle", js.FuncOf(app.jsRecycleCarryObj))
}

func (app *WasmClient) jsUnequipCarryObj(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.UnEquip,
		&c2t_obj.ReqUnEquip_data{UUID: id},
	)
	GetElementById(id).Call("blur")
	return nil
}
func (app *WasmClient) jsEquipCarryObj(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.Equip,
		&c2t_obj.ReqEquip_data{UUID: id},
	)
	GetElementById(id).Call("blur")
	return nil
}
func (app *WasmClient) jsDrinkPotion(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.DrinkPotion,
		&c2t_obj.ReqDrinkPotion_data{UUID: id},
	)
	GetElementById(id).Call("blur")
	return nil
}

func (app *WasmClient) jsReadScroll(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.ReadScroll,
		&c2t_obj.ReqReadScroll_data{UUID: id},
	)
	GetElementById(id).Call("blur")
	return nil
}

func (app *WasmClient) jsRecycleCarryObj(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.Recycle,
		&c2t_obj.ReqRecycle_data{UUID: id},
	)
	GetElementById(id).Call("blur")
	return nil
}
func (app *WasmClient) jsDropCarryObj(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.Drop,
		&c2t_obj.ReqDrop_data{UUID: id},
	)
	GetElementById(id).Call("blur")
	return nil
}

func (app *WasmClient) jsMove2Floor(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.MoveFloor,
		&c2t_obj.ReqMoveFloor_data{UUID: id},
	)
	GetElementById(id).Call("blur")
	return nil
}

func getChatMsg() string {
	msg := jsobj.GetTextValueFromInputText("chattext")
	msg = strings.TrimSpace(msg)
	if len(msg) > gameconst.MaxChatLen {
		msg = msg[:gameconst.MaxChatLen]
	}
	return msg
}

func (app *WasmClient) jsSendChat(this js.Value, args []js.Value) interface{} {
	msg := getChatMsg()
	go app.sendPacket(c2t_idcmd.Chat,
		&c2t_obj.ReqChat_data{Chat: msg})
	GetElementById("chatbutton").Call("blur")
	GetElementById("chattext").Call("blur")
	return nil
}

func AddEventListener(
	dst js.Value,
	evt string,
	fn func(this js.Value, args []js.Value) interface{}) {
	dst.Call("addEventListener", evt, js.FuncOf(fn))
}

func (app *WasmClient) registerKeyboardMouseEvent() {
	dst := js.Global().Get("window")
	AddEventListener(dst, "click", app.jsHandleMouseClick)
	AddEventListener(dst, "wheel", app.jsHandleMouseWheel)
	AddEventListener(dst, "mousemove", app.jsHandleMouseMove)
	AddEventListener(dst, "mousedown", app.jsHandleMouseDown)
	AddEventListener(dst, "mouseup", app.jsHandleMouseUp)
	AddEventListener(dst, "contextmenu", app.jsHandleContextMenu)
	AddEventListener(dst, "keydown", app.jsHandleKeyDown)
	AddEventListener(dst, "keypress", app.jsHandleKeyPress)
	AddEventListener(dst, "keyup", app.jsHandleKeyUp)
}

func (app *WasmClient) jsHandleMouseWheel(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	// evt.Call("preventDefault")

	deltaMode := evt.Get("deltaMode").Int()
	// deltaX deltaY deltaZ
	deltaX := evt.Get("deltaX").Int()
	deltaY := evt.Get("deltaY").Int()
	deltaZ := evt.Get("deltaZ").Int()

	jslog.Infof("mode %v x %v y %v z %v", deltaMode, deltaX, deltaY, deltaZ)
	switch deltaMode {
	case 0: // pixels
	case 1: // lines
	case 2: // ripagesght
	}
	return nil
}

func (app *WasmClient) jsHandleMouseClick(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")

	btn := evt.Get("button").Int()

	switch btn {
	case 0: // left
		app.makePathToMouseClick()
		return nil
	case 1: // wheel

	case 2: // right
	}
	return nil
}

func (app *WasmClient) makePathToMouseClick() {
	switch gameOptions.GetByIDBase("Viewport").State {
	case 0: // play viewpot mode
		cf := app.currentFloor()
		if cf == nil {
			jslog.Error("no current floor")
			return
		}
		tand := app.taNotiData
		if tand == nil || cf.FloorInfo.UUID != tand.FloorUUID {
			jslog.Error("invalid floor x,y")
			return
		}
		flX, flY := app.vp.mouseCursorFx, app.vp.mouseCursorFy

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
		if len(app.Path2dst) != 0 {
			autoPlayButton := autoActs.GetByIDBase("AutoPlay")
			if autoPlayButton.State == 0 {
				autoPlayButton.JSFn(js.Null(), nil)
			}
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

func (app *WasmClient) jsHandleMouseDown(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	// Never call,  relate focus , prevent key event listen
	// evt.Call("stopPropagation")
	// evt.Call("preventDefault")
	btn := evt.Get("button").Int()

	switch btn {
	case 2: // right click
		app.actByMouseRightDown()
	}
	return nil
}
func (app *WasmClient) actByMouseRightDown() {
	// follow mouse mode
	switch gameOptions.GetByIDBase("Viewport").State {
	case 0: // play viewpot mode
		autoPlayButton := autoActs.GetByIDBase("AutoPlay")
		if autoPlayButton.State == 0 {
			autoPlayButton.JSFn(js.Null(), nil)
		}
	case 1: // floor viewport mode
	}
	if app.ClientColtrolMode != clientcontroltype.FollowMouse {
		app.ClientColtrolMode = clientcontroltype.FollowMouse
		app.KeyDir = way9type.Center
		app.Path2dst = nil
		app.vp.ClearMovePath()
	} else {
		app.ClientColtrolMode = clientcontroltype.Keyboard
	}
}

func (app *WasmClient) jsHandleMouseUp(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")
	return nil
}
func (app *WasmClient) jsHandleContextMenu(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")
	return nil
}

func (app *WasmClient) jsHandleMouseMove(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")

	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Float()
	winH := win.Get("innerHeight").Float()
	jsMouse := ThreeJsNew("Vector2")
	jsMouse.Set("x", (evt.Get("clientX").Float()/winW)*2-1)
	jsMouse.Set("y", -(evt.Get("clientY").Float()/winH)*2+1)

	app.vp.mouseCursorFx, app.vp.mouseCursorFy = app.vp.FindRayCastingFxFy(jsMouse)

	app.actByMouseMove()
	return nil
}

func (app *WasmClient) actByMouseMove() {
	oldDir := app.MouseDir
	plx, ply := app.GetPlayerXY()
	app.MouseDir = way9type.RemoteDxDy2Way9(
		app.vp.mouseCursorFx-plx,
		app.vp.mouseCursorFy-ply,
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

///////////////////////////////////////////////////////
// keyboard handle

var jsInputTarget = js.Global().Get("body")

func (app *WasmClient) jsHandleKeyDown(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	if evt.Get("target").Equal(jsInputTarget) {
		evt.Call("stopPropagation")
		evt.Call("preventDefault")

		kcode := evt.Get("key").String()
		if kcode != "" {
			app.KeyboardPressedMap.KeyDown(kcode)
		}
		app.actByKeyPressMap(kcode)
	}
	return nil
}
func (app *WasmClient) jsHandleKeyPress(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	if evt.Get("target").Equal(jsInputTarget) {
		evt.Call("stopPropagation")
		evt.Call("preventDefault")

		kcode := evt.Get("key").String()
		app.actByKeyPressMap(kcode)
	}
	return nil
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
			app.vp.ClearMovePath()
			autoPlayButton := autoActs.GetByIDBase("AutoPlay")
			if autoPlayButton.State == 0 {
				autoPlayButton.JSFn(js.Null(), nil)
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
			app.vp.ClearMovePath()
			dir := app.KeyDir
			app.floorVPPosX += dir.Dx()
			app.floorVPPosY += dir.Dy()
			return true
		}
	}
	return false
}

func (app *WasmClient) jsHandleKeyUp(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	if evt.Get("target").Equal(jsInputTarget) {
		evt.Call("stopPropagation")
		evt.Call("preventDefault")

		kcode := evt.Get("key").String()
		// jslog.Infof("%v %v", evt, kcode)
		if kcode != "" {
			app.KeyboardPressedMap.KeyUp(kcode)
		}
		app.processKeyUpEvent(kcode)
	}
	return nil
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
		btn.JSFn(js.Null(), nil)
		return true
	}
	if btn := commandButtons.GetByKeyCode(kcode); btn != nil {
		btn.JSFn(js.Null(), nil)
		return true
	}
	if btn := gameOptions.GetByKeyCode(kcode); btn != nil {
		btn.JSFn(js.Null(), nil)
		return true
	}
	if btn := autoActs.GetByKeyCode(kcode); btn != nil {
		btn.JSFn(js.Null(), nil)
		return true
	}
	return false
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
