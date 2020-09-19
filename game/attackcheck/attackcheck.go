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

package attackcheck

import (
	"github.com/kasworld/go-abs"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/way9type"
)

func IsInLongAttack(x1, y1, x2, y2, w, h int) (bool, way9type.Way9Type) {
	absx := abs.Absi(x1 - x2)
	absy := abs.Absi(y1 - y2)
	if absx >= gameconst.AttackLongLen || absy >= gameconst.AttackLongLen {
		return false, way9type.Center
	}
	isWay9 := absx == 0 || absy == 0 || absx == absy
	dx, dy := way9type.CalcDxDyWrapped(x2-x1, y2-y1, w, h)
	way := way9type.RemoteDxDy2Way9(dx, dy)
	return isWay9, way
}
