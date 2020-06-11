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

package wasmclientgl

import (
	"github.com/kasworld/goguelike/lib/htmlbutton"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

var adminCmds = []c2t_idcmd.CommandID{
	c2t_idcmd.AdminTeleport,
	c2t_idcmd.AdminActiveObjCmd,
	c2t_idcmd.AdminActiveObjCmd,
	c2t_idcmd.AdminActiveObjCmd,
	c2t_idcmd.AdminFloorMove,
	c2t_idcmd.AdminFloorMove,
	c2t_idcmd.AdminActiveObjCmd,
}
var adminCommandButtons = htmlbutton.NewButtonGroup(" ",
	[]*htmlbutton.HTMLButton{
		htmlbutton.New("1", "AdminTeleport", []string{"Teleport"}, "teleport random", adminCmdTeleport, 0),
		htmlbutton.New("2", "IncExp", []string{"IncExp"}, "inc exp", adminCmdIncExp, 0),
		htmlbutton.New("3", "FloorBefore", []string{"FloorBefore"}, "before floor", adminCmdFloorBefore, 0),
		htmlbutton.New("4", "FloorNext", []string{"FloorNext"}, "next floor", adminCmdFloorNext, 0),
		htmlbutton.New("5", "Potion", []string{"Potion"}, "apply potion", adminCmdPotion, 0),
		htmlbutton.New("7", "Scroll", []string{"Scroll"}, "apply scroll", adminCmdScroll, 0),
		htmlbutton.New("8", "Condition", []string{"Condition"}, "set contidion", adminCmdCondition, 0),
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
	app.Focus2Canvas()
}
func adminCmdIncExp(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminActiveObjCmd,
		&c2t_obj.ReqAdminActiveObjCmd_data{
			Cmd: "IncExp",
			Arg: getChatMsg(),
		},
	)
	app.Focus2Canvas()
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
	app.Focus2Canvas()
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
	app.Focus2Canvas()
}

func adminCmdPotion(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminActiveObjCmd,
		&c2t_obj.ReqAdminActiveObjCmd_data{
			Cmd: "Potion",
			Arg: getChatMsg(),
		},
	)
	app.Focus2Canvas()
}

func adminCmdScroll(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminActiveObjCmd,
		&c2t_obj.ReqAdminActiveObjCmd_data{
			Cmd: "Scroll",
			Arg: getChatMsg(),
		},
	)
	app.Focus2Canvas()
}

func adminCmdCondition(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.AdminActiveObjCmd,
		&c2t_obj.ReqAdminActiveObjCmd_data{
			Cmd: "Condition",
			Arg: getChatMsg(),
		},
	)
	app.Focus2Canvas()
}
