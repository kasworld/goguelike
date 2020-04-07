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

package c2t_obj

import (
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/lib/g2id"
)

type ReqMeditate_data struct {
	Dummy uint8
}
type RspMeditate_data struct {
	Dummy uint8
}

type ReqKillSelf_data struct {
	Dummy uint8
}
type RspKillSelf_data struct {
	Dummy uint8
}

type ReqMove_data struct {
	Dir way9type.Way9Type
}
type RspMove_data struct {
	Dir way9type.Way9Type
}

type ReqAttack_data struct {
	Dir way9type.Way9Type
}
type RspAttack_data struct {
	Dummy uint8
}

type ReqPickup_data struct {
	G2ID g2id.G2ID
}
type RspPickup_data struct {
	Dummy uint8
}

type ReqDrop_data struct {
	G2ID g2id.G2ID
}
type RspDrop_data struct {
	Dummy uint8
}

type ReqEquip_data struct {
	G2ID g2id.G2ID
}
type RspEquip_data struct {
	Dummy uint8
}

type ReqUnEquip_data struct {
	G2ID g2id.G2ID
}
type RspUnEquip_data struct {
	Dummy uint8
}

type ReqDrinkPotion_data struct {
	G2ID g2id.G2ID
}
type RspDrinkPotion_data struct {
	Dummy uint8
}

type ReqReadScroll_data struct {
	G2ID g2id.G2ID
}
type RspReadScroll_data struct {
	Dummy uint8
}

type ReqRecycle_data struct {
	G2ID g2id.G2ID
}
type RspRecycle_data struct {
	Dummy uint8
}

type ReqEnterPortal_data struct {
	Dummy uint8
}
type RspEnterPortal_data struct {
	Dummy uint8
}

type ReqMoveFloor_data struct {
	G2ID g2id.G2ID
}
type RspMoveFloor_data struct {
	Dummy uint8
}

type ReqActTeleport_data struct {
	Dummy uint8
}
type RspActTeleport_data struct {
	Dummy uint8
}
