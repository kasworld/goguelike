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
)

func TileRooms64x32() []string {
	return []string{
		"AddRoom bgtile=Stone walltile=Wall terrace=false  x=0  y=0  w=16 h=16",
		"AddRoom bgtile=Sea walltile=Wall   terrace=false  x=16 y=0  w=16 h=16",
		"AddRoom bgtile=Sand walltile=Wall  terrace=false  x=32 y=0  w=16 h=16",
		"AddRoom bgtile=Magma walltile=Wall terrace=false  x=48 y=0  w=16 h=16",
		"AddRoom bgtile=Ice walltile=Wall   terrace=false  x=0  y=16 w=16 h=16",
		"AddRoom bgtile=Grass walltile=Wall terrace=false  x=16 y=16 w=16 h=16",
		"AddRoom bgtile=Swamp walltile=Wall terrace=false  x=32 y=16 w=16 h=16",
		"AddRoom bgtile=Soil walltile=Wall  terrace=false  x=48 y=16 w=16 h=16",
	}
}

func GogueLike() []string {
	return []string{
		"AddRoomsRand bgtile=Swamp walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Soil walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Stone walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Sand walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Sea walltile=Wall   terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Magma walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Ice walltile=Wall   terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Grass walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Tree,Soil walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Room walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Road walltile=Wall  terrace=false align=1 count=1 mean=8 stddev=4",
		"AddRoomsRand bgtile=Smoke walltile=Wall terrace=false align=1 count=1 mean=8 stddev=4",
		"ConnectRooms tile=Fog connect=1 allconnect=true diagonal=false",
		"ConnectRooms tile=Sand connect=1 allconnect=true diagonal=false",
	}
}

func Practice64x32() []string {
	return []string{
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=6 stddev=4",
		"ConnectRooms tile=Road connect=2 allconnect=true diagonal=false",
	}
}

func RogueLike80x43() []string {
	return []string{
		"AddRoomsRand bgtile=Soil walltile=Wall terrace=false align=1 count=12 mean=8 stddev=4",
		"ConnectRooms tile=Soil connect=2 allconnect=true diagonal=false",
	}
}

func Ghost80x43() []string {
	return []string{
		"AddRoomsRand bgtile=Smoke walltile=Window terrace=false align=1 count=12 mean=8 stddev=4",
		"ConnectRooms tile=Fog connect=2 allconnect=true diagonal=false",
	}
}

var allRoomTile = []string{
	"Room,Sand", "Soil", "Sand", "Stone", "Grass",
	"Ice", "Magma", "Swamp", "Sea", "Smoke,Sand",
}
var allRoadTile = []string{
	"Road,Stone", "Soil", "Sand", "Stone", "Tree,Grass", "Grass", "Fog,Stone",
}
var allWallTile = []string{
	"Wall", "Wall", "Window,Tree",
}

func RoguelikeRand(roomCount int, intnfn func(int) int) []string {
	rtn := make([]string, 0)
	for i := 0; i < roomCount; i++ {
		roomTile := allRoomTile[intnfn(len(allRoomTile))]
		wallTile := allWallTile[intnfn(len(allWallTile))]
		rtn = append(rtn, fmt.Sprintf(
			"AddRoomsRand bgtile=%v walltile=%v terrace=false align=1 count=1 mean=8 stddev=2",
			roomTile, wallTile))

	}
	roadTile := allRoadTile[intnfn(len(allRoadTile))]
	rtn = append(rtn, fmt.Sprintf(
		"ConnectRooms tile=%v connect=1 allconnect=true diagonal=false",
		roadTile))

	return rtn
}

func CityRoomsRand(roomCount int, intnfn func(int) int) []string {
	rtn := make([]string, 0)
	for i := 0; i < roomCount; i++ {
		roomTile := allRoomTile[intnfn(len(allRoomTile))]
		wallTile := allWallTile[intnfn(len(allWallTile))]
		rtn = append(rtn, fmt.Sprintf(
			"AddRoomsRand bgtile=%v walltile=%v terrace=false align=16 count=1 mean=16 stddev=2",
			roomTile, wallTile))

	}
	roadTile := allRoadTile[intnfn(len(allRoadTile))]
	rtn = append(rtn, fmt.Sprintf(
		"ConnectRooms tile=%v connect=1 allconnect=true diagonal=false",
		roadTile))

	return rtn
}

func CityRooms(floorW, floorH, roomW, roomH, roadW int, intnfn func(int) int) []string {
	rtn := make([]string, 0)

	for x := 0; x < floorW-roomW; x += roomW + roadW {
		for y := 0; y < floorH-roomH; y += roomH + roadW {
			roomTile := allRoomTile[intnfn(len(allRoomTile))]
			wallTile := allWallTile[intnfn(len(allWallTile))]
			rtn = append(rtn, fmt.Sprintf(
				"AddRoom bgtile=%v walltile=%v terrace=false x=%v y=%v w=%v h=%v",
				roomTile, wallTile, x, y, roomW, roomH,
			))
		}
	}
	for x := roomW; x < floorW; x += roomW + roadW {
		for dx := 0; dx < roadW; dx++ {
			roadTile := allRoadTile[intnfn(len(allRoadTile))]
			rtn = append(rtn, fmt.Sprintf(
				"TileVLine tile=%v x=%v y=%v h=%v",
				roadTile, x+dx, 0, floorH,
			))
		}
	}
	for y := roomH; y < floorH; y += roomH + roadW {
		for dy := 0; dy < roadW; dy++ {
			roadTile := allRoadTile[intnfn(len(allRoadTile))]
			rtn = append(rtn, fmt.Sprintf(
				"TileHLine tile=%v x=%v y=%v w=%v",
				roadTile, 0, y+dy, floorW,
			))
		}
	}

	return rtn
}
