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

type textGeoKey struct {
	Str  string
	Size float64
}

func (vp *Viewport) getTextGeometry(str string, size float64) js.Value {
	geo, exist := vp.textGeometryCache[textGeoKey{str, size}]
	if !exist {
		geo = vp.ThreeJsNew("TextGeometry", str,
			map[string]interface{}{
				"font":           vp.font_helvetiker_regular,
				"size":           size,
				"height":         5,
				"curveSegments":  12,
				"bevelEnabled":   true,
				"bevelThickness": 4,
				"bevelSize":      2,
				"bevelOffset":    0,
				"bevelSegments":  5,
			})
		vp.textGeometryCache[textGeoKey{str, size}] = geo
	}
	return geo
}

func (vp *Viewport) calcGeoMinMax(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("x").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("x").Float()
	return geoMin, geoMax
}
