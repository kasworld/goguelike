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
	"bytes"
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/kasworld/actjitter"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/clientcontroltype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/clientcookie"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/game/clientinitdata"
	"github.com/kasworld/goguelike/game/soundmap"
	"github.com/kasworld/goguelike/lib/clienttile"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/lib/htmlbutton"
	"github.com/kasworld/goguelike/lib/jskeypressmap"
	"github.com/kasworld/goguelike/lib/jsobj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_connwasm"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_pid2rspfn"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_version"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/textncount"
	"github.com/kasworld/gowasmlib/wrapspan"
	"github.com/kasworld/intervalduration"
)

var gInitData *clientinitdata.InitData
var gClientTile *clienttile.ClientTile

type WasmClient struct {
	// app info
	DoClose            func()
	systemMessage      textncount.TextNCountList
	KeyboardPressedMap *jskeypressmap.KeyPressMap
	Path2dst           [][2]int
	ClientColtrolMode  clientcontroltype.ClientControlType

	AOG2ID2AOClient       map[g2id.G2ID]*c2t_obj.ActiveObjClient
	CaObjG2ID2CaObjClient map[g2id.G2ID]interface{}
	G2ID2ClientFloor      map[g2id.G2ID]*clientfloor.ClientFloor
	FloorInfo             *c2t_obj.FloorInfo
	remainTurn2Rebirth    int

	// for net
	pid2recv             *c2t_pid2rspfn.PID2RspFn
	wsConn               *c2t_connwasm.Connection
	ServerJitter         *actjitter.ActJitter
	PingDur              time.Duration
	ServerClientTimeDiff time.Duration
	ClientJitter         *actjitter.ActJitter

	// from user input
	KeyDir   way9type.Way9Type
	MouseDir way9type.Way9Type

	// for turn
	waitObjList    bool
	needRefreshSet bool

	taNotiData     *c2t_obj.NotiVPTiles_data
	olNotiData     *c2t_obj.NotiObjectList_data
	lastOLNotiData *c2t_obj.NotiObjectList_data

	movePacketPerTurn int32
	actPacketPerTurn  int32
	lastEffBias       bias.Bias
	onFieldObj        *c2t_obj.FieldObjClient
	OverLoadRate      float64
	HPdiff            int
	SPdiff            int
	level             int

	// debug
	taNotiHeader c2t_packet.Header
	olNotiHeader c2t_packet.Header
	DispInterDur *intervalduration.IntervalDuration

	// for floor view mode
	floorVPPosX int
	floorVPPosY int

	vp *Viewport
}

// call after pageload
func InitPage() {
	// hide loading message
	js.Global().Get("document").Call("getElementById", "loadmsg").Get("style").Set("display", "none")

	gInitData = clientinitdata.New()
	gClientTile = clienttile.New()
	gameOptions = _gameopt // prevent compiler initialize loop

	app := &WasmClient{
		ServerJitter:     actjitter.New("Server"),
		G2ID2ClientFloor: make(map[g2id.G2ID]*clientfloor.ClientFloor),

		systemMessage:      make(textncount.TextNCountList, 0),
		KeyboardPressedMap: jskeypressmap.New(),
		DispInterDur:       intervalduration.New("Display"),
		ClientJitter:       actjitter.New("Client"),

		pid2recv: c2t_pid2rspfn.New(),
		DoClose:  func() { jslog.Errorf("Too early DoClose call") },
	}
	app.vp = NewViewport()
	js.Global().Call("requestAnimationFrame", js.FuncOf(app.drawCanvas))

	js.Global().Set("enterTower", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go app.enterTower(args[0].Int())
		return nil
	}))
	js.Global().Set("clearSession", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go clientcookie.ClearSession(args[0].Int())
		return nil
	}))
	app.registerJSButton()

	clientcookie.InitNickname()

	js.Global().Get("document").Call("getElementById", "centerinfo").Set("innerHTML",
		clientinitdata.MakeClientInfoHTML()+
			clientinitdata.MakeHelpFactionHTML()+
			MakeHelpInfoHTML()+
			clientinitdata.MakeHelpCarryObjectHTML()+
			clientinitdata.MakeHelpPotionHTML()+
			clientinitdata.MakeHelpScrollHTML()+
			clientinitdata.MakeHelpMoneyColorHTML()+
			clientinitdata.MakeHelpTileHTML()+
			clientinitdata.MakeHelpConditionHTML()+
			clientinitdata.MakeHelpFieldObjHTML())
	go func() {
		str := clientinitdata.LoadHighScoreHTML() +
			clientinitdata.MakeClientInfoHTML() +
			clientinitdata.MakeHelpFactionHTML() +
			MakeHelpInfoHTML() +
			clientinitdata.MakeHelpCarryObjectHTML() +
			clientinitdata.MakeHelpPotionHTML() +
			clientinitdata.MakeHelpScrollHTML() +
			clientinitdata.MakeHelpMoneyColorHTML() +
			clientinitdata.MakeHelpTileHTML() +
			clientinitdata.MakeHelpConditionHTML() +
			clientinitdata.MakeHelpFieldObjHTML()
		js.Global().Get("document").Call("getElementById", "centerinfo").Set("innerHTML", str)
	}()

	app.registerKeyboardMouseEvent()

	app.ResizeCanvas()

	win := js.Global().Get("window")
	win.Call("addEventListener", "resize", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			app.ResizeCanvas()
			return nil
		},
	))

	app.AOG2ID2AOClient = make(map[g2id.G2ID]*c2t_obj.ActiveObjClient)
	app.CaObjG2ID2CaObjClient = make(map[g2id.G2ID]interface{})
}

