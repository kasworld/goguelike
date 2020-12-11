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

package clientai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kasworld/actjitter"
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_connwsgorilla"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_pid2rspfn"
)

type ClientAI struct {
	log          *g2log.LogBase `prettystring:"hide"`
	sendRecvStop func()         `prettystring:"hide"`

	config    ClientAIConfig
	runResult error

	towerConn         *c2t_connwsgorilla.Connection
	ServiceInfo       *c2t_obj.ServiceInfo
	AccountInfo       *c2t_obj.AccountInfo
	TowerInfo         *c2t_obj.TowerInfo
	ViewportXYLenList findnear.XYLenList
	FloorInfoList     []*c2t_obj.FloorInfo
	CurrentFloor      *clientfloor.ClientFloor

	wg       *sync.WaitGroup
	pid2recv *c2t_pid2rspfn.PID2RspFn

	// turn data
	movePacketPerTurn     int32
	OLNotiData            *c2t_obj.NotiVPObjList_data
	playerActiveObjClient *c2t_obj.ActiveObjClient
	onFieldObj            *c2t_obj.FieldObjClient
	IsOverLoad            bool
	HPdiff                int
	SPdiff                int
	level                 int
	ServerClientTimeDiff  time.Duration
	ServerJitter          *actjitter.ActJitter
}

func New(config ClientAIConfig, l *g2log.LogBase) *ClientAI {
	cai := &ClientAI{
		config:            config,
		log:               l,
		ServerJitter:      actjitter.New("Server"),
		wg:                new(sync.WaitGroup),
		pid2recv:          c2t_pid2rspfn.New(),
		ViewportXYLenList: viewportdata.ViewportXYLenList,
	}
	cai.sendRecvStop = func() {
		cai.log.Error("Too early sendRecvStop call %v", cai)
	}
	cai.towerConn = c2t_connwsgorilla.New(10)
	return cai
}

func (cai *ClientAI) Cleanup() {
	cai.wg.Wait()
	if tc := cai.towerConn; tc != nil {
		tc.Cleanup()
	}
	cai.ServerJitter = nil
}

func (cai *ClientAI) Run(mainctx context.Context) {
	defer cai.Cleanup()

	ctx, closeCtx := context.WithCancel(mainctx)
	cai.sendRecvStop = closeCtx
	defer cai.sendRecvStop()

	if err := cai.towerConn.ConnectTo(cai.config.ConnectToTower); err != nil {
		cai.runResult = err
		cai.log.Error("%v", cai.runResult)
		return
	}
	cai.wg.Add(1)
	go func() {
		defer cai.wg.Done()
		err := cai.towerConn.Run(ctx,
			gameconst.ClientReadTimeoutSec*time.Second,
			gameconst.ClientWriteTimeoutSec*time.Second,
			c2t_gob.MarshalBodyFn,
			cai.handleRecvPacket,
			cai.handleSentPacket,
		)

		if err != nil {
			cai.runResult = err
			cai.log.Error("%v %v", cai, err)
		}
		cai.sendRecvStop()
	}()

	if err := cai.reqLogin(
		cai.config.SessionUUID,
		cai.config.Nickname,
		cai.config.Auth,
	); err != nil {
		cai.runResult = err
		cai.log.Error("%v", cai.runResult)
		return
	}

	timerPingTk := time.NewTicker(time.Second * gameconst.ServerPacketReadTimeOutSec / 2)
	defer timerPingTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerPingTk.C:
			cai.wg.Add(1)
			go func() {
				defer cai.wg.Done()
				err := cai.reqHeartbeat()
				if err != nil {
					cai.runResult = err
					cai.log.Error("%v", cai.runResult)
				}
			}()
		}
	}
}

func (cai *ClientAI) handleSentPacket(pk *c2t_packet.Packet) error {
	cai.log.TraceClient("sent %v", pk.Header)
	return nil
}

func (cai *ClientAI) handleRecvPacket(header c2t_packet.Header, body []byte) error {
	cai.log.TraceClient("recv %v", header)
	switch header.FlowType {
	default:
		return fmt.Errorf("Invalid packet type %v %v", header, body)
	case c2t_packet.Response:
		if err := cai.pid2recv.HandleRsp(header, body); err != nil {
			cai.sendRecvStop()
			cai.log.Fatal("%v %v %v %v", cai, header, body, err)
			return err
		}
	case c2t_packet.Notification:
		fn := DemuxNoti2ByteFnMap[header.Cmd]
		if err := fn(cai, header, body); err != nil {
			cai.sendRecvStop()
			cai.log.Fatal("%v %v %v %v", cai, header, body, err)
			return err
		}
	}
	return nil
}

func (cai *ClientAI) ReqWithRspFn(cmd c2t_idcmd.CommandID, body interface{},
	fn c2t_pid2rspfn.HandleRspFn) error {

	pid := cai.pid2recv.NewPID(fn)
	spk := c2t_packet.Packet{
		Header: c2t_packet.Header{
			Cmd:      uint16(cmd),
			ID:       pid,
			FlowType: c2t_packet.Request,
		},
		Body: body,
	}
	if err := cai.towerConn.EnqueueSendPacket(&spk); err != nil {
		cai.log.Error("End %s %v %+v %v",
			cai, spk.Header, spk.Body, err)
		cai.sendRecvStop()
		return fmt.Errorf("Send fail %s %v:%v %v",
			cai, cmd, pid, err)
	}
	return nil
}

func (cai *ClientAI) ReqWithRspFnWithAuth(cmd c2t_idcmd.CommandID, body interface{},
	fn c2t_pid2rspfn.HandleRspFn) error {
	if !cai.CanUseCmd(cmd) {
		return fmt.Errorf("Cmd not allowed %v", cmd)
	}
	return cai.ReqWithRspFn(cmd, body, fn)
}
