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

package towermake

import (
	"fmt"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/terraincmd"
	"github.com/kasworld/goguelike/game/terrain/paramconv"
	"github.com/kasworld/goguelike/lib/scriptparse"
)

func (fm *Floor) String() string {
	return fmt.Sprintf("%s[%vx%v]", fm.Name, fm.W, fm.H)
}

type Floor struct {
	rnd           *g2rand.G2Rand
	Name          string
	W, H          int
	PortalIDToUse int
	Script        []string
}

func NewFloor(name string, w, h int, ao, po int, turnBoost float64) *Floor {
	if w <= 0 || h <= 0 {
		fmt.Printf("%v w,h %v %v\n", name, w, h)
	}
	if ao < 0 || po < 0 {
		fmt.Printf("%v ao,po %v %v\n", name, ao, po)
	}
	fm := &Floor{
		rnd:    g2rand.New(),
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

func (fm *Floor) Appends(arg ...string) *Floor {
	fm.Script = append(fm.Script, arg...)
	return fm
}

func (fm *Floor) Appendf(format string, arg ...interface{}) *Floor {
	fm.Script = append(fm.Script,
		fmt.Sprintf(format, arg...),
	)
	return fm
}

func (fm *Floor) IsFinalizeTerrain() bool {
	for _, v := range fm.Script {
		if v == "FinalizeTerrain" {
			return true
		}
	}
	return false
}

func (fm *Floor) MakePortalIDStringInc() string {
	rtn := fmt.Sprintf("%v-%v", fm.Name, fm.PortalIDToUse)
	// inc portal id to use
	fm.PortalIDToUse++
	return rtn
}

// bidirection (in and out) portal
// suffix "InRoom" or "Rand" or " x=47 y=15"
func (fm *Floor) ConnectStairUp(suffix, suffix2 string, dstFloor *Floor) *Floor {
	srcID := fm.MakePortalIDStringInc()
	dstID := dstFloor.MakePortalIDStringInc()
	fm.Appendf(
		"AddPortal%[1]v display=StairUp acttype=PortalInOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortal%[1]v display=StairDn acttype=PortalInOut PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix2, dstID, srcID, fm.Name)
	return fm
}

// one way portal
// suffix "InRoom" or "Rand" or " x=47 y=15"
func (fm *Floor) ConnectPortalTo(suffix, suffix2 string, dstFloor *Floor) *Floor {
	srcID := fm.MakePortalIDStringInc()
	dstID := dstFloor.MakePortalIDStringInc()
	fm.Appendf(
		"AddPortal%[1]v display=PortalIn acttype=PortalIn PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortal%[1]v display=PortalOut acttype=PortalOut PortalID=%[2]v DstPortalID=%[3]v message=From%[4]v",
		suffix2, dstID, srcID, fm.Name)
	return fm
}

// one way auto activate portal
// suffix "InRoom" or "Rand" or " x=47 y=15"
func (fm *Floor) ConnectAutoInPortalTo(suffix, suffix2 string, dstFloor *Floor) *Floor {
	srcID := fm.MakePortalIDStringInc()
	dstID := dstFloor.MakePortalIDStringInc()
	fm.Appendf(
		"AddPortal%[1]v display=PortalAutoIn acttype=PortalAutoIn PortalID=%[2]v DstPortalID=%[3]v message=To%[4]v",
		suffix, srcID, dstID, dstFloor.Name,
	)
	dstFloor.Appendf(
		"AddPortal%[1]v display=PortalOut acttype=PortalOut PortalID=%[2]v DstPortalID=%[3]v message=From%[4]v",
		suffix2, dstID, srcID, fm.Name)
	return fm
}

// suffix "InRoom" or "Rand"
func (fm *Floor) AddTeleportIn(suffix string, count int) *Floor {
	if count <= 0 {
		fmt.Printf("%v AddTeleportIn count %v\n", fm, count)
		return fm
	}
	fm.Appendf(
		"AddTrapTeleports%[1]v DstFloor=%[2]v count=%[3]v message=Teleport",
		suffix, fm.Name, count)
	return fm
}

// suffix "InRoom" or "Rand"
func (fm *Floor) AddRecycler(suffix string, count int) *Floor {
	if count <= 0 {
		fmt.Printf("%v AddRecycler count %v\n", fm, count)
		return fm
	}
	fm.Appendf(
		"AddRecycler%[1]v display=Recycler count=%[2]v message=Recycle",
		suffix, count)
	return fm
}

// suffix "InRoom" or "Rand"
func (fm *Floor) AddTrapTeleportTo(suffix string, dstFloor *Floor) *Floor {
	fm.Appendf("AddTrapTeleports%[1]v DstFloor=%[2]v count=1 message=To%[2]v",
		suffix, dstFloor.Name)
	return fm
}

// suffix "InRoom" or "Rand"
func (fm *Floor) AddEffectTrap(suffix string, trapCount int) *Floor {
	if trapCount <= 0 {
		fmt.Printf("%v AddEffectTrap trapCount %v\n", fm, trapCount)
		return fm
	}

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
	return fm
}

// suffix "InRoom" or "Rand"
func (fm *Floor) AddAllEffectTrap(suffix string, countPerEffectTrapType int) *Floor {
	if countPerEffectTrapType <= 0 {
		fmt.Printf("%v AddAllEffectTrap countPerEffectTrapType %v\n", fm, countPerEffectTrapType)
		return fm
	}
	for j := 0; j < fieldobjacttype.FieldObjActType_Count; j++ {
		ft := fieldobjacttype.FieldObjActType(j)
		if fieldobjacttype.GetBuffByFieldObjActType(ft) == nil {
			continue
		}
		fm.Appendf("AddTraps%[1]v display=None acttype=%[2]v count=%[3]v message=%[2]v",
			suffix, ft, countPerEffectTrapType)
	}
	return fm
}

func (fm *Floor) CalcRoomCount() int {
	roomCount := 0
	for _, cmdline := range fm.Script {
		cmdstr, argLine := scriptparse.SplitCmdArgstr(cmdline, " ")
		if len(cmdstr) == 0 || cmdstr[0] == '#' {
			continue
		}
		cmd, exist := terraincmd.String2TerrainCmd(cmdstr)
		if !exist {
			fmt.Printf("unknown cmd %v", cmdstr)
			continue
		}
		_, name2value, err := scriptparse.Split2ListMap(argLine, " ", "=")
		if err != nil {
			fmt.Printf("Split2ListMap %v", err)
			continue
		}
		nameList, name2type, err := scriptparse.Split2ListMap(cmd.CommentString(), " ", ":")
		if err != nil {
			fmt.Printf("Split2ListMap %v", err)
			continue
		}
		ca := &scriptparse.CmdArgs{
			Type2ConvFn: paramconv.Type2ConvFn,
			Cmd:         cmdstr,
			Name2Value:  name2value,
			NameList:    nameList,
			Name2Type:   name2type,
		}
		switch cmd {
		case terraincmd.AddRoom, terraincmd.AddRoomMaze:
			roomCount++
		case terraincmd.AddRoomsRand:
			rc := 0
			if err := ca.SetArgByName("count", &rc); err == nil {
				roomCount += rc
			} else {
				fmt.Printf("SetArgByName %v", err)
			}
		}
	}
	return roomCount
}
