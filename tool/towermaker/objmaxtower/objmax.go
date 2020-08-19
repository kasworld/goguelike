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
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/tool/towermaker/floortemplate"
	"github.com/kasworld/goguelike/tool/towermaker/towermake"
)

func New(name string) *towermake.Tower {
	var rnd = g2rand.New()

	str := []string{
		"ResourceMazeWall resource=Fog   amount=500000  x=0 y=0 w=512 h=256 xn=80 yn=64 connerfill=true",
		"ResourceMazeWall resource=Water amount=1000000 x=2 y=2 w=512 h=256 xn=80 yn=64 connerfill=true",
		"ResourceMazeWall resource=Soil  amount=1000000 x=2 y=2 w=512 h=256 xn=80 yn=64 connerfill=true",
		"ResourceMazeWall resource=Fire  amount=1000000 x=4 y=4 w=512 h=256 xn=80 yn=64 connerfill=true",
		"ResourceMazeWall resource=Ice   amount=1000000 x=6 y=6 w=512 h=256 xn=80 yn=64 connerfill=true",
		"ResourceMazeWall resource=Plant amount=1000000 x=8 y=8 w=512 h=256 xn=80 yn=64  connerfill=true",
		"ResourceMazeWall resource=Plant amount=1000000 x=8 y=8 w=512 h=256 xn=80 yn=64  connerfill=true",
		"ResourceMazeWall resource=Stone amount=1000000 x=10 y=10 w=512 h=256 xn=80 yn=64 connerfill=true",
	}

	tw := towermake.New(name)

	fm := tw.Add("ObjMax0", 512, 256, 1024, 0, 1.0).Appends(
		floortemplate.CityRooms(512, 256, 12, 10, 5, rnd.Intn)...,
	)

	fm = tw.Add("ObjMax1", 512, 256, 1024, 0, 1.0).Appends(
		str...,
	)
	fm.Appends(
		"ResourceFillRect resource=Soil  amount=1  x=0 y=0  w=512 h=256",
	)
	fm.Appends(
		floortemplate.CityRoomsRand(1024, rnd.Intn)...,
	)

	fm = tw.Add("ObjMax2", 512, 256, 1024, 0, 1.0).Appends(
		str...,
	)
	fm.Appends(
		floortemplate.CityRoomsRand(1024, rnd.Intn)...,
	)

	for _, fm := range tw.GetList() {
		if !fm.IsFinalizeTerrain() {
			fm.Appends("FinalizeTerrain", "")
		}
	}
	for i := 0; i < 1024; i++ {
		tw.GetByName("ObjMax0").ConnectStairUp("Rand", "Rand", tw.GetByName("ObjMax1"))
		tw.GetByName("ObjMax1").ConnectStairUp("Rand", "Rand", tw.GetByName("ObjMax2"))
	}

	for _, fm := range tw.GetList() {
		suffix := "Rand"
		for i := 0; i < 256; i++ {
			fm.ConnectAutoInPortalTo(suffix, suffix, fm)
			fm.AddTrapTeleportTo(suffix, fm)
		}
		fm.AddAllEffectTrap(suffix, 512)
		fm.AddRecycler(suffix, 1024)
	}

	return tw
}
