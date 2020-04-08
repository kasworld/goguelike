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

import (
	"fmt"
	"sync"
	"time"
)

func (rs RangeStat) String() string {
	return fmt.Sprintf(
		"%s[hmin:%v smin:%v val:%v smax:%v hmax:%v total%v lap%v]",
		rs.name,
		rs.hardMin,
		rs.softMin,
		rs.currentVal,
		rs.softMax,
		rs.hardMax,
		rs.total,
		rs.lap,
	)
}

// range include min , include max
type RangeStat struct {
	mutex sync.RWMutex

	name string

	// for range
	hardMin    int
	softMin    int
	currentVal int
	softMax    int
	hardMax    int

	total statElement
	lap   statElement
}

func New(name string, min, max int) *RangeStat {
	if min >= max {
		return nil
	}
	if name == "" {
		name = "RangeStat"
	}
	rs := &RangeStat{
		name:       name,
		hardMin:    min,
		softMin:    min,
		currentVal: min,
		softMax:    max,
		hardMax:    max,
		total: statElement{
			decCount:  0,
			incCount:  0,
			beginTime: time.Now().UTC(),
		},
		lap: statElement{
			decCount:  0,
			incCount:  0,
			beginTime: time.Now().UTC(),
		},
	}
	return rs
}
