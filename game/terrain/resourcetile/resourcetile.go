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

package resourcetile

import (
	"bytes"
	"fmt"

	"github.com/kasworld/goguelike/enum/resourcetype"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
)

func (rt ResourceTile) String() string {
	var buf bytes.Buffer
	buf.WriteString("ResourceTile[")
	for i := 0; i < resourcetype.ResourceType_Count; i++ {
		fmt.Fprintf(&buf, "%v:%v ", i, resourcetype.ResourceType(i).String())
	}
	buf.WriteString("]")
	return buf.String()
}

type ResourceTile [resourcetype.ResourceType_Count]ResourceValue

// progress in cell generation
func (rt *ResourceTile) ActTurnGeneration(
	rndIntNFn func(n int) int) {
	rt.Move2To1(resourcetype.Ice, resourcetype.Fire, resourcetype.Water, 10)

	rt.Move2To1(resourcetype.Water, resourcetype.Fire, resourcetype.Fog, 10)

	rt.Move2To1(resourcetype.Plant, resourcetype.Fire, resourcetype.Soil, 10)

	rt.TurnsBy(resourcetype.Fog, resourcetype.Plant, resourcetype.Water, 10)

	rt.TurnsBy(resourcetype.Stone, resourcetype.Plant, resourcetype.Soil, 10)

	if rt[resourcetype.Fog] > 1000 {
		rt.TurnsTo2(resourcetype.Fog, resourcetype.Water, resourcetype.Fire, 2)
	}
}

func (rt1 *ResourceTile) AgeingWith(
	rt2 *ResourceTile,
	rndIntNFn func(n int) int) {
	heightdiff := (rt1[resourcetype.Soil] + rt1[resourcetype.Stone]) - (rt2[resourcetype.Soil] + rt2[resourcetype.Stone])

	if heightdiff > 0 {
		rt2.MoveTo(rt1, resourcetype.Water, 10)
	} else if heightdiff < 0 {
		rt1.MoveTo(rt2, resourcetype.Water, 10)
	}
	if heightdiff > 0 {
		rt1.MoveTo(rt2, resourcetype.Soil, 10)
	} else if heightdiff < 0 {
		rt2.MoveTo(rt1, resourcetype.Soil, 20)
	}

	rt1.ExchWith(rt2, resourcetype.Fire, 10)
	rt1.ExchWith(rt2, resourcetype.Ice, 20)

	rt1.MoveTo(rt2, resourcetype.Fog, 1)
	if rt1[resourcetype.Plant] > 1000000 {
		rt1.ExchWith(rt2, resourcetype.Plant, 3)
	} else {
		rt1.MoveTo(rt2, resourcetype.Plant, 1)
	}
}

// convert raw to ground and on tile
func (rt *ResourceTile) Resource2View() tile_flag.TileFlag {
	var tl tile_flag.TileFlag
	if rt[resourcetype.Fog] > 1000000 {
		tl.SetByTile(tile.Smoke)
	} else if rt[resourcetype.Fog] > 1000 {
		tl.SetByTile(tile.Fog)
	}

	v := rt.findMax()
	if v >= 0 {
		switch tt := tile.Tile(v); tt {
		default:
			tl.SetByTile(tile.Tile(v))

		case tile.Swamp:
			if rt[resourcetype.Water]/(1+rt[resourcetype.Soil]+rt[resourcetype.Stone]+rt[resourcetype.Plant]) > 1 {
				tl.SetByTile(tile.Sea)
			} else {
				tl.SetByTile(tile.Swamp)
				if rt[resourcetype.Plant] > 1000000 {
					tl.SetByTile(tile.Tree)
				}
			}

		case tile.Soil:
			if rt[resourcetype.Soil] > 1000 && rt[resourcetype.Stone] > 1000 {
				tl.SetByTile(tile.Sand)
			} else {
				tl.SetByTile(tile.Soil)
			}
			if rt[resourcetype.Plant] > 1000000 {
				tl.SetByTile(tile.Tree)
			}

		case tile.Stone:
			if rt[resourcetype.Stone] > 1000000 {
				tl.SetByTile(tile.Wall)
			}
			if rt[resourcetype.Soil] > 1000 && rt[resourcetype.Stone] > 1000 {
				tl.SetByTile(tile.Sand)
			} else {
				tl.SetByTile(tile.Stone)
			}
			if rt[resourcetype.Plant] > 1000000 {
				tl.SetByTile(tile.Tree)
			}

		case tile.Ice:
			tl.SetByTile(tile.Ice)
			if rt[resourcetype.Plant] > 1000000 {
				tl.SetByTile(tile.Tree)
			}

		case tile.Grass:
			if rt[resourcetype.Plant] > 1000000 {
				tl.SetByTile(tile.Tree)
			}
			tl.SetByTile(tile.Grass)

		case tile.Magma:
			if rt[resourcetype.Fire]/(1000+rt[resourcetype.Soil]) > 1 {
				tl.SetByTile(tile.Magma)
			} else {
				tl.SetByTile(tile.Soil)
			}
		}
	}

	// do clear bit by priority
	for i := tile.Tile_Count; i >= 0; i-- {
		tlt := tile.Tile(i)
		if tl.TestByTile(tlt) {
			tl.OverrideBits(tile.Tile(i))
			tl.SetByTile(tlt)
		}
	}

	return tl
}

var resource2Ground = map[resourcetype.ResourceType]tile.Tile{
	resourcetype.Water: tile.Swamp,
	resourcetype.Soil:  tile.Soil,
	resourcetype.Stone: tile.Stone,
	resourcetype.Ice:   tile.Ice,
	resourcetype.Plant: tile.Grass,
	resourcetype.Fire:  tile.Magma,
}

func (rt *ResourceTile) findMax() int {
	var max, rtni = ResourceValue(0), -1
	for ri, ti := range resource2Ground {
		if rt[ri] > max {
			max, rtni = rt[ri], int(ti)
		}
	}
	return rtni
}

func (rt *ResourceTile) FromRGBA(r, g, b, a uint32) {
	rt[resourcetype.Fire] = ResourceValue(r)
	rt[resourcetype.Soil] = ResourceValue(g)
	rt[resourcetype.Water] = ResourceValue(b)
	rt[resourcetype.Fog] = ResourceValue(a)
}
