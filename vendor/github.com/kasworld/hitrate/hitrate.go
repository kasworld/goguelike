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

package hitrate

import (
	"fmt"
)

type HitRate struct {
	miss, hit int
}

func (hr HitRate) String() string {
	return fmt.Sprintf("HitRate[%5.2f%% %v/%v]",
		hr.HitRate(), hr.hit, hr.hit+hr.miss)
}

func (hr *HitRate) Hit() {
	hr.hit++
}
func (hr *HitRate) Miss() {
	hr.miss++
}

func (hr HitRate) GetTotal() int {
	return hr.hit + hr.miss
}
func (hr HitRate) HitRate() float64 {
	return float64(hr.hit) * 100 / float64(hr.hit+hr.miss)
}