func (app *WasmClient) makeButtons() string {
	var buf bytes.Buffer
	adminCommandButtons.MakeButtonToolTipTop(&buf)
	gameOptions.MakeButtonToolTipTop(&buf)
	commandButtons.MakeButtonToolTipTop(&buf)
	autoActs.MakeButtonToolTipTop(&buf)
	return buf.String()
}

// link to enter tower button
func (app *WasmClient) enterTower(towerindex int) {
	gInitData.TowerIndex = towerindex

	jsdoc := js.Global().Get("document")
	jsobj.Hide(jsdoc.Call("getElementById", "titleform"))
	js.Global().Get("document").Call("getElementById", "centerinfo").Set("innerHTML", "")
	jsobj.Show(jsdoc.Call("getElementById", "cmdrow"))

	jsdoc.Call("getElementById", "leftinfo").Set("style",
		"color: white; position: fixed; top: 0; left: 0; overflow: hidden;")
	jsdoc.Call("getElementById", "rightinfo").Set("style",
		"color: white; position: fixed; top: 0; right: 0; overflow: hidden; text-align: right;")
	jsdoc.Call("getElementById", "centerinfo").Set("style",
		"color: white; position: fixed; top: 0%; left: 25%; overflow: hidden;")

	app.vp.hideTitle()
	app.Focus2Canvas()

	commandButtons.RegisterJSFn(app)
	autoActs.RegisterJSFn(app)
	gameOptions.RegisterJSFn(app)
	adminCommandButtons.RegisterJSFn(app)
	if err := gameOptions.SetFromURLArg(); err != nil {
		jslog.Errorf(err.Error())
	}
	jsdoc.Call("getElementById", "cmdbuttons").Set("innerHTML",
		app.makeButtons())

	app.reset2Default()
	// cmdToggleSound(app, gameOptions.GetByIDBase("Sound"))

	ctx, closeCtx := context.WithCancel(context.Background())
	app.DoClose = closeCtx
	defer app.DoClose()

	if err := app.NetInit(ctx); err != nil {
		jslog.Errorf("%v\n", err)
		return
	}
	defer app.Cleanup()
	app.systemMessage.Append(wrapspan.ColorTextf("yellow",
		"Welcome to Goguelike, %v!", gInitData.GetNickName()))
	// app.vp.NotiMessage.AppendTf(tcsInfo, "Welcome to Goguelike, %v!", gInitData.GetNickName())

	if gameconst.DataVersion != gInitData.ServiceInfo.DataVersion {
		jslog.Errorf("DataVersion mismatch client %v server %v",
			gameconst.DataVersion, gInitData.ServiceInfo.DataVersion)
	}
	if c2t_version.ProtocolVersion != gInitData.ServiceInfo.ProtocolVersion {
		jslog.Errorf("ProtocolVersion mismatch client %v server %v",
			c2t_version.ProtocolVersion, gInitData.ServiceInfo.ProtocolVersion)
	}

	// app.vp.ViewportPos2Index = gInitData.ViewportXYLenList.MakePos2Index()

	clientcookie.SetSession(towerindex, string(gInitData.AccountInfo.SessionG2ID), gInitData.AccountInfo.NickName)

	if gInitData.CanUseCmd(c2t_idcmd.AIPlay) {
		app.reqAIPlay(true)
	}
	for i, v := range adminCommandButtons.ButtonList {
		if !gInitData.CanUseCmd(adminCmds[i]) {
			v.Disable()
			v.Hide()
		}
	}

	go soundmap.Play("startsound")

	app.ResizeCanvas()

	timerPingTk := time.NewTicker(time.Second)
	defer timerPingTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerPingTk.C:
			go app.reqHeartbeat()
		}
	}
}

func (app *WasmClient) drawCanvas(this js.Value, args []js.Value) interface{} {
	defer func() {
		js.Global().Call("requestAnimationFrame", js.FuncOf(app.drawCanvas))
	}()
	act := app.DispInterDur.BeginAct()
	defer act.End()

	app.vp.Draw()

	return nil
}

func (app *WasmClient) ResizeCanvas() {
	if gInitData.AccountInfo == nil {
		app.vp.ResizeCanvas(true)
	} else {
		app.vp.ResizeCanvas(false)
		win := js.Global().Get("window")
		winH := win.Get("innerHeight").Int()
		ftsize := fmt.Sprintf("%vpx", winH/100)
		js.Global().Get("document").Call("getElementById", "body").Get("style").Set("font-size", ftsize)
		for _, v := range commandButtons.ButtonList {
			v.JSButton().Get("style").Set("font-size", ftsize)
		}
		for _, v := range autoActs.ButtonList {
			v.JSButton().Get("style").Set("font-size", ftsize)
		}
		for _, v := range gameOptions.ButtonList {
			v.JSButton().Get("style").Set("font-size", ftsize)
		}
		for _, v := range adminCommandButtons.ButtonList {
			v.JSButton().Get("style").Set("font-size", ftsize)
		}
		js.Global().Get("document").Call("getElementById", "chattext").Get("style").Set("font-size", ftsize)
		js.Global().Get("document").Call("getElementById", "chatbutton").Get("style").Set("font-size", ftsize)
	}
}

// for htmlbutton click fn
func btnFocus2Canvas(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	app.Focus2Canvas()
}

func (app *WasmClient) Focus2Canvas() {
	app.vp.Focus()
}
