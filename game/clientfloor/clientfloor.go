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

package clientfloor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/goguelike/game/tilearea4pathfind"
	"github.com/kasworld/goguelike/game/visitarea"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/wrapper"
)

type ClientFloor struct {
	FloorInfo *c2t_obj.FloorInfo
	Tiles     tilearea.TileArea `prettystring:"simple"`

	Visited        *visitarea.VisitArea `prettystring:"simple"`
	XWrapper       *wrapper.Wrapper     `prettystring:"simple"`
	YWrapper       *wrapper.Wrapper     `prettystring:"simple"`
	XWrapSafe      func(i int) int
	YWrapSafe      func(i int) int
	Tiles4PathFind *tilearea4pathfind.TileArea4PathFind `prettystring:"simple"`
	visitTime      time.Time                            `prettystring:"simple"`

	FieldObjPosMan *uuidposman.UUIDPosMan `prettystring:"simple"`
}

func New(FloorInfo *c2t_obj.FloorInfo) *ClientFloor {
	cf := ClientFloor{
		Tiles:     tilearea.New(FloorInfo.W, FloorInfo.H),
		Visited:   visitarea.NewVisitArea(FloorInfo),
		FloorInfo: FloorInfo,
		XWrapper:  wrapper.New(FloorInfo.W),
		YWrapper:  wrapper.New(FloorInfo.H),
	}
	cf.XWrapSafe = cf.XWrapper.GetWrapSafeFn()
	cf.YWrapSafe = cf.YWrapper.GetWrapSafeFn()
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	cf.FieldObjPosMan = uuidposman.New(FloorInfo.W, FloorInfo.H)

	return &cf
}

func (cf *ClientFloor) Cleanup() {
	cf.Tiles = nil
	if v := cf.Visited; v != nil {
		v.Cleanup()
	}
	cf.Visited = nil
	if t := cf.Tiles4PathFind; t != nil {
		t.Cleanup()
	}
	cf.Tiles4PathFind = nil
	if i := cf.FieldObjPosMan; i != nil {
		i.Cleanup()
	}
	cf.FieldObjPosMan = nil
}

func (cf *ClientFloor) Forget() {
	FloorInfo := cf.FloorInfo
	cf.Tiles = tilearea.New(FloorInfo.W, FloorInfo.H)
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	cf.Visited = visitarea.NewVisitArea(FloorInfo)
}

func (cf *ClientFloor) ReplaceFloorTiles(fta *c2t_obj.NotiFloorTiles_data) {
	cf.Tiles = fta.Tiles
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	for x, xv := range cf.Tiles {
		for y, yv := range xv {
			if yv != 0 {
				cf.Visited.CheckAndSetNolock(x, y)
			}
		}
	}
	return
}

func (cf *ClientFloor) String() string {
	return fmt.Sprintf("ClientFloor[%v %v %v %v]",
		cf.FloorInfo.Name,
		cf.Visited,
		cf.XWrapper.GetWidth(),
		cf.YWrapper.GetWidth(),
	)
}

func (cf *ClientFloor) UpdateFromViewportTile(
	vp *c2t_obj.NotiVPTiles_data,
	vpXYLenList findnear.XYLenList,
) error {

	if cf.FloorInfo.UUID != vp.FloorUUID {
		return fmt.Errorf("vptile data floor not match %v %v",
			cf.FloorInfo.UUID, vp.FloorUUID)

	}
	cf.Visited.UpdateByViewport2(vp.VPX, vp.VPY, vp.VPTiles, vpXYLenList)

	for i, v := range vpXYLenList {
		fx := cf.XWrapSafe(v.X + vp.VPX)
		fy := cf.YWrapSafe(v.Y + vp.VPY)
		if vp.VPTiles[i] != 0 {
			cf.Tiles[fx][fy] = vp.VPTiles[i]
		}
	}
	return nil
}

func (cf *ClientFloor) PosAddDir(x, y int, dir way9type.Way9Type) (int, int) {
	nextX := x + dir.Dx()
	nextY := y + dir.Dy()
	nextX = cf.XWrapper.Wrap(nextX)
	nextY = cf.YWrapper.Wrap(nextY)
	return nextX, nextY
}

func (cf *ClientFloor) FindMovableDir(x, y int, dir way9type.Way9Type) way9type.Way9Type {
	dirList := []way9type.Way9Type{
		dir,
		dir.TurnDir(1),
		dir.TurnDir(-1),
		dir.TurnDir(2),
		dir.TurnDir(-2),
	}
	if rand.Float64() >= 0.5 {
		dirList = []way9type.Way9Type{
			dir,
			dir.TurnDir(-1),
			dir.TurnDir(1),
			dir.TurnDir(-2),
			dir.TurnDir(2),
		}
	}
	for _, dir := range dirList {
		nextX, nextY := cf.PosAddDir(x, y, dir)
		if cf.Tiles[nextX][nextY].CharPlaceable() {
			return dir
		}
	}
	return way9type.Center
}

func (cf *ClientFloor) IsValidPos(x, y int) bool {
	return cf.XWrapper.IsIn(x) && cf.YWrapper.IsIn(y)
}

func (cf *ClientFloor) GetBias() bias.Bias {
	if cf.FloorInfo != nil {
		return cf.FloorInfo.Bias
	} else {
		return bias.Bias{}
	}
}

func (cf *ClientFloor) EnterFloor() {
	cf.visitTime = time.Now()
}

func (cf *ClientFloor) GetFieldObjAt(x, y int) *c2t_obj.FieldObjClient {
	po, ok := cf.FieldObjPosMan.Get1stObjAt(x, y).(*c2t_obj.FieldObjClient)
	if !ok {
		return nil
	}
	return po
}
