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

package leveldata

import (
	"math"

	"github.com/kasworld/goguelike/config/gameconst"
)

var attrib = [gameconst.MaxLevel + 1]struct {
	BaseExp     float64
	WeightLimit float64
	MaxHP       float64
	MaxSP       float64
	MaxAP       float64
	Sight       float64
}{}

func BaseExp(lv int) float64 {
	return attrib[lv].BaseExp
}
func WeightLimit(lv int) float64 {
	return attrib[lv].WeightLimit
}
func MaxHP(lv int) float64 {
	return attrib[lv].MaxHP
}
func MaxSP(lv int) float64 {
	return attrib[lv].MaxSP
}
func MaxAP(lv int) float64 {
	return attrib[lv].MaxAP
}

func Sight(lv int) float64 {
	return attrib[lv].Sight
}

func init() {
	for lv := range attrib {
		attrib[lv].BaseExp = calcBaseExpByLevel(lv)
		attrib[lv].WeightLimit = calcWeightLimitFromLevel(lv)
		attrib[lv].Sight = calcSightByLevel(lv)
		attrib[lv].MaxHP = gameconst.HPBase + gameconst.HPPerLevel*(float64(lv)-1)
		attrib[lv].MaxSP = gameconst.SPBase + gameconst.SPPerLevel*(float64(lv)-1)
		attrib[lv].MaxAP = gameconst.APBase + gameconst.APPerLevel*(float64(lv)-1)
		// fmt.Printf("%v %v\n", lv, attrib[lv])
	}
}

func CalcLevelFromExp(exp float64) float64 {
	lv := math.Sqrt(exp/1000) + 1
	if lv < 1 {
		lv = 1
	}
	if lv > gameconst.MaxLevel {
		lv = gameconst.MaxLevel
	}
	return lv
}

func calcBaseExpByLevel(lv int) float64 {
	if lv < 1 {
		lv = 1
	}
	return float64((lv-1)*(lv-1)) * 1000
}

func calcWeightLimitFromLevel(lv int) float64 {
	return gameconst.LvGram * (7 + float64(lv))
}

func calcSightByLevel(lvi int) float64 {
	lv := float64(lvi)
	return gameconst.SightBase + (lv-1)*gameconst.SightPerLevel
}
