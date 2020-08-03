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

	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/statusoptype"
	"github.com/kasworld/goguelike/game/cmd2floor"
	"github.com/kasworld/goguelike/game/cmd2tower"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
)

// from ReqChat
func AdminCmd(ao gamei.ActiveObjectI, cmdline string) c2t_error.ErrorCode {

	cmd, argLine := scriptparse.SplitCmdArgstr(cmdline, " ")
	if len(cmd) == 0 {
		return c2t_error.ActionCanceled
	}
	fnarg, exist := cmdFnArgs[cmd]
	if !exist {
		return c2t_error.ActionCanceled
	}
	_, name2value, err := scriptparse.Split2ListMap(argLine, " ", "=")
	if err != nil {
		return c2t_error.ActionCanceled
	}

	nameList, name2type, err := scriptparse.Split2ListMap(fnarg.format, " ", ":")
	if err != nil {
		return c2t_error.ActionCanceled
	}
	ca := &scriptparse.CmdArgs{
		Cmd:        cmd,
		Name2Value: name2value,
		NameList:   nameList,
		Name2Type:  name2type,
	}
	if err := fnarg.fn(ao, ca); err != nil {
		return c2t_error.ActionProhibited
	}
	return c2t_error.None
}

var cmdFnArgs = map[string]struct {
	fn     func(ao gamei.ActiveObjectI, ca *scriptparse.CmdArgs) error
	format string
}{
	"AddPotion":   {adminActiveObjCmd, ""},
	"AddScroll":   {adminActiveObjCmd, ""},
	"AddMoney":    {adminActiveObjCmd, ""},
	"AddEquip":    {adminActiveObjCmd, ""},
	"GetFloorMap": {adminActiveObjCmd, ""},
	"ForgetFloor": {adminActiveObjCmd, ""},
}

func init() {
	// verify format
	for _, v := range cmdFnArgs {
		_, _, err := scriptparse.Split2ListMap(v.format, " ", ":")
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
	}
}

func adminActiveObjCmd(ao gamei.ActiveObjectI, ca *scriptparse.CmdArgs) error {
	return ao.DoParsedAdminCmd(ca)
}

func (tw *Tower) bytesAPIFn_ReqAdminTowerCmd(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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
	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	spacket := &c2t_obj.RspAdminFloorCmd_data{}

	f := ao.GetCurrentFloor()
	if f == nil {
		return sendHeader, nil, fmt.Errorf("user not in floor %v", me)
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

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
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

	r, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	robj, ok := r.(*c2t_obj.ReqAdminTeleport_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", r)
	}

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return hd, nil, err
	}
	f := ao.GetCurrentFloor()
	if f == nil {
		return sendHeader, nil, fmt.Errorf("user not in floor %v", me)
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

// AdminAddExp  add arg to battle exp
func (tw *Tower) bytesAPIFn_ReqAdminAddExp(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminAddExp_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return sendHeader, nil, err
	}
	ao.AddBattleExp(float64(recvBody.Exp))
	sendBody := &c2t_obj.RspAdminAddExp_data{}
	return sendHeader, sendBody, nil
}

// AdminPotionEffect buff by arg potion type
func (tw *Tower) bytesAPIFn_ReqAdminPotionEffect(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminPotionEffect_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return sendHeader, nil, err
	}
	ao.GetBuffManager().Add(recvBody.Potion.String(), false, false,
		potiontype.GetBuffByPotionType(recvBody.Potion),
	)

	sendBody := &c2t_obj.RspAdminPotionEffect_data{}
	return sendHeader, sendBody, nil
}

// AdminScrollEffect buff by arg Scroll type
func (tw *Tower) bytesAPIFn_ReqAdminScrollEffect(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminScrollEffect_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return sendHeader, nil, err
	}
	ao.GetBuffManager().Add(recvBody.Scroll.String(), false, false,
		scrolltype.GetBuffByScrollType(recvBody.Scroll),
	)
	sendBody := &c2t_obj.RspAdminScrollEffect_data{}
	return sendHeader, sendBody, nil
}

// AdminCondition add arg condition for 100 turn
func (tw *Tower) bytesAPIFn_ReqAdminCondition(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminCondition_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	ao, err := tw.api_me2ao(me)
	if err != nil {
		return sendHeader, nil, err
	}
	buff2add := statusoptype.Repeat(100,
		statusoptype.OpArg{statusoptype.SetCondition, recvBody.Condition},
	)
	ao.GetBuffManager().Add(
		"admin"+recvBody.Condition.String(),
		true, true, buff2add)

	sendBody := &c2t_obj.RspAdminCondition_data{}
	return sendHeader, sendBody, nil
}

// AdminAddPotion add arg potion to inven
func (tw *Tower) bytesAPIFn_ReqAdminAddPotion(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminAddPotion_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddPotion_data{}
	return sendHeader, sendBody, nil
}

// AdminAddScroll add arg scroll to inven
func (tw *Tower) bytesAPIFn_ReqAdminAddScroll(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminAddScroll_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddScroll_data{}
	return sendHeader, sendBody, nil
}

// AdminAddMoney add arg money to inven
func (tw *Tower) bytesAPIFn_ReqAdminAddMoney(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminAddMoney_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddMoney_data{}
	return sendHeader, sendBody, nil
}

// AdminAddEquip add random equip to inven
func (tw *Tower) bytesAPIFn_ReqAdminAddEquip(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminAddEquip_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminAddEquip_data{}
	return sendHeader, sendBody, nil
}

// AdminForgetFloor forget current floor map
func (tw *Tower) bytesAPIFn_ReqAdminForgetFloor(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminForgetFloor_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminForgetFloor_data{}
	return sendHeader, sendBody, nil
}

// AdminFloorMap complete current floor map
func (tw *Tower) bytesAPIFn_ReqAdminFloorMap(
	me interface{}, hd c2t_packet.Header, rbody []byte) (
	c2t_packet.Header, interface{}, error) {
	robj, err := c2t_gob.UnmarshalPacket(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2t_obj.ReqAdminFloorMap_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	sendHeader := c2t_packet.Header{
		ErrorCode: c2t_error.None,
	}
	sendBody := &c2t_obj.RspAdminFloorMap_data{}
	return sendHeader, sendBody, nil
}
