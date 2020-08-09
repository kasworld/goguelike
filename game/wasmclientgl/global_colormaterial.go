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
	"sync"
	"syscall/js"
)

var gPoolColorMaterial = NewPoolColorMaterial()

type PoolColorMaterial struct {
	mutex    sync.Mutex
	poolData map[string][]js.Value
	newCount int
	getCount int
	putCount int
}

func NewPoolColorMaterial() *PoolColorMaterial {
	return &PoolColorMaterial{
		poolData: make(map[string][]js.Value),
	}
}

func (p *PoolColorMaterial) String() string {
	return fmt.Sprintf("PoolColorMaterial[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *PoolColorMaterial) Get(color string) js.Value {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn js.Value
	if l := len(p.poolData[color]); l > 0 {
		rtn = p.poolData[color][l-1]
		p.poolData[color] = p.poolData[color][:l-1]
		p.getCount++
	} else {
		rtn = newColorMaterial(color)
		p.newCount++
	}
	return rtn
}

func (p *PoolColorMaterial) Put(pb js.Value) {
	color := pb.Get("color").String()
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData[color] = append(p.poolData[color], pb)
	p.putCount++
}

func newColorMaterial(co string) js.Value {
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"color": co,
		},
	)
	// mat.Set("transparent", true)
	// mat.Set("opacity", 1)
	return mat
}
