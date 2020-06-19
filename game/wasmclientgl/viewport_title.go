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

func (vp *Viewport) hideTitle() {
	vp.scene.Call("remove", vp.jsoTitle)
	// SetPosition(vp.camera,
	// 	HelperSize/2, HelperSize/2, HelperSize,
	// )
	// vp.camera.Call("updateProjectionMatrix")
	// vp.scene.Call("remove", vp.light)
}

func (vp *Viewport) initTitle() {
	// init light
	vp.light = vp.ThreeJsNew("PointLight", 0xffffff, 1)
	SetPosition(vp.light,
		HelperSize,
		HelperSize,
		HelperSize,
	)
	vp.scene.Call("add", vp.light)

	// set title camera pos
	SetPosition(vp.camera,
		HelperSize/2, HelperSize/2, HelperSize,
	)
	vp.camera.Call("lookAt",
		vp.ThreeJsNew("Vector3",
			HelperSize/2, HelperSize/2, 0,
		),
	)
	vp.camera.Call("updateProjectionMatrix")

	vp.fontLoader.Call("load", "/fonts/helvetiker_regular.typeface.json",
		js.FuncOf(vp.fontLoaded),
	)
}

func (vp *Viewport) fontLoaded(this js.Value, args []js.Value) interface{} {
	vp.font_helvetiker_regular = args[0]
	str := "Goguelike"

	ftGeo := vp.getTextGeometry(str, 80)
	geoMin, geoMax := vp.calcGeoMinMax(ftGeo)

	co := vp.rnd.Uint32() & 0x00ffffff
	ftMat := vp.getColorMaterial(co)

	vp.jsoTitle = vp.ThreeJsNew("Mesh", ftGeo, ftMat)
	SetPosition(vp.jsoTitle,
		HelperSize/2-(geoMax-geoMin)/2,
		HelperSize/2,
		HelperSize/2,
	)
	vp.scene.Call("add", vp.jsoTitle)

	return nil
}
