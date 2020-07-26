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

package floormake

import (
	"fmt"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
)

type FloorMake struct {
	rnd           *g2rand.G2Rand
	Num           int
	Name          string
	W, H          int
	PortalIDToUse int
	Script        []string
}

func New(num int, name string, w, h int, ao, po int, turnBoost float64) *FloorMake {
	fm := &FloorMake{
		rnd:    g2rand.New(),
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
		"AddPortal%[1]v display=StairUp acttype=PortalInOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortal%[1]v display=StairDn acttype=PortalInOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, dstID, srcID, dstFloor.Name)
}

// one way portal
// suffix "InRoom" or "Rand"
func (fm *FloorMake) ConnectPortalTo(suffix string, dstFloor *FloorMake) {
	srcID := fm.MakePortalIDStringInc()
	dstID := dstFloor.MakePortalIDStringInc()
	fm.Appendf(
		"AddPortal%[1]v display=PortalIn acttype=PortalIn PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortal%[1]v display=PortalOut acttype=PortalOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, dstID, srcID, dstFloor.Name)
}

// one way auto activate portal
// suffix "InRoom" or "Rand"
func (fm *FloorMake) ConnectAutoInPortalTo(suffix string, dstFloor *FloorMake) {
	srcID := fm.MakePortalIDStringInc()
	dstID := dstFloor.MakePortalIDStringInc()
	fm.Appendf(
		"AddPortal%[1]v display=PortalAutoIn acttype=PortalAutoIn PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortal%[1]v display=PortalOut acttype=PortalOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
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
		dstFloor := floorList[fm.rnd.Intn(len(floorList))]
		fm.Appendf("AddTrapTeleports%[1]v DstFloor=%[2]v count=1 message=ToFloor%[2]v",
			suffix, dstFloor.Name)

		trapMade++
	}
}

// suffix "InRoom" or "Rand"
func (fm *FloorMake) AddEffectTrap(suffix string, trapCount int) {
	for trapMade := 0; trapMade < trapCount; {
		j := fm.rnd.Intn(fieldobjacttype.FieldObjActType_Count)
		ft := fieldobjacttype.FieldObjActType(j)
		if fieldobjacttype.GetBuffByFieldObjActType(ft) == nil {
			continue
		}
		fm.Appendf("AddTraps%[1]v display=None acttype=%[2]v count=1 message=%[2]v",
			suffix, ft)
		trapMade++
	}
}
