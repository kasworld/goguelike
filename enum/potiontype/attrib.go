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

package potiontype

import (
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/statusoptype"
	"github.com/kasworld/htmlcolors"
)

func (pt PotionType) Color24() htmlcolors.Color24 {
	return attrib[pt].Color
}

func (pt PotionType) Rune() string {
	return attrib[pt].Rune
}

var TotalPotionMakeRate int

func (pt PotionType) MakeRate() int {
	return attrib[pt].MakeRate
}

var attrib = [PotionType_Count]struct {
	Rune     string
	Color    htmlcolors.Color24
	MakeRate int
}{
	Empty:           {"!", htmlcolors.Black, 10},
	RecoverHP10:     {"!", htmlcolors.Red, 10},
	RecoverHPRate10: {"!", htmlcolors.Red, 5},
	RecoverHP50:     {"!", htmlcolors.Red, 10},
	RecoverHPRate50: {"!", htmlcolors.Red, 3},
	RecoverHP100:    {"!", htmlcolors.Red, 10},
	RecoverHPFull:   {"!", htmlcolors.Red, 1},
	BuffRecoverHP1:  {"!", htmlcolors.Red, 10},
	RecoverSP10:     {"!", htmlcolors.Lime, 10},
	RecoverSPRate10: {"!", htmlcolors.Lime, 5},
	RecoverSP50:     {"!", htmlcolors.Lime, 10},
	RecoverSPRate50: {"!", htmlcolors.Lime, 3},
	RecoverSP100:    {"!", htmlcolors.Lime, 10},
	RecoverSPFull:   {"!", htmlcolors.Lime, 1},
	BuffRecoverSP1:  {"!", htmlcolors.Lime, 10},
	BuffSight1:      {"!", htmlcolors.Blue, 10},
	BuffSight5:      {"!", htmlcolors.Blue, 5},
	BuffSightMax:    {"!", htmlcolors.Blue, 1},
}

func init() {
	sum := 0
	for _, v := range attrib {
		sum += v.MakeRate
	}
	TotalPotionMakeRate = sum
}

func GetBuffByPotionType(pt PotionType) []statusoptype.OpArg {
	return potion2BuffList[pt]
}

var potion2BuffList = [PotionType_Count][]statusoptype.OpArg{
	// one time buff
	RecoverHP10: []statusoptype.OpArg{
		{statusoptype.AddHP, 10.0},
	},
	RecoverHPRate10: []statusoptype.OpArg{
		{statusoptype.AddHPRate, 0.10},
	},
	RecoverSP10: []statusoptype.OpArg{
		{statusoptype.AddSP, 10.0},
	},
	RecoverSPRate10: []statusoptype.OpArg{
		{statusoptype.AddSPRate, 0.10},
	},
	RecoverHP50: []statusoptype.OpArg{
		{statusoptype.AddHP, 50.0},
	},
	RecoverHPRate50: []statusoptype.OpArg{
		{statusoptype.AddHPRate, 0.50},
	},
	RecoverSP50: []statusoptype.OpArg{
		{statusoptype.AddSP, 50.0},
	},
	RecoverSPRate50: []statusoptype.OpArg{
		{statusoptype.AddSPRate, 0.50},
	},
	RecoverHP100: []statusoptype.OpArg{
		{statusoptype.AddHP, 100.0},
	},
	RecoverSP100: []statusoptype.OpArg{
		{statusoptype.AddSP, 100.0},
	},
	RecoverHPFull: []statusoptype.OpArg{
		{statusoptype.AddHPRate, 1.00},
	},
	RecoverSPFull: []statusoptype.OpArg{
		{statusoptype.AddSPRate, 1.00},
	},

	// repeating buff
	BuffRecoverHP1: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.AddHP, 1.0},
	),
	BuffRecoverSP1: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.AddSP, 1.0},
	),

	BuffSight1: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.ModSight, 1.0},
	),

	BuffSight5: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.ModSight, 5.0},
	),
	BuffSightMax: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.ModSight, gameconst.SightXray},
	),
}

var AIRecycleMap = map[PotionType]bool{
	Empty:           true,
	RecoverHP10:     false,
	RecoverHPRate10: false,
	RecoverSP10:     false,
	RecoverSPRate10: false,
	RecoverHP50:     false,
	RecoverHPRate50: false,
	RecoverSP50:     false,
	RecoverSPRate50: false,
	RecoverHP100:    false,
	RecoverSP100:    false,
	RecoverHPFull:   false,
	RecoverSPFull:   false,
	BuffRecoverHP1:  false,
	BuffRecoverSP1:  false,
	BuffSight1:      false,
	BuffSight5:      false,
	BuffSightMax:    false,
}
