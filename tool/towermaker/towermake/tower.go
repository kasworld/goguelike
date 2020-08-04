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

package towermake

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kasworld/goguelike/game/towerscript"
)

type Tower struct {
	name   string
	byList []*Floor
	byName map[string]*Floor
}

func New(name string) *Tower {
	tw := &Tower{
		name:   name,
		byList: make([]*Floor, 0),
		byName: make(map[string]*Floor),
	}
	return tw
}

func (tw *Tower) Add(name string, w, h int, ao, po int, turnBoost float64) *Floor {
	fm := NewFloor(name, w, h, ao, po, turnBoost)
	return tw.AddFloor(name, fm)
}

func (tw *Tower) AddFloor(name string, fm *Floor) *Floor {
	tw.byList = append(tw.byList, fm)
	tw.byName[name] = fm
	return fm
}

func (tw *Tower) GetByName(name string) *Floor {
	return tw.byName[name]
}

func (tw *Tower) GetList() []*Floor {
	return tw.byList
}

func (tw *Tower) GetCount() int {
	return len(tw.byList)
}

func (tw *Tower) Save() error {
	twst := make(towerscript.TowerScript, 0)
	twst = append(twst, []string{
		fmt.Sprintf("# tower %v, made by tower maker", tw.name),
		fmt.Sprintf("# %s %s", filepath.Base(os.Args[0]), strings.Join(os.Args[1:], " ")),
	})
	for _, fm := range tw.byList {
		twst = append(twst, fm.Script)
	}
	twst = append(twst, []string{
		fmt.Sprintf("#tower %v, made by tower maker", tw.name),
	})
	err := twst.SaveJSON(tw.name + ".tower")
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	} else {
		fmt.Printf("tower %v made, total floor %v\n", tw.name, len(tw.byList))
	}
	return nil
}
