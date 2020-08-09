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

	"github.com/kasworld/goguelike/enum/tile_flag"
)

type Cursor3D struct {
	GeoInfo GeoInfo
	Mesh    [3]js.Value
}

func NewCursor3D() *Cursor3D {
	rtn := &Cursor3D{}
	geo := ThreeJsNew("PlaneGeometry", DstCellSize, DstCellSize)
	rtn.GeoInfo = GetGeoInfo(geo)
	for i, v := range [3]string{"#00ff00", "#0000ff", "#ff0000"} {
		mat := GetColorMaterialByCache(v)
		mat.Set("opacity", 0.5)
		rtn.Mesh[i] = ThreeJsNew("Mesh", geo, mat)
	}
	return rtn
}

func (aog *Cursor3D) SetFieldPosition(fx, fy int, tl tile_flag.TileFlag) {
	height := GetTile3DVisibleTopByCache(tl)
	aog.VisibleByTile(tl)
	for i := range aog.Mesh {
		SetPosition(
			aog.Mesh[i],
			float64(fx)*DstCellSize+DstCellSize/2,
			-float64(fy)*DstCellSize-DstCellSize/2,
			aog.GeoInfo.Len[2]/2+2+height,
		)
	}
}

func (aog *Cursor3D) VisibleByTile(tl tile_flag.TileFlag) {
	for i := range aog.Mesh {
		aog.Visible(i, false)
	}
	if !tl.CharPlaceable() {
		aog.Visible(2, true)
	} else {
		if tl.NoBattle() {
			aog.Visible(0, true)
		} else {
			aog.Visible(1, true)
		}
	}
}

func (ao3d *Cursor3D) Visible(i int, b bool) {
	ao3d.Mesh[i].Set("visible", b)
}
