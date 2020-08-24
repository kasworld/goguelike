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

package jskeypressmap

import (
	"sync"

	"github.com/kasworld/goguelike/enum/way9type"
)

type KeyPressMap struct {
	mutex              sync.RWMutex `prettystring:"hide"`
	KeyboardPressedMap map[string]bool
}

func New() *KeyPressMap {
	kpm := &KeyPressMap{
		KeyboardPressedMap: make(map[string]bool),
	}
	return kpm
}

func (kpm *KeyPressMap) KeyUp(kcode string) bool {
	kpm.mutex.Lock()
	defer kpm.mutex.Unlock()
	delete(kpm.KeyboardPressedMap, kcode)
	return true
}
func (kpm *KeyPressMap) KeyDown(kcode string) bool {
	kpm.mutex.Lock()
	defer kpm.mutex.Unlock()
	kpm.KeyboardPressedMap[kcode] = true
	return true
}

func (kpm *KeyPressMap) IsPressed(kcode string) bool {
	kpm.mutex.RLock()
	defer kpm.mutex.RUnlock()
	return kpm.KeyboardPressedMap[kcode]
}

func (kpm *KeyPressMap) SumMoveDxDy(
	key2dir map[string]way9type.Way9Type) (int, int) {

	kpm.mutex.RLock()
	defer kpm.mutex.RUnlock()
	dx, dy := 0, 0
	for k, d := range key2dir {
		if kpm.KeyboardPressedMap[k] {
			dx += d.Dx()
			dy += d.Dy()
		}
	}
	return dx, dy
}
