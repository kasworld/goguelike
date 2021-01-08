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

package bias

type BiasList []Bias

func (bl BiasList) Len() int { return len(bl) }
func (bl BiasList) Swap(i, j int) {
	bl[i], bl[j] = bl[j], bl[i]
}
func (bl BiasList) Less(i, j int) bool {
	for k := 0; k < 3; k++ {
		if bl[i][k] == bl[j][k] {
			continue
		}
		return bl[i][k] < bl[j][k]
	}
	return false
}
