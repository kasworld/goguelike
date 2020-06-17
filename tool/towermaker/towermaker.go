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

	"github.com/kasworld/goguelike/enum/fieldobjacttype"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/game/towerscript"
	"github.com/kasworld/goguelike/lib/g2log"
)

func main() {
	towername := flag.String("towername", "roguegen", "tower filename w/o .tower")
	floorcount := flag.Int("floorcount", 100, "floor count in tower")
	flag.Parse()

	fmt.Printf("make tower with towername:%v floorcount:%v\n", *towername, *floorcount)
	makeRogueTower(*towername, *floorcount)
}

func makeRogueTower(towerName string, floorCount int) {
	whList := []int{
		32, 64, 128,
	}
	rnd := g2rand.New()
	allRoomTile := []string{
		"Room", "Soil", "Sand", "Stone", "Grass", "Tree", "Ice", "Magma", "Swamp", "Sea", "Smoke",
	}
	allRoadTile := []string{
		"Road", "Soil", "Sand", "Stone", "Grass", "Tree", "Fog",
	}

	tw := make(towerscript.TowerScript, 0)

	for i := 0; i < floorCount; i++ {
		w := whList[rnd.Intn(len(whList))]
		h := whList[rnd.Intn(len(whList))]
		roomCount := w * h / 512
		if roomCount < 2 {
			roomCount = 2
		}
		roomTile := allRoomTile[rnd.Intn(len(allRoomTile))]
		roadTile := allRoadTile[rnd.Intn(len(allRoadTile))]
		fl := []string{
			fmt.Sprintf(
				"NewTerrain w=%v h=%v name=Floor%v ao=%v po=%v actturnboost=1.0",
				w, h, i, roomCount/2, 0),
			fmt.Sprintf("AddRoomsRand bgtile=%v walltile=Wall terrace=false align=1 count=%v mean=8 stddev=2 min=6",
				roomTile, roomCount),
			fmt.Sprintf("ConnectRooms tile=%v connect=1 allconnect=true diagonal=false",
				roadTile),
			"FinalizeTerrain",
			fmt.Sprintf(
				"AddPortalsInRoom display=StairDn acttype=PortalInOut PortalID=Floor%v-0 DstPortalID=Floor%v-1 message=Floor%v",
				i, wrapInt(i-1, floorCount), wrapInt(i-1, floorCount)),
			fmt.Sprintf(
				"AddPortalsInRoom display=StairUp acttype=PortalInOut PortalID=Floor%v-1 DstPortalID=Floor%v-0 message=Floor%v",
				i, wrapInt(i+1, floorCount), wrapInt(i+1, floorCount)),
			fmt.Sprintf("AddRecyclerInRoom display=Recycler count=%v message=Recycle", roomCount/2),
			fmt.Sprintf("AddTrapTeleportsInRoom DstFloor=Floor%v count=%v message=Teleport",
				i, roomCount/2),
		}

		for trapMade := 0; ; {
			dstFloor := rnd.Intn(floorCount)
			cmd := fmt.Sprintf("AddTrapTeleportsInRoom DstFloor=Floor%v count=%v message=ToFloor%v",
				dstFloor, 1, dstFloor)

			fl = append(fl, cmd)
			trapMade++
			if trapMade > roomCount/2 {
				break
			}
		}

		for trapMade := 0; ; {
			j := rnd.Intn(fieldobjacttype.FieldObjActType_Count)
			ft := fieldobjacttype.FieldObjActType(j)
			if fieldobjacttype.GetBuffByFieldObjActType(ft) == nil {
				continue
			}
			cmd := fmt.Sprintf("AddTrapsInRoom display=None acttype=%v count=1 message=%v",
				ft, ft)
			fl = append(fl, cmd)
			trapMade++
			if trapMade > roomCount {
				break
			}
		}
		tw = append(tw, fl)
	}

	err := tw.SaveJSON(towerName + ".tower")
	if err != nil {
		g2log.Error("%v", err)
	}
}

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func makePowerOf2(v int) int {
	return 1 << uint(bits.Len(uint(v-1)))
}
