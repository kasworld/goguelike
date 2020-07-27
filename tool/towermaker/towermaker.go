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

package main

import (
	"flag"
	"fmt"

	"github.com/kasworld/goguelike/game/towerscript"
	"github.com/kasworld/goguelike/tool/towermaker/floormake"
	"github.com/kasworld/goguelike/tool/towermaker/starttower"
)

func main() {
	towername := flag.String("towername", "roguegen", "tower filename w/o .tower")
	floorcount := flag.Int("floorcount", 100, "floor count in tower")
	flag.Parse()

	fmt.Printf("try make tower with towername:%v floorcount:%v\n", *towername, *floorcount)
	// floorList := roguetower.MakeRogueTower( *floorcount)
	floorList := starttower.MakeStartTower()

	Save(*towername, floorList)
}

func Save(towername string, floorList floormake.FloorList) {
	tw := make(towerscript.TowerScript, 0)
	tw = append(tw, []string{
		fmt.Sprintf("#tower %v, made by tower maker", towername),
	})
	for _, fm := range floorList {
		tw = append(tw, fm.Script)
	}
	err := tw.SaveJSON(towername + ".tower")
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("tower %v made, total floor %v\n", towername, len(floorList))
	}
}
