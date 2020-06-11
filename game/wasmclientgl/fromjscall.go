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
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/lib/jsobj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
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
		&c2t_obj.ReqUnEquip_data{G2ID: g2id.NewFromString(id)},
	)
	app.Focus2Canvas()
	return nil
}
func (app *WasmClient) jsEquipCarryObj(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.Equip,
		&c2t_obj.ReqEquip_data{G2ID: g2id.NewFromString(id)},
	)
	app.Focus2Canvas()
	return nil
}
func (app *WasmClient) jsDrinkPotion(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.DrinkPotion,
		&c2t_obj.ReqDrinkPotion_data{G2ID: g2id.NewFromString(id)},
	)
	app.Focus2Canvas()
	return nil
}

func (app *WasmClient) jsReadScroll(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.ReadScroll,
		&c2t_obj.ReqReadScroll_data{G2ID: g2id.NewFromString(id)},
	)
	app.Focus2Canvas()
	return nil
}

func (app *WasmClient) jsRecycleCarryObj(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.Recycle,
		&c2t_obj.ReqRecycle_data{G2ID: g2id.NewFromString(id)},
	)
	app.Focus2Canvas()
	return nil
}
func (app *WasmClient) jsDropCarryObj(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.Drop,
		&c2t_obj.ReqDrop_data{G2ID: g2id.NewFromString(id)},
	)
	app.Focus2Canvas()
	return nil
}

func (app *WasmClient) jsMove2Floor(this js.Value, args []js.Value) interface{} {
	id := strings.TrimSpace(args[0].String())
	go app.sendPacket(c2t_idcmd.MoveFloor,
		&c2t_obj.ReqMoveFloor_data{G2ID: g2id.NewFromString(id)},
	)
	app.Focus2Canvas()
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
	app.Focus2Canvas()
	return nil
}

func (app *WasmClient) registerKeyboardMouseEvent() {
	app.vp.AddEventListener("click", app.jsHandleMouseClickVP)
	app.vp.AddEventListener("mousemove", app.jsHandleMouseMoveVP)
	app.vp.AddEventListener("mousedown", app.jsHandleMouseDownVP)
	app.vp.AddEventListener("mouseup", app.jsHandleMouseUpVP)
	app.vp.AddEventListener("contextmenu", app.jsHandleContextMenu)
	app.vp.AddEventListener("keydown", app.jsHandleKeyDownVP)
	app.vp.AddEventListener("keypress", app.jsHandleKeyPressVP)
	app.vp.AddEventListener("keyup", app.jsHandleKeyUpVP)
}

func (app *WasmClient) jsHandleMouseClickVP(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")

	mouseX, mouseY := evt.Get("offsetX").Int(), evt.Get("offsetY").Int()
	btn := evt.Get("button").Int()

	switch btn {
	case 0: // left
		app.makePathToMouseClick(mouseX, mouseY)
		return nil
	case 1: // wheel

	case 2: // right
	}
	return nil
}

func (app *WasmClient) makePathToMouseClick(mouseX, mouseY int) {
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

func (app *WasmClient) jsHandleMouseDownVP(this js.Value, args []js.Value) interface{} {
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
	} else {
		app.ClientColtrolMode = clientcontroltype.Keyboard
	}
}

func (app *WasmClient) jsHandleMouseUpVP(this js.Value, args []js.Value) interface{} {
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

func (app *WasmClient) jsHandleMouseMoveVP(this js.Value, args []js.Value) interface{} {
	evt := args[0]
	evt.Call("stopPropagation")
	evt.Call("preventDefault")

	mouseX, mouseY := evt.Get("offsetX").Int(), evt.Get("offsetY").Int()
	app.actByMouseMove(mouseX, mouseY)
	return nil
}
func (app *WasmClient) actByMouseMove(mouseX, mouseY int) {
}

func (app *WasmClient) jsHandleKeyDownVP(this js.Value, args []js.Value) interface{} {
	return nil
}
func (app *WasmClient) jsHandleKeyPressVP(this js.Value, args []js.Value) interface{} {
	return nil
}

func (app *WasmClient) actByKeyPressMap(kcode string) bool {
	return false
}

func (app *WasmClient) jsHandleKeyUpVP(this js.Value, args []js.Value) interface{} {
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
