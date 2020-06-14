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
	"syscall/js"

	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (vp *Viewport) add2Scene(o *c2t_obj.ActiveObjClient) js.Value {
	if jso, exist := vp.jsSceneObjs[o.G2ID]; exist {
		SetPosition(jso, o.X, o.Y, 0)
		// jso.Get("rotation").Set("x", o.RotVt[0])
		// jso.Get("rotation").Set("y", o.RotVt[1])
		// jso.Get("rotation").Set("z", o.RotVt[2])
		return jso
	}
	geo := vp.getAOGeometry(o.Faction)
	mat := vp.getColorMaterial(uint32(o.Faction.Color24()))
	jso := vp.ThreeJsNew("Mesh", geo, mat)
	SetPosition(jso, o.X, o.Y, 0)
	// jso.Get("rotation").Set("x", o.RotVt[0])
	// jso.Get("rotation").Set("y", o.RotVt[1])
	// jso.Get("rotation").Set("z", o.RotVt[2])
	vp.scene.Call("add", jso)
	vp.jsSceneObjs[o.G2ID] = jso
	return jso
}

func (vp *Viewport) processRecvStageInfo(
	stageInfo *c2t_obj.NotiObjectList_data) {

	// bgPos := stageInfo.BackgroundPos
	// vp.background.Get("position").Set("x", bgPos[0])
	// vp.background.Get("position").Set("y", bgPos[1])

	// addUUID := make(map[string]bool)
	// for _, o := range stageInfo.ObjList {
	// 	vp.add2Scene(o)
	// 	addUUID[o.G2ID] = true
	// }
	// for id, jso := range vp.jsSceneObjs {
	// 	if !addUUID[id] {
	// 		vp.scene.Call("remove", jso)
	// 		delete(vp.jsSceneObjs, id)
	// 		if lt, exist := vp.lightCache[id]; exist {
	// 			vp.scene.Call("remove", lt.Helper)
	// 		}
	// 	}
	// }
}
