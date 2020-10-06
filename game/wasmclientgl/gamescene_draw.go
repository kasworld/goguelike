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
	"math"
	"time"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

// animate active obj
func (vp *GameScene) animateActiveObj(
	cf *clientfloor.ClientFloor, ao3d *ActiveObj3D,
	fx, fy int, dir way9type.Way9Type, frameProgress float64) {

	if dir != way9type.Center {
		dx, dy := dir.DxDy()
		tl := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
		shX := DstCellSize * frameProgress * float64(dx)
		shY := DstCellSize * frameProgress * float64(dy)
		shZ := CalcTile3DStepOn(tl)
		ao3d.SetFieldPosition(fx, fy, shX, shY, shZ)
	}
}

// fieldobj animate
// common to playview, floorview
func (vp *GameScene) animateFieldObj() {
	rad := time.Now().Sub(gInitData.TowerInfo.StartTime).Seconds()
	for _, fo := range vp.jsSceneFOs {
		fo.RotateZ(rad)
	}
}

func (vp *GameScene) animateDangerObj(frameProgress float64) {
	for _, dao3d := range vp.jsSceneDOs {
		rr := dao3d.Dao.AffectRate
		if rr < 0.2 {
			rr = 0.2
		}
		dao3d.ScaleX((1 - frameProgress) * rr)
		dao3d.ScaleY((1 - frameProgress) * rr)
	}
}

func (vp *GameScene) animateDoor() {
	rad := time.Now().Sub(gInitData.TowerInfo.StartTime).Seconds()

	meshs := vp.jsTile3DDarkMesh[tile.Door]
	for i := 0; i < vp.jsTile3DDarkCount[tile.Door]; i++ {
		oldmatrix := ThreeJsNew("Matrix4")
		meshs.Call("getMatrixAt", i, oldmatrix)
		v3pos := ThreeJsNew("Vector3")
		v3pos.Call("setFromMatrixPosition", oldmatrix)
		newmatrix := ThreeJsNew("Matrix4")
		newmatrix.Call("makeRotationZ", rad)
		newmatrix.Call("setPosition", v3pos)
		meshs.Call("setMatrixAt", i, newmatrix)
	}
	meshs.Get("instanceMatrix").Set("needsUpdate", true)

	meshs = vp.jsTile3DMesh[tile.Door]
	for i := 0; i < vp.jsTile3DCount[tile.Door]; i++ {
		oldmatrix := ThreeJsNew("Matrix4")
		meshs.Call("getMatrixAt", i, oldmatrix)
		v3pos := ThreeJsNew("Vector3")
		v3pos.Call("setFromMatrixPosition", oldmatrix)
		newmatrix := ThreeJsNew("Matrix4")
		newmatrix.Call("makeRotationZ", rad)
		newmatrix.Call("setPosition", v3pos)
		meshs.Call("setMatrixAt", i, newmatrix)
	}
	meshs.Get("instanceMatrix").Set("needsUpdate", true)
}

// tile scroll list of animate
// common to playview, floorview
func (vp *GameScene) animateTile(envBias bias.Bias) {
	for _, i := range []tile.Tile{
		tile.Swamp,
		tile.Soil,
		tile.Stone,
		tile.Sand,
		tile.Sea,
		tile.Magma,
		tile.Ice,
		tile.Grass,
		tile.Tree,
		// tile.Road,
		// tile.Room,
		// tile.Wall,
		// tile.Window,
		// tile.Door,
		tile.Fog,
		tile.Smoke,
	} {
		tilc := tile.TileScrollAttrib[i]
		shX := int(envBias[tilc.SrcXBiasAxis] * tilc.AniXSpeed)
		shY := int(envBias[tilc.SrcYBiasAxis] * tilc.AniYSpeed)
		gTile3D[i].DrawTexture(shX, shY)
		gTile3DDark[i].DrawTexture(shX, shY)
	}
}

