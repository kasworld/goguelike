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

package condition

import "github.com/kasworld/htmlcolors"

func (cn Condition) HideSelf() bool {
	return attrib[cn].hideSelf
}

func (cn Condition) HideOther() bool {
	return attrib[cn].hideOther
}

func (cn Condition) Color() htmlcolors.Color24 {
	return attrib[cn].color24
}

func (cn Condition) Rune() string {
	return attrib[cn].runeStr
}

// Probability effect probability
func (cn Condition) Probability() float64 {
	return attrib[cn].probability
}

var attrib = [Condition_Count]struct {
	runeStr     string
	probability float64
	hideSelf    bool // hide to client self ao
	hideOther   bool // hide to client other ao
	color24     htmlcolors.Color24
}{
	Blind:     {"bl", 1.00, false, false, htmlcolors.DarkRed},
	Invisible: {"iv", 1.00, false, false, htmlcolors.LemonChiffon},
	Burden:    {"bu", 1.00, false, false, htmlcolors.DeepPink},
	Float:     {"fl", 1.00, false, false, htmlcolors.Wheat},
	Greasy:    {"gr", 0.25, false, false, htmlcolors.PapayaWhip},
	Drunken:   {"dr", 1.00, false, false, htmlcolors.Plum},
	Sleep:     {"sl", 1.00, false, false, htmlcolors.LightCoral},
	Contagion: {"cn", 1.00, false, false, htmlcolors.DarkGreen},
	Slow:      {"sl", 1.00, false, false, htmlcolors.DarkBlue},
	Haste:     {"ha", 1.00, false, false, htmlcolors.LightBlue},
}
