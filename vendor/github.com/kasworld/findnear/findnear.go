// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// 2d tile space find functions
package findnear

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

var g_xylenList map[[2]int]XYLenList
var g_mutex sync.Mutex

func init() {
	g_xylenList = make(map[[2]int]XYLenList)
}

func G_XYLenListInfo() string {
	return fmt.Sprintf("%v", g_xylenList)
}

func NewXYLenList(xmax, ymax int) XYLenList {
	g_mutex.Lock()
	defer g_mutex.Unlock()

	xyll, exist := g_xylenList[[2]int{xmax, ymax}]
	if exist {
		return xyll
	}

	xyll = makeNewXYLenList(xmax, ymax)
	g_xylenList[[2]int{xmax, ymax}] = xyll
	return xyll
}

func makeNewXYLenList(xmax, ymax int) XYLenList {
	rtn := make(XYLenList, 0)
	for i := 0; i < xmax; i++ {
		for j := 0; j < ymax; j++ {
			x := i - xmax/2
			y := j - ymax/2
			rtn = append(rtn, XYLen{
				x, y,
				math.Sqrt(float64(x*x + y*y)),
			})
		}
	}
	sort.Sort(rtn)
	return rtn
}

type XYLen struct {
	X, Y int
	L    float64
}

func (xyl XYLen) String() string {
	return fmt.Sprintf("XYLen[%d %d %.2f]",
		xyl.X, xyl.Y, xyl.L,
	)
}

func (xyl XYLen) Sub(arg XYLen) XYLen {
	return XYLen{
		xyl.X - arg.X,
		xyl.Y - arg.Y,
		xyl.L - arg.L,
	}
}

func (xyl XYLen) IsSamePos(xyl2 XYLen) bool {
	return xyl.X == xyl2.X && xyl.Y == xyl2.Y
}

type XYLenList []XYLen

// func (s XYLenList) String() string {
// 	return fmt.Sprintf("XYLenList[%v]", len(s))
// }

func (s XYLenList) Len() int {
	return len(s)
}
func (s XYLenList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s XYLenList) Less(i, j int) bool {
	return s[i].L < s[j].L
}

func (pll XYLenList) SumLen() float64 {
	rtn := 0.0
	for _, v := range pll {
		rtn += v.L
	}
	return rtn
}

func (pll XYLenList) MakePos2Index() map[[2]int]int {
	rtn := make(map[[2]int]int)
	for i, v := range pll {
		rtn[[2]int{v.X, v.Y}] = i
	}
	return rtn
}

// search from center
type DoFn func(int, int) bool

func (pll XYLenList) FindAll(x, y int, fn DoFn) bool {
	return pll.Find(x, y, 0, len(pll), fn)
}

func (pll XYLenList) Find(x, y int, start, end int, fn DoFn) bool {
	if start > end || start < 0 || end > len(pll) {
		return false
	}
	for _, v := range pll[start:end] {
		if fn(x+v.X, y+v.Y) {
			return true
		}
	}
	return false
}
