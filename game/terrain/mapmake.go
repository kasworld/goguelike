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

package terrain

import (
	"fmt"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/game/terrain/corridor"
	"github.com/kasworld/goguelike/game/terrain/resourcetilearea"
	"github.com/kasworld/goguelike/game/terrain/roommanager"
	"github.com/kasworld/goguelike/game/terrain/scriptparse"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/wrapper"
)

var cmdFnArgs = map[string]struct {
	fn     func(tr *Terrain, ca *scriptparse.CmdArgs) error
	format string
}{
	"NewTerrain":             {cmdNewTerrain, "name:string w:int h:int ao:int po:int actturnboost:float"},
	"ResourceMazeWall":       {cmdResourceMazeWall, "resource:TileRsc_Type amount:int xn:int yn:int connerfill:bool"},
	"ResourceMazeWalk":       {cmdResourceMazeWalk, "resource:TileRsc_Type amount:int xn:int yn:int connerfill:bool"},
	"ResourceRand":           {cmdResourceRand, "resource:TileRsc_Type mean:int stddev:int repeat:int"},
	"Resource":               {cmdResource, "resource:TileRsc_Type amount:int x:int y:int"},
	"ResourceHLine":          {cmdResourceHLine, "resource:TileRsc_Type amount:int x:int w:int y:int"},
	"ResourceVLine":          {cmdResourceVLine, "resource:TileRsc_Type amount:int x:int y:int h:int"},
	"ResourceLine":           {cmdResourceLine, "resource:TileRsc_Type amount:int x1:int y1:int x2:int y2:int"},
	"ResourceRect":           {cmdResourceRect, "resource:TileRsc_Type amount:int x:int w:int y:int h:int"},
	"ResourceFillRect":       {cmdResourceFillRect, "resource:TileRsc_Type amount:int x:int w:int y:int h:int"},
	"ResourceFillEllipses":   {cmdResourceFillEllipses, "resource:TileRsc_Type amount:int x:int w:int y:int h:int"},
	"ResourceFromPNG":        {cmdResourceFromPNG, "name:string"},
	"ResourceAgeing":         {cmdAgeing, "initrun:int msper:int resetaftern:int"},
	"AddRoom":                {cmdAddRoom, "bgtile:Tile_Type walltile:Tile_Type terrace:bool x:int y:int w:int h:int"},
	"AddRoomMaze":            {cmdAddMazeRoom, "bgtile:Tile_Type walltile:Tile_Type terrace:bool x:int y:int w:int h:int xn:int yn:int connerfill:bool"},
	"AddRoomsRand":           {cmdAddRandRooms, "bgtile:Tile_Type walltile:Tile_Type terrace:bool align:int count:int mean:int stddev:int min:int"},
	"ConnectRooms":           {cmdConnectRooms, "tile:Tile_Type connect:int allconnect:bool diagonal:bool"},
	"FinalizeTerrain":        {cmdFinalizeTerrain, ""},
	"AddPortal":              {cmdAddPortal, "x:int y:int display:FieldObjDisplay_Type acttype:FieldObjAct_Type PortalID:string DstPortalID:string message:string"},
	"AddPortalRand":          {cmdAddPortalRand, "display:FieldObjDisplay_Type acttype:FieldObjAct_Type PortalID:string DstPortalID:string message:string"},
	"AddPortalInRoom":        {cmdAddPortalRandInRoom, "display:FieldObjDisplay_Type acttype:FieldObjAct_Type PortalID:string DstPortalID:string message:string"},
	"AddRecycler":            {cmdAddRecycler, "x:int y:int display:FieldObjDisplay_Type message:string"},
	"AddRecyclerRand":        {cmdAddRecyclerRand, "display:FieldObjDisplay_Type count:int message:string"},
	"AddRecyclerInRoom":      {cmdAddRecyclerRandInRoom, "display:FieldObjDisplay_Type count:int message:string"},
	"AddTrapTeleport":        {cmdAddTrapTeleport, "x:int y:int DstFloor:string message:string "},
	"AddTrapTeleportsRand":   {cmdAddTrapTeleportRand, "DstFloor:string count:int message:string"},
	"AddTrapTeleportsInRoom": {cmdAddTrapTeleportRandInRoom, "DstFloor:string count:int message:string"},
	"AddTrap":                {cmdAddTrap, "x:int y:int display:FieldObjDisplay_Type acttype:FieldObjAct_Type message:string"},
	"AddTrapsRand":           {cmdAddTrapRand, "display:FieldObjDisplay_Type acttype:FieldObjAct_Type count:int message:string"},
	"AddTrapsInRoom":         {cmdAddTrapRandInRoom, "display:FieldObjDisplay_Type acttype:FieldObjAct_Type count:int message:string"},
}

func init() {
	// verify format
	for _, v := range cmdFnArgs {
		_, n2v, err := scriptparse.Split2ListMap(v.format, " ", ":")
		for _, t := range n2v {
			_, exist := scriptparse.Type2ConvFn[t]
			if !exist {
				panic(fmt.Sprintf("unknown type %v %v", t, v.format))
			}
		}
		if err != nil {
			panic(err)
		}
	}
}

func (tr *Terrain) Execute1Cmdline(cmdline string) error {
	cmd, argLine := scriptparse.SplitCmdArgstr(cmdline, " ")
	if len(cmd) == 0 || cmd[0] == '#' {
		return nil
	}
	fnarg, exist := cmdFnArgs[cmd]
	if !exist {
		return fmt.Errorf("unknown cmd %v", cmd)
	}
	_, name2value, err := scriptparse.Split2ListMap(argLine, " ", "=")
	if err != nil {
		return err
	}
	nameList, name2type, err := scriptparse.Split2ListMap(fnarg.format, " ", ":")
	if err != nil {
		return err
	}
	ca := &scriptparse.CmdArgs{
		Cmd:        cmd,
		Name2Value: name2value,
		NameList:   nameList,
		Name2Type:  name2type,
	}
	return fnarg.fn(tr, ca)
}

func (tr *Terrain) execNewTerrain(
	name string, w, h int, aocount, pocount int, actturnboost float64) error {
	// if !isPowerOfTwo(w) || !isPowerOfTwo(h) {
	// 	return fmt.Errorf("w,h must power of 2, %v %v", w, h)
	// }
	tr.Xlen, tr.Ylen = w, h
	tr.ActiveObjCount = aocount
	tr.CarryObjCount = pocount

	tr.XWrapper = wrapper.New(tr.Xlen)
	tr.YWrapper = wrapper.New(tr.Ylen)
	tr.XWrap = tr.XWrapper.GetWrapFn()
	tr.YWrap = tr.YWrapper.GetWrapFn()

	tr.Name = name
	tr.ActTurnBoost = actturnboost
	tr.tileArea = tilearea.New(w, h)
	tr.resourceTileArea = resourcetilearea.New(w, h)
	tr.roomManager = roommanager.New(w, h)
	tr.corridorList = make([]*corridor.Corridor, 0)
	tr.foPosMan = uuidposman.New(tr.Xlen, tr.Ylen)

	tr.initCrpCache()
	tr.findList = findnear.NewXYLenList(w, h)
	return nil
}

func isPowerOfTwo(i int) bool {
	return (i & (i - 1)) == 0
}
