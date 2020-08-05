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

func SoilPlant64x64() []string {
	return []string{
		"ResourceFillRect resource=Soil amount=64   x=0  y=0  w=64 h=64",
		"ResourceMazeWall resource=Plant amount=2000000 xn=8  yn=8  connerfill=true",
		"ResourceMazeWall resource=Stone amount=1000000 xn=32 yn=32 connerfill=true",
		"ResourceAgeing initrun=1 msper=60000 resetaftern=60",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=4 mean=6 stddev=4 min=4",
		"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
	}
}

func SoilWater128x128() []string {
	return []string{
		"ResourceFillRect resource=Soil amount=1000000 x=0 w=128 y=0 h=128",
		"ResourceMazeWall resource=Water amount=1000000 xn=64 yn=64 connerfill=true",
		"ResourceAgeing initrun=1 msper=61000 resetaftern=1440",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=6 stddev=4 min=4",
		"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
	}
}

func SoilMagma128x128() []string {
	return []string{
		"ResourceFillRect resource=Soil amount=1 x=0  y=0  w=128 h=128",
		"ResourceMazeWall resource=Soil amount=500000 xn=64 yn=64 connerfill=true",
		"ResourceMazeWall resource=Fire amount=1000000 xn=64 yn=64 connerfill=true",
		"ResourceAgeing initrun=1 msper=62000 resetaftern=1440",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=8 mean=6 stddev=4 min=4",
		"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
	}
}

func SoilIce128x128() []string {
	return []string{
		"ResourceFillRect resource=Soil amount=256 x=0  y=0  w=128 h=128",
		"ResourceMazeWall  resource=Ice  amount=512 xn=64 yn=64 connerfill=true",
		"ResourceAgeing initrun=1 msper=63000 resetaftern=360",
		"AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=16 mean=6 stddev=4 min=4",
		"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
	}
}
