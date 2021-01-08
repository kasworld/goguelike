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

package factiontype

import (
	"github.com/kasworld/htmlcolors"
	"github.com/kasworld/wrapper"
)

func (ft FactionType) Rune() string {
	return attrib[ft].Rune
}

func (ft FactionType) Color24() htmlcolors.Color24 {
	return attrib[ft].Color24
}

func (ft FactionType) FactorBase() [3]float64 {
	return attrib[ft].FactorBase
}

var attrib = [FactionType_Count]struct {
	Rune       string
	FactorBase [3]float64
	Color24    htmlcolors.Color24
}{
	Black:           {"A", [3]float64{-1, -1, -1}, htmlcolors.Black},
	Maroon:          {"B", [3]float64{0, -1, -1}, htmlcolors.Maroon},
	Red:             {"C", [3]float64{1, -1, -1}, htmlcolors.Red},
	Green:           {"D", [3]float64{-1, 0, -1}, htmlcolors.Green},
	Olive:           {"E", [3]float64{0, 0, -1}, htmlcolors.Olive},
	DarkOrange:      {"F", [3]float64{1, 0, -1}, htmlcolors.DarkOrange},
	Lime:            {"G", [3]float64{-1, 1, -1}, htmlcolors.Lime},
	Chartreuse:      {"H", [3]float64{0, 1, -1}, htmlcolors.Chartreuse},
	Yellow:          {"I", [3]float64{1, 1, -1}, htmlcolors.Yellow},
	Navy:            {"J", [3]float64{-1, -1, 0}, htmlcolors.Navy},
	Purple:          {"K", [3]float64{0, -1, 0}, htmlcolors.Purple},
	DeepPink:        {"L", [3]float64{1, -1, 0}, htmlcolors.DeepPink},
	Teal:            {"N", [3]float64{-1, 0, 0}, htmlcolors.Teal},
	Salmon:          {"M", [3]float64{1, 0, 0}, htmlcolors.Salmon},
	SpringGreen:     {"O", [3]float64{-1, 1, 0}, htmlcolors.SpringGreen},
	LightGreen:      {"P", [3]float64{0, 1, 0}, htmlcolors.LightGreen},
	Khaki:           {"Q", [3]float64{1, 1, 0}, htmlcolors.Khaki},
	Blue:            {"R", [3]float64{-1, -1, 1}, htmlcolors.Blue},
	DarkViolet:      {"S", [3]float64{0, -1, 1}, htmlcolors.DarkViolet},
	Magenta:         {"T", [3]float64{1, -1, 1}, htmlcolors.Magenta},
	DodgerBlue:      {"U", [3]float64{-1, 0, 1}, htmlcolors.DodgerBlue},
	MediumSlateBlue: {"V", [3]float64{0, 0, 1}, htmlcolors.MediumSlateBlue},
	Violet:          {"W", [3]float64{1, 0, 1}, htmlcolors.Violet},
	Cyan:            {"X", [3]float64{-1, 1, 1}, htmlcolors.Cyan},
	Aquamarine:      {"Y", [3]float64{0, 1, 1}, htmlcolors.Aquamarine},
	White:           {"Z", [3]float64{1, 1, 1}, htmlcolors.White},
}

var Wraper = wrapper.New(FactionType_Count)
