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

package floor

import (
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (f *Floor) ToPacket_FloorInfo() *c2t_obj.FloorInfo {
	return &c2t_obj.FloorInfo{
		Name:       f.terrain.GetName(),
		W:          f.w,
		H:          f.h,
		Tiles:      f.terrain.GetTile2Discover(),
		Bias:       f.GetBias(),
		TurnPerSec: f.tower.Config().TurnPerSec / f.terrain.ActTurnBoost,
	}
}

func (f *Floor) ToPacket_NotiAgeing() *c2t_obj.NotiAgeing_data {
	return &c2t_obj.NotiAgeing_data{
		FloorName: f.GetName(),
	}
}
