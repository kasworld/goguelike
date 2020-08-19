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

	"github.com/kasworld/goguelike/enum/terraincmd"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/game/terrain/corridor"
	"github.com/kasworld/goguelike/game/terrain/paramconv"
	"github.com/kasworld/goguelike/game/terrain/resourcetilearea"
	"github.com/kasworld/goguelike/game/terrain/roommanager"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/wrapper"
)

var TerrainScriptFn = map[terraincmd.TerrainCmd]func(tr *Terrain, ca *scriptparse.CmdArgs) error{
	terraincmd.NewTerrain: cmdNewTerrain,

	terraincmd.ResourceMazeWall:     cmdResourceMazeWall,
	terraincmd.ResourceMazeWalk:     cmdResourceMazeWalk,
	terraincmd.ResourceRand:         cmdResourceRand,
	terraincmd.ResourceAt:           cmdResourceAt,
	terraincmd.ResourceHLine:        cmdResourceHLine,
	terraincmd.ResourceVLine:        cmdResourceVLine,
	terraincmd.ResourceLine:         cmdResourceLine,
	terraincmd.ResourceRect:         cmdResourceRect,
	terraincmd.ResourceFillRect:     cmdResourceFillRect,
	terraincmd.ResourceFillEllipses: cmdResourceFillEllipses,
	terraincmd.ResourceFromPNG:      cmdResourceFromPNG,
	terraincmd.ResourceAgeing:       cmdAgeing,

	terraincmd.AddRoom:      cmdAddRoom,
	terraincmd.AddRoomMaze:  cmdAddMazeRoom,
	terraincmd.AddRoomsRand: cmdAddRandRooms,
	terraincmd.ConnectRooms: cmdConnectRooms,

	// operations to tileLayer
	terraincmd.TileMazeWall:     cmdTileMazeWall,
	terraincmd.TileMazeWalk:     cmdTileMazeWalk,
	terraincmd.TileAt:           cmdTileAt,
	terraincmd.TileHLine:        cmdTileHLine,
	terraincmd.TileVLine:        cmdTileVLine,
	terraincmd.TileLine:         cmdTileLine,
	terraincmd.TileRect:         cmdTileRect,
	terraincmd.TileFillRect:     cmdTileFillRect,
	terraincmd.TileFillEllipses: cmdTileFillEllipses,

	terraincmd.FinalizeTerrain: cmdFinalizeTerrain,

	terraincmd.AddPortal:              cmdAddPortal,
	terraincmd.AddPortalRand:          cmdAddPortalRand,
	terraincmd.AddPortalInRoom:        cmdAddPortalRandInRoom,
	terraincmd.AddRecycler:            cmdAddRecycler,
	terraincmd.AddRecyclerRand:        cmdAddRecyclerRand,
	terraincmd.AddRecyclerInRoom:      cmdAddRecyclerRandInRoom,
	terraincmd.AddTrapTeleport:        cmdAddTrapTeleport,
	terraincmd.AddTrapTeleportsRand:   cmdAddTrapTeleportRand,
	terraincmd.AddTrapTeleportsInRoom: cmdAddTrapTeleportRandInRoom,
	terraincmd.AddTrap:                cmdAddTrap,
	terraincmd.AddTrapsRand:           cmdAddTrapRand,
	terraincmd.AddTrapsInRoom:         cmdAddTrapRandInRoom,
}

func init() {
	// verify format
	for i := 0; i < terraincmd.TerrainCmd_Count; i++ {
		format := terraincmd.TerrainCmd(i).CommentString()
		_, n2v, err := scriptparse.Split2ListMap(format, " ", ":")
		for _, t := range n2v {
			_, exist := paramconv.Type2ConvFn[t]
			if !exist {
				panic(fmt.Sprintf("unknown type %v %v", t, format))
			}
		}
		if err != nil {
			panic(err)
		}
	}
}

func (tr *Terrain) Execute1Cmdline(cmdline string) error {
	cmdstr, argLine := scriptparse.SplitCmdArgstr(cmdline, " ")
	if len(cmdstr) == 0 || cmdstr[0] == '#' {
		return nil
	}
	cmd, exist := terraincmd.String2TerrainCmd(cmdstr)
	if !exist {
		return fmt.Errorf("unknown cmd %v", cmd)
	}
	fn, exist := TerrainScriptFn[cmd]
	if !exist {
		return fmt.Errorf("unknown cmd %v", cmd)
	}
	_, name2value, err := scriptparse.Split2ListMap(argLine, " ", "=")
	if err != nil {
		return err
	}
	nameList, name2type, err := scriptparse.Split2ListMap(cmd.CommentString(), " ", ":")
	if err != nil {
		return err
	}
	ca := &scriptparse.CmdArgs{
		Type2ConvFn: paramconv.Type2ConvFn,
		Cmd:         cmdstr,
		Name2Value:  name2value,
		NameList:    nameList,
		Name2Type:   name2type,
	}
	return fn(tr, ca)
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
	tr.serviceTileArea = tilearea.New(w, h)
	tr.tileLayer = tilearea.New(w, h)
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
