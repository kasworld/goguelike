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

// Package activebuff timed buff in activeobject
package activebuff

import (
	"sync"

	"github.com/kasworld/goguelike/enum/statusoptype"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type BuffManager struct {
	Mutex    sync.RWMutex  `prettystring:"hide"`
	BuffList []*ActiveBuff `prettystring:"simple"`
}

func New() *BuffManager {
	bm := &BuffManager{
		BuffList: make([]*ActiveBuff, 0),
	}
	return bm
}

// Add return true if replace
func (bm *BuffManager) Add(name string, clearOnRebirth bool, replaceSameName bool, tb []statusoptype.OpArg) bool {
	afs := &ActiveBuff{
		Name:           name,
		ClearOnRebirth: clearOnRebirth,
		Buff:           tb,
	}
	bm.Mutex.Lock()
	defer bm.Mutex.Unlock()

	if replaceSameName {
		for i, v := range bm.BuffList {
			if v != nil && v.Name == name {
				bm.BuffList[i] = afs
				return true
			}
		}
	}
	for i, v := range bm.BuffList {
		if v == nil {
			bm.BuffList[i] = afs
			return false
		}
	}
	bm.BuffList = append(bm.BuffList, afs)
	return false
}

func (bm *BuffManager) Exist(name string) bool {
	for _, v := range bm.BuffList {
		if v != nil && v.Name == name {
			return true
		}
	}
	return false
}

func (bm *BuffManager) ClearOnRebirth() {
	bm.Mutex.Lock()
	defer bm.Mutex.Unlock()
	for i, v := range bm.BuffList {
		if v == nil {
			continue
		}
		if v.ClearOnRebirth {
			bm.BuffList[i] = nil
		}
	}
}

func (bm *BuffManager) GetOpArgListToApply() []statusoptype.OpArg {
	rtn := make([]statusoptype.OpArg, 0)
	bm.Mutex.Lock()
	defer bm.Mutex.Unlock()
	for i, v := range bm.BuffList {
		if v == nil {
			continue
		}
		oparg, remain := v.GetOpArgToApply()
		if !remain {
			bm.BuffList[i] = nil
			continue
		}
		rtn = append(rtn, oparg)
	}
	return rtn
}

func (bm *BuffManager) ToPacket_ActiveObjBuffList() []*c2t_obj.ActiveObjBuff {
	rtn := make([]*c2t_obj.ActiveObjBuff, 0)
	bm.Mutex.RLock()
	defer bm.Mutex.RUnlock()
	for _, v := range bm.BuffList {
		if v == nil || v.IsEnded() {
			continue
		}
		rtn = append(rtn, v.ToPacket_ActiveObjBuff())
	}
	return rtn
}
