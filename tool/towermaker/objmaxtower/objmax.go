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
	stairCount := 1024
	roomCount := 512

	str := []string{
		fmt.Sprintf(
			"ResourceMazeWall resource=Fog   amount=500000  x=0 y=0 w=%v h=%v xn=80 yn=64 connerfill=true",
			floorW, floorH),
		fmt.Sprintf(
			"ResourceMazeWall resource=Water amount=1000000 x=2 y=2 w=%v h=%v xn=80 yn=64 connerfill=true",
			floorW, floorH),
		fmt.Sprintf(
			"ResourceMazeWall resource=Soil  amount=1000000 x=2 y=2 w=%v h=%v xn=80 yn=64 connerfill=true",
			floorW, floorH),
		fmt.Sprintf(
			"ResourceMazeWall resource=Fire  amount=1000000 x=4 y=4 w=%v h=%v xn=80 yn=64 connerfill=true",
			floorW, floorH),
		fmt.Sprintf(
			"ResourceMazeWall resource=Ice   amount=1000000 x=6 y=6 w=%v h=%v xn=80 yn=64 connerfill=true",
			floorW, floorH),
		fmt.Sprintf(
			"ResourceMazeWall resource=Plant amount=1000000 x=8 y=8 w=%v h=%v xn=80 yn=64  connerfill=true",
			floorW, floorH),
		fmt.Sprintf(
			"ResourceMazeWall resource=Plant amount=1000000 x=8 y=8 w=%v h=%v xn=80 yn=64  connerfill=true",
			floorW, floorH),
		fmt.Sprintf(
			"ResourceMazeWall resource=Stone amount=1000000 x=10 y=10 w=%v h=%v xn=80 yn=64 connerfill=true",
			floorW, floorH),
	}

	tw := towermake.New(name)

	fm := tw.Add("ObjMax0", floorW, floorH, aoCount, 0, 1.0).Appends(
		floortemplate.CityRooms(floorW, floorH, 12, 10, 5, rnd.Intn)...,
	)

	fm = tw.Add("ObjMax1", floorW, floorH, aoCount, 0, 1.0).Appends(
		str...,
	)
	fm.Appends(
		fmt.Sprintf("ResourceFillRect resource=Soil  amount=1  x=0 y=0  w=%v h=%v", floorW, floorH),
	)
	fm.Appends(
		floortemplate.CityRoomsRand(roomCount, rnd.Intn)...,
	)

	fm = tw.Add("ObjMax2", floorW, floorH, aoCount, 0, 1.0).Appends(
		str...,
	)
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
		for i := 0; i < roomCount; i++ {
			fm.ConnectAutoInPortalTo(suffix, suffix, fm)
			fm.AddTrapTeleportTo(suffix, fm)
		}
		fm.AddAllEffectTrap(suffix, roomCount/8)
		fm.AddRecycler(suffix, stairCount)
	}

	return tw
}
