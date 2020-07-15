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

package wasmclientgl

import (
	"syscall/js"

	"github.com/kasworld/goguelike/config/gameconst"
)

type RaycastPlane struct {
	Mesh js.Value
}

func NewRaycastPlane() *RaycastPlane {
	w := gameconst.ClientViewPortW * DstCellSize
	h := gameconst.ClientViewPortH * DstCellSize
	Mat := ThreeJsNew("MeshBasicMaterial")
	Geo := ThreeJsNew("PlaneGeometry", w, h)
	Mesh := ThreeJsNew("Mesh", Geo, Mat)
	return &RaycastPlane{
		Mesh: Mesh,
	}
}

func (pl *RaycastPlane) MoveCenterTo(fx, fy int) {
	SetPosition(pl.Mesh,
		fx*DstCellSize,
		-fy*DstCellSize,
		0,
	)
}
