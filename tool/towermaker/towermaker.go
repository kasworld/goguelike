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

package main

import (
	"flag"
	"fmt"
	"math/bits"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/game/towerscript"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/tool/towermaker/floormake"
)

func main() {
	towername := flag.String("towername", "roguegen", "tower filename w/o .tower")
	floorcount := flag.Int("floorcount", 100, "floor count in tower")
	flag.Parse()

	fmt.Printf("make tower with towername:%v floorcount:%v\n", *towername, *floorcount)
	makeRogueTower(*towername, *floorcount)
}

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func makePowerOf2(v int) int {
	return 1 << uint(bits.Len(uint(v-1)))
}

func makeRogueTower(towerName string, floorCount int) {
	var rnd = g2rand.New()
	var whList = []int{
		32, 64, 128,
	}
	floorList := make([]*floormake.FloorMake, floorCount)
	for i := range floorList {
		w := whList[rnd.Intn(len(whList))]
		h := whList[rnd.Intn(len(whList))]
		roomCount := w * h / 512
		if roomCount < 2 {
			roomCount = 2
		}
		floorList[i] = floormake.New(
			i,
			fmt.Sprintf("Floor%v", i),
			w, h,
			roomCount/2, 0,
			1.0,
		)
	}

	for _, fm := range floorList {
		fm.MakeRoguelike()
	}

	for _, fm := range floorList {
		roomCount := fm.W * fm.H / 512
		fm.ConnectStairUp("InRoom", floorList[wrapInt(fm.Num+1, floorCount)])
		fm.AddRecycler("InRoom", roomCount/2)
		fm.AddTeleportIn("InRoom", roomCount/2)
		fm.AddTeleportTrapOut("InRoom", floorList, roomCount/2)
		fm.AddEffectTrap("InRoom", roomCount/2)
	}

	tw := make(towerscript.TowerScript, 0)
	for i := 0; i < floorCount; i++ {
		tw = append(tw, floorList[i].Script)
	}
	err := tw.SaveJSON(towerName + ".tower")
	if err != nil {
		g2log.Error("%v", err)
	}
}
