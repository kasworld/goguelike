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

package fieldobject

import (
	"math"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/lib/lineofsight"
)

// GetLineAttack calc dangerobj wingcount * line
func (fo *FieldObject) GetLineAttack() []findnear.XYLenList {
	rtn := make([]findnear.XYLenList, fo.WingCount)
	wingdeg := 360.0 / float64(fo.WingCount)
	for wing := 0; wing < fo.WingCount; wing++ {
		rad := (float64(wing)*wingdeg + float64(fo.Degree)) / 180.0 * math.Pi
		dx := float64(fo.WingLen) * math.Cos(rad)
		dy := float64(fo.WingLen) * math.Sin(rad)
		rtn[wing] = lineofsight.MakePosLenList(0.5, 0.5, dx+0.5, dy+0.5).ToCellLenList()
	}
	return rtn
}
