// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tilearea4pathfind

import (
	"fmt"
	"sync"

	"github.com/kasworld/goguelike/game/tilearea"
	astar "github.com/kasworld/goguelike/lib/go-astar"
	"github.com/kasworld/wrapper"
)

type PFAreaI interface {
	GetXYLen() (int, int)
	GetTiles() tilearea.TileArea
}

type TileArea4PathFind struct {
	tiles4PathFindMutex sync.Mutex        `prettystring:"hide"`
	tiles4PathFind      [][]*pfTile       `prettystring:"simple"`
	tileArea            tilearea.TileArea `prettystring:"simple"`
	w                   int
	h                   int
	xWrapper            *wrapper.Wrapper `prettystring:"simple"`
	yWrapper            *wrapper.Wrapper `prettystring:"simple"`
	xWrap               func(i int) int
	yWrap               func(i int) int
}

func New(tr PFAreaI) *TileArea4PathFind {
	ta4pf := &TileArea4PathFind{
		tileArea: tr.GetTiles(),
	}
	ta4pf.w, ta4pf.h = tr.GetXYLen()
	ta4pf.xWrapper = wrapper.New(ta4pf.w)
	ta4pf.yWrapper = wrapper.New(ta4pf.h)
	ta4pf.xWrap = ta4pf.xWrapper.GetWrapFn()
	ta4pf.yWrap = ta4pf.yWrapper.GetWrapFn()
	ta4pf.tiles4PathFind = make([][]*pfTile, ta4pf.w)
	for i, _ := range ta4pf.tiles4PathFind {
		ta4pf.tiles4PathFind[i] = make([]*pfTile, ta4pf.h)
	}
	return ta4pf
}

func (ta4pf *TileArea4PathFind) Cleanup() {
	ta4pf.tileArea = nil
	for i, v := range ta4pf.tiles4PathFind {
		for j, _ := range v {
			v[j] = nil
		}
		ta4pf.tiles4PathFind[i] = nil
	}
	ta4pf.tiles4PathFind = nil
}

func (ta4pf *TileArea4PathFind) WrapXY(x, y int) (int, int) {
	return ta4pf.xWrap(x), ta4pf.yWrap(y)
}

func (ta4pf *TileArea4PathFind) String() string {
	return fmt.Sprintf("TileArea4PathFind[]")
}

func (ta4pf *TileArea4PathFind) newTile4PathFind(x, y int) *pfTile {
	x, y = ta4pf.WrapXY(x, y)
	ta4pf.tiles4PathFindMutex.Lock()
	defer ta4pf.tiles4PathFindMutex.Unlock()
	tpf := ta4pf.tiles4PathFind[x][y]
	if tpf != nil {
		// need for check same pointer
		return tpf
	}
	tpf = &pfTile{X: x, Y: y, Ta4pf: ta4pf}
	ta4pf.tiles4PathFind[x][y] = tpf
	return tpf
}

func (ta4pf *TileArea4PathFind) FindPath(dstx, dsty, srcx, srcy int, trylimit int) [][2]int {
	p, count := astar.Path2(
		ta4pf.newTile4PathFind(srcx, srcy),
		ta4pf.newTile4PathFind(dstx, dsty),
		trylimit,
		ta4pf.w+ta4pf.h,
	)

	if p != nil || count < trylimit {
		rtn := make([][2]int, len(p))
		for i, v := range p {
			vv := v.(*pfTile)
			rtn[len(p)-i-1] = [2]int{vv.X, vv.Y}
		}
		return rtn
	}
	return nil
}
