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

package visitarea

import (
	"fmt"
	"sync"
)

type ID2VisitArea struct {
	mutex sync.RWMutex `prettystring:"hide"`

	id2visit map[string]*VisitArea
}

func NewID2VisitArea() *ID2VisitArea {
	i2v := &ID2VisitArea{
		id2visit: make(map[string]*VisitArea),
	}
	return i2v
}

func (i2v *ID2VisitArea) GetList() []*VisitArea {
	i2v.mutex.RLock()
	defer i2v.mutex.RUnlock()

	rtn := make([]*VisitArea, 0, len(i2v.id2visit))
	for _, v := range i2v.id2visit {
		rtn = append(rtn, v)
	}
	return rtn
}

func (i2v *ID2VisitArea) GetIDList() []string {
	i2v.mutex.RLock()
	defer i2v.mutex.RUnlock()

	rtn := make([]string, 0, len(i2v.id2visit))
	for i := range i2v.id2visit {
		rtn = append(rtn, i)
	}
	return rtn
}

func (i2v *ID2VisitArea) GetByID(id string) (*VisitArea, bool) {
	i2v.mutex.RLock()
	defer i2v.mutex.RUnlock()
	r, exist := i2v.id2visit[id]
	return r, exist
}

func (i2v *ID2VisitArea) Add(fi floorI) *VisitArea {
	i2v.mutex.Lock()
	defer i2v.mutex.Unlock()
	va := NewVisitArea(fi)
	i2v.id2visit[fi.GetName()] = va
	return va
}

func (i2v *ID2VisitArea) Del(id string) {
	i2v.mutex.Lock()
	defer i2v.mutex.Unlock()

	delete(i2v.id2visit, id)
}

func (i2v *ID2VisitArea) Forget(uuid string) error {
	i2v.mutex.Lock()
	defer i2v.mutex.Unlock()
	r, exist := i2v.id2visit[uuid]
	if exist {
		r.Forget()
		return nil
	} else {
		return fmt.Errorf("not found %v", uuid)
	}
}
