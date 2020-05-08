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

	"github.com/kasworld/goguelike/game/wasmclient/soundmap"
	"github.com/kasworld/goguelike/lib/htmlbutton"
	"github.com/kasworld/gowasmlib/jslog"
)

var gameOptions *htmlbutton.HTMLButtonGroup

// prevent compiler initialize loop error
var _gameopt = htmlbutton.NewButtonGroup("Options",
	[]*htmlbutton.HTMLButton{
		htmlbutton.New("q", "LeftInfo", []string{"LeftInfoOff", "LeftInfoOn"},
			"show/hide left info", btnFocus2Canvas, 1),
		htmlbutton.New("w", "CenterInfo", []string{"HelpOff", "Highscore", "ClientInfo", "Help", "FactionInfo",
			"CarryObjectInfo", "PotionInfo", "ScrollInfo", "MoneyColor",
			"TileInfo", "ConditionInfo", "FieldObjInfo"},
			"rotate help info", cmdRotateCenterInfo, 0),
		htmlbutton.New("e", "RightInfo", []string{
			"RightInfoOff", "Message", "DebugInfo", "InvenList", "FieldObjList", "FloorList"},
			"Rotate right info", cmdRotateRightInfo, 1),
		htmlbutton.New("r", "Viewport", []string{"PlayVP", "FloorVP"},
			"play view / floor view", cmdToggleVPFloorPlay, 0),
		htmlbutton.New("t", "Zoom", []string{"Zoom0", "Zoom1", "Zoom2"},
			"Zoom viewport", cmdToggleZoom, 0),
		htmlbutton.New("y", "Sound", []string{"SoundOn", "SoundOff"},
			"Sound on/off", cmdToggleSound, 1),
	})

func (app *WasmClient) updateLeftInfo() {
	v := gameOptions.GetByIDBase("LeftInfo")
	infoobj := js.Global().Get("document").Call("getElementById", "leftinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", "")
	case 1: // leftinfo
		infoobj.Set("innerHTML",
			app.makeBaseInfoHTML()+app.makeBuffListHTML())
	}
}

func cmdRotateCenterInfo(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.updateCenterInfo()
	Focus2Canvas()
}

func (app *WasmClient) updateCenterInfo() {
	v := gameOptions.GetByIDBase("CenterInfo")
	infoobj := js.Global().Get("document").Call("getElementById", "centerinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", "")
	case 1: // highscore
		go func() {
			infoobj.Set("innerHTML", loadHighScoreHTML())
		}()
	case 2: // clientinfo
		infoobj.Set("innerHTML", makeClientInfoHTML())
	case 3: // helpinfo
		infoobj.Set("innerHTML", makeHelpInfoHTML())
	case 4: // faction
		infoobj.Set("innerHTML", makeHelpFactionHTML())
	case 5: // carryobj
		infoobj.Set("innerHTML", makeHelpCarryObjectHTML())
	case 6: // potion
		infoobj.Set("innerHTML", makeHelpPotionHTML())
	case 7: // scroll
		infoobj.Set("innerHTML", makeHelpScrollHTML())
	case 8: // Money color
		infoobj.Set("innerHTML", makeHelpMoneyColorHTML())
	case 9: // tile
		infoobj.Set("innerHTML", makeHelpTileHTML())
	case 10: // condition
		infoobj.Set("innerHTML", makeHelpConditionHTML())
	case 11: // fieldobj
		infoobj.Set("innerHTML", makeHelpFieldObjHTML())
	}
}

func cmdRotateRightInfo(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.updateRightInfo()
	Focus2Canvas()
}

func (app *WasmClient) updateRightInfo() {
	v := gameOptions.GetByIDBase("RightInfo")
	infoobj := js.Global().Get("document").Call("getElementById", "rightinfo")
	switch v.State {
	case 0: // Hide
		infoobj.Set("innerHTML", "")
	case 1: // Message
		app.systemMessage = app.systemMessage.GetLastN(DisplayLineLimit)
		infoobj.Set("innerHTML", app.systemMessage.ToHtmlStringRev())
	case 2: // DebugInfo
		infoobj.Set("innerHTML", app.makeDebugInfoHTML())
	case 3: // InvenList
		infoobj.Set("innerHTML", app.makeInvenInfoHTML())
	case 4: // FieldObjList
		infoobj.Set("innerHTML", app.makeFieldObjListHTML())
	case 5: // FloorList
		infoobj.Set("innerHTML", app.makeFloorListHTML())
	}
}

func cmdToggleVPFloorPlay(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	switch v.State {
	case 0: // play viewpot mode
	case 1: // floor viewport mode
		app.floorVPPosX, app.floorVPPosY = app.GetPlayerXY()
	}
	Focus2Canvas()
}

func cmdToggleZoom(obj interface{}, v *htmlbutton.HTMLButton) {
	gVP2d.Zoom(v.State) // set zoomstate , needrecalc
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.ResizeCanvas()
	app.systemMessage.Appendf("Zoom%v", v.State)
	gVP2d.NotiMessage.AppendTf(tcsInfo, "Zoom%v", v.State)
	Focus2Canvas()
}

func cmdToggleSound(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	if v.State == 0 {
		soundmap.SoundOn = true
		app.systemMessage.Append("SoundOn")
		gVP2d.NotiMessage.AppendTf(tcsInfo, "SoundOn")
	} else {
		soundmap.SoundOn = false
		app.systemMessage.Append("SoundOff")
		gVP2d.NotiMessage.AppendTf(tcsInfo, "SoundOff")
	}
	Focus2Canvas()
}
