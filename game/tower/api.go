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

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/game/activeobject"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/cmd2floor"
	"github.com/kasworld/goguelike/game/cmd2tower"
	"github.com/kasworld/goguelike/game/conndata"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_msgp"
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

	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
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

	if err := c2sc.GetAuthorCmdList().UpdateByAuthKey(robj.AuthKey); err != nil {
		return rhd, nil, err
	}

	connData := c2sc.GetConnData().(*conndata.ConnData)

	ss := tw.sessionManager.UpdateOrNew(
		robj.SessionG2ID,
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
			homeFloor = tw.GetFloorManager().GetRandomFloor()
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
			SessionG2ID:   g2id.NewFromString(connData.Session.GetUUID()),
			ActiveObjG2ID: g2id.NewFromString(connData.Session.ActiveObjUUID),
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

	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
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
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
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
	spacket := &c2t_obj.RspAchieveInfo_data{}

	for i, v := range ao.GetAchieveStat() {
		spacket.Achieve[i] = int64(v)
	}
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

func (tw *Tower) bytesAPIFn_ReqMeditate(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspMeditate_data{}

	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act: c2t_idcmd.Meditate,
	})
	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqKillSelf(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	spacket := &c2t_obj.RspKillSelf_data{}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act: c2t_idcmd.KillSelf,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqMove(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqMove_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	spacket := &c2t_obj.RspMove_data{}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act: c2t_idcmd.Move,
		Dir: robj.Dir,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqAttack(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAttack_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspAttack_data{}

	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act: c2t_idcmd.Attack,
		Dir: robj.Dir,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqPickup(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqPickup_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspPickup_data{}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act:  c2t_idcmd.Pickup,
		G2ID: robj.G2ID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqDrop(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqDrop_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspDrop_data{}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act:  c2t_idcmd.Drop,
		G2ID: robj.G2ID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqEquip(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqEquip_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspEquip_data{}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act:  c2t_idcmd.Equip,
		G2ID: robj.G2ID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqUnEquip(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqUnEquip_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspUnEquip_data{}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act:  c2t_idcmd.UnEquip,
		G2ID: robj.G2ID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqDrinkPotion(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqDrinkPotion_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspDrinkPotion_data{}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act:  c2t_idcmd.DrinkPotion,
		G2ID: robj.G2ID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqReadScroll(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqReadScroll_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspReadScroll_data{}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act:  c2t_idcmd.ReadScroll,
		G2ID: robj.G2ID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqRecycle(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqRecycle_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspRecycle_data{}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act:  c2t_idcmd.Recycle,
		G2ID: robj.G2ID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqEnterPortal(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	spacket := &c2t_obj.RspEnterPortal_data{}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act: c2t_idcmd.EnterPortal,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqMoveFloor(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
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
		FloorUUID: robj.G2ID.String(),
	}
	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqActTeleport(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	spacket := &c2t_obj.RspActTeleport_data{}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act: c2t_idcmd.ActTeleport,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqAdminTowerCmd(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAdminTowerCmd_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rspCh := make(chan c2t_error.ErrorCode, 1)
	tw.GetReqCh() <- &cmd2tower.AdminTowerCmd{
		ActiveObj:  ao,
		RecvPacket: robj,
		RspCh:      rspCh,
	}
	ec := <-rspCh
	return c2t_packet.Header{
		ErrorCode: ec,
	}, &c2t_obj.RspAdminTowerCmd_data{}, nil
}

func (tw *Tower) bytesAPIFn_ReqAdminFloorCmd(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAdminFloorCmd_data)
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
	spacket := &c2t_obj.RspAdminFloorCmd_data{}

	f := ao.GetCurrentFloor()
	if f == nil {
		return rhd, nil, fmt.Errorf("user not in floor %v", me)
	}
	rspCh := make(chan c2t_error.ErrorCode, 1)
	f.GetReqCh() <- &cmd2floor.APIAdminCmd2Floor{
		ActiveObj: ao,
		ReqPk:     robj,
		RspCh:     rspCh,
	}
	ec := <-rspCh
	return c2t_packet.Header{
		ErrorCode: ec,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqAdminActiveObjCmd(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAdminActiveObjCmd_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	return c2t_packet.Header{
		ErrorCode: ao.DoAdminCmd(robj.Cmd, robj.Arg),
	}, &c2t_obj.RspAdminActiveObjCmd_data{}, nil
}

func (tw *Tower) bytesAPIFn_ReqAdminFloorMove(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAdminFloorMove_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	rspCh := make(chan c2t_error.ErrorCode, 1)
	tw.GetReqCh() <- &cmd2tower.AdminFloorMove{
		ActiveObj:  ao,
		RecvPacket: robj,
		RspCh:      rspCh,
	}
	ec := <-rspCh
	return c2t_packet.Header{
		ErrorCode: ec,
	}, &c2t_obj.RspAdminFloorMove_data{}, nil
}

func (tw *Tower) bytesAPIFn_ReqAdminTeleport(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAdminTeleport_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}

	rhd := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	f := ao.GetCurrentFloor()
	if f == nil {
		return rhd, nil, fmt.Errorf("user not in floor %v", me)
	}
	rspCh := make(chan c2t_error.ErrorCode, 1)
	f.GetReqCh() <- &cmd2floor.APIAdminTeleport2Floor{
		ActiveObj: ao,
		ReqPk:     robj,
		RspCh:     rspCh,
	}
	ec := <-rspCh
	return c2t_packet.Header{
		ErrorCode: ec,
	}, &c2t_obj.RspAdminTeleport_data{}, nil
}

func (tw *Tower) bytesAPIFn_ReqAIPlay(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_msgp.UnmarshalPacket(hd, rbody)
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
