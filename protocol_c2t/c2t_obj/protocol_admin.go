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

package c2t_obj

import (
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
)

type ReqAdminFloorMove_data struct {
	Floor string
}
type RspAdminFloorMove_data struct {
	Dummy uint8
}

type ReqAdminTeleport_data struct {
	X int
	Y int
}
type RspAdminTeleport_data struct {
	Dummy uint8
}

type ReqAdminTowerCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminTowerCmd_data struct {
	Dummy uint8
}

type ReqAdminFloorCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminFloorCmd_data struct {
	Dummy uint8
}

type ReqAdminActiveObjCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminActiveObjCmd_data struct {
	Dummy uint8
}

// AdminAddExp  add arg to battle exp
type ReqAdminAddExp_data struct {
	Exp int
}

// AdminAddExp  add arg to battle exp
type RspAdminAddExp_data struct {
	Dummy uint8 // change as you need
}

// AdminPotionEffect buff by arg potion type
type ReqAdminPotionEffect_data struct {
	Potion potiontype.PotionType
}

// AdminPotionEffect buff by arg potion type
type RspAdminPotionEffect_data struct {
	Dummy uint8 // change as you need
}

// AdminScrollEffect buff by arg Scroll type
type ReqAdminScrollEffect_data struct {
	Scroll scrolltype.ScrollType
}

// AdminScrollEffect buff by arg Scroll type
type RspAdminScrollEffect_data struct {
	Dummy uint8 // change as you need
}

// AdminCondition add arg condition for 100 turn
type ReqAdminCondition_data struct {
	Condition condition.Condition
}

// AdminCondition add arg condition for 100 turn
type RspAdminCondition_data struct {
	Dummy uint8 // change as you need
}

// AdminAddPotion add arg potion to inven
type ReqAdminAddPotion_data struct {
	Potion potiontype.PotionType
}

// AdminAddPotion add arg potion to inven
type RspAdminAddPotion_data struct {
	Dummy uint8 // change as you need
}

// AdminAddScroll add arg scroll to inven
type ReqAdminAddScroll_data struct {
	Scroll scrolltype.ScrollType
}

// AdminAddScroll add arg scroll to inven
type RspAdminAddScroll_data struct {
	Dummy uint8 // change as you need
}

// AdminAddMoney add arg money to inven
type ReqAdminAddMoney_data struct {
	Money int
}

// AdminAddMoney add arg money to inven
type RspAdminAddMoney_data struct {
	Dummy uint8 // change as you need
}

// AdminAddEquip add random equip to inven
type ReqAdminAddEquip_data struct {
	Equip equipslottype.EquipSlotType
}

// AdminAddEquip add random equip to inven
type RspAdminAddEquip_data struct {
	Dummy uint8 // change as you need
}

// AdminForgetFloor forget current floor map
type ReqAdminForgetFloor_data struct {
	Dummy uint8 // change as you need
}

// AdminForgetFloor forget current floor map
type RspAdminForgetFloor_data struct {
	Dummy uint8 // change as you need
}

// AdminFloorMap complete current floor map
type ReqAdminFloorMap_data struct {
	Dummy uint8 // change as you need
}

// AdminFloorMap complete current floor map
type RspAdminFloorMap_data struct {
	Dummy uint8 // change as you need
}
