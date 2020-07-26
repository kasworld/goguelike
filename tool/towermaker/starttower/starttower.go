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

package starttower

import (
	"github.com/kasworld/goguelike/tool/towermaker/floormake"
)

func MakeStartTower() floormake.FloorList {
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

	floorList.GetByName("Practice").Appends(
		"AddPortalInRoom display=PortalIn acttype=PortalIn PortalID=Practice-0 DstPortalID=SoilPlant-3 message=ToSoilPlant",
		"AddPortalInRoom display=PortalOut acttype=PortalOut PortalID=Practice-1 DstPortalID=MazeWalk-1 message=FromMazeWalk",
	)
	floorList.GetByName("SoilPlant").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=SoilPlant-1 DstPortalID=FreeForAll-2 message=DownToFreeForAll",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=SoilPlant-2 DstPortalID=ManyPortals-1 message=UpToManyPortals",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=SoilPlant-3 DstPortalID=ManyPortals-102 message=FromManyPortals",
	)
	floorList.GetByName("ManyPortals").Appends(
		"AddPortalInRoom display=StairDn      acttype=PortalInOut  PortalID=ManyPortals-1   DstPortalID=SoilPlant-2 message=DownToSoilPlant",
		"AddPortalInRoom display=StairUp      acttype=PortalInOut  PortalID=ManyPortals-2   DstPortalID=SoilWater-1 message=UpToSoilWater",
		"AddPortal x=63 y=63 display=PortalIn acttype=PortalIn     PortalID=ManyPortals-101 DstPortalID=TileRooms-0 message=ToTileRooms",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-102 DstPortalID=SoilPlant-3 message=ToSoilPlant",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-104 DstPortalID=SoilWater-3 message=ToSoilWater",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-105 DstPortalID=SoilMagma-3 message=ToSoilMagma",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-106 DstPortalID=SoilIce-3 message=ToSoilIce",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-107 DstPortalID=MadeByImage-3 message=ToMadeByImage",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-108 DstPortalID=AgeingCity-3 message=ToAgeingCity",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-109 DstPortalID=AgeingField-3 message=ToAgeingField",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-110 DstPortalID=AgeingMaze-3 message=ToAgeingMaze",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-111 DstPortalID=RogueLike-3 message=ToRogueLike",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-112 DstPortalID=FreeForAll-3 message=ToFreeForAll",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-113 DstPortalID=GogueLike-3 message=ToGogueLike",
		"AddPortalInRoom display=PortalIn     acttype=PortalIn     PortalID=ManyPortals-114 DstPortalID=Ghost-3 message=ToGhost",
	)
	floorList.GetByName("SoilWater").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=SoilWater-1 DstPortalID=ManyPortals-2 message=DownToManyPortals",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=SoilWater-2 DstPortalID=SoilMagma-1 message=UpToSoilMagma",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=SoilWater-3 DstPortalID=ManyPortals-104 message=FromManyPortals",
	)
	floorList.GetByName("SoilMagma").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=SoilMagma-1 DstPortalID=SoilWater-2 message=DownToSoilWater",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=SoilMagma-2 DstPortalID=SoilIce-1 message=UpToSoilIce",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=SoilMagma-3 DstPortalID=ManyPortals-105 message=FromManyPortals",
	)
	floorList.GetByName("SoilIce").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=SoilIce-1 DstPortalID=SoilMagma-2 message=DownToSoilMagma",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=SoilIce-2 DstPortalID=MadeByImage-1 message=UpToMadeByImage",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=SoilIce-3 DstPortalID=ManyPortals-106 message=FromManyPortals",
	)
	floorList.GetByName("MadeByImage").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=MadeByImage-1 DstPortalID=SoilIce-2 message=DownToSoilIce",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=MadeByImage-2 DstPortalID=AgeingCity-1 message=UpToAgeingCity",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=MadeByImage-3 DstPortalID=ManyPortals-107 message=FromManyPortals",
	)
	floorList.GetByName("AgeingCity").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=AgeingCity-1 DstPortalID=MadeByImage-2 message=DownToMadeByImage",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=AgeingCity-2 DstPortalID=AgeingField-1 message=UpToAgeingField",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=AgeingCity-3 DstPortalID=ManyPortals-108 message=FromManyPortals",
	)
	floorList.GetByName("AgeingField").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=AgeingField-1 DstPortalID=AgeingCity-2 message=DownToAgeingCity",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=AgeingField-2 DstPortalID=AgeingMaze-1 message=UpToAgeingMaze",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=AgeingField-3 DstPortalID=ManyPortals-109 message=FromManyPortals",
	)
	floorList.GetByName("AgeingMaze").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=AgeingMaze-1 DstPortalID=AgeingField-2 message=DownToAgeingField",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=AgeingMaze-2 DstPortalID=RogueLike-1 message=UpToRogueLike",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=AgeingMaze-3 DstPortalID=ManyPortals-110 message=FromManyPortals",
	)
	floorList.GetByName("RogueLike").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=RogueLike-1 DstPortalID=AgeingMaze-2 message=DownToAgeingMaze",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=RogueLike-2 DstPortalID=GogueLike-1 message=UpToGogueLike",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=RogueLike-3 DstPortalID=ManyPortals-111 message=FromManyPortals",
	)
	floorList.GetByName("GogueLike").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=GogueLike-1 DstPortalID=RogueLike-2 message=DownToRogueLike",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=GogueLike-2 DstPortalID=Ghost-1 message=UpToGhost",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=GogueLike-3 DstPortalID=ManyPortals-113 message=FromManyPortals",
	)
	floorList.GetByName("Ghost").Appends(
		"AddPortalInRoom display=StairDn acttype=PortalInOut   PortalID=Ghost-1 DstPortalID=GogueLike-2 message=DownToGogueLike",
		"AddPortalInRoom display=StairUp acttype=PortalInOut   PortalID=Ghost-2 DstPortalID=FreeForAll-1 message=UpToFreeForAll",
		"AddPortalInRoom display=PortalOut acttype=PortalOut   PortalID=Ghost-3 DstPortalID=ManyPortals-114 message=FromManyPortals",
	)
	floorList.GetByName("FreeForAll").Appends(
		"AddPortalRand display=StairDn acttype=PortalInOut       PortalID=FreeForAll-1 DstPortalID=Ghost-2 message=DownToGhost",
		"AddPortalRand display=StairUp acttype=PortalInOut       PortalID=FreeForAll-2 DstPortalID=SoilPlant-1 message=UpToSoilPlant",
		"AddPortalRand display=PortalOut acttype=PortalOut       PortalID=FreeForAll-3 DstPortalID=ManyPortals-112 message=FromManyPortals",
	)
	floorList.GetByName("TileRooms").Appends(
		"AddPortalRand display=PortalOut acttype=PortalOut PortalID=TileRooms-0 DstPortalID=ManyPortals-101 message=FromManyPortals",
		"AddPortalRand display=PortalIn acttype=PortalIn PortalID=TileRooms-1 DstPortalID=PortalMaze-0 message=ToPortalMaze",
	)
	floorList.GetByName("PortalMaze").Appends(
		"AddPortal x=07 y=07 display=PortalOut acttype=PortalOut PortalID=PortalMaze-0 DstPortalID=TileRooms-1 message=FromTileRooms",
		"AddPortal x=55 y=23 display=PortalIn  acttype=PortalIn  PortalID=PortalMaze-1 DstPortalID=MazeRooms1-0 message=ToMazeRooms",

		"AddPortal x=04 y=04 display=PortalOut acttype=PortalOut    PortalID=PortalMaze-00-0 DstPortalID=PortalMaze-00-1 message=From",
		"AddPortal x=04 y=10 display=None      acttype=PortalAutoIn PortalID=PortalMaze-00-1 DstPortalID=PortalMaze-01-0 message=To",
		"AddPortal x=10 y=04 display=None      acttype=PortalAutoIn PortalID=PortalMaze-00-2 DstPortalID=PortalMaze-10-0 message=To",
		"AddPortal x=10 y=10 display=None      acttype=PortalAutoIn PortalID=PortalMaze-00-3 DstPortalID=PortalMaze-11-0 message=To",

		"AddPortal x=20 y=04 display=PortalOut acttype=PortalOut    PortalID=PortalMaze-10-0 DstPortalID=PortalMaze-10-1 message=From",
		"AddPortal x=20 y=10 display=None      acttype=PortalAutoIn PortalID=PortalMaze-10-1 DstPortalID=PortalMaze-11-0 message=To",
		"AddPortal x=26 y=04 display=None      acttype=PortalAutoIn PortalID=PortalMaze-10-2 DstPortalID=PortalMaze-20-0 message=To",
		"AddPortal x=26 y=10 display=None      acttype=PortalAutoIn PortalID=PortalMaze-10-3 DstPortalID=PortalMaze-21-0 message=To",

		"AddPortal x=36 y=04 display=PortalOut acttype=PortalOut    PortalID=PortalMaze-20-0 DstPortalID=PortalMaze-20-1 message=From",
		"AddPortal x=36 y=10 display=None      acttype=PortalAutoIn PortalID=PortalMaze-20-1 DstPortalID=PortalMaze-21-0 message=To",
		"AddPortal x=42 y=04 display=None      acttype=PortalAutoIn PortalID=PortalMaze-20-2 DstPortalID=PortalMaze-30-0 message=To",
		"AddPortal x=42 y=10 display=None      acttype=PortalAutoIn PortalID=PortalMaze-20-3 DstPortalID=PortalMaze-31-0 message=To",

		"AddPortal x=52 y=04 display=PortalOut acttype=PortalOut    PortalID=PortalMaze-30-0 DstPortalID=PortalMaze-30-1 message=From",
		"AddPortal x=52 y=10 display=None      acttype=PortalAutoIn PortalID=PortalMaze-30-1 DstPortalID=PortalMaze-31-0 message=To",
		"AddPortal x=58 y=04 display=None      acttype=PortalAutoIn PortalID=PortalMaze-30-2 DstPortalID=PortalMaze-00-0 message=To",
		"AddPortal x=58 y=10 display=None      acttype=PortalAutoIn PortalID=PortalMaze-30-3 DstPortalID=PortalMaze-01-0 message=To",

		"AddPortal x=04 y=20 display=PortalOut acttype=PortalOut    PortalID=PortalMaze-01-0 DstPortalID=PortalMaze-01-1 message=From",
		"AddPortal x=04 y=26 display=None      acttype=PortalAutoIn PortalID=PortalMaze-01-1 DstPortalID=PortalMaze-00-0 message=To",
		"AddPortal x=10 y=20 display=None      acttype=PortalAutoIn PortalID=PortalMaze-01-2 DstPortalID=PortalMaze-11-0 message=To",
		"AddPortal x=10 y=26 display=None      acttype=PortalAutoIn PortalID=PortalMaze-01-3 DstPortalID=PortalMaze-10-0 message=To",

		"AddPortal x=20 y=20 display=PortalOut acttype=PortalOut    PortalID=PortalMaze-11-0 DstPortalID=PortalMaze-11-1 message=From",
		"AddPortal x=20 y=26 display=None      acttype=PortalAutoIn PortalID=PortalMaze-11-1 DstPortalID=PortalMaze-10-0 message=To",
		"AddPortal x=26 y=20 display=None      acttype=PortalAutoIn PortalID=PortalMaze-11-2 DstPortalID=PortalMaze-21-0 message=To",
		"AddPortal x=26 y=26 display=None      acttype=PortalAutoIn PortalID=PortalMaze-11-3 DstPortalID=PortalMaze-20-0 message=To",

		"AddPortal x=36 y=20 display=PortalOut acttype=PortalOut    PortalID=PortalMaze-21-0 DstPortalID=PortalMaze-21-1 message=From",
		"AddPortal x=36 y=26 display=None      acttype=PortalAutoIn PortalID=PortalMaze-21-1 DstPortalID=PortalMaze-20-0 message=To",
		"AddPortal x=42 y=20 display=None      acttype=PortalAutoIn PortalID=PortalMaze-21-2 DstPortalID=PortalMaze-31-0 message=To",
		"AddPortal x=42 y=26 display=None      acttype=PortalAutoIn PortalID=PortalMaze-21-3 DstPortalID=PortalMaze-30-0 message=To",

		"AddPortal x=52 y=20 display=PortalOut acttype=PortalOut    PortalID=PortalMaze-31-0 DstPortalID=PortalMaze-31-1 message=From",
		"AddPortal x=52 y=26 display=None      acttype=PortalAutoIn PortalID=PortalMaze-31-1 DstPortalID=PortalMaze-30-0 message=To",
		"AddPortal x=58 y=20 display=None      acttype=PortalAutoIn PortalID=PortalMaze-31-2 DstPortalID=PortalMaze-01-0 message=To",
		"AddPortal x=58 y=26 display=None      acttype=PortalAutoIn PortalID=PortalMaze-31-3 DstPortalID=PortalMaze-00-0 message=To",
	)
	floorList.GetByName("MazeRooms1").Appends(
		"AddPortal x=15 y=15 display=PortalOut acttype=PortalOut PortalID=MazeRooms1-0 DstPortalID=PortalMaze-1 message=FromPortalMaze",
		"AddPortal x=47 y=15 display=PortalIn  acttype=PortalIn  PortalID=MazeRooms1-1 DstPortalID=MazeRooms2-0 message=ToMazeRooms2",
	)
	floorList.GetByName("MazeRooms2").Appends(
		"AddPortalInRoom display=PortalOut acttype=PortalOut PortalID=MazeRooms2-0 DstPortalID=MazeRooms1-1 message=FromMazeRooms1",
		"AddPortalInRoom display=PortalIn  acttype=PortalIn  PortalID=MazeRooms2-1 DstPortalID=MazeRooms3-0 message=ToMazeRooms3",
	)
	floorList.GetByName("MazeRooms3").Appends(
		"AddPortalInRoom display=PortalOut acttype=PortalOut PortalID=MazeRooms3-0 DstPortalID=MazeRooms2-1 message=FromMazeRooms2",
		"AddPortalInRoom display=PortalIn  acttype=PortalIn  PortalID=MazeRooms3-1 DstPortalID=MazeWalk-0 message=ToMazeWalk",
	)
	floorList.GetByName("MazeWalk").Appends(
		"AddPortalRand display=PortalOut acttype=PortalOut PortalID=MazeWalk-0 DstPortalID=MazeRooms3-1 message=FromMazeRooms3",
		"AddPortalRand display=PortalIn  acttype=PortalIn  PortalID=MazeWalk-1 DstPortalID=Practice-1 message=ToPractice",
	)

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
		fm.ConnectAutoInPortalTo(suffix, fm)
		fm.AddTrapTeleportTo(suffix, fm)
		fm.AddAllEffectTrap(suffix, 1)
	}

	return floorList
}
