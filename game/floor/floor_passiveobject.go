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

package floor

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/game/carryingobject"
	"github.com/kasworld/goguelike/game/gamei"
)

func (f *Floor) canCarryObjPlaceAt(x, y int) bool {
	return f.terrain.GetTiles()[x][y].CharPlaceable()
}

func (f *Floor) addNewRandCarryObj2Floor() error {
	var obj gamei.CarryingObjectI
	switch f.rnd.Intn(4) {
	case 0:
		obj = carryingobject.NewRandFactionEquipObj(f.GetName(), f.GetEnvBias().NearFaction(), f.rnd)
	case 1:
		n := f.rnd.Intn(potiontype.TotalPotionMakeRate)
		obj = carryingobject.NewPotionByMakeRate(n)

	case 2:
		n := f.rnd.Intn(scrolltype.TotalScrollMakeRate)
		obj = carryingobject.NewScrollByMakeRate(n)

	case 3:
		v := f.rnd.NormFloat64Range(100, 50)
		if v < 1 {
			v = 1
		}
		obj = carryingobject.NewMoney(v)
	default:
		return nil
	}
	for try := 5; try > 0; try-- {
		x, y := f.rnd.Intn(f.w), f.rnd.Intn(f.h)
		if f.canCarryObjPlaceAt(x, y) {
			return f.placeCarryObj2FloorAt(x, y, obj)
		}
	}
	return fmt.Errorf("fail to addNewRandCarryObj2Floor")
}

func (f *Floor) placeCarryObj2FloorAt(x, y int, po gamei.CarryingObjectI) error {
	po.SetRemainTurnInFloor()
	return f.poPosMan.AddToXY(po, x, y)
}
