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
	"math"

	"github.com/kasworld/goguelike/enum/way9type"
	astar "github.com/kasworld/goguelike/lib/go-astar"
)

type CorridorPather struct {
	X  int
	Y  int
	FF *Terrain
}

var dir2vt8 = [][2]int{
	way9type.North:     [2]int{0, -1},
	way9type.NorthEast: [2]int{1, -1},
	way9type.East:      [2]int{1, 0},
	way9type.SouthEast: [2]int{1, 1},
	way9type.South:     [2]int{0, 1},
	way9type.SouthWest: [2]int{-1, 1},
	way9type.West:      [2]int{-1, 0},
	way9type.NorthWest: [2]int{-1, -1},
}

var dir2vt4 = [][2]int{
	way9type.North: [2]int{0, -1},
	way9type.East:  [2]int{1, 0},
	way9type.South: [2]int{0, 1},
	way9type.West:  [2]int{-1, 0},
}

func (crp *CorridorPather) PathNeighbors() []astar.Pather {
	rtn := make([]astar.Pather, 0, 8)
	if crp.FF.crp8Way {
		for _, v := range dir2vt8 {
			x, y := crp.FF.WrapXY(crp.X+v[0], crp.Y+v[1])
			if crp.FF.corridorPlaceable(x, y) {
				rtn = append(rtn, crp.FF.newCRPather(x, y))
			}
		}
	} else {
		for _, v := range dir2vt4 {
			x, y := crp.FF.WrapXY(crp.X+v[0], crp.Y+v[1])
			if crp.FF.corridorPlaceable(x, y) {
				rtn = append(rtn, crp.FF.newCRPather(x, y))
			}
		}
	}
	return rtn
}

func (crp *CorridorPather) PathNeighborCost(to astar.Pather) float64 {
	toftp := to.(*CorridorPather)
	w, h := crp.FF.GetXYLen()
	contact, dir := way9type.CalcContactDirWrappedXY(
		crp.X, crp.Y, toftp.X, toftp.Y, w, h)
	if !contact {
		crp.FF.log.Warn("invalid neighbor %v %v", crp, toftp)
	}
	return dir.Len()
}

func (crp *CorridorPather) PathEstimatedCost(to astar.Pather) float64 {
	toftp := to.(*CorridorPather)
	w, h := crp.FF.GetXYLen()
	dx, dy := way9type.CalcDxDyWrapped(
		toftp.X-crp.X, toftp.Y-crp.Y, w, h)
	if crp.FF.crp8Way {
		return math.Sqrt(float64(dx*dx + dy*dy))
	} else {
		if dx < 0 {
			dx = -dx
		}
		if dy < 0 {
			dy = -dy
		}
		return float64(dx + dy)
	}
}

func (tr *Terrain) corridorPlaceable(x, y int) bool {
	if tr.roomManager.GetRoomByPos(x, y) != nil {
		return false
	}
	return true
}

func (tr *Terrain) newCRPather(x, y int) *CorridorPather {
	x, y = tr.WrapXY(x, y)
	crp := tr.crpCache[x][y]
	if crp != nil {
		return crp
	}
	crp = &CorridorPather{X: x, Y: y, FF: tr}
	tr.crpCache[x][y] = crp
	return crp
}

func (tr *Terrain) findCorridor(srcx, srcy, dstx, dsty int, trylimit int, way8 bool) [][2]int {
	tr.crp8Way = way8
	p, count := astar.Path2(
		tr.newCRPather(srcx, srcy),
		tr.newCRPather(dstx, dsty), trylimit, tr.Xlen*tr.Ylen)

	if p != nil || count < trylimit {
		rtn := make([][2]int, len(p))
		for i, v := range p {
			vv := v.(*CorridorPather)
			rtn[len(p)-i-1] = [2]int{vv.X, vv.Y}
		}
		return rtn
	}
	return nil
}

func (tr *Terrain) initCrpCache() {
	tr.crpCache = make([][]*CorridorPather, tr.Xlen)
	for i, _ := range tr.crpCache {
		tr.crpCache[i] = make([]*CorridorPather, tr.Ylen)
	}
}
