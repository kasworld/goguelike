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

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type Viewport struct {
	CanvasGL js.Value
	renderer js.Value
	jsMouse  js.Value

	// from client floor gl

	light     [3]js.Value // rgb light
	lightW    js.Value    // white light
	scene     js.Value
	camera    js.Value
	raycaster js.Value

	raycastPlane *RaycastPlane
	cursor       *Cursor3D

	jsSceneCOs map[string]*CarryObj3D  // in sight only  carryobj
	jsSceneAOs map[string]*ActiveObj3D // in sight only ao
	jsSceneFOs map[string]*FieldObj3D  // in clientview fieldobj

	// tile 3d instancedmesh
	// count = gameconst.ClientViewPortW * gameconst.ClientViewPortH
	jsTile3DMesh      [tile.Tile_Count]js.Value
	jsTile3DCount     [tile.Tile_Count]int // in use count
	jsTile3DDarkMesh  [tile.Tile_Count]js.Value
	jsTile3DDarkCount [tile.Tile_Count]int // in use count

}

func NewViewport() *Viewport {
	vp := &Viewport{
		jsSceneCOs: make(map[string]*CarryObj3D),
		jsSceneAOs: make(map[string]*ActiveObj3D),
		jsSceneFOs: make(map[string]*FieldObj3D),
	}

	vp.jsMouse = ThreeJsNew("Vector2")
	vp.renderer = ThreeJsNew("WebGLRenderer")
	vp.CanvasGL = vp.renderer.Get("domElement")
	GetElementById("canvasglholder").Call("appendChild", vp.CanvasGL)
	vp.CanvasGL.Set("tabindex", "1")

	vp.camera = ThreeJsNew("PerspectiveCamera", 50, 1, 0.1, HelperSize*2)
	vp.scene = ThreeJsNew("Scene")
	vp.raycaster = ThreeJsNew("Raycaster")

	// no need to add to scene for raycasting
	vp.raycastPlane = NewRaycastPlane()
	vp.scene.Call("add", vp.raycastPlane.Mesh)
	vp.raycastPlane.Mesh.Set("visible", false)

	vp.cursor = NewCursor3D()
	vp.cursor.ChangeTile(gClientTile.CursorTiles[0])
	vp.scene.Call("add", vp.cursor.Mesh)

	lightAm := ThreeJsNew("AmbientLight", 0x808080)
	vp.scene.Call("add", lightAm)

	vp.lightW = ThreeJsNew("PointLight", 0xffffff, 0.5)
	vp.scene.Call("add", vp.lightW)
	for i, co := range [3]uint32{0xff0000, 0x00ff00, 0x0000ff} {
		vp.light[i] = ThreeJsNew("PointLight", co, 0.5)
		SetPosition(vp.light[i],
			HelperSize/2, HelperSize/2, HelperSize/2,
		)
		vp.scene.Call("add", vp.light[i])
		lightHelper := ThreeJsNew("PointLightHelper", vp.light[i], 2)
		vp.scene.Call("add", lightHelper)
	}

	axisSize := HelperSize
	axisHelper := ThreeJsNew("AxesHelper", axisSize*DstCellSize)
	vp.scene.Call("add", axisHelper)

	SetPosition(vp.camera,
		HelperSize/2, HelperSize/2, HelperSize,
	)
	vp.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			HelperSize/2, HelperSize/2, 0,
		),
	)
	vp.camera.Call("updateProjectionMatrix")

	for i := 0; i < tile.Tile_Count; i++ {
		tlt := tile.Tile(i)

		mat := gTile3D[tlt].Mat
		geo := gTile3D[tlt].Geo
		mesh := ThreeJsNew("InstancedMesh", geo, mat,
			gameconst.ClientViewPortW*gameconst.ClientViewPortH)
		mesh.Set("count", 0)
		vp.scene.Call("add", mesh)
		vp.jsTile3DMesh[i] = mesh

		mat = gTile3DDark[tlt].Mat
		geo = gTile3DDark[tlt].Geo
		mesh = ThreeJsNew("InstancedMesh", geo, mat,
			gameconst.ClientViewPortW*gameconst.ClientViewPortH)
		mesh.Set("count", 0)
		vp.scene.Call("add", mesh)
		vp.jsTile3DDarkMesh[i] = mesh
	}

	return vp
}

func (vp *Viewport) Hide() {
	vp.CanvasGL.Get("style").Set("display", "none")
}
func (vp *Viewport) Show() {
	vp.CanvasGL.Get("style").Set("display", "initial")
}

func (vp *Viewport) Resize(w, h float64) {
	vp.CanvasGL.Call("setAttribute", "width", w)
	vp.CanvasGL.Call("setAttribute", "height", h)
	vp.renderer.Call("setSize", w, h)

	vp.camera.Set("aspect", w/h)
	vp.camera.Call("updateProjectionMatrix")

}

func (vp *Viewport) Zoom(zoom int) {
	vp.camera.Set("zoom", 1.0+float64(zoom)/2)
	vp.camera.Call("updateProjectionMatrix")
}

func (vp *Viewport) UpdateFromViewportTile(
	cf *clientfloor.ClientFloor,
	taNoti *c2t_obj.NotiVPTiles_data,
	olNoti *c2t_obj.NotiObjectList_data) error {

	if cf.FloorInfo.UUID != taNoti.FloorUUID {
		return fmt.Errorf("vptile data floor not match %v %v",
			cf.FloorInfo.UUID, taNoti.FloorUUID)

	}
	vp.makeClientTileView(cf, taNoti)
	vp.updateFieldObjInView(cf, taNoti.VPX, taNoti.VPY)
	vp.raycastPlane.MoveCenterTo(taNoti.VPX, taNoti.VPY)
	return nil
}

func (vp *Viewport) FindRayCastingFxFy() (int, int) {
	// update the picking ray with the camera and mouse position
	vp.raycaster.Call("setFromCamera", vp.jsMouse, vp.camera)

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
