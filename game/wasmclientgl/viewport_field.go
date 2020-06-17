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

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

// for tile texture wrap
type textureTileWrapInfo struct {
	Xcount int
	Ycount int
	WrapX  func(int) int
	WrapY  func(int) int
}

func (ttwi textureTileWrapInfo) CalcSrc(fx, fy int, shiftx, shifty float64) (int, int) {
	tx := fx % ttwi.Xcount
	ty := fy % ttwi.Ycount
	srcx := ttwi.WrapX(tx*CellSize + int(shiftx))
	srcy := ttwi.WrapY(ty*CellSize + int(shifty))
	return srcx, srcy
}

type ClientField struct {
	W   int // canvas width
	H   int // canvas height
	Cnv js.Value
	Ctx js.Value

	Tex  js.Value
	Mat  js.Value
	Geo  js.Value
	Mesh js.Value
}

func (vp *Viewport) NewClientField(fi *c2t_obj.FloorInfo) *ClientField {
	w := fi.W * CellSize
	h := fi.H * CellSize
	xRepeat := 3
	yRepeat := 3
	clFd := &ClientField{
		W: w,
		H: h,
	}
	clFd.Cnv = js.Global().Get("document").Call("createElement",
		"CANVAS")
	clFd.Ctx = clFd.Cnv.Call("getContext", "2d")
	if !clFd.Ctx.Truthy() {
		jslog.Errorf("fail to get context %v", fi)
		return nil
	}
	clFd.Ctx.Set("imageSmoothingEnabled", false)
	clFd.Cnv.Set("width", w)
	clFd.Cnv.Set("height", h)

	clFd.Tex = vp.ThreeJsNew("CanvasTexture", clFd.Cnv)
	clFd.Tex.Set("wrapS", vp.threejs.Get("RepeatWrapping"))
	clFd.Tex.Set("wrapT", vp.threejs.Get("RepeatWrapping"))
	clFd.Tex.Get("repeat").Set("x", xRepeat)
	clFd.Tex.Get("repeat").Set("y", yRepeat)
	clFd.Mat = vp.ThreeJsNew("MeshBasicMaterial",
		map[string]interface{}{
			"map": clFd.Tex,
		},
	)
	clFd.Geo = vp.ThreeJsNew("PlaneBufferGeometry",
		w*xRepeat, h*yRepeat)
	clFd.Mesh = vp.ThreeJsNew("Mesh", clFd.Geo, clFd.Mat)

	SetPosition(clFd.Mesh, w/2, -h/2, -10)

	clFd.Ctx.Set("font", fmt.Sprintf("%dpx sans-serif", CellSize))
	clFd.Ctx.Set("fillStyle", "gray")
	clFd.Ctx.Call("fillText", fi.Name, 100, 100)
	// clFd.Ctx.Call("fillRect", 0, 0, 10, 100)

	// clFd.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
	// 	0, 0, CellSize*16, CellSize*16,
	// 	0, 0, CellSize*16, CellSize*16)

	return clFd
}

func (vp *Viewport) ReplaceFloorTiles(cf *clientfloor.ClientFloor) {
	clFd := vp.floorG2ID2ClientField[cf.FloorInfo.G2ID]
	// jslog.Infof("UpdateClientField %v %v", cf, clFd)
	drawCtx := clFd.Ctx
	for fx, xv := range cf.Tiles {
		for fy, yv := range xv {
			vp.drawTileAt(cf, fx, fy, yv, drawCtx)
		}
	}
}

func (vp *Viewport) drawTileAt(
	cf *clientfloor.ClientFloor, fx, fy int, tl tile_flag.TileFlag, drawCtx js.Value) {
	dstX := fx * CellSize
	dstY := fy * CellSize
	diffbase := fx*5 + fy*3

	for i := 0; i < tile.Tile_Count; i++ {
		shX := 0.0
		shY := 0.0
		tlt := tile.Tile(i)
		if tl.TestByTile(tlt) {
			if vp.textureTileList[i] != nil {
				// texture tile
				tlic := vp.textureTileList[i]
				srcx, srcy := vp.textureTileWrapInfoList[i].CalcSrc(fx, fy, shX, shY)
				// wrap := textureTileWrapInfo[i]
				// tx := fx % wrap.Xcount
				// ty := fy % wrap.Ycount
				// srcx := wrap.WrapX(tx*CellSize + int(shX))
				// srcy := wrap.WrapY(ty*CellSize + int(shY))
				drawCtx.Call("drawImage", tlic.Cnv,
					srcx, srcy, CellSize, CellSize,
					dstX, dstY, CellSize, CellSize)

			} else if tlt == tile.Wall {
				// wall tile process
				tlList := gClientTile.FloorTiles[i]
				tilediff := vp.calcWallTileDiff(cf, fx, fy)
				ti := tlList[tilediff%len(tlList)]
				drawCtx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, CellSize, CellSize)
			} else if tlt == tile.Window {
				// window tile process
				tlList := gClientTile.FloorTiles[i]
				tlindex := 0
				if vp.checkWallAt(cf, fx, fy-1) && vp.checkWallAt(cf, fx, fy+1) { // n-s window
					tlindex = 1
				}
				ti := tlList[tlindex]
				drawCtx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, CellSize, CellSize)
			} else {
				// bitmap tile
				tlList := gClientTile.FloorTiles[i]
				tilediff := diffbase + int(shX)
				ti := tlList[tilediff%len(tlList)]
				drawCtx.Call("drawImage", gClientTile.TilePNG.Cnv,
					ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
					dstX, dstY, CellSize, CellSize)
			}
		}
	}
}

