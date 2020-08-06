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
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/tool/towermaker/floortemplate"
	"github.com/kasworld/goguelike/tool/towermaker/towermake"
)

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

var allRoomTile = []tile.Tile{
	tile.Room, tile.Soil, tile.Sand, tile.Stone, tile.Grass,
	tile.Tree, tile.Ice, tile.Magma, tile.Swamp, tile.Sea, tile.Smoke,
}
var allRoadTile = []tile.Tile{
	tile.Road, tile.Soil, tile.Sand, tile.Stone, tile.Grass, tile.Tree, tile.Fog,
}
var allWallTile = []tile.Tile{
	tile.Wall, tile.Window,
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
			fm := tw.Add(floorName, w, h, roomCount/2, 0, 1.0)
			for i := 0; i < roomCount; i++ {
				roomTile := allRoomTile[rnd.Intn(len(allRoomTile))]
				wallTile := allWallTile[rnd.Intn(len(allWallTile))]
				fm.Appendf(
					"AddRoomsRand bgtile=%v walltile=%v terrace=false align=1 count=1 mean=8 stddev=2 min=6",
					roomTile, wallTile)

			}
			roadTile := allRoadTile[rnd.Intn(len(allRoadTile))]
			fm.Appendf(
				"ConnectRooms tile=%v connect=1 allconnect=true diagonal=false",
				roadTile)
		case 1:
			tw.Add(floorName, 256, 256, 256, 0, 1.0).Appends(
				floortemplate.AgeingCity256x256()...,
			)

		case 3:
			tw.Add(floorName, 256, 256, 256, 0, 1.0).Appends(
				floortemplate.AgeingField256x256()...,
			)
		case 5:
			tw.Add(floorName, 256, 256, 256, 0, 1.0).Appends(
				floortemplate.AgeingMaze256x256()...,
			)
		case 7:
			tw.Add(floorName, 64, 64, 16, 0, 1.0).Appends(
				floortemplate.FreeForAll64x64()...,
			)
		case 9:
			tw.Add(floorName, 256, 256, 256, 0, 1.0).Appends(
				"ResourceFromPNG name=imagefloor.png",
				"ResourceRand resource=Plant mean=100000000 stddev=65535 repeat=256",
				"ResourceAgeing initrun=0 msper=64000 resetaftern=1440",
				"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=32 mean=6 stddev=4 min=4",
				"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
			)
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
