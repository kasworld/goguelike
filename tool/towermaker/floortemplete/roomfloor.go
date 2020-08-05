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
