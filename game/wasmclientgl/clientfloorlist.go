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

package wasmclientgl

import "sort"

type ClientFloorGLList []*ClientFloorGL

func (cfList ClientFloorGLList) Len() int { return len(cfList) }
func (cfList ClientFloorGLList) Swap(i, j int) {
	cfList[i], cfList[j] = cfList[j], cfList[i]
}
func (cfList ClientFloorGLList) Less(i, j int) bool {
	ao1 := cfList[i]
	ao2 := cfList[j]
	if ao1.visitTime == ao2.visitTime {
		return ao1.FloorInfo.Name < ao2.FloorInfo.Name
	}
	return ao1.visitTime.After(ao2.visitTime)
}
func (cfList ClientFloorGLList) Sort() {
	sort.Stable(cfList)
}