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

package authdata

import (
	"fmt"

	"github.com/kasworld/goguelike/protocol_c2t/c2t_authorize"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

func AddAdminKey(key string) error {
	var err error
	if _, exist := Authkey2Admin[key]; exist {
		err = fmt.Errorf("key %v exist, overwright", key)
	}
	Authkey2Admin[key] = [2][]string{
		[]string{"Login", "Admin"}, []string{"DelAfterLogin"},
	}
	return err
}

var allAuthorizationSet = map[string]*c2t_authorize.AuthorizedCmds{
	"PreLogin": c2t_authorize.NewByCmdIDList([]c2t_idcmd.CommandID{
		c2t_idcmd.Login,
	}),

	"DelAfterLogin": c2t_authorize.NewByCmdIDList([]c2t_idcmd.CommandID{
		c2t_idcmd.Login,
	}),

	"Login": c2t_authorize.NewByCmdIDList([]c2t_idcmd.CommandID{
		c2t_idcmd.Heartbeat,
		c2t_idcmd.Chat,
		c2t_idcmd.AchieveInfo,
		c2t_idcmd.Rebirth,
		c2t_idcmd.VisitFloorList,
		c2t_idcmd.Meditate,
		c2t_idcmd.KillSelf,
		c2t_idcmd.Move,
		c2t_idcmd.Attack,
		c2t_idcmd.AttackWide,
		c2t_idcmd.AttackLong,
		c2t_idcmd.Pickup,
		c2t_idcmd.Drop,
		c2t_idcmd.Equip,
		c2t_idcmd.UnEquip,
		c2t_idcmd.DrinkPotion,
		c2t_idcmd.ReadScroll,
		c2t_idcmd.Recycle,
		c2t_idcmd.EnterPortal,
		c2t_idcmd.MoveFloor,
		c2t_idcmd.ActTeleport,

		c2t_idcmd.AIPlay,
	}),
	"Admin": c2t_authorize.NewByCmdIDList([]c2t_idcmd.CommandID{
		c2t_idcmd.AdminTowerCmd,
		c2t_idcmd.AdminFloorCmd,
		c2t_idcmd.AdminActiveObjCmd,
		c2t_idcmd.AdminFloorMove,
		c2t_idcmd.AdminTeleport,
		c2t_idcmd.AdminAddExp,
		c2t_idcmd.AdminPotionEffect,
		c2t_idcmd.AdminScrollEffect,
		c2t_idcmd.AdminCondition,
		c2t_idcmd.AdminAddPotion,
		c2t_idcmd.AdminAddScroll,
		c2t_idcmd.AdminAddMoney,
		c2t_idcmd.AdminAddEquip,
		c2t_idcmd.AdminForgetFloor,
		c2t_idcmd.AdminFloorMap,
	}),
}

func NewPreLoginAuthorCmdIDList() *c2t_authorize.AuthorizedCmds {
	return allAuthorizationSet["PreLogin"].Duplicate()
}

func UpdateByAuthKey(acicl *c2t_authorize.AuthorizedCmds, key string) error {
	ag, exist := Authkey2Admin[key]
	if !exist {
		ag = [2][]string{[]string{"Login"}, []string{"DelAfterLogin"}}
	}
	// process include
	for _, authgroupname := range ag[0] {
		cmdidList := allAuthorizationSet[authgroupname]
		if cmdidList == nil {
			return fmt.Errorf("Can't Found authgroup %v", authgroupname)
		}
		acicl.Union(cmdidList)
	}
	// process exclude
	for _, authgroupname := range ag[1] {
		cmdidList := allAuthorizationSet[authgroupname]
		if cmdidList == nil {
			return fmt.Errorf("Can't Found authgroup %v", authgroupname)
		}
		acicl.SubIntersection(cmdidList)
	}
	return nil
}
