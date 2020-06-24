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

package g2id

import "testing"

func TestUUID_String(t *testing.T) {
	id := New()
	t.Logf("%v %v", id, len(id))
	id2 := NewFromString(id.String())
	t.Logf("%v", id2)
}

func TestUUID_Compare(t *testing.T) {
	id1 := New()
	id2 := id1
	t.Logf("%v == %v : %v", id1, id2, id1 == id2)
}

func TestUUID_map(t *testing.T) {
	idmap := make(map[UUID]bool)
	id1 := New()
	id2 := id1

	idmap[id1] = true
	idmap[id2] = false

	t.Logf("%v", idmap)
}
