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

package towerdata

type TowerData struct {
	TowerName      string
	ScriptFilename string
	TurnPerSec     float64
	AutoStart      bool
}

var Default = []TowerData{
	{"Roguelike1", "roguelike100", 1.0, false},
	{"Roguelike2", "roguelike100", 2.0, false},
	{"Roguelike3", "roguelike100", 3.0, true},
	{"Roguelike4", "roguelike100", 4.0, false},
	{"Roguelike5", "roguelike100", 5.0, false},
	{"Goguelike1", "start", 1.0, false},
	{"Goguelike2", "start", 2.0, false},
	{"Goguelike3", "start", 3.0, true},
	{"Goguelike4", "start", 4.0, false},
	{"Goguelike5", "start", 5.0, false},
}
