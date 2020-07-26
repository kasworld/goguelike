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

func wrapInt(v, l int) int {
	return (v%l + l) % l
}

func makePowerOf2(v int) int {
	return 1 << uint(bits.Len(uint(v-1)))
}

var whList = []int{
	32, 64, 128,
}
var allRoomTile = []string{
	"Room", "Soil", "Sand", "Stone", "Grass", "Tree", "Ice", "Magma", "Swamp", "Sea", "Smoke",
}
var allRoadTile = []string{
	"Road", "Soil", "Sand", "Stone", "Grass", "Tree", "Fog",
}

var rnd = g2rand.New()

func makeRogueTower(towerName string, floorCount int) {

	floorList := make([]*FloorMake, floorCount)
	for i := range floorList {
		w := whList[rnd.Intn(len(whList))]
		h := whList[rnd.Intn(len(whList))]
		roomCount := w * h / 512
		if roomCount < 2 {
			roomCount = 2
		}
		floorList[i] = NewFloorMake(
			i,
			fmt.Sprintf("Floor%v", i),
			w, h,
			roomCount/2, 0,
			1.0,
		)
	}

	for _, fm := range floorList {
		fm.MakeRoguelike()
	}

	for _, fm := range floorList {
		roomCount := fm.W * fm.H / 512
		fm.ConnectStairUp("InRoom", floorList[wrapInt(fm.Num+1, floorCount)])
		fm.AddRecycler("InRoom", roomCount/2)
		fm.AddTeleportIn("InRoom", roomCount/2)
		fm.AddTeleportTrapOut("InRoom", floorList, roomCount/2)
		fm.AddEffectTrap("InRoom", roomCount/2)
	}

	tw := make(towerscript.TowerScript, 0)
	for i := 0; i < floorCount; i++ {
		tw = append(tw, floorList[i].Script)
	}
	err := tw.SaveJSON(towerName + ".tower")
	if err != nil {
		g2log.Error("%v", err)
	}
}

type FloorMake struct {
	Num           int
	Name          string
	W, H          int
	PortalIDToUse int
	Script        []string
}

func NewFloorMake(num int, name string, w, h int, ao, po int, turnBoost float64) *FloorMake {
	fm := &FloorMake{
		Num:    num,
		Name:   name,
		W:      w,
		H:      h,
		Script: make([]string, 0),
	}
	fm.Appendf(
		"NewTerrain w=%v h=%v name=%v ao=%v po=%v actturnboost=%v",
		w, h, name, ao, po, turnBoost)
	return fm
}

func (fm *FloorMake) Appendf(format string, arg ...interface{}) {
	fm.Script = append(fm.Script,
		fmt.Sprintf(format, arg...),
	)
}

func (fm *FloorMake) MakeRoguelike() {
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

func (fm *FloorMake) MakePortalIDStringInc() string {
	rtn := fmt.Sprintf("%v-%v", fm.Name, fm.PortalIDToUse)
	// inc portal id to use
	fm.PortalIDToUse++
	return rtn
}

// bidirection (in and out) portal
// suffix "InRoom" or "Rand"
func (fm *FloorMake) ConnectStairUp(suffix string, dstFloor *FloorMake) {
	srcID := fm.MakePortalIDStringInc()
	dstID := dstFloor.MakePortalIDStringInc()
	fm.Appendf(
		"AddPortals%[1]v display=StairUp acttype=PortalInOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortals%[1]v display=StairDn acttype=PortalInOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, dstID, srcID, dstFloor.Name)
}

// one way portal
// suffix "InRoom" or "Rand"
func (fm *FloorMake) ConnectPortalTo(suffix string, dstFloor *FloorMake) {
	srcID := fm.MakePortalIDStringInc()
	dstID := dstFloor.MakePortalIDStringInc()
	fm.Appendf(
		"AddPortals%[1]v display=PortalIn acttype=PortalIn PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortals%[1]v display=PortalOut acttype=PortalOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, dstID, srcID, dstFloor.Name)
}

// one way auto activate portal
// suffix "InRoom" or "Rand"
func (fm *FloorMake) ConnectAutoInPortalTo(suffix string, dstFloor *FloorMake) {
	srcID := fm.MakePortalIDStringInc()
	dstID := dstFloor.MakePortalIDStringInc()
	fm.Appendf(
		"AddPortals%[1]v display=PortalAutoIn acttype=PortalAutoIn PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortals%[1]v display=PortalOut acttype=PortalOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, dstID, srcID, dstFloor.Name)
}

// suffix "InRoom" or "Rand"
func (fm *FloorMake) AddTeleportIn(suffix string, count int) {
	fm.Appendf(
		"AddTrapTeleports%[1]v DstFloor=%[2]v count=%[3]v message=Teleport",
		suffix, fm.Name, count)
}

// suffix "InRoom" or "Rand"
func (fm *FloorMake) AddRecycler(suffix string, count int) {
	fm.Appendf(
		"AddRecycler%[1]v display=Recycler count=%[2]v message=Recycle",
		suffix, count)
}

// suffix "InRoom" or "Rand"
func (fm *FloorMake) AddTeleportTrapOut(
	suffix string, floorList []*FloorMake, trapCount int) {
	for trapMade := 0; trapMade < trapCount; {
		dstFloor := floorList[rnd.Intn(len(floorList))]
		fm.Appendf("AddTrapTeleports%[1]v DstFloor=%[2]v count=1 message=ToFloor%[2]v",
			suffix, dstFloor.Name)

		trapMade++
	}
}

// suffix "InRoom" or "Rand"
func (fm *FloorMake) AddEffectTrap(suffix string, trapCount int) {
	for trapMade := 0; trapMade < trapCount; {
		j := rnd.Intn(fieldobjacttype.FieldObjActType_Count)
		ft := fieldobjacttype.FieldObjActType(j)
		if fieldobjacttype.GetBuffByFieldObjActType(ft) == nil {
			continue
		}
		fm.Appendf("AddTraps%[1]v display=None acttype=%[2]v count=1 message=%[2]v",
			suffix, ft)
		trapMade++
	}
}
