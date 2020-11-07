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

package floor

import (
	"fmt"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/game/terraini"
	"github.com/kasworld/goguelike/lib/uuidposmani"
	"github.com/kasworld/intervalduration"
)

func (f *Floor) Initialized() bool {
	return f.initialized
}

func (f *Floor) GetTower() gamei.TowerI {
	return f.tower
}

func (f *Floor) GetName() string {
	return f.terrain.GetName()
}

func (f *Floor) GetReqCh() chan<- interface{} {
	return f.recvRequestCh
}

func (f *Floor) ReqChState() string {
	return fmt.Sprintf("%v/%v",
		len(f.recvRequestCh), cap(f.recvRequestCh),
	)
}

func (f *Floor) GetInterDur() *intervalduration.IntervalDuration {
	return f.interDur
}

func (f *Floor) GetActTurn() int {
	return f.interDur.GetCount()
}

func (f *Floor) GetStatPacketObjOver() *actpersec.ActPerSec {
	return f.statPacketObjOver
}

func (f *Floor) GetCmdFloorActStat() *actpersec.ActPerSec {
	return f.floorCmdActStat
}

func (f *Floor) GetTerrain() terraini.TerrainI {
	return f.terrain
}

func (f *Floor) GetWidth() int {
	return f.w
}

func (f *Floor) GetHeight() int {
	return f.h
}

// for visitarea
func (f *Floor) VisitableCount() int {
	return f.terrain.Tile2Discover
}

func (f *Floor) GetBias() bias.Bias {
	return f.bias
}
func (f *Floor) GetEnvBias() bias.Bias {
	return f.bias.Add(f.tower.GetBias())
}

// for web info

func (f *Floor) TotalActiveObjCount() int {
	return f.aoPosMan.Count()
}
func (f *Floor) TotalCarryObjCount() int {
	return f.poPosMan.Count()
}

func (f *Floor) GetActiveObjPosMan() uuidposmani.UUIDPosManI {
	return f.aoPosMan
}
func (f *Floor) GetFieldObjPosMan() uuidposmani.UUIDPosManI {
	return f.foPosMan
}
func (f *Floor) GetCarryObjPosMan() uuidposmani.UUIDPosManI {
	return f.poPosMan
}

func (f *Floor) FindPath(dstx, dsty, srcx, srcy int, trylimit int) [][2]int {
	return f.terrain.FindPath(dstx, dsty, srcx, srcy, trylimit)
}
