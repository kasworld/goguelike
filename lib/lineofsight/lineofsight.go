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

package lineofsight

import (
	"math"
	"sort"

	"github.com/kasworld/findnear"
)

type PosLen struct {
	X, Y, L float64
}

func (pl *PosLen) UpdataLenFrom(srcx, srcy float64) PosLen {
	pl.L = math.Sqrt((srcx-pl.X)*(srcx-pl.X) + (srcy-pl.Y)*(srcy-pl.Y))
	return *pl
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
	rtn := make(PosLenList, 0)
	rtn = append(rtn, PosLen{X: srcx, Y: srcy}) // add first
	rtn = append(rtn, PosLen{X: dstx, Y: dsty}) // add last

	orix, oriy := srcx, srcy
	dx := dstx - srcx
	dy := dsty - srcy

	if dx > 0 { // srcx -> dstx
		for x := math.Ceil(srcx); x <= math.Floor(dstx); x++ {
			y := (x-orix)*dy/dx + oriy
			rtn = append(rtn, PosLen{X: x, Y: y})
		}
	} else if dx < 0 { // dstx -> srcx
		for x := math.Ceil(dstx); x <= math.Floor(srcx); x++ {
			y := (x-orix)*dy/dx + oriy
			rtn = append(rtn, PosLen{X: x, Y: y})
		}
	} else {
		// skip dx == 0
	}

	if dy > 0 { // srcy -> dsty
		for y := math.Ceil(srcy); y <= math.Floor(dsty); y++ {
			x := (y-oriy)*dx/dy + orix
			rtn = append(rtn, PosLen{X: x, Y: y})
		}
	} else if dy < 0 { // dsty -> srcy
		for y := math.Ceil(dsty); y <= math.Floor(srcy); y++ {
			x := (y-oriy)*dx/dy + orix
			rtn = append(rtn, PosLen{X: x, Y: y})
		}
	} else {
		// skip dy ==0
	}

	for i := range rtn { // calc len from start pos
		rtn[i].UpdataLenFrom(orix, oriy)
	}
	sort.Sort(rtn)
	return rtn.delDup()
}

func (pll PosLenList) delDup() PosLenList {
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

// ToCellLenList calc from-srclen to cell-cross-len
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
	return rtn
}
