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

package fieldobjacttype

import (
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/statusoptype"
	"github.com/kasworld/htmlcolors"
)

func (v FieldObjActType) Color24() htmlcolors.Color24 {
	return attrib[v].Color24
}

func (v FieldObjActType) TrapNoti() bool {
	return attrib[v].TrapNoti
}

func (v FieldObjActType) AutoTrigger() bool {
	return attrib[v].AutoTrigger
}

func (v FieldObjActType) TriggerRate() float64 {
	return attrib[v].TriggerRate
}

func (v FieldObjActType) SkipAOAct() bool {
	return attrib[v].SkipAOAct
}
func (v FieldObjActType) NeedTANoti() bool {
	return attrib[v].NeedTANoti
}

var attrib = [FieldObjActType_Count]struct {
	TrapNoti    bool
	AutoTrigger bool
	TriggerRate float64
	SkipAOAct   bool
	NeedTANoti  bool
	Color24     htmlcolors.Color24
}{
	None:            {false, false, 1.0, false, false, htmlcolors.Black},
	PortalInOut:     {false, false, 1.0, false, false, htmlcolors.Cyan},
	PortalIn:        {false, false, 1.0, false, false, htmlcolors.Cyan},
	PortalOut:       {false, false, 1.0, false, false, htmlcolors.Cyan},
	PortalAutoIn:    {false, true, 1.0, true, true, htmlcolors.Cyan},
	RecycleCarryObj: {false, false, 1.0, false, false, htmlcolors.Green},
	Teleport:        {true, true, 0.1, true, true, htmlcolors.Red},
	ForgetFloor:     {true, true, 0.2, false, true, htmlcolors.OrangeRed},
	ForgetOneFloor:  {true, true, 0.3, false, true, htmlcolors.OrangeRed},
	AlterFaction:    {true, true, 0.5, false, false, htmlcolors.Red},
	AllFaction:      {true, true, 0.5, false, false, htmlcolors.Red},
	Bleeding:        {true, true, 0.2, false, false, htmlcolors.DeepPink},
	Chilly:          {true, true, 0.2, false, false, htmlcolors.DeepPink},
	Blind:           {true, true, 0.2, false, false, htmlcolors.DeepPink},
	Invisible:       {true, true, 0.5, false, false, htmlcolors.DeepPink},
	Burden:          {true, true, 0.2, false, false, htmlcolors.DeepPink},
	Float:           {true, true, 0.3, false, false, htmlcolors.DeepPink},
	Greasy:          {true, true, 0.5, false, false, htmlcolors.DeepPink},
	Drunken:         {true, true, 0.5, false, false, htmlcolors.DeepPink},
	Sleepy:          {true, true, 0.1, false, false, htmlcolors.DeepPink},
}

// try act on fieldobj
var ClientData = [FieldObjActType_Count]struct {
	ActOn bool
	Text  string
}{
	None:            {true, ""},
	PortalInOut:     {true, "portal in/out"},
	PortalIn:        {true, "portal oneway"},
	PortalOut:       {true, "portal out only"},
	PortalAutoIn:    {false, "portal auto in oneway"},
	RecycleCarryObj: {true, "recycle carryobj to money"},
	Teleport:        {false, "teleport somewhere"},
	ForgetFloor:     {false, "forget current floor"},
	ForgetOneFloor:  {false, "forget some floor you visited"},
	AlterFaction:    {false, "change faction randomly"},
	AllFaction:      {false, "rotate all faction"},
	Blind:           {false, "sight 0"},
	Bleeding:        {false, "hp damage"},
	Invisible:       {false, "other cannot see you"},
	Burden:          {false, "overload limit reduced"},
	Chilly:          {false, "sp damage"},
	Float:           {false, "float in air"},
	Greasy:          {false, "greasy body"},
	Drunken:         {false, "random direction"},
	Sleepy:          {false, "cannot act"},
}

func GetBuffByFieldObjActType(at FieldObjActType) []statusoptype.OpArg {
	return foAct2BuffList[at]
}

var foAct2BuffList = [FieldObjActType_Count][]statusoptype.OpArg{
	// immediate effect
	AlterFaction: {
		statusoptype.OpArg{statusoptype.RndFaction, nil},
	},
	AllFaction: statusoptype.RepeatShift(260, 10,
		statusoptype.OpArg{statusoptype.IncFaction, 1},
	),

	ForgetFloor: {
		statusoptype.OpArg{statusoptype.ForgetFloor, nil},
	},
	ForgetOneFloor: {
		statusoptype.OpArg{statusoptype.ForgetOneFloor, nil},
	},

	// statusop debuff
	Bleeding: statusoptype.RepeatShift(200, 10,
		statusoptype.OpArg{statusoptype.AddHPRate, -0.05},
	),
	Chilly: statusoptype.RepeatShift(200, 10,
		statusoptype.OpArg{statusoptype.AddSPRate, -0.05},
	),

	// condition debuff
	Blind: statusoptype.RepeatShift(200, 2,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Blind},
	),
	Invisible: statusoptype.RepeatShift(200, 2,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Invisible},
	),
	Burden: statusoptype.RepeatShift(100, 1,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Burden},
	),
	Float: statusoptype.RepeatShift(200, 1,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Float},
	),
	Greasy: statusoptype.RepeatShift(400, 4,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Greasy},
	),
	Drunken: statusoptype.RepeatShift(200, 2,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Drunken},
	),
	Sleepy: statusoptype.RepeatShift(200, 4,
		statusoptype.OpArg{statusoptype.SetCondition, condition.Sleep},
	),
}
