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

package rangestat

import "time"

func (rs *RangeStat) isInSoftRange(v int) bool {
	return rs.softMin <= v && v <= rs.softMax
}

func (rs *RangeStat) SetSoftMax(v int) bool {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	if rs.hardMin >= v || v > rs.hardMax {
		return false
	}
	if rs.softMin >= v {
		return false
	}
	rs.softMax = v
	return true
}
func (rs *RangeStat) SetSoftMin(v int) bool {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	if rs.hardMin > v || v >= rs.hardMax {
		return false
	}
	if rs.softMax <= v {
		return false
	}
	rs.softMin = v
	return true
}

func (rs *RangeStat) CanInc() bool {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	if (rs.currentVal + 1) > rs.softMax {
		return false
	}
	return true
}

func (rs *RangeStat) Inc() bool {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	if (rs.currentVal + 1) > rs.softMax {
		return false
	}
	rs.currentVal++
	rs.total.incCount++
	rs.lap.incCount++
	return true
}

func (rs *RangeStat) CanDec() bool {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	if (rs.currentVal - 1) < rs.softMin {
		return false
	}
	return true
}
func (rs *RangeStat) Dec() bool {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	if (rs.currentVal - 1) < rs.softMin {
		return false
	}
	rs.currentVal--
	rs.total.decCount++
	rs.lap.decCount++
	return true
}

func (rs *RangeStat) GetCurrentVal() int {
	return rs.currentVal
}
func (rs *RangeStat) GetHardMin() int {
	return rs.hardMin
}
func (rs *RangeStat) GetHardMax() int {
	return rs.hardMax
}
func (rs *RangeStat) GetSoftMin() int {
	return rs.softMin
}
func (rs *RangeStat) GetSoftMax() int {
	return rs.softMax
}

func (rs *RangeStat) UpdateLap() {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	rs.lap = statElement{
		decCount:  0,
		incCount:  0,
		beginTime: time.Now().UTC(),
	}
}

func (rs *RangeStat) GetTotalInc() int {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.total.incCount
}
func (rs *RangeStat) GetTotalDec() int {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.total.decCount
}
