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

	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/lib/g2id"
)

func (vp *Viewport) getColorMaterial(co uint32) js.Value {
	mat, exist := vp.colorMaterialCache[co]
	if !exist {
		mat = vp.ThreeJsNew("MeshPhongMaterial",
			map[string]interface{}{
				"color": co,
				// "flatShading": true,
			},
		)
		vp.colorMaterialCache[co] = mat
	}
	return mat
}

func (vp *Viewport) getAOGeometry(gotype factiontype.FactionType) js.Value {
	geo, exist := vp.aoGeometryCache[gotype]
	if !exist {
		radius := 32
		geo = vp.ThreeJsNew("BoxGeometry", radius*2, radius*2, radius*2)
		vp.aoGeometryCache[gotype] = geo
	}
	return geo
}

func (vp *Viewport) getFieldMesh(
	floorG2ID g2id.G2ID, w, h int) js.Value {
	mesh, exist := vp.fieldMeshCache[floorG2ID]
	if !exist {
		flTex := vp.textureLoader.Call("load", "/tiles/Grass.png")
		flTex.Set("wrapS", vp.threejs.Get("RepeatWrapping"))
		flTex.Set("wrapT", vp.threejs.Get("RepeatWrapping"))
		flTex.Get("repeat").Set("x", 5)
		flTex.Get("repeat").Set("y", 5)
		bgMaterial := vp.ThreeJsNew("MeshBasicMaterial",
			map[string]interface{}{
				"map": flTex,
			},
		)
		bgGeo := vp.ThreeJsNew("PlaneBufferGeometry",
			StageSize*5, StageSize*5)
		mesh = vp.ThreeJsNew("Mesh", bgGeo, bgMaterial)
		vp.fieldMeshCache[floorG2ID] = mesh
		// jslog.Info(vp.fieldMeshCache[floorG2ID])
		SetPosition(mesh,
			StageSize/2, StageSize/2, -10,
		)
	}
	return mesh
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
