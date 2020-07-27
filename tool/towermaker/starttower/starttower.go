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
			"",
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
			"",
			"ResourceFillRect  resource=Soil  amount=1024 x=32 y=00 w=16 h=16",
			"ResourceFillRect  resource=Water amount=1024 x=32 y=16 w=16 h=16",
			"ResourceFillRect  resource=Stone amount=1024 x=48 y=00 w=16 h=16",
			"ResourceFillRect  resource=Fire  amount=2000 x=48 y=16 w=16 h=16",
			"",
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
		fm.Appends("FinalizeTerrain", "")
	}

	floorList.GetByName("SoilPlant").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("ManyPortals"))
	floorList.GetByName("ManyPortals").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("SoilWater"))
	floorList.GetByName("SoilWater").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("SoilMagma"))
	floorList.GetByName("SoilMagma").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("SoilIce"))
	floorList.GetByName("SoilIce").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("MadeByImage"))
	floorList.GetByName("MadeByImage").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("AgeingCity"))
	floorList.GetByName("AgeingCity").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("AgeingField"))
	floorList.GetByName("AgeingField").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("AgeingMaze"))
	floorList.GetByName("AgeingMaze").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("RogueLike"))
	floorList.GetByName("RogueLike").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("GogueLike"))
	floorList.GetByName("GogueLike").ConnectStairUp("InRoom", "InRoom", floorList.GetByName("Ghost"))
	floorList.GetByName("Ghost").ConnectStairUp("InRoom", "Rand", floorList.GetByName("FreeForAll"))
	floorList.GetByName("FreeForAll").ConnectStairUp("Rand", "InRoom", floorList.GetByName("SoilPlant"))

	floorList.GetByName("Practice").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("SoilPlant"))
	floorList.GetByName("TileRooms").ConnectPortalTo("Rand", " x=07 y=07", floorList.GetByName("PortalMaze"))
	floorList.GetByName("PortalMaze").ConnectPortalTo(" x=55 y=23", " x=15 y=15", floorList.GetByName("MazeRooms1"))
	floorList.GetByName("MazeRooms1").ConnectPortalTo(" x=47 y=15", "InRoom", floorList.GetByName("MazeRooms2"))
	floorList.GetByName("MazeRooms2").ConnectPortalTo("Rand", "Rand", floorList.GetByName("MazeRooms3"))
	floorList.GetByName("MazeRooms3").ConnectPortalTo("Rand", "Rand", floorList.GetByName("MazeWalk"))
	floorList.GetByName("MazeWalk").ConnectPortalTo("Rand", "InRoom", floorList.GetByName("Practice"))

	floorList.GetByName("ManyPortals").ConnectPortalTo(" x=63 y=63", "Rand", floorList.GetByName("TileRooms"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("SoilPlant"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("SoilWater"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("SoilMagma"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("SoilIce"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("MadeByImage"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("AgeingCity"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("AgeingField"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("AgeingMaze"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("RogueLike"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("GogueLike"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "InRoom", floorList.GetByName("Ghost"))
	floorList.GetByName("ManyPortals").ConnectPortalTo("InRoom", "Rand", floorList.GetByName("FreeForAll"))

	floorList.GetByName("PortalMaze").Appends(
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
		fm.ConnectAutoInPortalTo(suffix, suffix, fm)
		fm.AddTrapTeleportTo(suffix, fm)
		fm.AddAllEffectTrap(suffix, 1)
	}

	return floorList
}
