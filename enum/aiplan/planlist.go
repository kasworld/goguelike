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

package aiplan

import (
	"bytes"
	"fmt"
)

type PlanList []AIPlan

func (pl PlanList) String() string {
	var buf bytes.Buffer
	for _, v := range pl {
		fmt.Fprintf(&buf, "%s ", v)
	}
	return buf.String()
}

// get front plan
func (pl PlanList) GetCurrentPlan() AIPlan {
	return pl.GetLast()
}

func (pl PlanList) GetFront() AIPlan {
	return pl[0]
}
func (pl PlanList) GetLast() AIPlan {
	return pl[len(pl)-1]
}

func (pl PlanList) FindPos(p AIPlan) int {
	pos := -1
	for i, v := range pl {
		if v == p {
			pos = i
			break
		}
	}
	return pos
}

// move p to front
func (pl PlanList) Move2Front(p AIPlan) bool {
	pos := pl.FindPos(p)
	if pos != -1 {
		copy(pl[1:], pl[:pos])
		pl[0] = p
		return true
	}
	return false
}

// pop current and push back
func (pl PlanList) Front2Last() {
	p := pl[0]
	copy(pl, pl[1:])
	pl[len(pl)-1] = p
}

func (pl PlanList) Dup() PlanList {
	rtn := make(PlanList, len(pl))
	copy(rtn, pl)
	return rtn
}
