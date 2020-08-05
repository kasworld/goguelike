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

package floortemplete

func BGTile9Rooms128x128() []string {
	return []string{
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
	}
}

func AgeingCity256x256() []string {
	return []string{
		"ResourceFillRect resource=Soil amount=1 x=0  y=0  w=256 h=256",
		"ResourceRand resource=Ice   mean=1600000000 stddev=10000000 repeat=8",
		"ResourceRand resource=Water mean=100000000  stddev=10000000 repeat=256",
		"ResourceRand resource=Fire  mean=1600000000 stddev=10000000 repeat=8",
		"ResourceRand resource=Soil  mean=100000000  stddev=10000000 repeat=256",
		"ResourceRand resource=Plant mean=100000000  stddev=10000000 repeat=256",
		"ResourceAgeing initrun=10 msper=65000 resetaftern=1440",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=64 mean=8 stddev=4 min=4",
		"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
	}
}

func AgeingField256x256() []string {
	return []string{
		"ResourceFillRect resource=Soil amount=1 x=0  y=0  w=256 h=256",
		"ResourceRand resource=Ice   mean=1600000000 stddev=10000000 repeat=8",
		"ResourceRand resource=Water mean=100000000  stddev=10000000 repeat=256",
		"ResourceRand resource=Fire  mean=1600000000 stddev=10000000 repeat=8",
		"ResourceRand resource=Soil  mean=100000000  stddev=10000000 repeat=256",
		"ResourceRand resource=Plant mean=100000000  stddev=10000000 repeat=256",
		"ResourceAgeing initrun=10 msper=66000 resetaftern=1440",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=8 stddev=4 min=4",
		"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
	}
}

func AgeingMaze256x256() []string {
	return []string{
		"ResourceFillRect resource=Soil amount=1 x=0  y=0  w=256 h=256",
		"ResourceRand resource=Ice   mean=1600000000 stddev=10000000 repeat=8",
		"ResourceRand resource=Water mean=100000000  stddev=10000000 repeat=256",
		"ResourceRand resource=Fire  mean=1600000000 stddev=10000000 repeat=8",
		"ResourceRand resource=Soil  mean=100000000  stddev=10000000 repeat=256",
		"ResourceRand resource=Plant mean=100000000  stddev=10000000 repeat=256",
		"ResourceMazeWall resource=Stone amount=2000001 xn=32 yn=32 connerfill=false",
		"ResourceAgeing initrun=10 msper=67000 resetaftern=1440",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=10 stddev=4 min=4",
	}
}

func FreeForAll64x64() []string {
	return []string{
		"ResourceFillRect     resource=Soil amount=256 x=0  y=0  w=64 h=64",
		"ResourceFillEllipses resource=Plant amount=256 x=0  y=0  w=64 h=64",
		"ResourceFillEllipses resource=Water amount=256 x=8  y=8  w=48 h=48",
		"ResourceFillEllipses resource=Ice amount=256 x=16  y=16  w=32 h=32",
	}
}

// finalized floor
func PortalMaze64x32Finalized() []string {
	return []string{
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
	}
}
