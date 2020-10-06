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

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/tile"
)

type GameScene struct {
	renderer js.Value

	// from client floor gl

	lightRGB  [3]js.Value // rgb light
	lightW    js.Value    // white light
	scene     js.Value
	camera    js.Value
	raycaster js.Value

	raycastPlane *RaycastPlane
	cursor       *Cursor3D

	mouseCursorFx int
	mouseCursorFy int

	// player ao bars
	HP *ColorBar3D
	SP *ColorBar3D
	AP *ColorBar3D

	jsSceneCOs map[string]*CarryObj3D  // in sight only carryobj
	jsSceneAOs map[string]*ActiveObj3D // in sight only ao
	jsSceneDOs map[string]*DangerObj3D // in sight only dangerobj
	jsSceneFOs map[string]*FieldObj3D  // in clientview fieldobj

	jsSceneMovePathArrows map[[2]int]*ColorArrow3D // in clientview move path arrow

	// tile 3d instancedmesh
	// count = gameconst.ClientViewPortW * gameconst.ClientViewPortH
	jsTile3DMesh      [tile.Tile_Count]js.Value
	jsTile3DCount     [tile.Tile_Count]int // in use count
	jsTile3DDarkMesh  [tile.Tile_Count]js.Value
	jsTile3DDarkCount [tile.Tile_Count]int // in use count

}

func NewGameScene() *GameScene {
	vp := &GameScene{
		jsSceneCOs:            make(map[string]*CarryObj3D),
		jsSceneAOs:            make(map[string]*ActiveObj3D),
		jsSceneDOs:            make(map[string]*DangerObj3D),
		jsSceneFOs:            make(map[string]*FieldObj3D),
		jsSceneMovePathArrows: make(map[[2]int]*ColorArrow3D),

		HP: NewColorBar3D("red"),
		SP: NewColorBar3D("yellow"),
		AP: NewColorBar3D("lime"),
	}

	vp.renderer = ThreeJsNew("WebGLRenderer")
	vp.renderer.Get("shadowMap").Set("enabled", true)
	// vp.renderer.Get("shadowMap").Set("type", true)

	rendererDom := vp.renderer.Get("domElement")
	GetElementById("canvasglholder").Call("appendChild", rendererDom)

	vp.camera = ThreeJsNew("PerspectiveCamera", 50, 1, 0.1, HelperSize*2)
	vp.scene = ThreeJsNew("Scene")
	vp.raycaster = ThreeJsNew("Raycaster")

	// no need to add to scene for raycasting
	vp.raycastPlane = NewRaycastPlane()
	vp.scene.Call("add", vp.raycastPlane.Mesh)
	vp.raycastPlane.Mesh.Set("visible", false)

	vp.cursor = NewCursor3D()
	for i := range vp.cursor.Mesh {
		vp.scene.Call("add", vp.cursor.Mesh[i])
	}

	vp.scene.Call("add", vp.HP.Mesh)
	vp.scene.Call("add", vp.SP.Mesh)
	vp.scene.Call("add", vp.AP.Mesh)

	vp.lightW = ThreeJsNew("PointLight", 0xffffff, 0.5)
	vp.lightW.Set("castShadow", true)
	vp.lightW.Get("shadow").Get("camera").Set("near", HelperSize/2)
	vp.lightW.Get("shadow").Get("camera").Set("far", HelperSize*2)
	vp.scene.Call("add", vp.lightW)
	fogco := 0x404040
	// vp.scene.Set("fog", ThreeJsNew("FogExp2", fogco, 0.00025*1.5))
	vp.scene.Set("fog", ThreeJsNew("Fog", fogco, HelperSize*1.1, HelperSize*2))
	vp.scene.Set("background", ThreeJsNew("Color", fogco))
	vp.scene.Call("add", ThreeJsNew("AmbientLight", 0x808080))

	// vp.scene.Call("add", ThreeJsNew("HemisphereLight", 0x000000, fogco))
	// vp.scene.Call("add", ThreeJsNew("DirectionalLight", 0xffffff))

	for i, co := range [3]uint32{0xff0000, 0x00ff00, 0x0000ff} {
		vp.lightRGB[i] = ThreeJsNew("PointLight", co, 0.5)
		vp.lightRGB[i].Set("castShadow", true)
		vp.lightRGB[i].Get("shadow").Get("camera").Set("near", DstCellSize*8/2)
		vp.lightRGB[i].Get("shadow").Get("camera").Set("far", DstCellSize*8*2)
		geo := ThreeJsNew("SphereGeometry", DstCellSize/4, DstCellSize, DstCellSize)
		mat := ThreeJsNew("MeshStandardMaterial",
			map[string]interface{}{
				"emissive":          co,
				"emissiveIntensity": 1.0,
				"color":             co,
			},
		)
		mat.Set("transparent", true)
		mat.Set("opacity", 0.5)
		vp.lightRGB[i].Call("add", ThreeJsNew("Mesh", geo, mat))
		vp.scene.Call("add", vp.lightRGB[i])
	}

	axisSize := HelperSize
	axisHelper := ThreeJsNew("AxesHelper", axisSize*DstCellSize)
	vp.scene.Call("add", axisHelper)

	for i := 0; i < tile.Tile_Count; i++ {
		tlt := tile.Tile(i)

		mat := gTile3D[tlt].Mat
		geo := gTile3D[tlt].Geo
		mesh := ThreeJsNew("InstancedMesh", geo, mat,
			gameconst.ClientViewPortW*gameconst.ClientViewPortH)
		mesh.Set("castShadow", tile.TileShadowAttrib[i].Cast)
		mesh.Set("receiveShadow", tile.TileShadowAttrib[i].Receive)

		mesh.Set("count", 0)
		vp.scene.Call("add", mesh)
		vp.jsTile3DMesh[i] = mesh

		mat = gTile3DDark[tlt].Mat
		geo = gTile3DDark[tlt].Geo
		mesh = ThreeJsNew("InstancedMesh", geo, mat,
			gameconst.ClientViewPortW*gameconst.ClientViewPortH)
		// mesh.Set("castShadow", tile.TileShadowAttrib[i].Cast)
		// mesh.Set("receiveShadow", tile.TileShadowAttrib[i].Receive)
		mesh.Set("count", 0)
		vp.scene.Call("add", mesh)
		vp.jsTile3DDarkMesh[i] = mesh
	}

	return vp
}

func (vp *GameScene) Resize(w, h float64) {
	vp.renderer.Call("setSize", w, h)
	vp.camera.Set("aspect", w/h)
	vp.camera.Call("updateProjectionMatrix")
}

func (vp *GameScene) Zoom(zoom int) {
	vp.camera.Set("zoom", 1.0+float64(zoom)/2)
	vp.camera.Call("updateProjectionMatrix")
}

func (vp *GameScene) FindRayCastingFxFy(jsMouse js.Value) (int, int) {
	// update the picking ray with the camera and mouse position
	vp.raycaster.Call("setFromCamera", jsMouse, vp.camera)

	// calculate objects intersecting the picking ray
	intersects := vp.raycaster.Call(
		"intersectObject", vp.raycastPlane.Mesh)

	for i := 0; i < intersects.Length(); i++ {
		obj := intersects.Index(i)
		pos3 := obj.Get("point")
		x := pos3.Get("x").Float()
		y := pos3.Get("y").Float()
		fx := int(x / DstCellSize)
		fy := int(-y / DstCellSize)
		return fx, fy
		_ = fx
		_ = fy
		// jslog.Infof("pos fx:%v fy:%v", fx, fy)
	}
	return 0, 0
}
