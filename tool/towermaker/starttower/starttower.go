// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
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
	"fmt"

	"github.com/kasworld/goguelike/enum/decaytype"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/tool/towermaker/floortemplate"
	"github.com/kasworld/goguelike/tool/towermaker/towermake"
)

func New(name string) *towermake.Tower {
	var rnd = g2rand.New()

	tw := towermake.New(name)
	tw.Add("Practice", 64, 32, 0.7).Appends(
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=6 stddev=4",
		"ConnectRooms tile=Road connect=2 allconnect=true diagonal=false",
	)
	tw.Add("SoilPlant", 64, 64, 1.0).Appends(
		floortemplate.SoilPlant64x64()...,
	).Appendf("ActiveObjectsRand count=%v", 16)

	tw.Add("ManyPortals", 128, 128, 1.0).Appends(
		floortemplate.BGTile9Rooms128x128()...,
	).Appendf("ActiveObjectsRand count=%v", 64)

	tw.Add("SoilWater", 128, 128, 1.0).Appends(
		floortemplate.SoilWater128x128()...,
	).Appendf("ActiveObjectsRand count=%v", 64)

	tw.Add("SoilMagma", 128, 128, 1.0).Appends(
		floortemplate.SoilMagma128x128()...,
	).Appendf("ActiveObjectsRand count=%v", 64)

	tw.Add("SoilIce", 128, 128, 1.0).Appends(
		floortemplate.SoilIce128x128()...,
	).Appendf("ActiveObjectsRand count=%v", 64)

	tw.Add("MadeByImage", 256, 256, 1.0).Appends(
		"ResourceFromPNG name=imagefloor.png",
		"ResourceRand resource=Plant mean=100000000 stddev=65535 repeat=256",
		"ResourceAgeing initrun=0 msper=64000 resetaftern=1440",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=32 mean=6 stddev=4",
		"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
	).Appendf("ActiveObjectsRand count=%v", 256)

	tw.Add("AgeingCity", 256, 256, 1.0).Appends(
		floortemplate.AgeingCitySize(256)...,
	).Appendf("ActiveObjectsRand count=%v", 256)

	tw.Add("AgeingField", 512, 512, 1.0).Appends(
		floortemplate.AgeingFieldSize(512)...,
	).Appendf("ActiveObjectsRand count=%v", 512)

	tw.Add("AgeingMaze", 1024, 1024, 1.0).Appends(
		floortemplate.AgeingMazeSize(1024)...,
	).Appendf("ActiveObjectsRand count=%v", 1024)

	fm := tw.Add("BedTown", 256, 256, 1.0).Appends(
		floortemplate.CityRooms(256, 256, 11, 11, 5, rnd.Intn)...,
	).Appendf("ActiveObjectsRand count=%v", 128)

	fm = tw.Add("ResourceMazeFill", 256, 256, 1.0).Appends(
		fmt.Sprintf("ResourceFillRect resource=Soil amount=1 x=0 y=0 w=%v h=%v", 256, 256),
	).Appendf("ActiveObjectsRand count=%v", 128)
	fm.Appends(
		floortemplate.MixedResourceMaze(256, 256)...,
	)
	fm.Appends(
		floortemplate.CityRoomsRand(512, rnd.Intn)...,
	)

	fm = tw.Add("ResourceMaze", 256, 256, 1.0).Appends(
		floortemplate.MixedResourceMaze(256, 256)...,
	).Appendf("ActiveObjectsRand count=%v", 128)
	fm.Appends(
		floortemplate.CityRoomsRand(512, rnd.Intn)...,
	)

	lhSize := gameconst.ViewPortW * 10
	fm = tw.Add("MovingDanger", lhSize, lhSize, 1.0).Appendf(
		"ResourceFillRect resource=Soil amount=1 x=0 y=0 w=%v h=%v",
		lhSize, lhSize,
	)
	fm.Appends(
		floortemplate.MixedResourceMaze(lhSize, lhSize)...,
	)

	tw.Add("RogueLike", 80, 43, 1.0).Appends(
		"AddRoomsRand bgtile=Soil walltile=Wall terrace=false align=1 count=12 mean=8 stddev=4",
		"ConnectRooms tile=Soil connect=2 allconnect=true diagonal=false",
	).Appendf("ActiveObjectsRand count=%v", 16)

	tw.Add("GogueLike", 80, 43, 1.0).Appends(
		floortemplate.GoguelikeRand(12, rnd.Intn)...,
	).Appendf("ActiveObjectsRand count=%v", 32)

	tw.Add("Ghost", 80, 43, 1.0).Appends(
		"AddRoomsRand bgtile=Smoke walltile=Window terrace=false align=1 count=12 mean=8 stddev=4",
		"ConnectRooms tile=Fog connect=2 allconnect=true diagonal=false",
	).Appendf("ActiveObjectsRand count=%v", 16)

	tw.Add("FreeForAll", 64, 64, 1.0).Appends(
		floortemplate.FreeForAll64x64()...,
	).Appendf("ActiveObjectsRand count=%v", 16)

	tw.Add("TileRooms", 64, 32, 1.5).Appends(
		"AddRoom bgtile=Stone walltile=Wall terrace=false  x=0  y=0  w=16 h=16",
		"AddRoom bgtile=Sea walltile=Wall   terrace=false  x=16 y=0  w=16 h=16",
		"AddRoom bgtile=Sand walltile=Wall  terrace=false  x=32 y=0  w=16 h=16",
		"AddRoom bgtile=Magma walltile=Wall terrace=false  x=48 y=0  w=16 h=16",
		"AddRoom bgtile=Ice walltile=Wall   terrace=false  x=0  y=16 w=16 h=16",
		"AddRoom bgtile=Grass walltile=Wall terrace=false  x=16 y=16 w=16 h=16",
		"AddRoom bgtile=Swamp walltile=Wall terrace=false  x=32 y=16 w=16 h=16",
		"AddRoom bgtile=Soil walltile=Wall  terrace=false  x=48 y=16 w=16 h=16",
	)

	tw.Add("PortalMaze", 64, 32, 1.5).Appends(
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
		"FinalizeTerrain",
		"",
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
	tw.Add("MazeRooms1", 64, 32, 1.5).Appends(
		"ResourceFillRect resource=Soil amount=64  x=0  y=0  w=64 h=32",
		"AddRoomMaze bgtile=Room walltile=Wall terrace=false  x=0   y=0   w=33 h=31 xn=16 yn=15 connerfill=true",
		"AddRoomMaze bgtile=Road walltile=Wall terrace=false  x=40  y=8  w=17 h=17 xn=8 yn=8 connerfill=true",
	)
	tw.Add("MazeRooms2", 64, 32, 1.5).Appends(
		"ResourceFillRect resource=Soil amount=64  x=0  y=0  w=64 h=32",
		"AddRoomMaze bgtile=Stone walltile=Wall terrace=false  x=0  y=0  w=15 h=15 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Sea walltile=Wall   terrace=false  x=16 y=0  w=15 h=15 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Sand walltile=Wall  terrace=false  x=32 y=0  w=15 h=15 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Magma walltile=Wall terrace=false  x=48 y=0  w=15 h=15 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Ice walltile=Wall   terrace=false  x=0  y=16 w=15 h=15 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Grass walltile=Wall terrace=false  x=16 y=16 w=15 h=15 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Swamp walltile=Wall terrace=false  x=32 y=16 w=15 h=15 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Tree walltile=Wall  terrace=false  x=48 y=16 w=15 h=15 xn=7 yn=7 connerfill=true",
	)
	tw.Add("MazeRooms3", 64, 32, 1.5).Appends(
		"ResourceFillRect resource=Soil amount=64  x=0  y=0  w=64 h=32",
		"AddRoomMaze bgtile=Stone walltile=Wall terrace=false  x=0  y=0  w=16 h=16 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Sea walltile=Wall   terrace=false  x=16 y=0  w=16 h=16 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Sand walltile=Wall  terrace=false  x=32 y=0  w=16 h=16 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Magma walltile=Wall terrace=false  x=48 y=0  w=16 h=16 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Ice walltile=Wall   terrace=false  x=0  y=16 w=16 h=16 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Grass walltile=Wall terrace=false  x=16 y=16 w=16 h=16 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Swamp walltile=Wall terrace=false  x=32 y=16 w=16 h=16 xn=7 yn=7 connerfill=true",
		"AddRoomMaze bgtile=Tree walltile=Wall  terrace=false  x=48 y=16 w=16 h=16 xn=7 yn=7 connerfill=true",
	)
	tw.Add("MazeWalk", 64, 64, 2.0).Appends(
		"ResourceMazeWalk resource=Soil amount=64 x=0 y=0 w=64 h=64 xn=16 yn=16 connerfill=true",
	)

	for _, fm := range tw.GetList() {
		if !fm.IsFinalizeTerrain() {
			fm.Appends("FinalizeTerrain", "")
		}
	}

	tw.GetByName("ManyPortals").ConnectPortalTo("Rand", "Rand", tw.GetByName("TileRooms"))
	tw.GetByName("TileRooms").ConnectPortalTo("Rand", " x=07 y=07", tw.GetByName("PortalMaze"))
	tw.GetByName("PortalMaze").ConnectPortalTo(" x=55 y=23", " x=15 y=15", tw.GetByName("MazeRooms1"))
	tw.GetByName("MazeRooms1").ConnectPortalTo(" x=47 y=15", "InRoom", tw.GetByName("MazeRooms2"))
	tw.GetByName("MazeRooms2").ConnectPortalTo("Rand", "Rand", tw.GetByName("MazeRooms3"))
	tw.GetByName("MazeRooms3").ConnectPortalTo("Rand", "Rand", tw.GetByName("MazeWalk"))
	tw.GetByName("MazeWalk").ConnectPortalTo("Rand", "Rand", tw.GetByName("ManyPortals"))

	tw.GetByName("Practice").ConnectStairUp("InRoom", "InRoom", tw.GetByName("SoilPlant"))
	tw.GetByName("SoilPlant").ConnectStairUp("InRoom", "InRoom", tw.GetByName("AgeingCity"))
	tw.GetByName("AgeingCity").ConnectStairUp("InRoom", "InRoom", tw.GetByName("BedTown"))
	tw.GetByName("BedTown").ConnectStairUp("InRoom", "InRoom", tw.GetByName("RogueLike"))

	tw.GetByName("RogueLike").ConnectStairUp("InRoom", "InRoom", tw.GetByName("SoilWater"))
	tw.GetByName("SoilWater").ConnectStairUp("InRoom", "InRoom", tw.GetByName("AgeingField"))
	tw.GetByName("AgeingField").ConnectStairUp("InRoom", "InRoom", tw.GetByName("ResourceMazeFill"))
	tw.GetByName("ResourceMazeFill").ConnectStairUp("InRoom", "Rand", tw.GetByName("MovingDanger"))
	tw.GetByName("MovingDanger").ConnectStairUp("Rand", "InRoom", tw.GetByName("GogueLike"))

	tw.GetByName("GogueLike").ConnectStairUp("InRoom", "InRoom", tw.GetByName("SoilMagma"))
	tw.GetByName("SoilMagma").ConnectStairUp("InRoom", "InRoom", tw.GetByName("AgeingMaze"))
	tw.GetByName("AgeingMaze").ConnectStairUp("InRoom", "InRoom", tw.GetByName("ResourceMaze"))
	tw.GetByName("ResourceMaze").ConnectStairUp("InRoom", "InRoom", tw.GetByName("Ghost"))

	tw.GetByName("Ghost").ConnectStairUp("InRoom", "Rand", tw.GetByName("SoilIce"))
	tw.GetByName("SoilIce").ConnectStairUp("InRoom", "InRoom", tw.GetByName("MadeByImage"))
	tw.GetByName("MadeByImage").ConnectStairUp("InRoom", "Rand", tw.GetByName("FreeForAll"))
	tw.GetByName("FreeForAll").ConnectStairUp("Rand", "InRoom", tw.GetByName("Practice"))

	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("Practice"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("SoilPlant"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("SoilWater"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("SoilMagma"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("SoilIce"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("MadeByImage"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("AgeingCity"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("AgeingField"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("AgeingMaze"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("RogueLike"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("GogueLike"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("Ghost"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("FreeForAll"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("MovingDanger"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("BedTown"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("ResourceMazeFill"))
	tw.GetByName("ManyPortals").ConnectStairUp("Rand", "Rand", tw.GetByName("ResourceMaze"))

	fm = tw.GetByName("MovingDanger")
	perturn := 10
	lhCount := 0
	gwCount := 0
	lhlen := gameconst.ViewPortW / 2

	for x := 0; x < fm.W; x += lhlen * 2 {
		for y := 0; y < fm.H; y += lhlen * 2 {
			decay := decaytype.DecayType(rnd.Intn(decaytype.DecayType_Count))
			fm.Appendf(
				"AddRotateLineAttack x=%v y=%v display=RotateLineAttack winglen=%v wingcount=1 degree=0 perturn=%v decay=%v message=RotDanger1",
				x, y, lhlen, perturn, decay,
			)
			perturn = -perturn
			lhCount++
		}
	}
	for x := lhlen; x < fm.W; x += lhlen * 2 {
		for y := lhlen; y < fm.H; y += lhlen * 2 {
			decay := decaytype.DecayType(rnd.Intn(decaytype.DecayType_Count))
			fm.Appendf(
				"AddRotateLineAttack x=%v y=%v display=RotateLineAttack winglen=%v wingcount=2 degree=0 perturn=%v decay=%v message=RotDanger2",
				x, y, lhlen/2, perturn, decay,
			)
			perturn = -perturn
			gwCount++
		}
	}
	fm.Appendf(
		"AddMine%v display=None decay=%v count=%v message=Mine",
		"Rand", "Decrease", (lhCount+gwCount)/3,
	)
	fm.Appendf(
		"AddMine%v display=None decay=%v count=%v message=Mine",
		"Rand", "Even", (lhCount+gwCount)/3,
	)
	fm.Appendf(
		"AddMine%v display=None decay=%v count=%v message=Mine",
		"Rand", "Increase", (lhCount+gwCount)/3,
	)

	for _, fm := range tw.GetList() {
		noRoomFloorNames := map[string]bool{
			"FreeForAll":   true,
			"MovingDanger": true,
			"PortalMaze":   true,
			"MazeRooms1":   true,
			"MazeRooms2":   true,
			"MazeRooms3":   true,
			"MazeWalk":     true,
		}
		suffix := "InRoom"
		if noRoomFloorNames[fm.Name] {
			suffix = "Rand"
		}

		roomCount := fm.CalcRoomCount()
		// fmt.Printf("%v Room %v\n", fm, roomCount)
		recycleCount := fm.W * fm.H / gameconst.ViewPortWH / 8
		if recycleCount < 2 {
			recycleCount = 2
		}

		fm.ConnectAutoInPortalTo(suffix, suffix, fm)
		fm.AddTrapTeleportTo(suffix, fm)
		for i := 0; i < recycleCount/8; i++ {
			fm.ConnectAutoInPortalTo("Rand", "Rand", fm)
			fm.AddTrapTeleportTo("Rand", fm)
		}
		fm.AddAllEffectTrap(suffix, 1)
		if count := recycleCount / 32; count > 0 {
			fm.AddAllEffectTrap("Rand", count)
		}

		if recycleCount > roomCount {
			fm.AddRecycler(suffix, roomCount)
		} else {
			fm.AddRecycler(suffix, recycleCount)
		}
		if recycleCount-roomCount > 0 {
			fm.AddRecycler("Rand", recycleCount-roomCount)
		}
		for j := 0; j < decaytype.DecayType_Count; j++ {
			decay := decaytype.DecayType(j)
			fm.Appendf(
				"AddRotateLineAttack%v display=RotateLineAttack winglen=%v wingcount=1 degree=0 perturn=10 decay=%v count=%v message=RotDanger1",
				"Rand", lhlen, decay, 1,
			)
			fm.Appendf(
				"AddRotateLineAttack%v display=RotateLineAttack winglen=%v wingcount=2 degree=0 perturn=10 decay=%v count=%v message=RotDanger2",
				"Rand", lhlen/2, decay, 1,
			)
			fm.Appendf(
				"AddMine%v display=None decay=%v count=%v message=Mine",
				"Rand", decay, 1,
			)
		}

	}
	return tw
}
