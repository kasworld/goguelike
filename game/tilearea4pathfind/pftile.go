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

package tilearea4pathfind

import (
	"fmt"
	"math"

	"github.com/kasworld/goguelike/enum/way9type"
	astar "github.com/kasworld/goguelike/lib/go-astar"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

var _ astar.Pather = &pfTile{}

var dir2vt = [][2]int{
	way9type.Center:    [2]int{0, 0},
	way9type.North:     [2]int{0, -1},
	way9type.NorthEast: [2]int{1, -1},
	way9type.East:      [2]int{1, 0},
	way9type.SouthEast: [2]int{1, 1},
	way9type.South:     [2]int{0, 1},
	way9type.SouthWest: [2]int{-1, 1},
	way9type.West:      [2]int{-1, 0},
	way9type.NorthWest: [2]int{-1, -1},
}

type pfTile struct {
	X     int
	Y     int
	Ta4pf *TileArea4PathFind
}

func (pft *pfTile) PathNeighbors() []astar.Pather {
	rtn := make([]astar.Pather, 0, 8)
	ta := pft.Ta4pf.tileArea
	fx, fy := pft.X, pft.Y
	for _, v := range dir2vt[1:] {
		x, y := pft.Ta4pf.WrapXY(v[0]+fx, v[1]+fy)
		if ta[x][y].CharPlaceable() {
			rtn = append(rtn, pft.Ta4pf.newTile4PathFind(x, y))
		}
	}
	return rtn
}

func (pft *pfTile) PathNeighborCost(to astar.Pather) float64 {
	toftp := to.(*pfTile)
	contact, dir := way9type.CalcContactDirWrappedXY(
		pft.X, pft.Y,
		toftp.X, toftp.Y,
		pft.Ta4pf.w, pft.Ta4pf.h)
	if !contact {
		panic(fmt.Sprintf("invalid neighbor %v %v", pft, toftp))
	}
	if dir == way9type.Center {
		return 0
	}
	tl := pft.Ta4pf.tileArea[pft.X][pft.Y]
	hp, sp := tl.ActHPSP(c2t_idcmd.Move)
	return (1 + hp + sp) * dir.Len()
}

func (pft *pfTile) PathEstimatedCost(to astar.Pather) float64 {
	toftp := to.(*pfTile)
	dx, dy := way9type.CalcDxDyWrapped(
		toftp.X-pft.X, toftp.Y-pft.Y,
		pft.Ta4pf.w, pft.Ta4pf.h,
	)
	return math.Sqrt(float64(dx*dx + dy*dy))
}
