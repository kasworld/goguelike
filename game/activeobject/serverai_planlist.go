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

package activeobject

import (
	"bytes"
	"fmt"

	"github.com/kasworld/goguelike/enum/aiplan"
	"github.com/kasworld/goguelike/enum/aotype"
)

type planObj struct {
	InitFn func(sai *ServerAIState) int
	ActFn  func(sai *ServerAIState) bool
}

var aoType2aiPlan = [...]planList{
	aotype.System: planList{
		aiplan.Chat,
		aiplan.StrollAround,
		aiplan.Move2Dest,
		aiplan.Revenge,
		aiplan.UsePortal,
		aiplan.RechargeSafe,
		aiplan.RechargeCan,
		aiplan.PickupCarryObj,
		aiplan.Equip,
		aiplan.Attack,
		aiplan.MoveStraight3,
		aiplan.MoveStraight5,
	},
	aotype.User: planList{
		aiplan.StrollAround,
		aiplan.Move2Dest,
		aiplan.Revenge,
		aiplan.UsePortal,
		aiplan.MoveToRecycler,
		aiplan.RechargeSafe,
		aiplan.RechargeCan,
		aiplan.PickupCarryObj,
		aiplan.Equip,
		aiplan.UsePotion,
		aiplan.Attack,
		aiplan.MoveStraight3,
		aiplan.MoveStraight5,
	},
}

type planList []aiplan.AIPlan

func (pl planList) String() string {
	var buf bytes.Buffer
	for _, v := range pl {
		fmt.Fprintf(&buf, "%s ", v)
	}
	return buf.String()
}

// get front plan
func (pl planList) GetCurrentPlan() aiplan.AIPlan {
	return pl.getLast()
}

func (pl planList) getFront() aiplan.AIPlan {
	return pl[0]
}
func (pl planList) getLast() aiplan.AIPlan {
	return pl[len(pl)-1]
}

func (pl planList) findPos(p aiplan.AIPlan) int {
	pos := -1
	for i, v := range pl {
		if v == p {
			pos = i
			break
		}
	}
	return pos
}

// move p to front
func (pl planList) move2Front(p aiplan.AIPlan) bool {
	pos := pl.findPos(p)
	if pos != -1 {
		copy(pl[1:], pl[:pos])
		pl[0] = p
		return true
	}
	return false
}

// pop current and push back
func (pl planList) front2Last() {
	p := pl[0]
	copy(pl, pl[1:])
	pl[len(pl)-1] = p
}

func (pl planList) dup() planList {
	rtn := make(planList, len(pl))
	copy(rtn, pl)
	return rtn
}
