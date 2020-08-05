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

	"github.com/kasworld/goguelike/enum/tile"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/tool/towermaker/towermake"
)

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func makePowerOf2(v int) int {
	return 1 << uint(bits.Len(uint(v-1)))
}

func MakeRogueTower(name string, floorCount int) *towermake.Tower {
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
		tw.Add(
			fmt.Sprintf("Floor%v", i),
			w, h,
			roomCount/2, 0,
			1.0,
		)
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
	for _, fm := range tw.GetList() {
		roomCount := fm.W * fm.H / 512
		if roomCount < 2 {
			roomCount = 2
		}
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
		fm.Appends("FinalizeTerrain", "")
	}

	for i, fm := range tw.GetList() {
		roomCount := fm.W * fm.H / 512
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
