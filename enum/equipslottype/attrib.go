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

package equipslottype

import (
	"github.com/kasworld/goguelike/config/gameconst"
)

func (et EquipSlotType) Attack() bool {
	return attrib[et].Attack
}

func (et EquipSlotType) Defence() bool {
	return attrib[et].Defence
}

func (et EquipSlotType) Rune() string {
	return attrib[et].Rune
}

func (et EquipSlotType) Names() []string {
	return attrib[et].Names
}

func (et EquipSlotType) Materials() []string {
	return attrib[et].Materials
}

func (et EquipSlotType) BaseLen() float64 {
	return float64(gameconst.ActiveObjBaseBiasLen) * attrib[et].CreateBiasRate
}

var attrib = [EquipSlotType_Count]struct {
	Attack         bool
	Defence        bool
	CreateBiasRate float64 // bias rate of gameconst.ActiveObjBaseBiasLen when create
	Rune           string
	Names          []string
	Materials      []string
}{
	Weapon: {true, false, .5, "/",
		[]string{
			"Dagger", "Mace", "Blade", "Spear", "Hammer", "Club", "Sword", "Pike", "Sickle",
			"Knife", "Razor", "Trident", "Axe", "Staff", "Machete", "Rapier", "Scimitar",
			"Trident", "Cutlass", "Flail", "Gladius", "Bardiche", "Falchion", "Flail",
			"Glaive", "Guisarme", "Halberd", "Lance", "hook", "Pickaxe", "Scythe",
			"Nunchaku", "Sabre", "Whip", "Estoc", "Flindbar", "Fauchard", "Harpoon",
		},
		[]string{
			"Wood", "Bone",
			"Stone", "Crystal", "Diamond",
			"Copper", "Brass", "Bronze", "Iron", "Steel", "Duralumin", "Silver", "Gold", "Platinum", "Titanium",
		},
	},
	Shield: {false, true, .25, "(",
		[]string{
			"Buckler", "LightShield", "MediumShield", "HeavyShield", "TowerShield",
		},
		[]string{
			"Wood", "Bone", "Hide",
			"Stone",
			"Copper", "Brass", "Bronze", "Iron", "Steel", "Duralumin", "Silver", "Gold", "Platinum", "Titanium",
		},
	},
	Helmet: {false, true, 0.125, "^",
		[]string{
			"Circlet", "Crown", "Hat", "Helm", "Hood", "Headband",
		},
		[]string{
			"Wood", "Bone", "Hide", "Silk",
			"Stone",
			"Copper", "Brass", "Bronze", "Iron", "Steel", "Duralumin", "Silver", "Gold", "Platinum", "Titanium",
		},
	},
	Armor: {false, true, .25, "%",
		[]string{
			"Armor", "Coat", "Tunic", "Shirt", "Bodysuit", "Mail", "Chainmail", "Bandedmail",
			"Splintmail", "Fieldplate", "Breastplate", "Halfplate", "Fullplate",
		},
		[]string{
			"Wood", "Bone", "Hide", "Silk",
			"Stone",
			"Copper", "Brass", "Bronze", "Iron", "Steel", "Duralumin", "Silver", "Gold", "Platinum", "Titanium",
		},
	},
	Gauntlet: {false, true, 0.125, "=",
		[]string{
			"Gauntlet",
		},
		[]string{
			"Wood", "Bone", "Hide", "Silk",
			"Stone",
			"Copper", "Brass", "Bronze", "Iron", "Steel", "Duralumin", "Silver", "Gold", "Platinum", "Titanium",
		},
	},
	Footwear: {false, true, 0.125, ",,",
		[]string{
			"Boots", "Sandals", "Shoes", "Slippers",
		},
		[]string{
			"Wood", "Bone", "Hide", "Silk",
			"Stone",
			"Copper", "Brass", "Bronze", "Iron", "Steel", "Duralumin", "Silver", "Gold", "Platinum", "Titanium",
		},
	},
	Ring: {true, false, 0.125, "o",
		[]string{
			"Ring",
		},
		[]string{
			"Wood", "Bone",
			"Stone", "Coral", "Crystal", "Ruby", "Diamond",
			"Copper", "Brass", "Bronze", "Iron", "Steel", "Duralumin", "Silver", "Gold", "Platinum", "Titanium",
		},
	},
	Amulet: {true, false, 0.125, "+",
		[]string{
			"Amulet",
		},
		[]string{
			"Wood", "Bone",
			"Stone", "Coral", "Crystal", "Ruby", "Diamond",
			"Copper", "Brass", "Bronze", "Iron", "Steel", "Duralumin", "Silver", "Gold", "Platinum", "Titanium",
		},
	},
}
