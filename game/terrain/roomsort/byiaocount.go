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

package roomsort

import (
	"sort"

	"github.com/kasworld/goguelike/game/terrain/room"
)

type ByRecyclerCount []*room.Room

func (rl ByRecyclerCount) Len() int { return len(rl) }
func (rl ByRecyclerCount) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}
func (rl ByRecyclerCount) Less(i, j int) bool {
	r1 := rl[i]
	r2 := rl[j]
	if r1.RecyclerCount == r2.RecyclerCount {
		return r1.PortalCount < r2.PortalCount
	}
	return r1.RecyclerCount < r2.RecyclerCount
}
func (rl ByRecyclerCount) Sort() {
	sort.Sort(rl)
}

type ByPortalCount []*room.Room

func (rl ByPortalCount) Len() int { return len(rl) }
func (rl ByPortalCount) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}
func (rl ByPortalCount) Less(i, j int) bool {
	r1 := rl[i]
	r2 := rl[j]
	if r1.PortalCount == r2.PortalCount {
		return r1.RecyclerCount < r2.RecyclerCount
	}
	return r1.PortalCount < r2.PortalCount
}
func (rl ByPortalCount) Sort() {
	sort.Sort(rl)
}

type ByTrapCount []*room.Room

func (rl ByTrapCount) Len() int { return len(rl) }
func (rl ByTrapCount) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}
func (rl ByTrapCount) Less(i, j int) bool {
	r1 := rl[i]
	r2 := rl[j]
	if r1.TrapCount == r2.TrapCount {
		return r1.RecyclerCount < r2.RecyclerCount
	}
	return r1.TrapCount < r2.TrapCount
}
func (rl ByTrapCount) Sort() {
	sort.Sort(rl)
}

type ByRotateLineAttackCount []*room.Room

func (rl ByRotateLineAttackCount) Len() int { return len(rl) }
func (rl ByRotateLineAttackCount) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}
func (rl ByRotateLineAttackCount) Less(i, j int) bool {
	r1 := rl[i]
	r2 := rl[j]
	if r1.RotateLineAttackCount == r2.RotateLineAttackCount {
		return r1.RecyclerCount < r2.RecyclerCount
	}
	return r1.RotateLineAttackCount < r2.RotateLineAttackCount
}
func (rl ByRotateLineAttackCount) Sort() {
	sort.Sort(rl)
}

type ByMineCount []*room.Room

func (rl ByMineCount) Len() int { return len(rl) }
func (rl ByMineCount) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}
func (rl ByMineCount) Less(i, j int) bool {
	r1 := rl[i]
	r2 := rl[j]
	if r1.MineCount == r2.MineCount {
		return r1.RecyclerCount < r2.RecyclerCount
	}
	return r1.MineCount < r2.MineCount
}
func (rl ByMineCount) Sort() {
	sort.Sort(rl)
}
