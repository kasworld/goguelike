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

package gameconst

import "math"

// game const
const (

	// packet viewport size
	ViewPortW  = 19
	ViewPortH  = 19
	ViewPortWH = ViewPortW * ViewPortH

	// client gl viewport tile size
	ClientViewPortW = ViewPortW * 2
	ClientViewPortH = ViewPortH * 2

	// int packet
	ActiveObjCountInViewportLimit = 100
	CarryObjCountInViewportLimit  = 100
	FieldObjCountInViewportLimit  = 100

	// turn to del in floor CarryObj
	CarryingObjectLifeTurnInFloor = 400

	InitCarryObjEquipCount = 2
	InitPotionCount        = 4
	InitScrollCount        = 1
	InitGoldMean           = 100

	DropCarryObjonActiveObjDeathRate = 0.5

	TowerBaseBiasLen     = 100
	FloorBaseBiasLen     = 100
	ActiveObjBaseBiasLen = 100.0

	MaxLevel = 100

	HPBase     = 100.0
	HPPerLevel = 20.0

	SPBase     = 100.0
	SPPerLevel = 10.0

	SightBase     = 10
	SightXray     = math.MaxFloat32
	SightBlock    = SightXray / ViewPortWH
	SightPerLevel = 0.5

	RebirthHPRate = 0.8
	RebirthSPRate = 0.8

	ActiveObjRebirthWaitTurn = 30

	// carryingobject weight, price
	MoneyGram     = 0.1
	EquipABSGram  = 100.0
	EquipABSValue = 100.0
	PotionGram    = 100.0
	PotionValue   = 100.0
	ScrollGram    = 100.0
	ScrollValue   = 100.0

	LvGram = ActiveObjBaseBiasLen/4*EquipABSGram + PotionGram*2 + ScrollGram*1 + MoneyGram*10000

	CarryObjRecycleRate = 0.5

	MaxChatLen = 80
)

// activeobject experience constant
const (
	// explore experience

	// find floor
	ActiveObjExp_VisitFloor = 500
	// discover all floor per tile
	ActiveObjExp_CompleteFloorPerTile = 1
	// discover a tile
	ActiveObjExp_DiscoverTile = 1

	// battle experience
	// kill ao per level
	ActiveObjExp_KillLevel = 10
	// deal damage ao per damage
	ActiveObjExp_Damage = 1
	// reduce battle exp when die
	ActiveObjExp_DieRate = 0.99
)
