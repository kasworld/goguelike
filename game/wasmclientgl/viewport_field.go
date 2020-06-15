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

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
)

var wrapInfo = [tile.Tile_Count]struct {
	Xcount int
	Ycount int
	WrapX  func(int) int
	WrapY  func(int) int
}{}

type ClientField struct {
	W   int
	H   int
	Cnv js.Value
	Ctx js.Value

	Tex  js.Value
	Mat  js.Value
	Geo  js.Value
	Mesh js.Value
}

func (vp *Viewport) drawFloorTiles(
	drawCtx js.Value,
	cf *clientfloor.ClientFloor,
	vpData *c2t_obj.NotiVPTiles_data,
	ViewportXYLenList findnear.XYLenList,
) {
	for i, v := range ViewportXYLenList {
		fx := cf.XWrapSafe(v.X + vpData.VPX)
		fy := cf.YWrapSafe(v.Y + vpData.VPY)
		tl := vpData.VPTiles[i]
		dstX := fx * CellSize
		dstY := fy * CellSize
		diffbase := fx*5 + fy*3
		for i := 0; i < tile.Tile_Count; i++ {
			shX := 0.0
			shY := 0.0
			tlt := tile.Tile(i)
			if tl.TestByTile(tlt) {
				if vp.TileImgCnvList[i] != nil {
					// texture tile
					tlic := vp.TileImgCnvList[i]
					wrap := wrapInfo[i]
					tx := fx % wrap.Xcount
					ty := fy % wrap.Ycount
					srcx := wrap.WrapX(tx*CellSize + int(shX))
					srcy := wrap.WrapY(ty*CellSize + int(shY))
					drawCtx.Call("drawImage", tlic.Cnv,
						srcx, srcy, CellSize, CellSize,
						dstX, dstY, CellSize, CellSize)

				} else if tlt == tile.Wall {
					// wall tile process
					tlList := vp.clientTile.FloorTiles[i]
					tilediff := vp.calcWallTileDiff(cf, fx, fy)
					ti := tlList[tilediff%len(tlList)]
					drawCtx.Call("drawImage", vp.clientTile.TilePNG.Cnv,
						ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
						dstX, dstY, CellSize, CellSize)
				} else if tlt == tile.Window {
					// window tile process
					tlList := vp.clientTile.FloorTiles[i]
					tlindex := 0
					if vp.checkWallAt(cf, fx, fy-1) && vp.checkWallAt(cf, fx, fy+1) { // n-s window
						tlindex = 1
					}
					ti := tlList[tlindex]
					drawCtx.Call("drawImage", vp.clientTile.TilePNG.Cnv,
						ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
						dstX, dstY, CellSize, CellSize)
				} else {
					// bitmap tile
					tlList := vp.clientTile.FloorTiles[i]
					tilediff := diffbase + int(shX)
					ti := tlList[tilediff%len(tlList)]
					drawCtx.Call("drawImage", vp.clientTile.TilePNG.Cnv,
						ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
						dstX, dstY, CellSize, CellSize)
				}
			}
		}
	}
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

func (vp *Viewport) getFieldMesh(
	floorG2ID g2id.G2ID, w, h int) *ClientField {
	clFd, exist := vp.floorG2ID2ClientField[floorG2ID]
	if !exist {
		clFd = &ClientField{
			W: w,
			H: h,
		}
		clFd.Cnv = js.Global().Get("document").Call("createElement",
			"CANVAS")
		clFd.Ctx = clFd.Cnv.Call("getContext", "2d")
		if !clFd.Ctx.Truthy() {
			jslog.Errorf("fail to get context %v", floorG2ID)
			return nil
		}
		clFd.Ctx.Set("imageSmoothingEnabled", false)
		clFd.Cnv.Set("width", w)
		clFd.Cnv.Set("height", h)

		clFd.Tex = vp.ThreeJsNew("CanvasTexture", clFd.Cnv)
		clFd.Tex.Set("wrapS", vp.threejs.Get("RepeatWrapping"))
		clFd.Tex.Set("wrapT", vp.threejs.Get("RepeatWrapping"))
		clFd.Tex.Get("repeat").Set("x", 5)
		clFd.Tex.Get("repeat").Set("y", 5)
		clFd.Mat = vp.ThreeJsNew("MeshBasicMaterial",
			map[string]interface{}{
				"map": clFd.Tex,
			},
		)
		clFd.Geo = vp.ThreeJsNew("PlaneBufferGeometry",
			StageSize*5, StageSize*5)
		clFd.Mesh = vp.ThreeJsNew("Mesh", clFd.Geo, clFd.Mat)
		vp.floorG2ID2ClientField[floorG2ID] = clFd
		// jslog.Info(vp.floorG2ID2ClientField[floorG2ID])
	}
	return clFd
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

func (vp *Viewport) initField() {
	clFd := vp.getFieldMesh(
		g2id.New(), 256, 256,
	)
	SetPosition(clFd.Mesh, StageSize/2, StageSize/2, -10)
	// clFd.Ctx.Call("clearRect", 0, 0, w, h)
	clFd.Ctx.Set("fillStyle", "gray")
	clFd.Ctx.Call("fillRect", 0, 0, 10, 100)

	vp.scene.Call("add", clFd.Mesh)
}
