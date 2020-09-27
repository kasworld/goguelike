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

// Package terrain static terrain include generationstep
package terrain

import (
	"fmt"
	"sync/atomic"

	"github.com/kasworld/findnear"
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/terrain/corridor"
	"github.com/kasworld/goguelike/game/terrain/resourcetilearea"
	"github.com/kasworld/goguelike/game/terrain/roommanager"
	"github.com/kasworld/goguelike/game/terrain/viewportcache"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/goguelike/game/tilearea4pathfind"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/wrapper"
)

func (tr *Terrain) String() string {
	return fmt.Sprintf("Terrain[%v Seed:%v (%v %v) ]",
		tr.Name,
		tr.Seed,
		tr.Xlen, tr.Ylen,
	)
}

type Terrain struct {
	dataDir       string
	terrainScript []string
	rnd           *g2rand.G2Rand `prettystring:"hide"`
	log           *g2log.LogBase `prettystring:"hide"`

	oriTiles resourcetilearea.ResourceTileArea `prettystring:"simple"`

	viewportCache *viewportcache.ViewportCache         `prettystring:"simple"`
	ta4ff         *tilearea4pathfind.TileArea4PathFind `prettystring:"simple"`

	ageingCount int64
	inAgeing    int32

	// convert serviceTileArea
	resourceTileArea resourcetilearea.ResourceTileArea `prettystring:"simple"`

	// service tile area
	serviceTileArea tilearea.TileArea `prettystring:"simple"`

	// need override to tilearea after ageing
	tileLayer tilearea.TileArea `prettystring:"simple"`

	// tile layer relate data
	roomManager  *roommanager.RoomManager `prettystring:"simple"`
	corridorList []*corridor.Corridor     `prettystring:"simple"`
	crp8Way      bool                     // corridor connect 4way or 8way

	// must clear after terrain made
	crpCache [][]*CorridorPather `prettystring:"simple"`
	findList findnear.XYLenList  `prettystring:"simple"`

	foPosMan *uuidposman.UUIDPosMan `prettystring:"simple"`

	Xlen     int
	Ylen     int
	XWrapper *wrapper.Wrapper `prettystring:"simple"`
	YWrapper *wrapper.Wrapper `prettystring:"simple"`
	XWrap    func(i int) int
	YWrap    func(i int) int

	Name              string
	Seed              int64
	ActTurnBoost      float64
	ActiveObjCount    int
	CarryObjCount     int
	MSPerAgeing       int64
	ResetAfterNAgeing int64
	Tile2Discover     int
}

func New(seed int64, script []string, dataDir string, l *g2log.LogBase) *Terrain {
	tr := &Terrain{
		dataDir:       dataDir,
		terrainScript: script,
		log:           l,
	}
	tr.viewportCache = viewportcache.New(tr)
	tr.Seed = seed
	tr.rnd = g2rand.NewWithSeed(seed)
	return tr
}
func (tr *Terrain) Cleanup() {
	tr.log.TraceService("Start Cleanup Terrain %v", tr.Name)
	defer func() { tr.log.TraceService("End Cleanup Terrain %v", tr.Name) }()
	tr.foPosMan.Cleanup()
}

func (tr *Terrain) Init() error {
	for i, cmdline := range tr.terrainScript {
		if err := tr.Execute1Cmdline(cmdline); err != nil {
			tr.log.Fatal(
				"fail to exec terrain %v %v:%v %v", tr, i, cmdline, err)
			return err
		}
	}
	if tr.Name == "" {
		return nil // skip no name terrain
	}
	tr.oriTiles = tr.resourceTileArea.Dup()
	return nil
}

func (tr *Terrain) Ageing() error {
	if tr.MSPerAgeing == 0 {
		return fmt.Errorf("Not ageable terrain %v", tr)
	}
	var err error
	ageingLimit := tr.GetResetAfterNAgeing()
	if ageingLimit == 0 || ageingLimit > tr.ageingCount {
		err = tr.AgeingNoCheck()
	} else {
		err = tr.ResetAgeing()
	}
	return err
}

func (tr *Terrain) AgeingNoCheck() error {
	if atomic.CompareAndSwapInt32(&tr.inAgeing, 0, 1) {
		defer atomic.AddInt32(&tr.inAgeing, -1)
		rta := tr.GetRcsTiles().Dup().Ageing(tr.rnd.Intn, 1)
		tr.resourceTileArea = rta
		tr.renderServiceTileArea()
		tr.ageingCount++
		return nil
	} else {
		return fmt.Errorf("skip AgeingNoCheck, do next time %v", tr)
	}
}

func (tr *Terrain) ResetAgeing() error {
	if atomic.CompareAndSwapInt32(&tr.inAgeing, 0, 1) {
		defer atomic.AddInt32(&tr.inAgeing, -1)
		rta := tr.GetOriRcsTiles().Dup()
		tr.resourceTileArea = rta
		tr.renderServiceTileArea()
		tr.ageingCount = 0
		return nil
	} else {
		return fmt.Errorf("skip ResetAgeing, do next time %v", tr)
	}
}

func (tr *Terrain) renderServiceTileArea() {
	tr.resource2View()
	tr.tileLayer2SeviceTileArea()
	tr.openBlockedDoor()
	tr.Tile2Discover = tr.serviceTileArea.CalcNotEmptyTileCount()
	tr.viewportCache.Reset()
	tr.ta4ff = tilearea4pathfind.New(tr.GetTiles())
	for _, o := range tr.foPosMan.GetAllList() {
		fo, ok := o.(*fieldobject.FieldObject)
		if !ok {
			tr.log.Fatal("not *fieldobject.FieldObject %v", o)
			continue
		}
		x, y, exist := tr.foPosMan.GetXYByUUID(o.GetUUID())
		if !exist {
			tr.log.Fatal("fieldobj not found %v", o)
			continue
		}
		if !tr.serviceTileArea[x][y].CharPlaceable() {
			tr.log.Fatal("fieldobj placed at NonCharPlaceable tile %v", fo)
		}
	}
}
