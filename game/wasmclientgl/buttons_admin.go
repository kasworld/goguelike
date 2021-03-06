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

package wasmclientgl

import (
	"strconv"

	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/lib/htmlbutton"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

var adminCmds = []c2t_idcmd.CommandID{
	c2t_idcmd.AdminTeleport,
	c2t_idcmd.AdminFloorMove,
	c2t_idcmd.AdminFloorMove,
	c2t_idcmd.AdminForgetFloor,
	c2t_idcmd.AdminFloorMap,
	c2t_idcmd.AdminAddExp,
	c2t_idcmd.AdminPotionEffect,
	c2t_idcmd.AdminScrollEffect,
	c2t_idcmd.AdminCondition,
	c2t_idcmd.AdminAddMoney,
	c2t_idcmd.AdminAddPotion,
	c2t_idcmd.AdminAddScroll,
	c2t_idcmd.AdminAddEquip,
}

var adminCommandButtons = htmlbutton.NewButtonGroup(" ",
	[]*htmlbutton.HTMLButton{
		htmlbutton.New("1", "AdminTeleport", []string{"Teleport"}, "teleport random", adminCmdTeleport, 0),
		htmlbutton.New("2", "FloorBefore", []string{"FloorBefore"}, "before floor", adminCmdFloorBefore, 0),
		htmlbutton.New("3", "FloorNext", []string{"FloorNext"}, "next floor", adminCmdFloorNext, 0),
		htmlbutton.New("4", "FloorForget", []string{"FloorForget"}, "forget floor", adminCmdForgetFloor, 0),
		htmlbutton.New("5", "FloorMap", []string{"FloorMap"}, "floor map complete", adminCmdFloorMap, 0),
		htmlbutton.New("6", "AddExp", []string{"AddExp"}, "+/- exp", adminCmdAddExp, 0),
		htmlbutton.New("7", "PotionBuff", []string{"PotionBuff"}, "apply potion buff", adminCmdPotionEffect, 0),
		htmlbutton.New("8", "ScrollBuff", []string{"ScrollBuff"}, "apply scroll buff", adminCmdScrollEffect, 0),
		htmlbutton.New("9", "Condition", []string{"Condition"}, "set contidion", adminCmdCondition, 0),
		htmlbutton.New("0", "AddMoney", []string{"AddMoney"}, "add money", adminCmdAddMoney, 0),
		htmlbutton.New("-", "Potion", []string{"Potion"}, "add potion", adminCmdAddPotion, 0),
		htmlbutton.New("=", "Scroll", []string{"Scroll"}, "add scroll", adminCmdAddScroll, 0),
		htmlbutton.New("", "Equip", []string{"Equip"}, "add equip by slot", adminCmdAddEquip, 0),
	})

func adminCmdTeleport(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminTeleport,
		&c2t_obj.ReqAdminTeleport_data{
			X: 0,
			Y: 0,
		},
	)
	v.Blur()
}
func adminCmdFloorBefore(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminFloorMove,
		&c2t_obj.ReqAdminFloorMove_data{
			Floor: "Before",
		},
	)
	v.Blur()
}
func adminCmdFloorNext(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminFloorMove,
		&c2t_obj.ReqAdminFloorMove_data{
			Floor: "Next",
		},
	)
	v.Blur()
}

func adminCmdForgetFloor(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminForgetFloor,
		&c2t_obj.ReqAdminForgetFloor_data{},
	)
	v.Blur()
}

func adminCmdFloorMap(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminFloorMap,
		&c2t_obj.ReqAdminFloorMap_data{},
	)
	v.Blur()
}

func adminCmdAddExp(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	expstr := getChatMsg()
	i64, err := strconv.ParseInt(expstr, 0, 64)
	if err != nil {
		jslog.Errorf("invalid AddExp %v", expstr)
	} else {
		go app.sendPacket(c2t_idcmd.AdminAddExp,
			&c2t_obj.ReqAdminAddExp_data{
				Exp: int(i64),
			},
		)
	}
	v.Blur()
}

func adminCmdPotionEffect(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	args := getChatMsg()
	pt, exist := potiontype.String2PotionType(args)
	if exist {
		go app.sendPacket(c2t_idcmd.AdminPotionEffect,
			&c2t_obj.ReqAdminPotionEffect_data{
				Potion: pt,
			},
		)
	} else {
		jslog.Errorf("unknown admin potion %v", args)
	}
	v.Blur()
}

func adminCmdScrollEffect(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	args := getChatMsg()
	pt, exist := scrolltype.String2ScrollType(args)
	if exist {
		go app.sendPacket(c2t_idcmd.AdminScrollEffect,
			&c2t_obj.ReqAdminScrollEffect_data{
				Scroll: pt,
			},
		)
	} else {
		jslog.Errorf("unknown admin scroll %v", args)
	}

	v.Blur()
}

func adminCmdCondition(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	args := getChatMsg()
	cond, exist := condition.String2Condition(args)
	if exist {
		go app.sendPacket(c2t_idcmd.AdminCondition,
			&c2t_obj.ReqAdminCondition_data{
				Condition: cond,
			},
		)
	} else {
		jslog.Errorf("unknown admin condition %v", args)
	}

	v.Blur()
}

func adminCmdAddMoney(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	expstr := getChatMsg()
	i64, err := strconv.ParseInt(expstr, 0, 64)
	if err != nil {
		jslog.Errorf("invalid AddMoney %v", expstr)
	} else {
		go app.sendPacket(c2t_idcmd.AdminAddMoney,
			&c2t_obj.ReqAdminAddMoney_data{
				Money: int(i64),
			},
		)
	}
	v.Blur()
}

func adminCmdAddPotion(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	args := getChatMsg()
	pt, exist := potiontype.String2PotionType(args)
	if exist {
		go app.sendPacket(c2t_idcmd.AdminAddPotion,
			&c2t_obj.ReqAdminAddPotion_data{
				Potion: pt,
			},
		)
	} else {
		jslog.Errorf("unknown admin potion %v", args)
	}
	v.Blur()
}

func adminCmdAddScroll(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	args := getChatMsg()
	pt, exist := scrolltype.String2ScrollType(args)
	if exist {
		go app.sendPacket(c2t_idcmd.AdminAddScroll,
			&c2t_obj.ReqAdminAddScroll_data{
				Scroll: pt,
			},
		)
	} else {
		jslog.Errorf("unknown admin scroll %v", args)
	}

	v.Blur()
}

func adminCmdAddEquip(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	ft := app.olNotiData.ActiveObj.Bias.NearFaction()
	args := getChatMsg()
	pt, exist := equipslottype.String2EquipSlotType(args)
	if exist {
		go app.sendPacket(c2t_idcmd.AdminAddEquip,
			&c2t_obj.ReqAdminAddEquip_data{
				Faction: ft,
				Equip:   pt,
			},
		)
	} else {
		jslog.Errorf("unknown admin scroll %v", args)
	}
	v.Blur()
}
