// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// log destination group
package logdestinationgroup

import (
	"fmt"
	"sync"

	"github.com/kasworld/log/logdestinationi"
)

type LogDestinationGroup struct {
	mutex sync.RWMutex

	outputs  []logdestinationi.LogDestinationI
	name2out map[string]logdestinationi.LogDestinationI
}

func New() *LogDestinationGroup {
	return &LogDestinationGroup{
		name2out: make(map[string]logdestinationi.LogDestinationI),
	}
}

func (lw *LogDestinationGroup) Write(b []byte) error {
	for _, o := range lw.outputs[:] {
		if err := o.Write(b); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (lw *LogDestinationGroup) AddDestination(o logdestinationi.LogDestinationI) bool {
	lw.mutex.Lock()
	defer lw.mutex.Unlock()

	if _, ok := lw.name2out[o.Name()]; !ok {
		lw.name2out[o.Name()] = o
		lw.outputs = append(lw.outputs, o)
		return true
	}
	return false
}

func (lw *LogDestinationGroup) DelDestination(o logdestinationi.LogDestinationI) bool {
	lw.mutex.Lock()
	defer lw.mutex.Unlock()

	if _, ok := lw.name2out[o.Name()]; !ok {
		return false
	}
	delete(lw.name2out, o.Name())
	for i, v := range lw.outputs {
		if v == o {
			lw.outputs = append(lw.outputs[:i], lw.outputs[i+1:]...)
			return true
		}
	}
	return false
}

func (lw *LogDestinationGroup) GetDestinations() []logdestinationi.LogDestinationI {
	return lw.outputs[:]
}
