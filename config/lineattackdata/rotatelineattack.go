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

package lineattackdata

import (
	"math"
	"sync"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/lib/lineofsight"
)

// wingCache hold cache winglen wingangle to line
var wingCache = struct {
	data  map[int][360]findnear.XYLenList
	mutex sync.Mutex
}{
	data: make(map[int][360]findnear.XYLenList),
}

func UpdateCache360Line(winglen int) {
	wingCache.mutex.Lock()
	defer wingCache.mutex.Unlock()
	if _, exist := wingCache.data[winglen]; exist {
		return
	}
	var rtn [360]findnear.XYLenList
	for deg := range rtn {
		rad := float64(deg) / 180.0 * math.Pi
		dx := float64(winglen) * math.Cos(rad)
		dy := float64(winglen) * math.Sin(rad)
		rtn[deg] = lineofsight.MakePosLenList(0.5, 0.5, dx+0.5, dy+0.5).ToCellLenList()
	}
	wingCache.data[winglen] = rtn
}

func GetWingLines(winglen int) [360]findnear.XYLenList {
	return wingCache.data[winglen]
}
