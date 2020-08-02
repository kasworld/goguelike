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
