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

package actjitter

import (
	"fmt"
	"sync"
	"time"
)

type ActJitter struct {
	mutex       sync.RWMutex `prettystring:"hide"`
	name        string
	startTime   time.Time `prettystring:"simple"`
	lastActTime time.Time `prettystring:"simple"`
	count       int64
	lastJitter  float64
	lastDur     time.Duration
	avgDur      time.Duration
}

func (jt *ActJitter) String() string {
	name := jt.name
	if jt.name == "" {
		name = "ActJitter"
	}
	return fmt.Sprintf(
		"%s[Count:%v Avg:%v Last[%v %4.2f%%]",
		name, jt.count, jt.GetAvg(), jt.lastDur, jt.lastJitter)
}

func New(name string) *ActJitter {
	jt := &ActJitter{
		name:        name,
		startTime:   time.Now(),
		lastActTime: time.Now(),
		count:       0,
	}
	return jt
}

func (jt *ActJitter) GetAvg() time.Duration {
	return jt.avgDur
}

func (jt *ActJitter) Act() float64 {
	jt.mutex.Lock()
	defer jt.mutex.Unlock()

	return jt.actByValue(time.Now())
}

func (jt *ActJitter) ActByValue(actTime time.Time) float64 {
	jt.mutex.Lock()
	defer jt.mutex.Unlock()

	return jt.actByValue(actTime)
}

func (jt *ActJitter) GetInFrameProgress() float64 {
	if jt.count == 0 {
		return 0.5
	}
	thisDur := time.Now().Sub(jt.lastActTime).Seconds()
	avg := jt.GetAvg().Seconds()
	rtn := thisDur / avg
	if rtn > 1 {
		rtn = 1
	}
	return rtn
}

func (jt *ActJitter) GetInFrameProgress2() float64 {
	if jt.count == 0 {
		return 0
	}
	thisDur := time.Now().Sub(jt.lastActTime).Seconds()
	avg := jt.GetAvg().Seconds()
	rtn := thisDur / avg
	return rtn
}

func (jt *ActJitter) actByValue(actTime time.Time) float64 {
	if jt.count == 0 {
		jt.lastActTime = actTime
		jt.count++
		jt.avgDur = actTime.Sub(jt.startTime)
		return jt.lastJitter
	}

	jt.count++
	thisDur := actTime.Sub(jt.lastActTime)
	oldAvg := jt.GetAvg()
	jt.avgDur = (oldAvg + thisDur) / 2
	jt.lastDur = thisDur
	jt.lastJitter = float64(jt.lastDur-jt.avgDur) * 100 / float64(oldAvg)
	jt.lastActTime = actTime
	return jt.lastJitter
}

func (jt *ActJitter) GetLastJitter() float64 {
	jt.mutex.RLock()
	defer jt.mutex.RUnlock()
	return jt.lastJitter
}

func (jt *ActJitter) GetLastDur() time.Duration {
	jt.mutex.RLock()
	defer jt.mutex.RUnlock()
	return jt.lastDur
}
