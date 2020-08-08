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

// Package towerwebjsconn manage websocket connection to tower
package wasmclient

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"syscall/js"
	"time"

	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/clientcookie"
	"github.com/kasworld/goguelike/lib/jsobj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_connwasm"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_pid2rspfn"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/wasmcookie"
	"github.com/kasworld/gowasmlib/wrapspan"
)

func GetQuery() url.Values {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	if err != nil {
		jslog.Errorf("%v", err)
	}
	return u.Query()
}

func (app *WasmClient) NetInit(ctx context.Context) error {
	app.wsConn = c2t_connwasm.New(
		gInitData.GetConnectToTowerURL(),
		c2t_gob.MarshalBodyFn,
		app.handleRecvPacket,
		app.handleSentPacket)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err := app.wsConn.Connect(ctx, &wg)
		if err != nil {
			jslog.Errorf("wsConn.Connect err %v", err)
			app.systemMessage.Append(
				wrapspan.ColorTextf("red", "tower disconnected %v", err))
			app.DoClose()
		}
	}()
	authkey := GetQuery().Get("authkey")
	nick := jsobj.GetTextValueFromInputText("nickname")
	ck := wasmcookie.GetMap()
	sessionkey := string(ck[clientcookie.SessionKeyName(gInitData.TowerIndex)])
	wg.Wait()

	wg.Add(1)
	app.ReqWithRspFn(
		c2t_idcmd.Login,
		&c2t_obj.ReqLogin_data{
			SessionUUID: sessionkey,
			NickName:    nick,
			AuthKey:     authkey,
		},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspLogin_data)
			gInitData.ServiceInfo = rpk.ServiceInfo
			gInitData.AccountInfo = rpk.AccountInfo
			wg.Done()
			return nil
		},
	)
	wg.Wait()
	return nil
}

func (app *WasmClient) Cleanup() {
	app.wsConn.SendRecvStop()
}

func (app *WasmClient) handleSentPacket(header c2t_packet.Header) error {
	return nil
}

func (app *WasmClient) handleRecvPacket(header c2t_packet.Header, body []byte) error {
	robj, err := c2t_gob.UnmarshalPacket(header, body)
	if err != nil {
		return err
	}
	switch header.FlowType {
	default:
		return fmt.Errorf("Invalid packet type %v %v", header, robj)
	case c2t_packet.Response:
		if err := app.pid2recv.HandleRsp(header, robj); err != nil {
			return err
		}
	case c2t_packet.Notification:
		fn := ProcessRecvObjNotiFnMap[header.Cmd]
		if err := fn(app, header, robj); err != nil {
			return err
		}
	}
	return nil
}

func (app *WasmClient) ReqWithRspFn(cmd c2t_idcmd.CommandID, body interface{},
	fn c2t_pid2rspfn.HandleRspFn) error {

	pid := app.pid2recv.NewPID(fn)
	spk := c2t_packet.Packet{
		Header: c2t_packet.Header{
			Cmd:      uint16(cmd),
			ID:       pid,
			FlowType: c2t_packet.Request,
		},
		Body: body,
	}
	if err := app.wsConn.EnqueueSendPacket(spk); err != nil {
		app.wsConn.SendRecvStop()
		return fmt.Errorf("Send fail %s %v:%v %v", app, cmd, pid, err)
	}
	return nil
}

func (app *WasmClient) ReqWithRspFnWithAuth(cmd c2t_idcmd.CommandID, body interface{},
	fn c2t_pid2rspfn.HandleRspFn) error {
	if !gInitData.CanUseCmd(cmd) {
		return fmt.Errorf("Cmd not allowed %v", cmd)
	}
	return app.ReqWithRspFn(cmd, body, fn)
}

func (app *WasmClient) reqAIPlay(onoff bool) error {
	return app.ReqWithRspFnWithAuth(
		c2t_idcmd.AIPlay,
		&c2t_obj.ReqAIPlay_data{onoff},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}

func (app *WasmClient) reqAchieveInfo() error {
	return app.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspAchieveInfo_data)
			app.systemMessage.Append(wrapspan.ColorText("Gold",
				"== Achievement == "))
			for i, v := range rpk.Achieve {
				app.systemMessage.Append(wrapspan.ColorTextf("Gold",
					"%v : %v ", achievetype.AchieveType(i).String(), v))
			}
			return nil
		},
	)
}

func (app *WasmClient) reqHeartbeat() error {
	return app.ReqWithRspFnWithAuth(
		c2t_idcmd.Heartbeat,
		&c2t_obj.ReqHeartbeat_data{
			Time: time.Now(),
		},
		func(hd c2t_packet.Header, rsp interface{}) error {
			rpk := rsp.(*c2t_obj.RspHeartbeat_data)
			pingDur := time.Now().Sub(rpk.Time)
			app.PingDur = (app.PingDur + pingDur) / 2
			return nil
		},
	)
}

func (app *WasmClient) sendPacket(cmd c2t_idcmd.CommandID, arg interface{}) {
	if cmd.NeedTurn() != 0 { // is act?
		atomic.AddInt32(&app.actPacketPerTurn, 1)
	}
	app.ReqWithRspFnWithAuth(
		cmd, arg,
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}

func (app *WasmClient) sendMovePacketByInput(tryDir way9type.Way9Type) bool {
	cf := app.currentFloor()
	playerX, playerY := app.GetPlayerXY()
	moveDir := cf.FindMovableDir(playerX, playerY, tryDir)
	if moveDir != way9type.Center {
		atomic.AddInt32(&app.movePacketPerTurn, 1)
		go app.sendPacket(c2t_idcmd.Move,
			&c2t_obj.ReqMove_data{Dir: moveDir},
		)
		return true
	} else if tryDir != way9type.Center {
		app.systemMessage.Appendf(
			"Cannot move to %v", tryDir.String())
	}
	return false
}
