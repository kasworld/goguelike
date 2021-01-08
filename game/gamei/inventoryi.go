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

package gamei

import "github.com/kasworld/goguelike/enum/equipslottype"

type InventoryI interface {
	GetEquipSlot() [equipslottype.EquipSlotType_Count]EquipObjI
	GetTypeList() (
		[]EquipObjI, []PotionI, []ScrollI)
	GetEquipList() []EquipObjI
	GetPotionList() []PotionI
	GetScrollList() []ScrollI
	GetTotalWeight() float64

	RemoveByUUID(poid string) CarryingObjectI
	GetByUUID(poid string) CarryingObjectI
	AddToBag(po CarryingObjectI) error

	RecycleCarryObjByID(poid string) (float64, error)

	AddToWallet(po MoneyI) error
	SubFromWallet(po MoneyI) error
	GetWalletValue() float64

	EquipFromBagByUUID(id string) error
	UnEquipToBagByUUID(id string) (EquipObjI, error)
}
