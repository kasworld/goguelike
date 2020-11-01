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

package floor4clientman

import (
	"fmt"
	"sync"

	"github.com/kasworld/goguelike/game/floor4client"
	"github.com/kasworld/goguelike/game/visitarea"
)

type Floor4ClientMan struct {
	mutex sync.RWMutex `prettystring:"hide"`

	name2fc map[string]*floor4client.Floor4Client
}

func New() *Floor4ClientMan {
	f4cm := &Floor4ClientMan{
		name2fc: make(map[string]*floor4client.Floor4Client),
	}
	return f4cm
}

func (f4cm *Floor4ClientMan) GetList() []*floor4client.Floor4Client {
	f4cm.mutex.RLock()
	defer f4cm.mutex.RUnlock()

	rtn := make([]*floor4client.Floor4Client, 0, len(f4cm.name2fc))
	for _, v := range f4cm.name2fc {
		rtn = append(rtn, v)
	}
	return rtn
}

func (f4cm *Floor4ClientMan) GetNameList() []string {
	f4cm.mutex.RLock()
	defer f4cm.mutex.RUnlock()

	rtn := make([]string, 0, len(f4cm.name2fc))
	for i := range f4cm.name2fc {
		rtn = append(rtn, i)
	}
	return rtn
}

func (f4cm *Floor4ClientMan) GetByName(name string) (*floor4client.Floor4Client, bool) {
	f4cm.mutex.RLock()
	defer f4cm.mutex.RUnlock()
	r, exist := f4cm.name2fc[name]
	return r, exist
}

func (f4cm *Floor4ClientMan) Add(fi visitarea.FloorI) *floor4client.Floor4Client {
	f4cm.mutex.Lock()
	defer f4cm.mutex.Unlock()
	f4cm.name2fc[fi.GetName()] = floor4client.New(fi)
	return f4cm.name2fc[fi.GetName()]
}

func (f4cm *Floor4ClientMan) Del(id string) {
	f4cm.mutex.Lock()
	defer f4cm.mutex.Unlock()

	delete(f4cm.name2fc, id)
}

func (f4cm *Floor4ClientMan) Forget(name string) error {
	f4cm.mutex.Lock()
	defer f4cm.mutex.Unlock()
	r, exist := f4cm.name2fc[name]
	if exist {
		r.Forget()
		return nil
	} else {
		return fmt.Errorf("not found %v", name)
	}
}
