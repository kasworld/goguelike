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

package aoid2activeobject

import (
	"fmt"
	"sync"

	"github.com/kasworld/goguelike/game/gamei"
)

func (aoman *ActiveObjID2ActiveObject) String() string {
	return fmt.Sprintf("%s[%v]", aoman.name, aoman.Count())
}

type ActiveObjID2ActiveObject struct {
	mutex     sync.RWMutex `prettystring:"hide"`
	name      string
	aouuid2ao map[string]gamei.ActiveObjectI
}

func New(name string) *ActiveObjID2ActiveObject {
	if name == "" {
		name = "ActiveObjID2ActiveObject"
	}
	return &ActiveObjID2ActiveObject{
		name:      name,
		aouuid2ao: make(map[string]gamei.ActiveObjectI),
	}
}

func (aoman *ActiveObjID2ActiveObject) Cleanup() {
	aoman.mutex.Lock()
	defer aoman.mutex.Unlock()
	for _, ao := range aoman.aouuid2ao {
		ao.Cleanup()
	}
}

func (aoman *ActiveObjID2ActiveObject) Count() int {
	return len(aoman.aouuid2ao)
}

func (aoman *ActiveObjID2ActiveObject) Add(ao gamei.ActiveObjectI) error {
	aoman.mutex.Lock()
	defer aoman.mutex.Unlock()
	if _, exist := aoman.aouuid2ao[ao.GetUUID()]; exist {
		return fmt.Errorf("ao exist %v", ao)
	}
	aoman.aouuid2ao[ao.GetUUID()] = ao
	return nil
}

func (aoman *ActiveObjID2ActiveObject) DelByUUID(aouuid string) (gamei.ActiveObjectI, error) {
	aoman.mutex.Lock()
	defer aoman.mutex.Unlock()
	ao, exist := aoman.aouuid2ao[aouuid]
	if !exist {
		return nil, fmt.Errorf("aouuid not found %v", aouuid)
	}
	delete(aoman.aouuid2ao, aouuid)
	return ao, nil
}

func (aoman *ActiveObjID2ActiveObject) GetByUUID(aouuid string) (gamei.ActiveObjectI, bool) {
	aoman.mutex.RLock()
	defer aoman.mutex.RUnlock()
	if ao, exist := aoman.aouuid2ao[aouuid]; !exist {
		return nil, false
	} else {
		return ao, true
	}
}

func (aoman *ActiveObjID2ActiveObject) GetAllList() []gamei.ActiveObjectI {
	aoman.mutex.RLock()
	defer aoman.mutex.RUnlock()
	rtn := make([]gamei.ActiveObjectI, len(aoman.aouuid2ao))
	i := 0
	for _, ao := range aoman.aouuid2ao {
		rtn[i] = ao
		i++
	}
	return rtn
}
