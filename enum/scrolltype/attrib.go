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

package scrolltype

import (
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/statusoptype"
	"github.com/kasworld/htmlcolors"
)

func (st ScrollType) Color24() htmlcolors.Color24 {
	return attrib[st].Color
}

func (pt ScrollType) Rune() string {
	return attrib[pt].Rune
}

var TotalScrollMakeRate int

func (pt ScrollType) MakeRate() int {
	return attrib[pt].MakeRate
}

var attrib = [ScrollType_Count]struct {
	Rune     string
	Color    htmlcolors.Color24
	MakeRate int
}{
	Empty:                  {"📜", htmlcolors.PaleGreen, 10},
	FloorMap:               {"📜", htmlcolors.LimeGreen, 1},
	Teleport:               {"📜", htmlcolors.DarkSeaGreen, 5},
	FactionRnd:             {"📜", htmlcolors.Green, 5},
	FactionNext:            {"📜", htmlcolors.Green, 5},
	FactionBorn:            {"📜", htmlcolors.Green, 5},
	FactionBlack:           {"📜", htmlcolors.Black, 5},
	FactionMaroon:          {"📜", htmlcolors.Maroon, 5},
	FactionRed:             {"📜", htmlcolors.Red, 5},
	FactionGreen:           {"📜", htmlcolors.Green, 5},
	FactionOlive:           {"📜", htmlcolors.Olive, 5},
	FactionDarkOrange:      {"📜", htmlcolors.DarkOrange, 5},
	FactionLime:            {"📜", htmlcolors.Lime, 5},
	FactionChartreuse:      {"📜", htmlcolors.Chartreuse, 5},
	FactionYellow:          {"📜", htmlcolors.Yellow, 5},
	FactionNavy:            {"📜", htmlcolors.Navy, 5},
	FactionPurple:          {"📜", htmlcolors.Purple, 5},
	FactionDeepPink:        {"📜", htmlcolors.DeepPink, 5},
	FactionTeal:            {"📜", htmlcolors.Teal, 5},
	FactionSalmon:          {"📜", htmlcolors.Salmon, 5},
	FactionSpringGreen:     {"📜", htmlcolors.SpringGreen, 5},
	FactionLightGreen:      {"📜", htmlcolors.LightGreen, 5},
	FactionKhaki:           {"📜", htmlcolors.Khaki, 5},
	FactionBlue:            {"📜", htmlcolors.Blue, 5},
	FactionDarkViolet:      {"📜", htmlcolors.DarkViolet, 5},
	FactionMagenta:         {"📜", htmlcolors.Magenta, 5},
	FactionDodgerBlue:      {"📜", htmlcolors.DodgerBlue, 5},
	FactionMediumSlateBlue: {"📜", htmlcolors.MediumSlateBlue, 5},
	FactionViolet:          {"📜", htmlcolors.Violet, 5},
	FactionCyan:            {"📜", htmlcolors.Cyan, 5},
	FactionAquamarine:      {"📜", htmlcolors.Aquamarine, 5},
	FactionWhite:           {"📜", htmlcolors.White, 5},
	BiasNeg:                {"📜", htmlcolors.GreenYellow, 5},
	BiasRotateR:            {"📜", htmlcolors.GreenYellow, 5},
	BiasRotateL:            {"📜", htmlcolors.GreenYellow, 5},
}

func init() {
	sum := 0
	for _, v := range attrib {
		sum += v.MakeRate
	}
	TotalScrollMakeRate = sum
}

func GetBuffByScrollType(st ScrollType) []statusoptype.OpArg {
	return scroll2BuffList[st]
}

var scroll2BuffList = [ScrollType_Count][]statusoptype.OpArg{
	FactionRnd: []statusoptype.OpArg{
		{statusoptype.RndFaction, 0},
	},
	FactionNext: []statusoptype.OpArg{
		{statusoptype.IncFaction, 0},
	},
	FactionBorn: []statusoptype.OpArg{
		{statusoptype.ResetFaction, 0},
	},
	BiasNeg: []statusoptype.OpArg{
		{statusoptype.NegBias, 0},
	},
	BiasRotateR: []statusoptype.OpArg{
		{statusoptype.RotateBiasRight, 0},
	},
	BiasRotateL: []statusoptype.OpArg{
		{statusoptype.RotateBiasLeft, 0},
	},

	FactionBlack: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Black},
	},
	FactionMaroon: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Maroon},
	},
	FactionRed: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Red},
	},
	FactionGreen: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Green},
	},
	FactionOlive: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Olive},
	},
	FactionDarkOrange: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.DarkOrange},
	},
	FactionLime: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Lime},
	},
	FactionChartreuse: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Chartreuse},
	},
	FactionYellow: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Yellow},
	},
	FactionNavy: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Navy},
	},
	FactionPurple: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Purple},
	},
	FactionDeepPink: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.DeepPink},
	},
	FactionTeal: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Teal},
	},
	FactionSalmon: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Salmon},
	},
	FactionSpringGreen: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.SpringGreen},
	},
	FactionLightGreen: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.LightGreen},
	},
	FactionKhaki: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Khaki},
	},
	FactionBlue: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Blue},
	},
	FactionDarkViolet: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.DarkViolet},
	},
	FactionMagenta: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Magenta},
	},
	FactionDodgerBlue: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.DodgerBlue},
	},
	FactionMediumSlateBlue: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.MediumSlateBlue},
	},
	FactionViolet: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Violet},
	},
	FactionCyan: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Cyan},
	},
	FactionAquamarine: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.Aquamarine},
	},
	FactionWhite: []statusoptype.OpArg{
		{statusoptype.SetFaction, factiontype.White},
	},
}

var AIRecycleMap = map[ScrollType]bool{
	Empty:    true,
	FloorMap: false,
	Teleport: false,

	FactionRnd:             true,
	FactionNext:            true,
	FactionBorn:            true,
	FactionBlack:           false,
	FactionMaroon:          false,
	FactionRed:             false,
	FactionGreen:           false,
	FactionOlive:           false,
	FactionDarkOrange:      false,
	FactionLime:            false,
	FactionChartreuse:      false,
	FactionYellow:          false,
	FactionNavy:            false,
	FactionPurple:          false,
	FactionDeepPink:        false,
	FactionTeal:            false,
	FactionSalmon:          false,
	FactionSpringGreen:     false,
	FactionLightGreen:      false,
	FactionKhaki:           false,
	FactionBlue:            false,
	FactionDarkViolet:      false,
	FactionMagenta:         false,
	FactionDodgerBlue:      false,
	FactionMediumSlateBlue: false,
	FactionViolet:          false,
	FactionCyan:            false,
	FactionAquamarine:      false,
	FactionWhite:           false,
	BiasNeg:                true,
	BiasRotateR:            true,
	BiasRotateL:            true,
}
