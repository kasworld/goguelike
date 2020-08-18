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
		"ResourceFillRect resource=Soil  amount=1  x=0 y=0  w=800 h=640",
		"ResourceMazeWall resource=Fog  amount=500000   x=0 y=0 w=800 h=640 xn=128 yn=128 connerfill=true",
		"ResourceMazeWall resource=Water amount=1000000 x=0 y=0 w=800 h=640 xn=63 yn=63 connerfill=true",
		"ResourceMazeWall resource=Fire  amount=1000000 x=0 y=0 w=800 h=640 xn=67 yn=67 connerfill=true",
		"ResourceMazeWall resource=Ice   amount=1000000 x=0 y=0 w=800 h=640 xn=69 yn=69 connerfill=true",
		"ResourceMazeWall resource=Plant amount=2000000 x=0 y=0 w=800 h=640 xn=71 yn=71  connerfill=true",
		"ResourceMazeWall resource=Stone amount=1000000 x=0 y=0 w=800 h=640 xn=73 yn=73 connerfill=true",
	}

	tw := towermake.New(name)
	fm := tw.Add("ObjMax1", 800, 640, 8192, 0, 1.0).Appends(
		str...,
	)
	fm.Appends(
		floortemplate.CityRoomsRand(1024, rnd.Intn)...,
	)

	fm = tw.Add("ObjMax2", 800, 640, 8192, 0, 1.0).Appends(
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
		tw.GetByName("ObjMax1").ConnectStairUp("Rand", "Rand", tw.GetByName("ObjMax2"))
	}

	for _, fm := range tw.GetList() {
		suffix := "Rand"
		for i := 0; i < 2048; i++ {
			fm.ConnectAutoInPortalTo(suffix, suffix, fm)
			fm.AddTrapTeleportTo(suffix, fm)
		}
		fm.AddAllEffectTrap(suffix, 512)
		fm.AddRecycler(suffix, 4096)
	}

	return tw
}
