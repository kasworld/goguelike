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

	"github.com/kasworld/goguelike/enum/factiontype"
)

func (vp *Viewport) getColorMaterial(co uint32) js.Value {
	mat, exist := vp.colorMaterialCache[co]
	if !exist {
		mat = vp.ThreeJsNew("MeshPhongMaterial",
			map[string]interface{}{
				"color": co,
				// "flatShading": true,
			},
		)
		vp.colorMaterialCache[co] = mat
	}
	return mat
}

func (vp *Viewport) getAOGeometry(gotype factiontype.FactionType) js.Value {
	geo, exist := vp.aoGeometryCache[gotype]
	if !exist {
		radius := 32
		geo = vp.ThreeJsNew("BoxGeometry", radius*2, radius*2, radius*2)
		vp.aoGeometryCache[gotype] = geo
	}
	return geo
}
