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
	"math"
	"syscall/js"
)

const StageSize = 1024.0

func (vp *Viewport) makeGridHelper(
	co uint32,
	x, y, z float64,
	lookat js.Value,
) js.Value {
	helper := vp.ThreeJsNew("GridHelper",
		StageSize, 10, co, 0x404040)
	helper.Get("position").Set("x", x)
	helper.Get("position").Set("y", y)
	helper.Get("position").Set("z", z)
	helper.Get("geometry").Call("rotateX", math.Pi/2)
	helper.Call("lookAt", lookat)
	return helper
}

func (vp *Viewport) initHelpers() {
	min := 0.0
	max := StageSize
	mid := StageSize / 2
	center := vp.ThreeJsNew("Vector3",
		mid, mid, mid,
	)

	// y min
	vp.scene.Call("add", vp.makeGridHelper(
		0x0000ff, mid, min, mid, center,
	))

	// y max
	vp.scene.Call("add", vp.makeGridHelper(
		0xffff00, mid, max, mid, center,
	))

	// x min
	vp.scene.Call("add", vp.makeGridHelper(
		0xff0000, min, mid, mid, center,
	))

	// x max
	vp.scene.Call("add", vp.makeGridHelper(
		0x00ffff, max, mid, mid, center,
	))

	// z min
	vp.scene.Call("add", vp.makeGridHelper(
		0x00ff00, mid, mid, min, center,
	))

	// z max
	vp.scene.Call("add", vp.makeGridHelper(
		0xff00ff, mid, mid, max, center,
	))

	box3 := vp.ThreeJsNew("Box3",
		vp.ThreeJsNew("Vector3", 0, 0, 0),
		vp.ThreeJsNew("Vector3", StageSize, StageSize, StageSize),
	)
	helper := vp.ThreeJsNew("Box3Helper", box3, 0xffffff)
	vp.scene.Call("add", helper)
}
