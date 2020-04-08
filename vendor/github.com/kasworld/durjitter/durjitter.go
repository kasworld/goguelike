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

// Package durjitter calc jitter without targetDuration
package durjitter

import (
	"fmt"
	"sync"
	"time"
)

type DurJitter struct {
	mutex sync.RWMutex `webformhide:"" stringformhide:""`
	name  string

	count int
	sum   time.Duration

	lastDur    time.Duration
	lastJitter float64
}

func (jt *DurJitter) String() string {
	name := jt.name
	if name == "" {
		name = "DurJitter"
	}
	return fmt.Sprintf("%s[avg:%v jitter:%4.2f%%]",
		name, jt.GetAvg(), jt.lastJitter)
}

func New(name string) *DurJitter {
	jt := &DurJitter{
		name: name,
	}
	return jt
}

func (jt *DurJitter) GetCount() int {
	return jt.count
}

func (jt *DurJitter) GetAvg() time.Duration {
	if jt.count == 0 {
		return 0
	}
	return jt.sum / time.Duration(jt.count)
}

func (jt *DurJitter) Add(dur time.Duration) float64 {
	jt.mutex.Lock()
	defer jt.mutex.Unlock()

	if jt.count == 0 {
		jt.sum += dur
		jt.count++
		return jt.lastJitter
	}

	jt.lastJitter = float64(dur-jt.GetAvg()) * 100 / float64(jt.GetAvg())
	jt.sum += dur
	jt.lastDur = dur
	jt.count++
	return jt.lastJitter
}

func (jt *DurJitter) GetLastDuration() time.Duration {
	jt.mutex.RLock()
	defer jt.mutex.RUnlock()
	return jt.lastDur
}

func (jt *DurJitter) GetLastJitter() float64 {
	jt.mutex.RLock()
	defer jt.mutex.RUnlock()
	return jt.lastJitter
}
