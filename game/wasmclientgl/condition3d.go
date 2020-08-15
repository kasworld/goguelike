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
	"fmt"
	"math"
	"sync"
	"syscall/js"

	"github.com/kasworld/goguelike/enum/condition"
)

var gPoolCondition3D = NewPoolCondition3D()

type PoolCondition3D struct {
	mutex    sync.Mutex
	poolData map[condition.Condition][]*Condition3D
	newCount int
	getCount int
	putCount int
}

func NewPoolCondition3D() *PoolCondition3D {
	return &PoolCondition3D{
		poolData: make(map[condition.Condition][]*Condition3D),
	}
}

func (p *PoolCondition3D) String() string {
	return fmt.Sprintf("PoolCondition3D[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolCondition3D) Get(cn condition.Condition) *Condition3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *Condition3D
	if l := len(p.poolData[cn]); l > 0 {
		rtn = p.poolData[cn][l-1]
		p.poolData[cn] = p.poolData[cn][:l-1]
		p.getCount++
	} else {
		rtn = NewCondition3D(cn)
		p.newCount++
	}
	return rtn
}

func (p *PoolCondition3D) Put(pb *Condition3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[pb.Condition] = append(p.poolData[pb.Condition], pb)
	p.putCount++
}

type Condition3D struct {
	Condition condition.Condition
	Tex       js.Value
	GeoInfo   GeoInfo
	Mesh      js.Value
}

func NewCondition3D(cn condition.Condition) *Condition3D {
	colorstr := cn.Color().ToHTMLColorString()
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"color": colorstr,
		},
	)
	mat.Set("transparent", true)
	mat.Set("opacity", 0.9)

	geo := ThreeJsNew("CylinderGeometry", 1, 1, DstCellSize)
	geo.Call("rotateZ", math.Pi/2)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &Condition3D{
		Condition: cn,
		GeoInfo:   GetGeoInfo(geo),
		Mesh:      mesh,
	}
}
