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

package dangerobject

import (
	"fmt"
)

type DOPool struct {
	poolData []*DangerObject
	newCount int
	getCount int
	putCount int
}

func NewDOPool() *DOPool {
	return &DOPool{
		poolData: make([]*DangerObject, 0),
	}
}

func (p *DOPool) String() string {
	return fmt.Sprintf("DOPool[%v new:%v get:%v put:%v]",
		len(p.poolData), p.newCount, p.getCount, p.putCount,
	)
}

func (p *DOPool) Get() *DangerObject {
	var rtn *DangerObject
	if l := len(p.poolData); l > 0 {
		rtn = p.poolData[l-1]
		p.poolData = p.poolData[:l-1]
		p.getCount++
	} else {
		rtn = nil
		p.newCount++
	}
	return rtn
}

func (p *DOPool) Put(pb *DangerObject) {
	p.poolData = append(p.poolData, pb)
	p.putCount++
}
