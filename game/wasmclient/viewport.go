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
	"fmt"

	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"

	"syscall/js"

	"github.com/kasworld/goguelike/enum/clientcontroltype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/wasmclient/htmlbutton"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/wrapspan"
)

// for htmlbutton click fn
func btnFocus2Canvas(obj interface{}, v *htmlbutton.HTMLButton) {
	Focus2Canvas()
}

func Focus2Canvas() {
	gVP2d.Focus()
}

func (app *WasmClient) ResizeCanvas() {

	if gInitData.AccountInfo == nil {
		gVP2d.DrawTitle()
	} else {
		gVP2d.Resize()
		ftsize := fmt.Sprintf("%vpx", gVP2d.RefSize/4)
		js.Global().Get("document").Call("getElementById", "body").Get("style").Set("font-size", ftsize)
		for _, v := range commandButtons.ButtonList {
			js.Global().Get("document").Call("getElementById", v.JSID()).Get("style").Set("font-size", ftsize)
		}
		for _, v := range autoActs.ButtonList {
			js.Global().Get("document").Call("getElementById", v.JSID()).Get("style").Set("font-size", ftsize)
		}
		for _, v := range gameOptions.ButtonList {
			js.Global().Get("document").Call("getElementById", v.JSID()).Get("style").Set("font-size", ftsize)
		}
		for _, v := range adminCommandButtons.ButtonList {
			js.Global().Get("document").Call("getElementById", v.JSID()).Get("style").Set("font-size", ftsize)
		}
		js.Global().Get("document").Call("getElementById", "chattext").Get("style").Set("font-size", ftsize)
		js.Global().Get("document").Call("getElementById", "chatbutton").Get("style").Set("font-size", ftsize)

		js.Global().Get("document").Call("getElementById", "leftinfo").Set("style",
			"color: white; position: fixed; top: 0; left: 0; overflow: hidden;")
		js.Global().Get("document").Call("getElementById", "rightinfo").Set("style",
			"color: white; position: fixed; top: 0; right: 0; overflow: hidden; text-align: right;")
		js.Global().Get("document").Call("getElementById", "centerinfo").Set("style",
			"color: white; position: fixed; top: 0%; left: 25%; overflow: hidden;")

	}
}

func (app *WasmClient) drawCanvas(this js.Value, args []js.Value) interface{} {
	if app.waitObjList {
		// app.systemMessage.Append(
		// 	wrapspan.ColorText("OrangeRed", "waiting object list"))
		app.needRefreshSet = true
		return nil
	}
	defer func() {
		js.Global().Call("requestAnimationFrame", js.FuncOf(app.drawCanvas))
	}()

	cf := app.currentFloor()
	if cf == nil {
		jslog.Warn("no currentfloor")
		return nil
	}
	if cf.FloorInfo == nil {
		jslog.Warn("no floorinfo")
		return nil
	}
	if app.taNotiData == nil {
		jslog.Warn("no app.taNotiData")
		return nil
	}
	if cf.FloorInfo.G2ID != app.taNotiData.FloorG2ID {
		app.systemMessage.Append(wrapspan.ColorText("OrangeRed", "changeing floor..."))
		return nil
	}
	actMoveDir := app.getMoveDirByControlMode()
	frameProgress := app.ClientJitter.GetInFrameProgress2()
	envBias := app.GetEnvBias()
	act := app.DispInterDur.BeginAct()
	if gameOptions.GetByIDBase("Viewport").State == 0 { // play viewport mode
		scrollDir := app.getScrollDir()
		if err := gVP2d.DrawPlayVP(
			frameProgress, envBias, app.TowerBias(), cf.GetBias(),
			gInitData.AccountInfo.ActiveObjG2ID, cf,
			scrollDir, actMoveDir,
			app.taNotiData,
			app.olNotiData, app.lastOLNotiData,
			app.Path2dst,
		); err != nil {
			jslog.Errorf("%v", err)
		}
	} else {
		if err := gVP2d.DrawFloorVP(
			envBias,
			cf, app.floorVPPosX, app.floorVPPosY, actMoveDir); err != nil {
			jslog.Errorf("%v", err)
		}
	}
	act.End()
	return nil
}

func (app *WasmClient) getMoveDirByControlMode() way9type.Way9Type {
	dir := way9type.Center
	if app.ClientColtrolMode == clientcontroltype.FollowMouse {
		dir = app.MouseDir
	}
	if app.ClientColtrolMode == clientcontroltype.Keyboard {
		dir = app.KeyDir
	}
	return dir
}

func (app *WasmClient) getScrollDir() way9type.Way9Type {
	scrollDir := way9type.Center
	if pao, exist := app.AOG2ID2AOClient[gInitData.AccountInfo.ActiveObjG2ID]; exist {
		if pao.Act == c2t_idcmd.Move && pao.Result == c2t_error.None {
			scrollDir = pao.Dir
		}
	}
	return scrollDir
}
