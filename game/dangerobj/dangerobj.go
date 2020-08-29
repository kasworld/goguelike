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

package dangerobj

import (
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/findnear"
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

type DangerObj struct {
	Owner             uuidposman.UUIDPosI // id of owner ( ao, floor? , fieldobj?)
	RemainTargetCount int                 // end life when 0
	RemainTurn        int                 // remain turn to affect
	TargetTileLen     int                 // effect len in targetTiles
	TargetTiles       findnear.XYLenList  // absolute pos, cross len
}

// IDPosI interface
func (p *DangerObj) GetUUID() string {
	return p.owner.GetUUID()
}

func NewAOAttact(attacker uuidposman.UUIDPosI, dstx, dsty int) *DangerObj {
	return &DangerObj{
		Owner: attacker,
		RemainTargetCount:1,
		RemainTurn: 1,
		TargetTileLen:1,
		TargetTiles: []findnear.XYLenList{ findnear.XYLen{
			X:dstx, Y: dsty, L: 1,
		}}
	}
}
