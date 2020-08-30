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

package dangertype

import (
	"github.com/kasworld/findnear"
	"github.com/kasworld/htmlcolors"
)

/*
attactype ?
None make empty error
BasicAttack
ThrowAttack

target type
None
FirstTaget stop at 1st target
Split distribute to all target
Full full to all target

targer area
None
OneTile
MoveAlongWithLineOfTile
TileArea xylenlist

distance
None
ContactTile
RangedTile

*/

func (dt DangerType) TargetCount() int {
	return attrib[dt].targetCount
}
func (dt DangerType) Turn2Live() int {
	return attrib[dt].turn2Live
}
func (dt DangerType) Color24() htmlcolors.Color24 {
	return attrib[dt].color
}
func (dt DangerType) TargetArea() findnear.XYLenList {
	return attrib[dt].targetArea
}

var attrib = [DangerType_Count]struct {
	targetCount int
	turn2Live   int
	color       htmlcolors.Color24
	targetArea  findnear.XYLenList // + pos + dir
}{
	None:        {0, 0, htmlcolors.Black, findnear.XYLenList{}},
	BasicAttack: {1, 1, htmlcolors.Red, findnear.XYLenList{findnear.XYLen{0, 0, 1}}},
}
