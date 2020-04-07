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

package floor

import (
	"time"

	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (f *Floor) ToPacket_FloorInfo() *c2t_obj.FloorInfo {
	return &c2t_obj.FloorInfo{
		Name:       f.terrain.GetName(),
		G2ID:       g2id.NewFromString(f.uuid),
		W:          f.w,
		H:          f.h,
		Tiles:      f.terrain.GetTile2Discover(),
		Bias:       f.GetBias(),
		TurnPerSec: f.tower.Config().TurnPerSec / f.terrain.ActTurnBoost,
	}
}

func (f *Floor) ToPacket_NotiAgeing() *c2t_obj.NotiAgeing_data {
	return &c2t_obj.NotiAgeing_data{
		G2ID: g2id.NewFromString(f.uuid),
	}
}

func (f *Floor) ToPacket_NotiObjectList(
	turnTime time.Time,
	cache *CacheVPIXYOList,
	x, y int, sight float64) *c2t_obj.NotiObjectList_data {

	x, y = f.terrain.WrapXY(x, y)
	vpixyolists := cache.GetAtByCache(x, y)

	sightMat := f.terrain.GetViewportCache().GetByCache(x, y)
	aOs := f.makeViewportActiveObjs2(vpixyolists[0], sightMat, float32(sight))
	pOs := f.makeViewportCarryObjs2(vpixyolists[1], sightMat, float32(sight))
	fOs := f.makeViewportFieldObjs2(vpixyolists[2], sightMat, float32(sight))

	return &c2t_obj.NotiObjectList_data{
		Time:          turnTime,
		FloorG2ID:     g2id.NewFromString(f.uuid),
		ActiveObjList: aOs,
		CarryObjList:  pOs,
		FieldObjList:  fOs,
	}
}

func (f *Floor) ToPacket_NotiTileArea(
	x, y int, sight float64) *c2t_obj.NotiVPTiles_data {

	x, y = f.terrain.WrapXY(x, y)
	sightMat := f.terrain.GetViewportCache().GetByCache(x, y)
	cstiles := f.makeViewportTiles2(x, y, sightMat, float32(sight))

	return &c2t_obj.NotiVPTiles_data{
		FloorG2ID: g2id.NewFromString(f.uuid),
		VPX:       x,
		VPY:       y,
		VPTiles:   cstiles,
	}
}
