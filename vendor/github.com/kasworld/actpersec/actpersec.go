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

// Package actpersec calc act per sec , lap and total
package actpersec

import (
	"fmt"
	"sync"
	"time"
)

type ActPerSec struct {
	mutex sync.RWMutex

	totalCount int
	beginTime  time.Time

	lastLapCount    int
	lastLapDuration time.Duration

	currentLapCount     int
	currentLapBeginTime time.Time
}

func (a ActPerSec) String() string {
	now := time.Now()
	totalDur := now.Sub(a.beginTime)
	lapCount := a.lastLapCount + a.currentLapCount
	lapDur := a.lastLapDuration + now.Sub(a.currentLapBeginTime)
	return fmt.Sprintf(
		"total[%v %4.2f/s] lap[%v %4.2f/s]",
		a.totalCount, float64(a.totalCount)/totalDur.Seconds(),
		lapCount, float64(lapCount)/lapDur.Seconds(),
	)
}

func New() *ActPerSec {
	now := time.Now()
	r := &ActPerSec{
		totalCount:          0,
		beginTime:           now,
		currentLapBeginTime: now,
	}
	return r
}

func (a *ActPerSec) GetTotalCount() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.totalCount
}

func (a *ActPerSec) PerSec() float64 {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	now := time.Now()
	totalDur := a.beginTime.Sub(now)
	perSec := float64(a.totalCount) / totalDur.Seconds()
	return perSec
}

func (a *ActPerSec) LapPerSec() float64 {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	now := time.Now()
	lapCount := a.lastLapCount + a.currentLapCount
	lapDur := a.lastLapDuration + now.Sub(a.currentLapBeginTime)
	perSec := float64(lapCount) / lapDur.Seconds()
	return perSec
}

func (a *ActPerSec) UpdateLap() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	now := time.Now()
	a.lastLapCount = a.currentLapCount
	a.lastLapDuration = now.Sub(a.currentLapBeginTime)
	a.currentLapCount = 0
	a.currentLapBeginTime = now

}

func (a *ActPerSec) Inc() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.totalCount++
	a.currentLapCount++
}

func (a *ActPerSec) Add(n int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.totalCount += n
	a.currentLapCount += n
}