func (vp *Viewport) UpdateClientField(
	cf *clientfloor.ClientFloor,
	vpData *c2t_obj.NotiVPTiles_data,
	ViewportXYLenList findnear.XYLenList,
) {
	clFd := vp.floorG2ID2ClientField[cf.FloorInfo.G2ID]
	// jslog.Infof("UpdateClientField %v %v", cf, clFd)
	drawCtx := clFd.Ctx
	for i, v := range ViewportXYLenList {
		fx := cf.XWrapSafe(v.X + vpData.VPX)
		fy := cf.YWrapSafe(v.Y + vpData.VPY)
		tl := vpData.VPTiles[i]
		vp.drawTileAt(cf, fx, fy, tl, drawCtx)
	}

	// move camera, light
	cameraX := vpData.VPX * CellSize
	cameraY := -vpData.VPY * CellSize
	SetPosition(vp.light,
		cameraX, cameraY, CellSize*2,
	)
	SetPosition(vp.camera,
		cameraX, cameraY, HelperSize,
	)
	vp.camera.Call("lookAt",
		vp.ThreeJsNew("Vector3",
			cameraX, cameraY, 0,
		),
	)

	clFd.Tex.Set("needsUpdate", true)
}

func (vp *Viewport) calcWallTileDiff(cf *clientfloor.ClientFloor, flx, fly int) int {
	rtn := 0
	if vp.checkWallAt(cf, flx, fly-1) {
		rtn |= 1
	}
	if vp.checkWallAt(cf, flx+1, fly) {
		rtn |= 1 << 1
	}
	if vp.checkWallAt(cf, flx, fly+1) {
		rtn |= 1 << 2
	}
	if vp.checkWallAt(cf, flx-1, fly) {
		rtn |= 1 << 3
	}
	return rtn
}

func (vp *Viewport) checkWallAt(cf *clientfloor.ClientFloor, flx, fly int) bool {
	flx = cf.XWrapSafe(flx)
	fly = cf.YWrapSafe(fly)
	tl := cf.Tiles[flx][fly]
	return tl.TestByTile(tile.Wall) ||
		tl.TestByTile(tile.Door) ||
		tl.TestByTile(tile.Window)
}

func (vp *Viewport) ChangeToClientField(cf *clientfloor.ClientFloor) {
	clFd := vp.floorG2ID2ClientField[cf.FloorInfo.G2ID]
	for _, v := range vp.floorG2ID2ClientField {
		vp.scene.Call("remove", v.Mesh)
	}
	vp.scene.Call("add", clFd.Mesh)
}

// var groundTexture = loader.load( 'textures/terrain/grasslight-big.jpg' );
// groundTexture.wrapS = groundTexture.wrapT = THREE.RepeatWrapping;
// groundTexture.repeat.set( 25, 25 );
// groundTexture.anisotropy = 16;
// groundTexture.encoding = THREE.sRGBEncoding;
// var groundMaterial = new THREE.MeshLambertMaterial( { map: groundTexture } );
// var mesh = new THREE.Mesh( new THREE.PlaneBufferGeometry( 20000, 20000 ), groundMaterial );
// mesh.position.y = - 250;
// mesh.rotation.x = - Math.PI / 2;
// mesh.receiveShadow = true;
// scene.add( mesh );

// func (vp *Viewport) initField() {
// 	clFd := vp.getFieldMesh(
// 		g2id.New(), 256, 256,
// 	)
// 	SetPosition(clFd.Mesh, HelperSize/2, HelperSize/2, -10)
// 	// clFd.Ctx.Call("clearRect", 0, 0, w, h)
// 	clFd.Ctx.Set("fillStyle", "gray")
// 	clFd.Ctx.Call("fillRect", 0, 0, 10, 100)

// 	vp.scene.Call("add", clFd.Mesh)
// }
