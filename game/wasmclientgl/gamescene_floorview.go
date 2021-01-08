// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
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
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/clientfloor"
)

// floorview frame update
func (vp *GameScene) UpdateFloorViewFrame(
	cf *clientfloor.ClientFloor, vpx, vpy int, envBias bias.Bias) {

	vp.makeClientTile4FloorView(cf, vpx, vpy)
	vp.updateFieldObjInView(cf, vpx, vpy)
	vp.animateFieldObj()
	vp.animateTile(envBias)
	vp.moveCameraLight(
		cf, vpx, vpy,
		0, 0,
		envBias,
	)
	vp.renderer.Call("render", vp.scene, vp.camera)
}

// make floor tiles for floorview
func (vp *GameScene) makeClientTile4FloorView(
	cf *clientfloor.ClientFloor, vpx, vpy int) {
	for i := 0; i < tile.Tile_Count; i++ {
		vp.jsTile3DCount[i] = 0     // clear use count
		vp.jsTile3DDarkCount[i] = 0 // clear use count
	}
	matrix := ThreeJsNew("Matrix4")
	for _, v := range gXYLenListView {
		fx := v.X + vpx
		fy := v.Y + vpy
		newTile := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		for ti := 0; ti < tile.Tile_Count; ti++ {
			if !newTile.TestByTile(tile.Tile(ti)) {
				continue
			}
			matrix.Call("setPosition",
				gTile3D[ti].MakePosVector3(fx, fy),
			)
			vp.jsTile3DMesh[ti].Call("setMatrixAt",
				vp.jsTile3DCount[ti], matrix)
			vp.jsTile3DCount[ti]++
		}
	}
	for ti := 0; ti < tile.Tile_Count; ti++ {
		vp.jsTile3DMesh[ti].Set("count", vp.jsTile3DCount[ti])
		vp.jsTile3DMesh[ti].Get("instanceMatrix").Set("needsUpdate", true)
		vp.jsTile3DDarkMesh[ti].Set("count", vp.jsTile3DDarkCount[ti])
		vp.jsTile3DDarkMesh[ti].Get("instanceMatrix").Set("needsUpdate", true)
	}
}
