// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objmaxtower

import (
	"fmt"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/tool/towermaker/floortemplate"
	"github.com/kasworld/goguelike/tool/towermaker/towermake"
)

func New(name string) *towermake.Tower {
	var rnd = g2rand.New()

	floorW, floorH := 512, 256
	aoCount := 1024
	stairCount := 128
	roomCount := 512

	tw := towermake.New(name)

	fm := tw.Add("ObjMax0", floorW, floorH, 1.0).Appends(
		floortemplate.CityRooms(floorW, floorH, 12, 10, 5, rnd.Intn)...,
	)
	fm.Appendf("ActiveObjectsRand count=%v", aoCount)

	fm = tw.Add("ObjMax1", floorW, floorH, 1.0).Appends(
		floortemplate.MixedResourceMaze(floorW, floorH)...,
	)
	fm.Appendf("ActiveObjectsRand count=%v", aoCount)
	fm.Appends(
		fmt.Sprintf("ResourceFillRect resource=Soil  amount=1  x=0 y=0  w=%v h=%v", floorW, floorH),
	)
	fm.Appends(
		floortemplate.CityRoomsRand(roomCount, rnd.Intn)...,
	)

	fm = tw.Add("ObjMax2", floorW, floorH, 1.0).Appends(
		floortemplate.MixedResourceMaze(floorW, floorH)...,
	)
	fm.Appendf("ActiveObjectsRand count=%v", aoCount)
	fm.Appends(
		floortemplate.CityRoomsRand(roomCount, rnd.Intn)...,
	)

	for _, fm := range tw.GetList() {
		if !fm.IsFinalizeTerrain() {
			fm.Appends("FinalizeTerrain", "")
		}
	}
	for i := 0; i < stairCount; i++ {
		tw.GetByName("ObjMax0").ConnectStairUp("Rand", "Rand", tw.GetByName("ObjMax1"))
		tw.GetByName("ObjMax1").ConnectStairUp("Rand", "Rand", tw.GetByName("ObjMax2"))
		tw.GetByName("ObjMax2").ConnectStairUp("Rand", "Rand", tw.GetByName("ObjMax0"))
	}

	for _, fm := range tw.GetList() {
		suffix := "Rand"
		for i := 0; i < roomCount/8; i++ {
			fm.ConnectAutoInPortalTo(suffix, suffix, fm)
		}
		for i := 0; i < roomCount; i++ {
			fm.AddTrapTeleportTo(suffix, fm)
		}
		fm.AddAllEffectTrap(suffix, roomCount/8)
		fm.AddRecycler(suffix, stairCount)
	}

	return tw
}
