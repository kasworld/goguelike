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

package activeobject

import (
	"fmt"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/game/carryingobject"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
)

func (ao *ActiveObject) DoAdminCmd(cmd, args string) c2t_error.ErrorCode {
	ao.log.AdminAudit("$c %v %v", ao, cmd, args)

	ec := c2t_error.None

	switch cmd {
	default:
		ao.log.Error("unknown admin ao cmd %v %v", cmd, args)
		ec = c2t_error.ActionProhibited

		// case "IncExp":
		// 	f64, err := strconv.ParseFloat(args, 64)
		// 	if err != nil {
		// 		ao.log.Error("invalid IncExp arg %v %v", cmd, args)
		// 		ec = c2t_error.ActionProhibited
		// 	} else {
		// 		ao.AddBattleExp(f64)
		// 	}

		// case "Potion":
		// 	pt, exist := potiontype.String2PotionType(args)
		// 	if exist {
		// 		ao.buffManager.Add(pt.String(), false, false,
		// 			potiontype.GetBuffByPotionType(pt),
		// 		)
		// 	} else {
		// 		ao.log.Error("unknown admin potion arg %v %v", cmd, args)
		// 		ec = c2t_error.ActionProhibited
		// 	}

		// case "Scroll":
		// 	pt, exist := scrolltype.String2ScrollType(args)
		// 	if exist {
		// 		ao.buffManager.Add(pt.String(), false, false,
		// 			scrolltype.GetBuffByScrollType(pt),
		// 		)
		// 	} else {
		// 		ao.log.Error("unknown admin potion arg %v %v", cmd, args)
		// 		ec = c2t_error.ActionProhibited
		// 	}

		// case "Condition":
		// 	cond, exist := condition.String2Condition(args)
		// 	if exist {
		// 		buff2add := statusoptype.Repeat(100,
		// 			statusoptype.OpArg{statusoptype.SetCondition, cond},
		// 		)
		// 		ao.buffManager.Add("admin"+args, true, true, buff2add)
		// 	}
	}

	return ec
}

func (ao *ActiveObject) DoParsedAdminCmd(ca *scriptparse.CmdArgs) error {
	ao.log.AdminAudit("%v %v", ao, ca)
	switch ca.Cmd {
	default:
		return fmt.Errorf("unknown cmd %v", ca.Cmd)
	case "AddPotion":
		ao.addRandPotion(1)
	case "AddScroll":
		ao.addRandScroll(1)
	case "AddMoney":
		v := ao.rnd.NormFloat64Range(gameconst.InitGoldMean, gameconst.InitGoldMean/2)
		if v < 1 {
			v = 1
		}
		ao.inven.AddToWallet(carryingobject.NewMoney(v))
		ao.achieveStat.Add(achievetype.MoneyGet, float64(v))
	case "AddEquip":
		ao.addRandCarryObjEquip(ao.nickName, 1)

	case "ForgetFloor":
		err := ao.forgetFloorByUUID(ao.currrentFloor.GetUUID())
		if err != nil {
			ao.log.Error("%v", err)
		}

	case "GetFloorMap":
		err := ao.makeFloorComplete(ao.currrentFloor)
		if err != nil {
			ao.log.Error("%v", err)
		}
	}
	return nil
}
