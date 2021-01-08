// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tower

import (
	"context"
	"fmt"
	"time"

	"github.com/kasworld/goguelike/protocol_t2g/t2g_conntcp"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_error"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_idcmd"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_json"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_packet"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_pid2rspfn"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_statapierror"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_statcallapi"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_statnoti"
)

// service const
const (
	// for client
	readTimeoutSec  = 6 * time.Second
	writeTimeoutSec = 3 * time.Second
)

type Conn2Ground struct {
	addr         string
	c2sc         *t2g_conntcp.Connection
	sendRecvStop func()
	apistat      *t2g_statcallapi.StatCallAPI
	pid2statobj  *t2g_statcallapi.PacketID2StatObj
	notistat     *t2g_statnoti.StatNotification
	errstat      *t2g_statapierror.StatAPIError
	pid2recv     *t2g_pid2rspfn.PID2RspFn
	connected    bool
}

func NewConn2Ground(addr string) *Conn2Ground {
	c2g := &Conn2Ground{
		addr:        addr,
		apistat:     t2g_statcallapi.New(),
		pid2statobj: t2g_statcallapi.NewPacketID2StatObj(),
		notistat:    t2g_statnoti.New(),
		errstat:     t2g_statapierror.New(),
		pid2recv:    t2g_pid2rspfn.New(),
	}
	c2g.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call\n")
	}
	c2g.c2sc = t2g_conntcp.New(10)
	return c2g
}

func (c2g *Conn2Ground) IsConnected() bool {
	return c2g.connected
}

func (c2g *Conn2Ground) Run(mainctx context.Context) error {
	if err := c2g.c2sc.ConnectTo(c2g.addr); err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	ctx, stopFn := context.WithCancel(mainctx)
	c2g.sendRecvStop = stopFn
	defer c2g.sendRecvStop()
	c2g.connected = true
	err := c2g.c2sc.Run(ctx,
		readTimeoutSec, writeTimeoutSec,
		t2g_json.MarshalBodyFn,
		c2g.handleRecvPacket,
		c2g.handleSentPacket,
	)
	c2g.connected = false
	return err
}

func (c2g *Conn2Ground) handleSentPacket(pk *t2g_packet.Packet) error {
	if err := c2g.apistat.AfterSendReq(pk.Header); err != nil {
		return err
	}
	return nil
}

func (c2g *Conn2Ground) handleRecvPacket(header t2g_packet.Header, body []byte) error {
	switch header.FlowType {
	default:
		return fmt.Errorf("Invalid packet type %v %v", header, body)
	case t2g_packet.Notification:
		// noti stat
		c2g.notistat.Add(header)
		//process noti here
		// robj, err := t2g_json.UnmarshalPacket(header, body)

	case t2g_packet.Response:
		// error stat
		c2g.errstat.Inc(t2g_idcmd.CommandID(header.Cmd), header.ErrorCode)
		// api stat
		if err := c2g.apistat.AfterRecvRsp(header); err != nil {
			fmt.Printf("%v %v\n", c2g, err)
			return err
		}
		psobj := c2g.pid2statobj.Get(header.ID)
		if psobj == nil {
			return fmt.Errorf("no statobj for %v", header.ID)
		}
		psobj.CallServerEnd(header.ErrorCode == t2g_error.None)
		c2g.pid2statobj.Del(header.ID)

		// process response
		if err := c2g.pid2recv.HandleRsp(header, body); err != nil {
			return err
		}

	}
	return nil
}

func (c2g *Conn2Ground) ReqWithRspFn(cmd t2g_idcmd.CommandID, body interface{},
	fn t2g_pid2rspfn.HandleRspFn) error {

	pid := c2g.pid2recv.NewPID(fn)
	spk := t2g_packet.Packet{
		Header: t2g_packet.Header{
			Cmd:      uint16(cmd),
			ID:       pid,
			FlowType: t2g_packet.Request,
		},
		Body: body,
	}

	// add api stat
	psobj, err := c2g.apistat.BeforeSendReq(spk.Header)
	if err != nil {
		return nil
	}
	c2g.pid2statobj.Add(spk.Header.ID, psobj)

	if err := c2g.c2sc.EnqueueSendPacket(&spk); err != nil {
		fmt.Printf("End %v %v %v\n", c2g, spk, err)
		c2g.sendRecvStop()
		return fmt.Errorf("Send fail %v %v", c2g, err)
	}
	return nil
}
