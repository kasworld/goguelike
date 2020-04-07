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
	Black:           {"가", [3]float64{-1, -1, -1}, htmlcolors.Black},
	Maroon:          {"나", [3]float64{0, -1, -1}, htmlcolors.Maroon},
	Red:             {"다", [3]float64{1, -1, -1}, htmlcolors.Red},
	Green:           {"라", [3]float64{-1, 0, -1}, htmlcolors.Green},
	Olive:           {"마", [3]float64{0, 0, -1}, htmlcolors.Olive},
	DarkOrange:      {"바", [3]float64{1, 0, -1}, htmlcolors.DarkOrange},
	Lime:            {"석", [3]float64{-1, 1, -1}, htmlcolors.Lime},
	Chartreuse:      {"아", [3]float64{0, 1, -1}, htmlcolors.Chartreuse},
	Yellow:          {"정", [3]float64{1, 1, -1}, htmlcolors.Yellow},
	Navy:            {"차", [3]float64{-1, -1, 0}, htmlcolors.Navy},
	Purple:          {"카", [3]float64{0, -1, 0}, htmlcolors.Purple},
	DeepPink:        {"타", [3]float64{1, -1, 0}, htmlcolors.DeepPink},
	Teal:            {"파", [3]float64{-1, 0, 0}, htmlcolors.Teal},
	Salmon:          {"하", [3]float64{1, 0, 0}, htmlcolors.Salmon},
	SpringGreen:     {"강", [3]float64{-1, 1, 0}, htmlcolors.SpringGreen},
	LightGreen:      {"야", [3]float64{0, 1, 0}, htmlcolors.LightGreen},
	Khaki:           {"어", [3]float64{1, 1, 0}, htmlcolors.Khaki},
	Blue:            {"현", [3]float64{-1, -1, 1}, htmlcolors.Blue},
	DarkViolet:      {"오", [3]float64{0, -1, 1}, htmlcolors.DarkViolet},
	Magenta:         {"요", [3]float64{1, -1, 1}, htmlcolors.Magenta},
	DodgerBlue:      {"우", [3]float64{-1, 0, 1}, htmlcolors.DodgerBlue},
	MediumSlateBlue: {"유", [3]float64{0, 0, 1}, htmlcolors.MediumSlateBlue},
	Violet:          {"으", [3]float64{1, 0, 1}, htmlcolors.Violet},
	Cyan:            {"이", [3]float64{-1, 1, 1}, htmlcolors.Cyan},
	Aquamarine:      {"김", [3]float64{0, 1, 1}, htmlcolors.Aquamarine},
	White:           {"원", [3]float64{1, 1, 1}, htmlcolors.White},
}

var Wraper = wrapper.New(FactionType_Count)
