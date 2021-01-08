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

package floortemplate

import "fmt"

func AgeingFloorWH(w, h int) []string {
	return []string{
		fmt.Sprintf("ResourceFillRect resource=Soil amount=1 x=0  y=0  w=%v h=%v", w, h),
		fmt.Sprintf("ResourceRand resource=Ice   mean=1600000000 stddev=10000000 repeat=%v", w*h/256/32),
		fmt.Sprintf("ResourceRand resource=Water mean=100000000  stddev=10000000 repeat=%v", w*h/256),
		fmt.Sprintf("ResourceRand resource=Fire  mean=1600000000 stddev=10000000 repeat=%v", w*h/256/32),
		fmt.Sprintf("ResourceRand resource=Soil  mean=100000000  stddev=10000000 repeat=%v", w*h/256),
		fmt.Sprintf("ResourceRand resource=Plant mean=100000000  stddev=10000000 repeat=%v", w*h/256),
		"ResourceAgeing initrun=10 msper=65000 resetaftern=1440",
	}
}

func AgeingCitySize(size int) []string {
	return append(
		AgeingFloorWH(size, size),
		[]string{
			fmt.Sprintf("AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=%v mean=8 stddev=4", size*size/256/4),
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		}...)
}

func AgeingFieldSize(size int) []string {
	return append(
		AgeingFloorWH(size, size),
		[]string{
			fmt.Sprintf("AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=%v mean=8 stddev=4", size*size/256/32),
			"ConnectRooms tile=Road connect=1 allconnect=false diagonal=true",
		}...)
}

func AgeingMazeSize(size int) []string {
	return append(
		AgeingFloorWH(size, size),
		[]string{
			fmt.Sprintf("ResourceMazeWall resource=Stone amount=2000001 x=0 y=0 w=%v h=%v xn=%v yn=%v connerfill=false",
				size, size, size/8, size/8),
			fmt.Sprintf("AddRoomsRand bgtile=Room walltile=Wall terrace=false align=1 count=%v mean=8 stddev=4", size*size/256/32),
		}...)
}
