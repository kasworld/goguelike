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
	Empty:             {"🍾", htmlcolors.Black, 10},
	MinorHealing:      {"🍾", htmlcolors.Red, 10},
	MinorHeal:         {"🍾", htmlcolors.Red, 5},
	MajorHealing:      {"🍾", htmlcolors.Red, 10},
	MajorHeal:         {"🍾", htmlcolors.Red, 3},
	GreatHealing:      {"🍾", htmlcolors.Red, 10},
	CompleteHeal:      {"🍾", htmlcolors.Red, 1},
	MinorSpanHealing:  {"🍾", htmlcolors.Red, 10},
	MinorActing:       {"🍾", htmlcolors.Lime, 10},
	MinorAct:          {"🍾", htmlcolors.Lime, 5},
	MajorActing:       {"🍾", htmlcolors.Lime, 10},
	MajorAct:          {"🍾", htmlcolors.Lime, 3},
	GreatActing:       {"🍾", htmlcolors.Lime, 10},
	CompleteAct:       {"🍾", htmlcolors.Lime, 1},
	MinorSpanActing:   {"🍾", htmlcolors.Lime, 10},
	MinorSpanVision:   {"🍾", htmlcolors.Blue, 10},
	MajorSpanVision:   {"🍾", htmlcolors.Blue, 5},
	PerfectSpanVision: {"🍾", htmlcolors.Blue, 1},
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
	MinorHealing: []statusoptype.OpArg{
		{statusoptype.AddHP, 10.0},
	},
	MinorHeal: []statusoptype.OpArg{
		{statusoptype.AddHPRate, 0.10},
	},
	MinorActing: []statusoptype.OpArg{
		{statusoptype.AddSP, 10.0},
	},
	MinorAct: []statusoptype.OpArg{
		{statusoptype.AddSPRate, 0.10},
	},
	MajorHealing: []statusoptype.OpArg{
		{statusoptype.AddHP, 50.0},
	},
	MajorHeal: []statusoptype.OpArg{
		{statusoptype.AddHPRate, 0.50},
	},
	MajorActing: []statusoptype.OpArg{
		{statusoptype.AddSP, 50.0},
	},
	MajorAct: []statusoptype.OpArg{
		{statusoptype.AddSPRate, 0.50},
	},
	GreatHealing: []statusoptype.OpArg{
		{statusoptype.AddHP, 100.0},
	},
	GreatActing: []statusoptype.OpArg{
		{statusoptype.AddSP, 100.0},
	},
	CompleteHeal: []statusoptype.OpArg{
		{statusoptype.AddHPRate, 1.00},
	},
	CompleteAct: []statusoptype.OpArg{
		{statusoptype.AddSPRate, 1.00},
	},

	// repeating buff
	MinorSpanHealing: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.AddHP, 1.0},
	),
	MinorSpanActing: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.AddSP, 1.0},
	),

	MinorSpanVision: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.ModSight, 1.0},
	),

	MajorSpanVision: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.ModSight, 5.0},
	),
	PerfectSpanVision: statusoptype.Repeat(300,
		statusoptype.OpArg{statusoptype.ModSight, gameconst.SightXray},
	),
}

var AIRecycleMap = map[PotionType]bool{
	Empty:             true,
	MinorHealing:      false,
	MinorHeal:         false,
	MinorActing:       false,
	MinorAct:          false,
	MajorHealing:      false,
	MajorHeal:         false,
	MajorActing:       false,
	MajorAct:          false,
	GreatHealing:      false,
	GreatActing:       false,
	CompleteHeal:      false,
	CompleteAct:       false,
	MinorSpanHealing:  false,
	MinorSpanActing:   false,
	MinorSpanVision:   false,
	MajorSpanVision:   false,
	PerfectSpanVision: false,
}
