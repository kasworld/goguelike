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

package corridor

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/tile"
)

type Corridor struct {
	Tile tile.Tile
	P    [][2]int
}

func (cr Corridor) String() string {
	l := len(cr.P)
	var st, ed [2]int
	if l > 1 {
		st = cr.P[0]
		ed = cr.P[l-1]
	}
	return fmt.Sprintf("Corridor[%v %v %v-%v]",
		cr.Tile, l, st, ed)
}
