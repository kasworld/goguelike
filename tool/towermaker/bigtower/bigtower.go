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

package bigtower

import (
	"fmt"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/tool/towermaker/floortemplate"
	"github.com/kasworld/goguelike/tool/towermaker/towermake"
)

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func New(name string, floorCount int) *towermake.Tower {
	var rnd = g2rand.New()
	var whList = []int{
		32, 64, 128,
	}
	tw := towermake.New(name)
	for i := 0; i < floorCount; i++ {
		floorName := fmt.Sprintf("Floor%v", i)
		switch i % 10 {
		default: // roguelike floor
			w := whList[rnd.Intn(len(whList))]
			h := whList[rnd.Intn(len(whList))]
			roomCount := w * h / 512
			if roomCount < 2 {
				roomCount = 2
			}
			fm := tw.Add(floorName, w, h, 1.0)
			fm.Appendf("ActiveObjectsRand count=%v", roomCount/2)
			fm.Appends(
				floortemplate.RoguelikeRand(roomCount, rnd.Intn)...,
			)
		case 1:
			fm := tw.Add(floorName, 256, 256, 1.0).Appends(
				floortemplate.AgeingCity256x256()...,
			)
			fm.Appendf("ActiveObjectsRand count=%v", 256)

		case 3:
			fm := tw.Add(floorName, 256, 256, 1.0).Appends(
				floortemplate.AgeingField256x256()...,
			)
			fm.Appendf("ActiveObjectsRand count=%v", 256)
		case 5:
			fm := tw.Add(floorName, 256, 256, 1.0).Appends(
				floortemplate.AgeingMaze256x256()...,
			)
			fm.Appendf("ActiveObjectsRand count=%v", 256)
		case 7:
			fm := tw.Add(floorName, 64, 64, 1.0).Appends(
				floortemplate.FreeForAll64x64()...,
			)
			fm.Appendf("ActiveObjectsRand count=%v", 16)

		case 9:
			fm := tw.Add(floorName, 256, 256, 1.0).Appends(
				"ResourceFromPNG name=imagefloor.png",
				"ResourceRand resource=Plant mean=100000000 stddev=65535 repeat=256",
				"ResourceAgeing initrun=0 msper=64000 resetaftern=1440",
				"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=32 mean=6 stddev=4",
				"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
			)
			fm.Appendf("ActiveObjectsRand count=%v", 256)
		}
	}

	for _, fm := range tw.GetList() {
		if !fm.IsFinalizeTerrain() {
			fm.Appends("FinalizeTerrain", "")
		}
	}

	for i, fm := range tw.GetList() {
		roomCount := fm.CalcRoomCount()
		suffix1 := "InRoom"
		if roomCount == 0 {
			suffix1 = "Rand"
		}
		suffix2 := "InRoom"
		if tw.GetList()[wrapInt(i+1, floorCount)].CalcRoomCount() == 0 {
			suffix2 = "Rand"
		}
		fm.ConnectStairUp(suffix1, suffix2, tw.GetList()[wrapInt(i+1, floorCount)])

		fm.AddRecycler(suffix1, roomCount/2)
		fm.AddTeleportIn(suffix1, roomCount/2)
		for trapMade := 0; trapMade < roomCount/2; trapMade++ {
			dstFloor := tw.GetList()[rnd.Intn(tw.GetCount())]
			fm.AddTrapTeleportTo(suffix1, dstFloor)
		}
		fm.AddEffectTrap(suffix1, roomCount/2)
	}
	return tw
}
