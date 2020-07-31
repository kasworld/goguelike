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

	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
)

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
