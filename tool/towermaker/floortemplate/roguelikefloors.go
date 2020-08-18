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

package floortemplate

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/tile"
)

func GogueLike() []string {
	return []string{
		"AddRoomsRand bgtile=Swamp walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Soil walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Stone walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Sand walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Sea walltile=Wall   terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Magma walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Ice walltile=Wall   terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Grass walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Tree walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Room walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Road walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"AddRoomsRand bgtile=Smoke walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4 min=4",
		"ConnectRooms tile=Fog connect=1 allconnect=true diagonal=false",
		"ConnectRooms tile=Sand connect=1 allconnect=true diagonal=false",
	}
}

func Practice64x32() []string {
	return []string{
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=6 stddev=4 min=4",
		"ConnectRooms tile=Road connect=2 allconnect=true diagonal=false",
	}
}

func RogueLike80x43() []string {
	return []string{
		"AddRoomsRand bgtile=Soil walltile=Wall terrace=false align=1 count=12 mean=8 stddev=4 min=4",
		"ConnectRooms tile=Soil connect=2 allconnect=true diagonal=false",
	}
}

func Ghost80x43() []string {
	return []string{
		"AddRoomsRand bgtile=Smoke walltile=Window terrace=false align=1 count=12 mean=8 stddev=4 min=4",
		"ConnectRooms tile=Fog connect=2 allconnect=true diagonal=false",
	}
}

func RoguelikeRand(roomCount int, intnfn func(int) int) []string {
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
	rtn := make([]string, 0)
	for i := 0; i < roomCount; i++ {
		roomTile := allRoomTile[intnfn(len(allRoomTile))]
		wallTile := allWallTile[intnfn(len(allWallTile))]
		rtn = append(rtn, fmt.Sprintf(
			"AddRoomsRand bgtile=%v walltile=%v terrace=false align=1 count=1 mean=8 stddev=2 min=6",
			roomTile, wallTile))

	}
	roadTile := allRoadTile[intnfn(len(allRoadTile))]
	rtn = append(rtn, fmt.Sprintf(
		"ConnectRooms tile=%v connect=1 allconnect=true diagonal=false",
		roadTile))

	return rtn
}

func CityRoomsRand(roomCount int, intnfn func(int) int) []string {
	var allRoomTile = []tile.Tile{
		tile.Room, tile.Soil, tile.Sand, tile.Stone, tile.Grass,
		tile.Tree, tile.Ice, tile.Magma, tile.Swamp, tile.Sea, tile.Smoke,
	}
	var allRoadTile = []tile.Tile{
		tile.Road, tile.Soil, tile.Sand, tile.Stone, tile.Grass, tile.Tree, tile.Fog,
	}
	var allWallTile = []tile.Tile{
		tile.Wall, tile.Wall, tile.Window,
	}
	rtn := make([]string, 0)
	for i := 0; i < roomCount; i++ {
		roomTile := allRoomTile[intnfn(len(allRoomTile))]
		wallTile := allWallTile[intnfn(len(allWallTile))]
		rtn = append(rtn, fmt.Sprintf(
			"AddRoomsRand bgtile=%v walltile=%v terrace=false align=16 count=1 mean=16 stddev=2 min=8",
			roomTile, wallTile))

	}
	roadTile := allRoadTile[intnfn(len(allRoadTile))]
	rtn = append(rtn, fmt.Sprintf(
		"ConnectRooms tile=%v connect=1 allconnect=true diagonal=false",
		roadTile))

	return rtn
}
