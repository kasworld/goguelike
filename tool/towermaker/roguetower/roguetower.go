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

package roguetower

import (
	"fmt"
	"math/bits"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/tool/towermaker/floortemplate"
	"github.com/kasworld/goguelike/tool/towermaker/towermake"
)

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func makePowerOf2(v int) int {
	return 1 << uint(bits.Len(uint(v-1)))
}

func New(name string, floorCount int) *towermake.Tower {
	var rnd = g2rand.New()
	var whList = []int{
		32, 64, 128,
	}
	tw := towermake.New(name)
	for i := 0; i < floorCount; i++ {
		w := whList[rnd.Intn(len(whList))]
		h := whList[rnd.Intn(len(whList))]
		roomCount := w * h / 512
		if roomCount < 2 {
			roomCount = 2
		}
		fm := tw.Add(fmt.Sprintf("Floor%v", i), w, h, 1.0)
		fm.Appendf("ActiveObjectsRand count=%v", roomCount/2)
		fm.Appends(
			floortemplate.RoguelikeRand(roomCount, rnd.Intn)...,
		)
	}

	for _, fm := range tw.GetList() {
		if !fm.IsFinalizeTerrain() {
			fm.Appends("FinalizeTerrain", "")
		}
	}

	for i, fm := range tw.GetList() {
		roomCount := fm.CalcRoomCount()
		fm.ConnectStairUp("InRoom", "InRoom", tw.GetList()[wrapInt(i+1, floorCount)])
		fm.AddRecycler("InRoom", roomCount/2)
		fm.AddTeleportIn("InRoom", roomCount/2)

		for trapMade := 0; trapMade < roomCount/2; trapMade++ {
			dstFloor := tw.GetList()[rnd.Intn(tw.GetCount())]
			fm.AddTrapTeleportTo("InRoom", dstFloor)
		}
		fm.AddEffectTrap("InRoom", roomCount/2)
	}
	return tw
}
