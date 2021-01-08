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

package c2t_idcmd

func (cmd CommandID) SleepCancel() bool {
	return attrib[cmd].sleepCancel
}

func (cmd CommandID) NeedTurn() float64 {
	return attrib[cmd].needTurn
}

var attrib = [CommandID_Count]struct {
	sleepCancel bool
	needTurn    float64
}{
	Invalid:     {false, 0},
	Login:       {false, 0},
	Heartbeat:   {false, 0},
	Chat:        {false, 0},
	AchieveInfo: {false, 0},

	Rebirth:        {false, 0},
	MoveFloor:      {false, 1}, // need check need turn
	AIPlay:         {false, 0},
	VisitFloorList: {false, 0},

	Meditate:    {false, 1},
	KillSelf:    {false, 1},
	Move:        {true, 1},
	Attack:      {true, 1.5},
	AttackWide:  {true, 3},
	AttackLong:  {true, 3},
	Pickup:      {true, 1},
	Drop:        {true, 1},
	Equip:       {true, 1},
	UnEquip:     {true, 1},
	DrinkPotion: {true, 1},
	ReadScroll:  {true, 1},
	Recycle:     {true, 1},
	EnterPortal: {true, 1},
	ActTeleport: {false, 1},

	AdminTowerCmd:     {false, 0},
	AdminFloorCmd:     {false, 0},
	AdminActiveObjCmd: {false, 0},
	AdminFloorMove:    {false, 0},
	AdminTeleport:     {false, 0},

	AdminAddExp:       {false, 0},
	AdminPotionEffect: {false, 0},
	AdminScrollEffect: {false, 0},
	AdminCondition:    {false, 0},
	AdminAddPotion:    {false, 0},
	AdminAddScroll:    {false, 0},
	AdminAddMoney:     {false, 0},
	AdminAddEquip:     {false, 0},
	AdminForgetFloor:  {false, 0},
	AdminFloorMap:     {false, 0},
}
