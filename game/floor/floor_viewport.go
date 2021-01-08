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
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/game/activeobject"
	"github.com/kasworld/goguelike/game/dangerobject"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/uuidposmani"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (f *Floor) makeViewportTiles2(centerX, centerY int, sightMat *viewportdata.ViewportSight2,
	sight float32) *viewportdata.ViewportTileArea2 {

	tiles := f.terrain.GetTiles()
	xWrap := f.terrain.GetXWrapper().GetWrapFn()
	yWrap := f.terrain.GetYWrapper().GetWrapFn()

	var rtn viewportdata.ViewportTileArea2
	for i, v := range viewportdata.ViewportXYLenList {
		if sightMat[i] <= sight {
			// make viewport tiles
			fx := xWrap(centerX + v.X)
			fy := yWrap(centerY + v.Y)
			rtn[i] = tiles[fx][fy]
		}
	}
	return &rtn
}

func (f *Floor) makeViewportActiveObjs2(
	vpixyolist []uuidposmani.VPIXYObj,
	sightMat *viewportdata.ViewportSight2, sight float32) []*c2t_obj.ActiveObjClient {

	maxobj := gameconst.ActiveObjCountInViewportLimit
	rtn := make([]*c2t_obj.ActiveObjClient, 0, len(vpixyolist))
	for _, v := range vpixyolist {
		if sightMat[v.I] >= sight {
			continue
		}
		ww := v.O.(*activeobject.ActiveObject)
		if ww.AOTurnData.Condition.TestByCondition(condition.Invisible) {
			// skip invisible ao
			continue
		}
		rtn = append(rtn, ww.ToPacket_ActiveObjClient(v.X, v.Y))
		maxobj--
		if maxobj < 0 {
			f.statPacketObjOver.Inc()
			break
		}
	}
	return rtn
}

func (f *Floor) makeViewportCarryObjs2(
	vpixyolist []uuidposmani.VPIXYObj,
	sightMat *viewportdata.ViewportSight2, sight float32) []*c2t_obj.CarryObjClientOnFloor {

	maxobj := gameconst.CarryObjCountInViewportLimit
	rtn := make([]*c2t_obj.CarryObjClientOnFloor, 0, len(vpixyolist))
	for _, v := range vpixyolist {
		if sightMat[v.I] >= sight {
			continue
		}
		ww := v.O.(gamei.CarryingObjectI)
		rtn = append(rtn, ww.ToPacket_CarryObjClientOnFloor(v.X, v.Y))
		maxobj--
		if maxobj < 0 {
			f.statPacketObjOver.Inc()
			break
		}
	}
	return rtn
}

func (f *Floor) makeViewportFieldObjs2(
	vpixyolist []uuidposmani.VPIXYObj,
	sightMat *viewportdata.ViewportSight2, sight float32) []*c2t_obj.FieldObjClient {

	maxobj := gameconst.FieldObjCountInViewportLimit
	rtn := make([]*c2t_obj.FieldObjClient, 0, len(vpixyolist))
	for _, v := range vpixyolist {
		if sightMat[v.I] >= sight {
			continue
		}
		ww := v.O.(*fieldobject.FieldObject)
		if ww.DisplayType == fieldobjdisplaytype.None {
			continue
		}
		rtn = append(rtn, ww.ToPacket_FieldObjClient(v.X, v.Y))
		maxobj--
		if maxobj < 0 {
			f.statPacketObjOver.Inc()
			break
		}
	}
	return rtn
}

func (f *Floor) makeViewportDangerObjs2(
	vpixyolist []uuidposmani.VPIXYObj,
	sightMat *viewportdata.ViewportSight2, sight float32) []*c2t_obj.DangerObjClient {

	maxobj := gameconst.DangerObjCountInViewportLimit
	rtn := make([]*c2t_obj.DangerObjClient, 0, len(vpixyolist))
	for _, v := range vpixyolist {
		if sightMat[v.I] >= sight {
			continue
		}
		ww := v.O.(*dangerobject.DangerObject)
		rtn = append(rtn, ww.ToPacket_DangerObjClient(v.X, v.Y))
		maxobj--
		if maxobj < 0 {
			f.statPacketObjOver.Inc()
			break
		}
	}
	return rtn
}
