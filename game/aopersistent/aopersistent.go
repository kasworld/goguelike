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

package aopersistent

import (
	"time"

	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/game/bias"
)

type AOPersistent struct {
	UUID string // sessionid? aouuid?

	Create   time.Time `prettystring:"simple"` // first
	NickName string
	Faction  factiontype.FactionType // born

	Update time.Time `prettystring:"simple"` // save time
	Exp    float64
	Wealth float64
	Bias   bias.Bias `prettystring:"simple"` // current
}
