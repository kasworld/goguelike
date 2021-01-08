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

package resourcetile

import "github.com/kasworld/goguelike/enum/resourcetype"

type ResourceValue int32

func (t *ResourceTile) Move2To1(s1, s2, d resourcetype.ResourceType, rate ResourceValue) {
	amount := t[s2]
	if t[s1] < amount {
		amount = t[s1]
	}
	amount /= rate
	t[s1] -= amount
	t[s2] -= amount
	t[d] += amount
}

func (t *ResourceTile) TurnsBy(s, w, d resourcetype.ResourceType, rate ResourceValue) {
	amount := t[w]
	if t[s] < amount {
		amount = t[s]
	}
	amount /= rate
	t[s] -= amount
	t[d] += amount
}

func (t *ResourceTile) TurnsTo(s, d resourcetype.ResourceType, rate ResourceValue) {
	amount := t[s] / rate
	t[s] -= amount
	t[d] += amount
}

func (t *ResourceTile) TurnsTo2(s, d1, d2 resourcetype.ResourceType, rate ResourceValue) {
	amount := t[s] / rate
	t[s] -= amount
	t[d1] += amount
	t[d2] += amount
}

func (t *ResourceTile) MakeTo(s, d resourcetype.ResourceType, rate ResourceValue) {
	amount := t[s] / rate
	t[d] += amount
}

func (t *ResourceTile) Decay(s resourcetype.ResourceType, rate ResourceValue) {
	amount := t[s] / rate
	t[s] -= amount
}

func (t1 *ResourceTile) MoveTo(t2 *ResourceTile, ele resourcetype.ResourceType, rate ResourceValue) {
	amount := t1[ele] / rate
	t1[ele] -= amount
	t2[ele] += amount
}

func (t1 *ResourceTile) ExchWith(t2 *ResourceTile, ele resourcetype.ResourceType, rate ResourceValue) {
	amount := (t1[ele] - t2[ele]) / rate
	t1[ele] -= amount
	t2[ele] += amount
}
