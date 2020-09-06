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
	"github.com/kasworld/goguelike/tool/towermaker/noroomtower"
	"github.com/kasworld/goguelike/tool/towermaker/objmaxtower"
	"github.com/kasworld/goguelike/tool/towermaker/roguetower"
	"github.com/kasworld/goguelike/tool/towermaker/sighttower"
	"github.com/kasworld/goguelike/tool/towermaker/starttower"
)

func main() {
	towername := flag.String("towername", "",
		"all,start,roguelike,big,sight,objmax,noroom tower to make")
	floorcount := flag.Int("floorcount", 100, "roguelike,big tower floor count")
	flag.Parse()

	switch *towername {
	default:
		fmt.Println("invalid option")
		flag.PrintDefaults()

	case "all":
		starttower.New("start").Save()
		sighttower.New("sight").Save()
		objmaxtower.New("objmax").Save()
		roguetower.New(fmt.Sprintf("roguelike%v", *floorcount), *floorcount).Save()
		bigtower.New(fmt.Sprintf("big%v", *floorcount), *floorcount).Save()

	case "start":
		fmt.Printf("try make tower with towername:%v\n", *towername)
		starttower.New(*towername).Save()

	case "noroom":
		fmt.Printf("try make tower with towername:%v\n", *towername)
		noroomtower.New(*towername).Save()

	case "sight":
		fmt.Printf("try make tower with towername:%v\n", *towername)
		sighttower.New(*towername).Save()

	case "objmax":
		fmt.Printf("try make tower with towername:%v\n", *towername)
		objmaxtower.New(*towername).Save()

	case "roguelike":
		twname := fmt.Sprintf("%v%v", *towername, *floorcount)
		fmt.Printf("try make tower with towername:%v floorcount:%v\n", twname, *floorcount)
		roguetower.New(twname, *floorcount).Save()

	case "big":
		twname := fmt.Sprintf("%v%v", *towername, *floorcount)
		fmt.Printf("try make tower with towername:%v floorcount:%v\n", twname, *floorcount)
		bigtower.New(twname, *floorcount).Save()

	}

}
