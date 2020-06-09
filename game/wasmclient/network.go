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
	"syscall/js"

	"github.com/kasworld/goguelike/game/clientcookie"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/lib/jsobj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_connwasm"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_msgp"
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
		c2t_msgp.MarshalBodyFn,
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
	sessionkey := g2id.G2ID(ck[clientcookie.SessionKeyName(gInitData.TowerIndex)])
	wg.Wait()

	wg.Add(1)
	app.ReqWithRspFn(
		c2t_idcmd.Login,
		&c2t_obj.ReqLogin_data{
			SessionG2ID: sessionkey,
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
	robj, err := c2t_msgp.UnmarshalPacket(header, body)
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

// implement C2SConnectI
func (app *WasmClient) CheckAPI(hd c2t_packet.Header) error {
	if hd.ErrorCode == c2t_error.None {
		return nil
	}
	return hd.ErrorCode
}
