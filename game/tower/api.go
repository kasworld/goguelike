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

package tower

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/kasworld/goguelike/config/authdata"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/game/activeobject"
	"github.com/kasworld/goguelike/game/cmd2tower"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/conndata"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_serveconnbyte"
	"github.com/kasworld/version"
)

func (tw *Tower) bytesAPIFn_ReqInvalid(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	rhd := c2t_packet.Header{}
	spacket := &c2t_obj.RspInvalid_data{}
	return rhd, spacket, fmt.Errorf("invalid packet")
}

func (tw *Tower) bytesAPIFn_ReqLogin(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	c2sc, ok := me.(*c2t_serveconnbyte.ServeConnByte)
	if !ok {
		panic(fmt.Sprintf("invalid me not c2t_serveconnbyte.ServeConnByte %#v", me))
	}

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqLogin_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}

	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}

	if err := authdata.UpdateByAuthKey(c2sc.GetAuthorCmdList(), robj.AuthKey); err != nil {
		return rhd, nil, err
	}

	connData := c2sc.GetConnData().(*conndata.ConnData)

	ss := tw.sessionManager.UpdateOrNew(
		robj.SessionUUID,
		connData.RemoteAddr,
		robj.NickName)

	if oldc2sc := tw.connManager.Get(ss.ConnUUID); oldc2sc != nil {
		oldc2sc.Disconnect()
		// wait
		trycount := 10
		for tw.connManager.Get(ss.ConnUUID) != nil && trycount > 0 {
			runtime.Gosched()
			time.Sleep(time.Millisecond * 100)
			trycount--
		}
	}
	if tw.connManager.Get(ss.ConnUUID) != nil {
		tw.log.Fatal("old connection online %v", ss)
		return rhd, nil, err
	}
	ss.ConnUUID = connData.UUID
	connData.Session = ss

	oldAO, exist := tw.id2aoSuspend.GetByUUID(connData.Session.ActiveObjUUID)
	if exist {
		// connect to exist ao
		oldAO.Resume(c2sc)
		rspCh := make(chan error, 1)
		tw.GetReqCh() <- &cmd2tower.ActiveObjResumeTower{
			ActiveObj: oldAO,
			RspCh:     rspCh,
		}
		err = <-rspCh
	} else {
		// new ao
		var homeFloor gamei.FloorI
		if strings.HasPrefix(robj.NickName, "MC_") {
			homeFloor = tw.GetFloorManager().GetFloorByIndex(tw.rnd.Intn(tw.GetFloorManager().GetFloorCount()))
		} else {
			homeFloor = tw.GetFloorManager().GetStartFloor()
		}
		newAO := activeobject.NewUserActiveObj(
			homeFloor,
			connData.Session.NickName,
			tw.log,
			tw.towerAchieveStat,
			c2sc)
		connData.Session.ActiveObjUUID = newAO.GetUUID()
		rspCh := make(chan error, 1)
		tw.GetReqCh() <- &cmd2tower.ActiveObjEnterTower{
			ActiveObj: newAO,
			RspCh:     rspCh,
		}
		err = <-rspCh
	}

	if err != nil {
		return rhd, nil, err
	} else {
		acinfo := &c2t_obj.AccountInfo{
			SessionUUID:   connData.Session.GetUUID(),
			ActiveObjUUID: connData.Session.ActiveObjUUID,
			NickName:      connData.Session.NickName,
			CmdList:       *c2sc.GetAuthorCmdList(),
		}
		return rhd, &c2t_obj.RspLogin_data{
			ServiceInfo: tw.serviceInfo,
			AccountInfo: acinfo,
		}, nil
	}
}

func (tw *Tower) bytesAPIFn_ReqHeartbeat(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqHeartbeat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}

	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	spacket := &c2t_obj.RspHeartbeat_data{
		Time: robj.Time,
	}

	return rhd, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqChat(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqChat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	robj.Chat = strings.TrimSpace(robj.Chat)
	if !version.IsRelease() && len(robj.Chat) > 1 && robj.Chat[0] == '/' {
		return c2t_packet.Header{
			ErrorCode: AdminCmd(ao, robj.Chat[1:]),
		}, &c2t_obj.RspChat_data{}, nil
	} else {
		if len(robj.Chat) > gameconst.MaxChatLen {
			robj.Chat = robj.Chat[:gameconst.MaxChatLen]
		}
		ao.SetChat(robj.Chat)
		return rhd, &c2t_obj.RspChat_data{}, nil
	}
}

func (tw *Tower) bytesAPIFn_ReqAchieveInfo(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	spacket := &c2t_obj.RspAchieveInfo_data{
		AchieveStat:   *ao.GetAchieveStat(),
		PotionStat:    *ao.GetPotionStat(),
		ScrollStat:    *ao.GetScrollStat(),
		FOActStat:     *ao.GetFieldObjActStat(),
		AOActionStat:  *ao.GetActStat(),
		ConditionStat: *ao.GetConditionStat(),
	}

	// for i, v := range ao.GetAchieveStat() {
	// 	spacket.Achieve[i] = int64(v)
	// }
	return rhd, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqRebirth(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	spacket := &c2t_obj.RspRebirth_data{}

	err = ao.TryRebirth()
	if err != nil {
		rhd.ErrorCode = c2t_error.ActionProhibited
		tw.log.Error("%v", err)
	}

	return rhd, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqMoveFloor(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqMoveFloor_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspMoveFloor_data{}

	tw.GetReqCh() <- &cmd2tower.FloorMove{
		ActiveObj: ao,
		FloorName: robj.UUID,
	}
	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqAIPlay(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAIPlay_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	if err := ao.DoAIOnOff(robj.On); err != nil {
		tw.log.Error("fail to AIOn %v %v", me)
	}
	return rhd, &c2t_obj.RspAIPlay_data{}, nil
}
