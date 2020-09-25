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
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/decaytype"
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
			fm := tw.Add(floorName, 256, 256, 1.0).Appends(
				floortemplate.CityRooms(256, 256, 11, 11, 5, rnd.Intn)...,
			)
			fm.Appendf("ActiveObjectsRand count=%v", 256)

		case 9:
			fm := tw.Add(floorName, 256, 256, 1.0).Appends(
				fmt.Sprintf("ResourceFillRect resource=Soil amount=1 x=0 y=0 w=%v h=%v", 256, 256),
			)
			fm.Appends(
				floortemplate.MixedResourceMaze(256, 256)...,
			)
			fm.Appends(
				floortemplate.CityRoomsRand(512, rnd.Intn)...,
			)
			fm.Appendf("ActiveObjectsRand count=%v", 256)
		}
	}

	for _, fm := range tw.GetList() {
		if !fm.IsFinalizeTerrain() {
			fm.Appends("FinalizeTerrain", "")
		}
	}

	lhlen := gameconst.ViewPortW / 2

	for i, fm := range tw.GetList() {
		roomCount := fm.CalcRoomCount()
		recycleCount := fm.W * fm.H / gameconst.ViewPortWH / 8
		if recycleCount < 2 {
			recycleCount = 2
		}

		suffix1 := "InRoom"
		if roomCount == 0 {
			suffix1 = "Rand"
		}

		fm2 := tw.GetList()[wrapInt(i+1, floorCount)]
		suffix2 := "InRoom"
		if fm2.CalcRoomCount() == 0 {
			suffix2 = "Rand"
		}
		fm.ConnectStairUp(suffix1, suffix2, fm2)

		fm2 = tw.GetList()[rnd.Intn(tw.GetCount())]
		suffix2 = "InRoom"
		if fm2.CalcRoomCount() == 0 {
			suffix2 = "Rand"
		}
		fm.ConnectPortalTo(suffix1, suffix2, fm2)
		fm.ConnectAutoInPortalTo(suffix1, suffix1, fm)
		fm.AddTrapTeleportTo(suffix1, fm)
		for j := 0; j < recycleCount/8; j++ {
			fm2 = tw.GetList()[rnd.Intn(tw.GetCount())]
			suffix2 = "InRoom"
			if fm2.CalcRoomCount() == 0 {
				suffix2 = "Rand"
			}
			fm.ConnectPortalTo(suffix1, suffix2, fm2)

			fm.ConnectAutoInPortalTo("Rand", "Rand", fm)
			fm.AddTrapTeleportTo("Rand", fm)
		}
		fm.AddAllEffectTrap(suffix1, 1)
		if count := recycleCount / 32; count > 0 {
			fm.AddAllEffectTrap("Rand", count)
		}

		if recycleCount > roomCount {
			fm.AddRecycler(suffix1, roomCount)
		} else {
			fm.AddRecycler(suffix1, recycleCount)
		}
		if recycleCount-roomCount > 0 {
			fm.AddRecycler("Rand", recycleCount-roomCount)
		}
		for j := 0; j < decaytype.DecayType_Count; j++ {
			decay := decaytype.DecayType(j)
			fm.Appendf(
				"AddRotateLineAttack%v display=RotateLineAttack winglen=%v wingcount=1 degree=0 perturn=10 decay=%v count=%v message=RotDanger1",
				"Rand", lhlen, decay, 1,
			)
			fm.Appendf(
				"AddRotateLineAttack%v display=RotateLineAttack winglen=%v wingcount=2 degree=0 perturn=10 decay=%v count=%v message=RotDanger2",
				"Rand", lhlen/2, decay, 1,
			)
			fm.Appendf(
				"AddMine%v display=None decay=%v count=%v message=Mine",
				"Rand", decay, 1,
			)
		}
	}
	return tw
}