// common to playview, floorview
func (vp *GameScene) moveCameraLight(
	cf *clientfloor.ClientFloor, vpx, vpy int, frameProgress float64,
	scrollDir way9type.Way9Type, envBias bias.Bias) {

	sx, sy := CalcShiftDxDy(frameProgress)
	scrollDx := -scrollDir.Dx() * sx
	scrollDy := scrollDir.Dy() * sy

	// move camera, light
	cameraX := float64(vpx*DstCellSize + scrollDx)
	cameraY := float64(-vpy*DstCellSize + scrollDy)
	cameraR := float64(HelperSize)

	angle := gameOptions.GetByIDBase("Angle").State
	zoom := gameOptions.GetByIDBase("Zoom").State
	anglezoom := [3][3]float64{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	carmeraDy := anglezoom[angle][zoom] * DstCellSize

	cameraRad := []float64{
		math.Pi / 2,
		math.Pi / 3,
		math.Pi / 4,
	}[gameOptions.GetByIDBase("Angle").State]

	envBias = envBias.MakeAbsSumTo(1)
	cx := float64(cf.XWrapper.GetWidth()) * DstCellSize / 2
	cy := float64(cf.YWrapper.GetWidth()) * DstCellSize / 2
	r := cx
	if cy < cx {
		r = cy
	}
	r *= 0.7
	for i := range vp.lightRGB {
		rad := envBias[i] * 2 * math.Pi
		x := cx + r*math.Sin(rad)
		y := cy + r*math.Cos(rad)
		SetPosition(vp.lightRGB[i],
			x, -y, DstCellSize*8,
		)
	}
	SetPosition(vp.lightW,
		cameraX, cameraY, cameraR,
	)

	SetPosition(vp.camera,
		cameraX,
		cameraY-cameraR*math.Cos(cameraRad)-carmeraDy,
		cameraR*math.Sin(cameraRad),
	)
	vp.camera.Call("lookAt",
		ThreeJsNew("Vector3",
			cameraX,
			cameraY-carmeraDy,
			0,
		),
	)
}

// add fo to clientview by vp.FieldObjPosMan
func (vp *GameScene) updateFieldObjInView(
	cf *clientfloor.ClientFloor, vpx, vpy int) {
	addFOuuid := make(map[string]bool)
	for _, v := range gXYLenListView {
		fx := v.X + vpx
		fy := v.Y + vpy
		obj := cf.FieldObjPosMan.Get1stObjAt(fx, fy) // wrap in fn
		if obj == nil {
			continue
		}
		fo := obj.(*c2t_obj.FieldObjClient)
		fo3d, exist := vp.jsSceneFOs[obj.GetUUID()]
		if !exist {
			// add new fieldobj
			fo3d = gPoolFieldObj3D.Get(fo.ActType, fo.DisplayType)
			vp.jsSceneFOs[obj.GetUUID()] = fo3d
			vp.scene.Call("add", fo3d.Mesh)

			fo3d.Label = gPoolLabel3D.Get(fo.Message)
			vp.scene.Call("add", fo3d.Label.Mesh)
		}
		if !addFOuuid[obj.GetUUID()] {
			tl := cf.Tiles[cf.XWrapSafe(fx)][cf.YWrapSafe(fy)]
			shZ := CalcTile3DStepOn(tl)
			fo3d.SetFieldPosition(fx, fy, shZ)
			addFOuuid[obj.GetUUID()] = true
			fo3d.Label.SetFieldPosition(fx, fy,
				0, 0, fo3d.GeoInfo.Len[2]+2+shZ)
		} else {
			// same fo in gXYLenListView
			// field too small
			// skip farther fo
		}
	}

	// del removed obj
	for id, fo3d := range vp.jsSceneFOs {
		if !addFOuuid[id] {
			vp.scene.Call("remove", fo3d.Mesh)
			gPoolFieldObj3D.Put(fo3d)
			delete(vp.jsSceneFOs, id)

			vp.scene.Call("remove", fo3d.Label.Mesh)
			gPoolLabel3D.Put(fo3d.Label)
			fo3d.Label = nil
		}
	}
}

func (vp *GameScene) ClearMovePath() {
	for pos, ar3d := range vp.jsSceneMovePathArrows {
		vp.scene.Call("remove", ar3d.Mesh)
		gPoolColorArrow3D.Put(ar3d)
		delete(vp.jsSceneMovePathArrows, pos)
	}
}

func (vp *GameScene) makeMovePathInView(
	cf *clientfloor.ClientFloor,
	vpx, vpy int,
	path2dst [][2]int) {

	addAr3Duuid := make(map[[2]int]bool)
	if len(path2dst) > 0 {
		w, h := cf.XWrapper.GetWidth(), cf.YWrapper.GetWidth()
		for i, pos := range path2dst[:len(path2dst)-1] {
			ar3d, exist := vp.jsSceneMovePathArrows[pos]
			if !exist {
				ar3d = gPoolColorArrow3D.Get("#808080")
				vp.jsSceneMovePathArrows[pos] = ar3d
				vp.scene.Call("add", ar3d.Mesh)
			}
			addAr3Duuid[pos] = true

			// change tile for dir
			pos2 := path2dst[i+1]
			dx, dy := way9type.CalcDxDyWrapped(
				pos2[0]-pos[0], pos2[1]-pos[1],
				w, h,
			)
			diri := way9type.RemoteDxDy2Way9(dx, dy)
			ar3d.SetDir(diri)
			x, y := CalcAroundPos(w, h, vpx, vpy, pos[0], pos[1])
			tl := cf.Tiles[cf.XWrapSafe(x)][cf.YWrapSafe(y)]
			shZ := CalcTile3DStepOn(tl)
			ar3d.SetFieldPosition(x, y, 0, 0, shZ)
		}
		// add last ?
	}

	for pos, ar3d := range vp.jsSceneMovePathArrows {
		if !addAr3Duuid[pos] {
			vp.scene.Call("remove", ar3d.Mesh)
			gPoolColorArrow3D.Put(ar3d)
			delete(vp.jsSceneMovePathArrows, pos)
		}
	}
}
