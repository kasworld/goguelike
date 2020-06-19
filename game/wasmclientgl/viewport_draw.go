// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wasmclientgl

import (
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (vp *Viewport) processNotiObjectList(
	olNoti *c2t_obj.NotiObjectList_data) {
	clFd := vp.floorG2ID2ClientField[olNoti.FloorG2ID]
	shY := int(-float64(clFd.CellSize) * 0.8)
	addUUID := make(map[g2id.G2ID]bool)
	for _, o := range olNoti.ActiveObjList {
		if jso, exist := vp.jsSceneObjs[o.G2ID]; exist {
			SetPosition(
				jso,
				o.X*clFd.CellSize,
				-o.Y*clFd.CellSize+shY,
				0)
		} else {
			geo := vp.getTextGeometry(
				o.Faction.String()[:2],
				clFd.CellSize/2,
			)
			mat := vp.getColorMaterial(uint32(o.Faction.Color24()))
			jso := vp.ThreeJsNew("Mesh", geo, mat)
			SetPosition(
				jso,
				o.X*clFd.CellSize,
				-o.Y*clFd.CellSize+shY,
				0)
			vp.scene.Call("add", jso)
			vp.jsSceneObjs[o.G2ID] = jso
		}
		addUUID[o.G2ID] = true
	}
	for id, jso := range vp.jsSceneObjs {
		if !addUUID[id] {
			vp.scene.Call("remove", jso)
			delete(vp.jsSceneObjs, id)
		}
	}
}
