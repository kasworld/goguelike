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

	if dx > 0 {
		for x := math.Ceil(srcx); x <= math.Floor(dstx); x++ {
			y := x * dy / dx
			rtn = append(rtn, PosLen{x, y, math.Sqrt((orix-x)*(orix-x) + (oriy-y)*(oriy-y))})
		}
	} else if dx < 0 {
		for x := math.Floor(srcx); x >= math.Ceil(dstx); x-- {
			y := x * dy / dx
			rtn = append(rtn, PosLen{x, y, math.Sqrt((orix-x)*(orix-x) + (oriy-y)*(oriy-y))})
		}
	} else {
		// skip dx == 0
	}

	if dy > 0 {
		for y := math.Ceil(srcy); y <= math.Floor(dsty); y++ {
			x := y * dx / dy
			rtn = append(rtn, PosLen{x, y, math.Sqrt((orix-x)*(orix-x) + (oriy-y)*(oriy-y))})
		}
	} else if dy < 0 {
		for y := math.Floor(srcy); y >= math.Ceil(dsty); y-- {
			x := y * dx / dy
			rtn = append(rtn, PosLen{x, y, math.Sqrt((orix-x)*(orix-x) + (oriy-y)*(oriy-y))})
		}
	} else {
		// skip dy ==0
	}

	rtn = append(rtn, last)

	sort.Sort(rtn)
	if len(rtn) < 2 {
		return rtn
	}
	// del dup
	rtn2 := make(PosLenList, 0, len(rtn))
	rtn2 = append(rtn2, rtn[0])
	for _, v := range rtn[1:] {
		if rtn2[len(rtn2)-1] == v {
			continue
		}
		rtn2 = append(rtn2, v)
	}
	return rtn2
}

// fromsrclen to in square len
func (pll PosLenList) ToCellLenList() findnear.XYLenList {
	if len(pll) == 0 {
		return nil
	}
	// calc diff
	rtn := make(findnear.XYLenList, 0, len(pll))
	for i, v := range pll[:len(pll)-1] {
		next := pll[i+1]
		rtn = append(rtn, findnear.XYLen{
			X: int(math.Floor((v.X + next.X) / 2)),
			Y: int(math.Floor((v.Y + next.Y) / 2)),
			L: next.L - v.L,
		})
	}
	// calc last?
	// last := pll[len(pll)-1]
	// rtn = append(rtn, findnear.XYLen{
	// 	X: int(math.Floor((last.X))),
	// 	Y: int(math.Floor((last.Y))),
	// 	L: 0,
	// })
	return rtn
}

func CalcXYLenListLine(x1, y1, x2, y2 int) (findnear.XYLenList, error) {
	return MakePosLenList(
		// float64(x1), float64(y1), float64(x2), float64(y2),
		float64(x1)+0.5, float64(y1)+0.5, float64(x2)+0.5, float64(y2)+0.5,
	).ToCellLenList(), nil
}
