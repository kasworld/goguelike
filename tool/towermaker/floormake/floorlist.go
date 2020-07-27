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

package floormake

import (
	"fmt"

	"github.com/kasworld/goguelike/game/towerscript"
)

type FloorList []*FloorMake

func (fl FloorList) GetByName(name string) *FloorMake {
	for _, fm := range fl {
		if fm.Name == name {
			return fm
		}
	}
	return nil
}

func (fl FloorList) Save(towername string) error {
	tw := make(towerscript.TowerScript, 0)
	for _, fm := range fl {
		tw = append(tw, fm.Script)
	}
	tw = append(tw, []string{
		fmt.Sprintf("#tower %v, made by tower maker", towername),
	})
	err := tw.SaveJSON(towername + ".tower")
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	} else {
		fmt.Printf("tower %v made, total floor %v\n", towername, len(fl))
	}
	return nil
}
