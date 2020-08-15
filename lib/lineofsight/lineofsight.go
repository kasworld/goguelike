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

package lineofsight

import (
	"math"
	"sort"

	"github.com/kasworld/findnear"
)

type PosLen struct {
	X, Y, L float64
}

type PosLenList []PosLen

func (pll PosLenList) Len() int {
	return len(pll)
}
func (pll PosLenList) Swap(i, j int) {
	pll[i], pll[j] = pll[j], pll[i]
}
func (pll PosLenList) Less(i, j int) bool {
	return pll[i].L < pll[j].L
}

func MakePosLenList(srcx, srcy, dstx, dsty float64) PosLenList {
	first := PosLen{srcx, srcy, 0}
	last := PosLen{dstx, dsty,
		math.Sqrt((srcx-dstx)*(srcx-dstx) + (srcy-dsty)*(srcy-dsty))}

	orix, oriy := srcx, srcy
	rtn := make(PosLenList, 0)
	rtn = append(rtn, first)
	dx := dstx - srcx
	dy := dsty - srcy

	if srcx < dstx {
		for x := math.Ceil(srcx); x < dstx; x++ {
			y := x * dy / dx
			l := math.Sqrt((orix-x)*(orix-x) + (oriy-y)*(oriy-y))
			cp := PosLen{x, y, l}
			rtn = append(rtn, cp)
		}
	} else if srcx > dstx {
		for x := math.Floor(srcx); x > dstx; x-- {
			y := x * dy / dx
			l := math.Sqrt((orix-x)*(orix-x) + (oriy-y)*(oriy-y))
			cp := PosLen{x, y, l}
			rtn = append(rtn, cp)
		}
	} else {
		// skip srcx == dstx
	}

	if srcy < dsty {
		for y := math.Ceil(srcy); y < dsty; y++ {
			x := y * dx / dy
			l := math.Sqrt((orix-x)*(orix-x) + (oriy-y)*(oriy-y))
			cp := PosLen{x, y, l}
			rtn = append(rtn, cp)
		}
	} else if srcy > dsty {
		for y := math.Floor(srcy); y > dsty; y-- {
			x := y * dx / dy
			l := math.Sqrt((orix-x)*(orix-x) + (oriy-y)*(oriy-y))
			cp := PosLen{x, y, l}
			rtn = append(rtn, cp)
		}
	} else {
		// skip srcy == dsty
	}

	rtn = append(rtn, last)

	sort.Sort(rtn)
	return rtn
}

// remove dup pos
func (pll PosLenList) DelDup() PosLenList {
	if len(pll) < 2 {
		return pll
	}
	// del dup
	rtn := make(PosLenList, 0, len(pll))
	rtn = append(rtn, pll[0])
	for _, v := range pll[1:] {
		if rtn[len(rtn)-1] == v {
			continue
		}
		rtn = append(rtn, v)
	}
	return rtn
}

// fromsrclen to in square len
func (pll PosLenList) ToCellLenList() findnear.XYLenList {
	if len(pll) == 0 {
		return nil
	}
	// calc diff
	rtn := make(findnear.XYLenList, 0, len(pll))
	rtn = append(rtn, findnear.XYLen{
		X: int(math.Floor(pll[0].X)),
		Y: int(math.Floor(pll[0].Y)),
		L: pll[0].L,
	})
	for i, v := range pll[1:] {
		// loop i == 0, v = pll[1] ...
		last := pll[i]
		rtn = append(rtn, findnear.XYLen{
			X: int(math.Floor((v.X + last.X) / 2)),
			Y: int(math.Floor((v.Y + last.Y) / 2)),
			L: v.L - last.L,
		})
	}
	return rtn
}

func CalcXYLenListLine(x1, y1, x2, y2 int) (findnear.XYLenList, error) {
	return MakePosLenList(
		float64(x1)+0.5, float64(y1)+0.5, float64(x2)+0.5, float64(y2)+0.5,
	).DelDup().ToCellLenList(), nil
}
