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

package main

import (
	"flag"
	"fmt"
	"math/bits"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/game/towerscript"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/tool/towermaker/floormake"
)

func main() {
	towername := flag.String("towername", "roguegen", "tower filename w/o .tower")
	floorcount := flag.Int("floorcount", 100, "floor count in tower")
	flag.Parse()

	fmt.Printf("make tower with towername:%v floorcount:%v\n", *towername, *floorcount)
	// makeRogueTower(*towername, *floorcount)
	makeStartTower(*towername)
}

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func makePowerOf2(v int) int {
	return 1 << uint(bits.Len(uint(v-1)))
}

func makeRogueTower(towerName string, floorCount int) {
	var rnd = g2rand.New()
	var whList = []int{
		32, 64, 128,
	}
	floorList := make(floormake.FloorList, floorCount)
	for i := range floorList {
		w := whList[rnd.Intn(len(whList))]
		h := whList[rnd.Intn(len(whList))]
		roomCount := w * h / 512
		if roomCount < 2 {
			roomCount = 2
		}
		floorList[i] = floormake.New(
			fmt.Sprintf("Floor%v", i),
			w, h,
			roomCount/2, 0,
			1.0,
		)
	}

	var allRoomTile = []string{
		"Room", "Soil", "Sand", "Stone", "Grass", "Tree", "Ice", "Magma", "Swamp", "Sea", "Smoke",
	}
	var allRoadTile = []string{
		"Road", "Soil", "Sand", "Stone", "Grass", "Tree", "Fog",
	}
	for _, fm := range floorList {
		roomCount := fm.W * fm.H / 512
		if roomCount < 2 {
			roomCount = 2
		}
		roomTile := allRoomTile[rnd.Intn(len(allRoomTile))]
		roadTile := allRoadTile[rnd.Intn(len(allRoadTile))]
		fm.Appendf(
			"AddRoomsRand bgtile=%v walltile=Wall terrace=false align=1 count=%v mean=8 stddev=2 min=6",
			roomTile, roomCount)
		fm.Appendf(
			"ConnectRooms tile=%v connect=1 allconnect=true diagonal=false",
			roadTile)
		fm.Appendf("FinalizeTerrain")
	}

	for i, fm := range floorList {
		roomCount := fm.W * fm.H / 512
		fm.ConnectStairUp("InRoom", floorList[wrapInt(i+1, floorCount)])
		fm.AddRecycler("InRoom", roomCount/2)
		fm.AddTeleportIn("InRoom", roomCount/2)
		fm.AddTeleportTrapOut("InRoom", floorList, roomCount/2)
		fm.AddEffectTrap("InRoom", roomCount/2)
	}

	tw := make(towerscript.TowerScript, 0)
	for _, fm := range floorList {
		tw = append(tw, fm.Script)
	}
	err := tw.SaveJSON(towerName + ".tower")
	if err != nil {
		g2log.Error("%v", err)
	}
}

