// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package aoscore score data for ground server
package aoscore

import "sort"

type ByExp []*ActiveObjScore

func (aoList ByExp) Len() int { return len(aoList) }
func (aoList ByExp) Swap(i, j int) {
	aoList[i], aoList[j] = aoList[j], aoList[i]
}
func (aoList ByExp) Less(i, j int) bool {
	ao1 := aoList[i]
	ao2 := aoList[j]
	return ao1.Exp > ao2.Exp
}
func (aoList ByExp) Sort() {
	sort.Stable(aoList)
}
