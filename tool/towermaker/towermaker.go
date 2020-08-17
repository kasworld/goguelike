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

	"github.com/kasworld/goguelike/tool/towermaker/bigtower"
	"github.com/kasworld/goguelike/tool/towermaker/objmaxtower"
	"github.com/kasworld/goguelike/tool/towermaker/roguetower"
	"github.com/kasworld/goguelike/tool/towermaker/sighttower"
	"github.com/kasworld/goguelike/tool/towermaker/starttower"
)

func main() {
	towername := flag.String("towername", "",
		"starting,roguelike,big,sight,objmax tower to make")
	floorcount := flag.Int("floorcount", 100, "roguelike,big tower floor count")
	flag.Parse()

	switch *towername {
	default:
		fmt.Println("invalid option")
		flag.PrintDefaults()

	case "starting":
		fmt.Printf("try make tower with towername:%v\n", *towername)
		tower := starttower.New(*towername)
		tower.Save()

	case "sight":
		fmt.Printf("try make tower with towername:%v\n", *towername)
		tower := sighttower.New(*towername)
		tower.Save()

	case "objmax":
		fmt.Printf("try make tower with towername:%v\n", *towername)
		tower := objmaxtower.New(*towername)
		tower.Save()

	case "roguelike":
		twname := fmt.Sprintf("%v%v", *towername, *floorcount)
		fmt.Printf("try make tower with towername:%v floorcount:%v\n", twname, *floorcount)
		tower := roguetower.New(twname, *floorcount)
		tower.Save()

	case "big":
		twname := fmt.Sprintf("%v%v", *towername, *floorcount)
		fmt.Printf("try make tower with towername:%v floorcount:%v\n", twname, *floorcount)
		tower := bigtower.New(twname, *floorcount)
		tower.Save()

	}

}
