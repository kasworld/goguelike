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

	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/gowasmlib/jslog"
)

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
