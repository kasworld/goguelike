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

	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
)

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
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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

func (tw *Tower) bytesAPIFn_ReqAttackWide(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAttackWide_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspAttackWide_data{}

	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act: c2t_idcmd.AttackWide,
		Dir: robj.Dir,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqAttackLong(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAttackLong_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	spacket := &c2t_obj.RspAttackLong_data{}

	ao.SetReq2Handle(&aoactreqrsp.Act{
		Act: c2t_idcmd.AttackLong,
		Dir: robj.Dir,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqPickup(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
		UUID: robj.UUID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqDrop(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
		UUID: robj.UUID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqEquip(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
		UUID: robj.UUID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqUnEquip(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
		UUID: robj.UUID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqDrinkPotion(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
		UUID: robj.UUID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqReadScroll(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
		UUID: robj.UUID,
	})

	return c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}, spacket, nil
}

func (tw *Tower) bytesAPIFn_ReqRecycle(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
		UUID: robj.UUID,
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
