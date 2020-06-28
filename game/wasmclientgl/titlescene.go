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

type TitleScene struct {
	camera   js.Value
	light    js.Value
	scene    js.Value
	jsoTitle js.Value
}

func NewTitleScene() *TitleScene {
	ts := &TitleScene{}
	ts.camera = ThreeJsNew("PerspectiveCamera", 60, 1, 1, HelperSize*2)
	ts.scene = ThreeJsNew("Scene")
	ts.light = ThreeJsNew("PointLight", 0xffffff, 1)
	SetPosition(ts.light,
		HelperSize,
		HelperSize,
		HelperSize,
	)
	ts.scene.Call("add", ts.light)

	axisHelper := ThreeJsNew("AxesHelper", HelperSize)
	ts.scene.Call("add", axisHelper)

	// set title camera pos
	SetPosition(ts.camera,
		HelperSize/2, HelperSize/2, HelperSize,
	)
	ts.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			HelperSize/2, HelperSize/2, 0,
		),
	)
	ts.camera.Call("updateProjectionMatrix")
	return ts
}

func (ts *TitleScene) Resize(w, h float64) {
	ts.camera.Set("aspect", w/h)
	ts.camera.Call("updateProjectionMatrix")
}

func (ts *TitleScene) addTitle() {
	str := "Goguelike"
	ftGeo := GetTextGeometryByCache(str, 80)
	geoMin, geoMax := CalcGeoMinMaxX(ftGeo)
	co := gRnd.Uint32() & 0x00ffffff
	ftMat := GetColorMaterialByCache(co)
	ts.jsoTitle = ThreeJsNew("Mesh", ftGeo, ftMat)
	SetPosition(ts.jsoTitle,
		HelperSize/2-(geoMax-geoMin)/2,
		HelperSize/2,
		HelperSize/2,
	)
	ts.scene.Call("add", ts.jsoTitle)
}
