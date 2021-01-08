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

package aoturndata

import (
	"github.com/kasworld/goguelike/enum/condition_flag"
	"github.com/kasworld/goguelike/game/bias"
)

type ActiveObjTurnData struct {
	NonBattleExp float64                      // from visitarea
	TotalExp     float64                      // battleexp + nonbattleexp
	Level        float64                      // from total exp
	Sight        float64                      // from level + buff
	LoadRate     float64                      // from inven /  level
	AttackBias   bias.Bias                    `prettystring:"simple"` // from level + inven equip
	DefenceBias  bias.Bias                    `prettystring:"simple"` // from level + inven equip
	HPMax        float64                      // from level + def bias sum
	SPMax        float64                      // from level + atk bias sum
	Condition    condition_flag.ConditionFlag // current condition
}
