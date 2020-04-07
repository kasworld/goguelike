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

package wasmclient

import (
	"github.com/kasworld/goguelike/game/wasmclient/htmlbutton"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/wrapspan"
)

var commandButtons = htmlbutton.NewButtonGroup("Commands",
	[]*htmlbutton.HTMLButton{
		{"a", "KillSelf", []string{"KillSelf"}, "Kill self", cmdKillSelf, 0},
		{"s", "ShowAchieve", []string{"ShowAchieve"}, "Show Achievement", cmdShowAchieve, 0},
		{"d", "EnterPortal", []string{"EnterPortal"}, "Enter portal", cmdEnterPortal, 0},
		{"f", "Teleport", []string{"Teleport"}, "Teleport random in floor", cmdTeleport, 0},
		{"g", "Rebirth", []string{"Rebirth"}, "Rebirth", cmdRebirth, 0},
	})

func cmdKillSelf(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.KillSelf,
		&c2t_obj.ReqKillSelf_data{},
	)
	Focus2Canvas()
}

func cmdShowAchieve(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.reqAchieveInfo()
	Focus2Canvas()
}

func cmdEnterPortal(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.EnterPortal,
		&c2t_obj.ReqEnterPortal_data{},
	)
	Focus2Canvas()
}

func cmdTeleport(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	go app.sendPacket(c2t_idcmd.ActTeleport,
		&c2t_obj.ReqActTeleport_data{},
	)
	Focus2Canvas()
}

func cmdRebirth(obj interface{}, v *htmlbutton.HTMLButton) {
	app, ok := obj.(*WasmClient)
	if !ok {
		jslog.Errorf("obj not app %v", obj)
		return
	}
	if app.olNotiData.ActiveObj.HP <= 0 {
		go app.sendPacket(c2t_idcmd.Rebirth,
			&c2t_obj.ReqRebirth_data{},
		)
	} else {
		app.systemMessage.Append(
			wrapspan.ColorText("OrangeRed",
				"no need to rebirth"))
	}
	Focus2Canvas()
}
