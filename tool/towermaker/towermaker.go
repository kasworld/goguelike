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

	"github.com/kasworld/goguelike/tool/towermaker/starttower"
)

func main() {
	towername := flag.String("towername", "towermaker", "tower filename w/o .tower")
	floorcount := flag.Int("floorcount", 100, "floor count in tower")
	flag.Parse()

	fmt.Printf("try make tower with towername:%v floorcount:%v\n", *towername, *floorcount)
	// tower := roguetower.MakeRogueTower(*towername, *floorcount)
	tower := starttower.MakeStartTower(*towername)

	tower.Save()
}
