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
	"fmt"
	"syscall/js"
)

type TitleScene struct {
	camera   js.Value
	light    js.Value
	scene    js.Value
	jsoTitle js.Value
}

func NewTitleScene() *TitleScene {
	ts := &TitleScene{}
	ts.camera = ThreeJsNew("PerspectiveCamera", 60, 1, 1, HelperSize)
	ts.scene = ThreeJsNew("Scene")
	ts.light = ThreeJsNew("PointLight", 0xffffff, 1)
	SetPosition(ts.light,
		0,
		0,
		HelperSize/2,
	)
	ts.scene.Call("add", ts.light)

	axisHelper := ThreeJsNew("AxesHelper", HelperSize)
	ts.scene.Call("add", axisHelper)

	// set title camera pos
	SetPosition(ts.camera,
		0, 0, HelperSize/2,
	)
	ts.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			0, 0, 0,
		),
	)
	ts.camera.Set("zoom", 3.0)
	ts.camera.Call("updateProjectionMatrix")
	return ts
}

func (ts *TitleScene) Resize(w, h float64) {
	ts.camera.Set("aspect", w/h)
	ts.camera.Call("updateProjectionMatrix")
}

func (ts *TitleScene) addTitle() {
	str := "Goguelike-GL"
	geo := MakeTitleTextGeometry(str, 80)
	geoInfo := GetGeoInfo(geo)
	co := gRnd.Uint32() & 0x00ffffff
	mat := gPoolColorMaterial.Get(fmt.Sprintf("#%06x", co))
	mat.Set("opacity", 1)
	ts.jsoTitle = ThreeJsNew("Mesh", geo, mat)
	SetPosition(ts.jsoTitle,
		0-geoInfo.Len[0]/2,
		0,
		0,
	)
	ts.scene.Call("add", ts.jsoTitle)
}

func MakeTitleTextGeometry(str string, size float64) js.Value {
	geo := ThreeJsNew("TextGeometry", str,
		map[string]interface{}{
			"font":           gFont_helvetiker_regular,
			"size":           size,
			"height":         5,
			"curveSegments":  size / 3,
			"bevelEnabled":   true,
			"bevelThickness": size / 8,
			"bevelSize":      size / 16,
			"bevelOffset":    0,
			"bevelSegments":  size / 8,
		})
	return geo
}
