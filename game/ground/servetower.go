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

package ground

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/kasworld/goguelike/protocol_t2g/t2g_authorize"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_error"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_json"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_obj"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_packet"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_serveconnbyte"
	"github.com/kasworld/uuidstr"
)

func (grd *Ground) listenTowerTCP(ctx context.Context) {
	ListenRPCFromTowerPort := fmt.Sprintf(":%v", grd.sconfig.GroundRPCPort)
	tcpaddr, err := net.ResolveTCPAddr("tcp", ListenRPCFromTowerPort)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	defer listener.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			listener.SetDeadline(time.Now().Add(time.Duration(1 * time.Second)))
			conn, err := listener.AcceptTCP()
			if err != nil {
				operr, ok := err.(*net.OpError)
				if ok && operr.Timeout() {
					continue
				}
				fmt.Printf("error %#v\n", err)
			} else {
				go grd.serveTCPClient(ctx, conn)
			}
		}
	}
}

const (
	sendBufferSize  = 10
	readTimeoutDur  = 6 * time.Second
	writeTimeoutDur = 3 * time.Second
)

func (grd *Ground) serveTCPClient(ctx context.Context, conn *net.TCPConn) {

	connID := uuidstr.New()
	c2sc := t2g_serveconnbyte.New(
		connID,
		sendBufferSize,
		t2g_authorize.NewAllSet(),
		grd.SendStat, grd.RecvStat,
		grd.demuxReq2BytesAPIFnMap)

	// add to conn manager
	grd.connManager.Add(connID, c2sc)

	// start client service
	c2sc.StartServeTCP(ctx, conn,
		readTimeoutDur, writeTimeoutDur,
		t2g_json.MarshalBodyFn)
	conn.Close()
	// connection cleanup here

	// del from conn manager
	grd.connManager.Del(connID)
}

func (grd *Ground) bytesAPIFn_ReqInvalid(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error) {

	robj, err := t2g_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*t2g_obj.ReqInvalid_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody
	t2gc, ok := me.(*t2g_serveconnbyte.ServeConnByte)
	if !ok {
		return hd, nil, fmt.Errorf("me type miss match %v", me)
	}
	_ = t2gc
	sendBody := &t2g_obj.RspInvalid_data{}
	return hd, sendBody, fmt.Errorf("invalid packet")
}

func (grd *Ground) bytesAPIFn_ReqRegister(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error) {
	robj, err := t2g_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*t2g_obj.ReqRegister_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody
	t2gc, ok := me.(*t2g_serveconnbyte.ServeConnByte)
	if !ok {
		return hd, nil, fmt.Errorf("me type miss match %v", me)
	}
	_ = t2gc
	err = grd.twMan.Register(recvBody)
	hd.ErrorCode = t2g_error.None
	sendBody := &t2g_obj.RspRegister_data{}
	return hd, sendBody, err
}

func (grd *Ground) bytesAPIFn_ReqHeartbeat(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error) {
	robj, err := t2g_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*t2g_obj.ReqHeartbeat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody
	t2gc, ok := me.(*t2g_serveconnbyte.ServeConnByte)
	if !ok {
		return hd, nil, fmt.Errorf("me type miss match %v", me)
	}
	_ = t2gc
	err = grd.twMan.Heartbeat(recvBody)
	hd.ErrorCode = t2g_error.None
	sendBody := &t2g_obj.RspHeartbeat_data{}
	return hd, sendBody, err
}

func (grd *Ground) bytesAPIFn_ReqHighScore(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error) {
	robj, err := t2g_json.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*t2g_obj.ReqHighScore_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody
	t2gc, ok := me.(*t2g_serveconnbyte.ServeConnByte)
	if !ok {
		return hd, nil, fmt.Errorf("me type miss match %v", me)
	}
	_ = t2gc
	go grd.AddActiveObj2HighScoreAndSort(recvBody.ActiveObjScore)
	hd.ErrorCode = t2g_error.None
	sendBody := &t2g_obj.RspHighScore_data{}
	return hd, sendBody, nil
}