func makeStartTower(towerName string) {
	floorList := floormake.FloorList{
		floormake.New("Practice", 32, 32, 0, 0, 0.7).Appends(
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=4 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Road connect=2 allconnect=true diagonal=false",
		),
		floormake.New("SoilPlant", 64, 64, 16, 0, 1.0).Appends(
			"ResourceFillRect resource=Soil amount=64   x=0  y=0  w=64 h=64",
			"ResourceMazeWall resource=Plant amount=2000000 xn=8  yn=8  connerfill=true",
			"ResourceMazeWall resource=Stone amount=1000000 xn=32 yn=32 connerfill=true",
			"ResourceAgeing initrun=1 msper=60000 resetaftern=60",
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=4 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		),
		floormake.New("ManyPortals", 128, 128, 64, 0, 1.0).Appends(
			"ResourceFillRect resource=Soil   amount=255    x=0  y=0  w=128 h=128",
			"ResourceRect resource=Plant  amount=1000001 x=0  y=0  w=128 h=128",
			"ResourceLine resource=Stone amount=1000001 x1=30 y1=30 x2=96 y2=96",
			"ResourceLine resource=Stone amount=1000001 x1=31 y1=97 x2=97 y2=31",

			"AddRoom bgtile=Ice walltile=Wall terrace=false    x=8  y=8  w=24 h=24",
			"AddRoom bgtile=Sand walltile=Wall terrace=false   x=8  y=96 w=24 h=24",
			"AddRoom bgtile=Magma walltile=Wall terrace=false  x=96 y=8  w=24 h=24",
			"AddRoom bgtile=Swamp walltile=Wall terrace=false  x=96 y=96 w=24 h=24",
			"AddRoom bgtile=Smoke walltile=Wall terrace=false  x=53 y=8  w=24 h=24",
			"AddRoom bgtile=Sea walltile=Wall terrace=false    x=53 y=96 w=24 h=24",
			"AddRoom bgtile=Road walltile=Wall terrace=false   x=8  y=53 w=24 h=24",
			"AddRoom bgtile=Stone walltile=Wall terrace=false  x=96 y=53 w=24 h=24",
			"AddRoom bgtile=Room walltile=Wall terrace=false   x=53 y=53 w=24 h=24",
		),
		floormake.New("SoilWater", 128, 128, 64, 0, 1.0).Appends(
			"ResourceFillRect resource=Soil amount=1000000 x=0 w=128 y=0 h=128",
			"ResourceMazeWall resource=Water amount=1000000 xn=64 yn=64 connerfill=true",
			"ResourceAgeing initrun=1 msper=61000 resetaftern=1440",
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		),
		floormake.New("SoilMagma", 128, 128, 64, 0, 1.0).Appends(
			"ResourceFillRect resource=Soil amount=1 x=0  y=0  w=128 h=128",
			"ResourceMazeWall resource=Soil amount=500000 xn=64 yn=64 connerfill=true",
			"ResourceMazeWall resource=Fire amount=1000000 xn=64 yn=64 connerfill=true",
			"ResourceAgeing initrun=1 msper=62000 resetaftern=1440",
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		),
		floormake.New("SoilIce", 128, 128, 64, 0, 1.0).Appends(
			"ResourceFillRect resource=Soil amount=256 x=0  y=0  w=128 h=128",
			"ResourceMazeWall  resource=Ice  amount=512 xn=64 yn=64 connerfill=true",
			"ResourceAgeing initrun=1 msper=63000 resetaftern=360",
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=16 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		),
		floormake.New("MadeByImage", 256, 256, 256, 0, 1.0).Appends(
			"ResourceFromPNG name=imagefloor.png",
			"ResourceRand resource=Plant mean=100000000 stddev=65535 repeat=256",
			"ResourceAgeing initrun=0 msper=64000 resetaftern=1440",
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=32 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		),
		floormake.New("AgeingCity", 256, 256, 256, 0, 1.0).Appends(
			"ResourceFillRect resource=Soil amount=1 x=0  y=0  w=256 h=256",
			"ResourceRand resource=Ice   mean=1600000000 stddev=10000000 repeat=8",
			"ResourceRand resource=Water mean=100000000  stddev=10000000 repeat=256",
			"ResourceRand resource=Fire  mean=1600000000 stddev=10000000 repeat=8",
			"ResourceRand resource=Soil  mean=100000000  stddev=10000000 repeat=256",
			"ResourceRand resource=Plant mean=100000000  stddev=10000000 repeat=256",
			"ResourceAgeing initrun=10 msper=65000 resetaftern=1440",
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=64 mean=8 stddev=4 min=4",
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		),
		floormake.New("AgeingField", 256, 256, 256, 0, 1.0).Appends(
			"ResourceFillRect resource=Soil amount=1 x=0  y=0  w=256 h=256",
			"ResourceRand resource=Ice   mean=1600000000 stddev=10000000 repeat=8",
			"ResourceRand resource=Water mean=100000000  stddev=10000000 repeat=256",
			"ResourceRand resource=Fire  mean=1600000000 stddev=10000000 repeat=8",
			"ResourceRand resource=Soil  mean=100000000  stddev=10000000 repeat=256",
			"ResourceRand resource=Plant mean=100000000  stddev=10000000 repeat=256",
			"ResourceAgeing initrun=10 msper=66000 resetaftern=1440",
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=8 stddev=4 min=4",
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		),
		floormake.New("AgeingMaze", 256, 256, 256, 0, 1.0).Appends(
			"ResourceFillRect resource=Soil amount=1 x=0  y=0  w=256 h=256",
			"ResourceRand resource=Ice   mean=1600000000 stddev=10000000 repeat=8",
			"ResourceRand resource=Water mean=100000000  stddev=10000000 repeat=256",
			"ResourceRand resource=Fire  mean=1600000000 stddev=10000000 repeat=8",
			"ResourceRand resource=Soil  mean=100000000  stddev=10000000 repeat=256",
			"ResourceRand resource=Plant mean=100000000  stddev=10000000 repeat=256",
			"ResourceMazeWall resource=Stone amount=2000001 xn=32 yn=32 connerfill=false",
			"ResourceAgeing initrun=10 msper=67000 resetaftern=1440",
			"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=10 stddev=4 min=4",
		),
		floormake.New("RogueLike", 64, 32, 16, 0, 1.0).Appends(
			"AddRoomsRand bgtile=Soil walltile=Wall terrace=false align=1 count=4 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Soil connect=2 allconnect=true diagonal=false",
		),
		floormake.New("GogueLike", 64, 32, 32, 0, 1.0).Appends(
			"AddRoomsRand bgtile=Swamp walltile=Wall terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Soil walltile=Wall  terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Stone walltile=Wall terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Sand walltile=Wall  terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Sea walltile=Wall   terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Magma walltile=Wall terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Ice walltile=Wall   terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Grass walltile=Wall terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Tree walltile=Wall  terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Room walltile=Wall  terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Road walltile=Wall  terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"AddRoomsRand bgtile=Smoke walltile=Wall terrace=false align=1 count=1 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Fog connect=1 allconnect=true diagonal=false",
			"ConnectRooms tile=Sand connect=1 allconnect=true diagonal=false",
		),
		floormake.New("Ghost", 64, 32, 16, 0, 1.0).Appends(
			"AddRoomsRand bgtile=Smoke walltile=Wall terrace=false align=1 count=4 mean=6 stddev=4 min=4",
			"ConnectRooms tile=Fog connect=2 allconnect=true diagonal=false",
		),
		floormake.New("FreeForAll", 64, 64, 16, 0, 1.0).Appends(
			"ResourceFillRect     resource=Soil amount=256 x=0  y=0  w=64 h=64",
			"ResourceFillEllipses resource=Plant amount=256 x=0  y=0  w=64 h=64",
			"ResourceFillEllipses resource=Water amount=256 x=16  y=16  w=32 h=32",
		),
		floormake.New("TileRooms", 64, 32, 0, 0, 1.5).Appends(
			"AddRoom bgtile=Stone walltile=Wall terrace=false  x=0  y=0  w=16 h=16",
			"AddRoom bgtile=Sea walltile=Wall   terrace=false  x=16 y=0  w=16 h=16",
			"AddRoom bgtile=Sand walltile=Wall  terrace=false  x=32 y=0  w=16 h=16",
			"AddRoom bgtile=Magma walltile=Wall terrace=false  x=48 y=0  w=16 h=16",
			"AddRoom bgtile=Ice walltile=Wall   terrace=false  x=0  y=16 w=16 h=16",
			"AddRoom bgtile=Grass walltile=Wall terrace=false  x=16 y=16 w=16 h=16",
			"AddRoom bgtile=Swamp walltile=Wall terrace=false  x=32 y=16 w=16 h=16",
			"AddRoom bgtile=Soil walltile=Wall  terrace=false  x=48 y=16 w=16 h=16",
		),
		floormake.New("PortalMaze", 64, 32, 0, 0, 1.5).Appends(
			"ResourceFillRect  resource=Soil  amount=1024 x=00 y=00 w=16 h=16",
			"ResourceFillRect  resource=Plant amount=1024 x=00 y=16 w=16 h=16",
			"ResourceFillRect  resource=Stone amount=1024 x=16 y=00 w=16 h=16",
			"ResourceFillRect  resource=Ice   amount=1024 x=16 y=16 w=16 h=16",

			"ResourceFillRect  resource=Soil  amount=1024 x=32 y=00 w=16 h=16",
			"ResourceFillRect  resource=Water amount=1024 x=32 y=16 w=16 h=16",
			"ResourceFillRect  resource=Stone amount=1024 x=48 y=00 w=16 h=16",
			"ResourceFillRect  resource=Fire  amount=2000 x=48 y=16 w=16 h=16",

			"ResourceHLine resource=Stone amount=1000001 x=0  y=15 w=64",
			"ResourceHLine resource=Stone amount=1000001 x=0  y=31 w=64",
			"ResourceVLine resource=Stone amount=1000001 x=15 y=0  h=32",
			"ResourceVLine resource=Stone amount=1000001 x=31 y=0  h=32",
			"ResourceVLine resource=Stone amount=1000001 x=47 y=0  h=32",
			"ResourceVLine resource=Stone amount=1000001 x=63 y=0  h=32",
		),
		floormake.New("MazeRooms1", 64, 32, 0, 0, 1.5).Appends(
			"ResourceFillRect resource=Soil amount=64  x=0  y=0  w=64 h=32",
			"AddRoomMaze bgtile=Room walltile=Wall terrace=false  x=0   y=0   w=33 h=31 xn=16 yn=15 connerfill=true",
			"AddRoomMaze bgtile=Road walltile=Wall terrace=false  x=40  y=8  w=17 h=17 xn=8 yn=8 connerfill=true",
		),
		floormake.New("MazeRooms2", 64, 32, 0, 0, 1.5).Appends(
			"ResourceFillRect resource=Soil amount=64  x=0  y=0  w=64 h=32",
			"AddRoomMaze bgtile=Stone walltile=Wall terrace=false  x=0  y=0  w=15 h=15 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Sea walltile=Wall   terrace=false  x=16 y=0  w=15 h=15 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Sand walltile=Wall  terrace=false  x=32 y=0  w=15 h=15 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Magma walltile=Wall terrace=false  x=48 y=0  w=15 h=15 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Ice walltile=Wall   terrace=false  x=0  y=16 w=15 h=15 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Grass walltile=Wall terrace=false  x=16 y=16 w=15 h=15 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Swamp walltile=Wall terrace=false  x=32 y=16 w=15 h=15 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Tree walltile=Wall  terrace=false  x=48 y=16 w=15 h=15 xn=7 yn=7 connerfill=true",
		),
		floormake.New("MazeRooms3", 64, 32, 0, 0, 1.5).Appends(
			"ResourceFillRect resource=Soil amount=64  x=0  y=0  w=64 h=32",
			"AddRoomMaze bgtile=Stone walltile=Wall terrace=false  x=0  y=0  w=16 h=16 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Sea walltile=Wall   terrace=false  x=16 y=0  w=16 h=16 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Sand walltile=Wall  terrace=false  x=32 y=0  w=16 h=16 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Magma walltile=Wall terrace=false  x=48 y=0  w=16 h=16 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Ice walltile=Wall   terrace=false  x=0  y=16 w=16 h=16 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Grass walltile=Wall terrace=false  x=16 y=16 w=16 h=16 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Swamp walltile=Wall terrace=false  x=32 y=16 w=16 h=16 xn=7 yn=7 connerfill=true",
			"AddRoomMaze bgtile=Tree walltile=Wall  terrace=false  x=48 y=16 w=16 h=16 xn=7 yn=7 connerfill=true",
		),
		floormake.New("MazeWalk", 64, 64, 0, 0, 2.0).Appends(
			"ResourceMazeWalk resource=Soil amount=64 xn=16 yn=16 connerfill=true",
		),
	}

	for _, fm := range floorList {
		fm.Appends("FinalizeTerrain")
	}

	for _, fm := range floorList {
		randList := map[string]bool{
			"FreeForAll": true,
			"PortalMaze": true,
			"MazeRooms1": true,
			"MazeRooms2": true,
			"MazeRooms3": true,
			"MazeWalk":   true,
		}
		suffix := "InRoom"
		if randList[fm.Name] {
			suffix = "Rand"
		}
		fm.AddRecycler(suffix, 4)
		fm.AddTeleportIn(suffix, 1)
		fm.AddAllEffectTrap(suffix, 1)
	}

	tw := make(towerscript.TowerScript, 0)
	for _, fm := range floorList {
		tw = append(tw, fm.Script)
	}
	err := tw.SaveJSON(towerName + ".tower")
	if err != nil {
		g2log.Error("%v", err)
	}
}
