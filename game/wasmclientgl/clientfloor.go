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

package wasmclientgl

import (
	"fmt"
	"math/rand"
	"syscall/js"
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

type ClientFloorGL struct {
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

	PlaneTile  *PlaneLayer
	PlaneSight *PlaneLayer

	camera           js.Value
	light            js.Value
	scene            js.Value
	jsSceneObjs      map[string]js.Value   // in sight only ao, carryobj
	jsSceneTreeObjs  map[[2]int][]js.Value // tile 3d tree at 9 pos
	jsSceneGrassObjs map[[2]int][]js.Value // tile 3d grass at 9 pos
}

func NewClientFloorGL(fi *c2t_obj.FloorInfo) *ClientFloorGL {
	cf := ClientFloorGL{
		Tiles:            tilearea.New(fi.W, fi.H),
		Visited:          visitarea.NewVisitArea(fi),
		FloorInfo:        fi,
		XWrapper:         wrapper.New(fi.W),
		YWrapper:         wrapper.New(fi.H),
		jsSceneObjs:      make(map[string]js.Value),
		jsSceneTreeObjs:  make(map[[2]int][]js.Value),
		jsSceneGrassObjs: make(map[[2]int][]js.Value),
	}
	cf.XWrapSafe = cf.XWrapper.GetWrapSafeFn()
	cf.YWrapSafe = cf.YWrapper.GetWrapSafeFn()
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	cf.FieldObjPosMan = uuidposman.New(fi.W, fi.H)

	cf.PlaneTile = NewPlaneLayer(fi, -1)
	cf.PlaneSight = NewPlaneLayer(fi, 1)

	cf.camera = ThreeJsNew("PerspectiveCamera", 60, 1, 1, HelperSize*2)
	cf.scene = ThreeJsNew("Scene")
	cf.light = ThreeJsNew("PointLight", 0xffffff, 1)
	SetPosition(cf.light,
		HelperSize,
		HelperSize,
		HelperSize,
	)
	cf.scene.Call("add", cf.light)

	axisHelper := ThreeJsNew("AxesHelper", HelperSize)
	cf.scene.Call("add", axisHelper)

	// set title camera pos
	SetPosition(cf.camera,
		HelperSize/2, HelperSize/2, HelperSize,
	)
	cf.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			HelperSize/2, HelperSize/2, 0,
		),
	)
	cf.camera.Call("updateProjectionMatrix")

	cf.scene.Call("add", cf.PlaneTile.Mesh)
	cf.scene.Call("add", cf.PlaneSight.Mesh)

	return &cf
}

func (cf *ClientFloorGL) Cleanup() {
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

func (cf *ClientFloorGL) Forget() {
	FloorInfo := cf.FloorInfo
	cf.Tiles = tilearea.New(FloorInfo.W, FloorInfo.H)
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	cf.Visited = visitarea.NewVisitArea(FloorInfo)
}

func (cf *ClientFloorGL) ReplaceFloorTiles(fta *c2t_obj.NotiFloorTiles_data) {
	for x, xv := range fta.Tiles {
		for y, yv := range xv {
			if yv != 0 {
				cf.Visited.CheckAndSetNolock(x, y)
				cf.drawTileAt(x, y, yv) // must before cf.tiles update
			}
		}
	}
	cf.PlaneTile.Tex.Set("needsUpdate", true)
	cf.Tiles = fta.Tiles
	cf.Tiles4PathFind = tilearea4pathfind.New(cf.Tiles)
	return
}

func (cf *ClientFloorGL) String() string {
	return fmt.Sprintf("ClientFloorGL[%v %v %v %v]",
		cf.FloorInfo.Name,
		cf.Visited,
		cf.XWrapper.GetWidth(),
		cf.YWrapper.GetWidth(),
	)
}

func (cf *ClientFloorGL) UpdateFromViewportTile(
	taNoti *c2t_obj.NotiVPTiles_data,
	ViewportXYLenList findnear.XYLenList,
) error {

	if cf.FloorInfo.UUID != taNoti.FloorUUID {
		return fmt.Errorf("vptile data floor not match %v %v",
			cf.FloorInfo.UUID, taNoti.FloorUUID)

	}
	cf.Visited.UpdateByViewport2(taNoti.VPX, taNoti.VPY, taNoti.VPTiles)

	for i, v := range ViewportXYLenList {
		fx := cf.XWrapSafe(v.X + taNoti.VPX)
		fy := cf.YWrapSafe(v.Y + taNoti.VPY)
		if taNoti.VPTiles[i] != 0 {
			tl := taNoti.VPTiles[i]
			cf.drawTileAt(fx, fy, tl) // must before cf.tiles update
			cf.Tiles[fx][fy] = taNoti.VPTiles[i]
		}
	}
	cf.PlaneTile.Tex.Set("needsUpdate", true)
	return nil
}

func (cf *ClientFloorGL) PosAddDir(x, y int, dir way9type.Way9Type) (int, int) {
	nextX := x + dir.Dx()
	nextY := y + dir.Dy()
	nextX = cf.XWrapper.Wrap(nextX)
	nextY = cf.YWrapper.Wrap(nextY)
	return nextX, nextY
}

func (cf *ClientFloorGL) FindMovableDir(x, y int, dir way9type.Way9Type) way9type.Way9Type {
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

func (cf *ClientFloorGL) IsValidPos(x, y int) bool {
	return cf.XWrapper.IsIn(x) && cf.YWrapper.IsIn(y)
}

func (cf *ClientFloorGL) GetBias() bias.Bias {
	if cf.FloorInfo != nil {
		return cf.FloorInfo.Bias
	} else {
		return bias.Bias{}
	}
}

func (cf *ClientFloorGL) EnterFloor() {
	cf.visitTime = time.Now()
	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Float()
	winH := win.Get("innerHeight").Float()
	cf.Resize(winW, winH)
}

func (cf *ClientFloorGL) Resize(w, h float64) {
	cf.camera.Set("aspect", w/h)
	cf.camera.Call("updateProjectionMatrix")
}

func (cf *ClientFloorGL) GetFieldObjAt(x, y int) *c2t_obj.FieldObjClient {
	po, ok := cf.FieldObjPosMan.Get1stObjAt(x, y).(*c2t_obj.FieldObjClient)
	if !ok {
		return nil
	}
	return po
}
