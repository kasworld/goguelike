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
	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/enum/tile_flag"
)

func (tr *Terrain) Search1stByXYLenList(
	xylenlist findnear.XYLenList,
	sx, sy int,
	filterfn func(t *tile_flag.TileFlag) bool) (int, int, bool) {

	for _, v := range xylenlist {
		x, y := tr.WrapXY(sx+v.X, sy+v.Y)
		if filterfn(&tr.serviceTileArea[x][y]) {
			return x, y, true
		}
	}
	return 0, 0, false
}
