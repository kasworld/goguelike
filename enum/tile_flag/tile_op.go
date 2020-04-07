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

package tile_flag

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tileoptype"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

func (t TileFlag) Meditateable() bool {
	hp, sp := t.ActHPSP(c2t_idcmd.Meditate)
	return hp == 0 && sp == 0
}

func (t TileFlag) Empty() bool {
	return t == 0
}

type TileTypeValue struct {
	T tileoptype.TileOpType
	V interface{}
}

func (t *TileFlag) Op(rv TileTypeValue) error {
	switch rv.V.(type) {
	case TileFlag:
		v := rv.V.(TileFlag)
		switch rv.T {
		case tileoptype.SetBits:
			t.SetByTileFlag(v)
		case tileoptype.ClearBits:
			t.ClearByTileFlag(v)
		default:
			return fmt.Errorf("Invalid op %v", rv)
		}
	case tile.Tile:
		v := rv.V.(tile.Tile)
		switch rv.T {
		case tileoptype.SetBit:
			t.SetByTile(v)
		case tileoptype.ClearBit:
			t.ClearByTile(v)
		case tileoptype.OverrideBits:
			t.OverrideBits(v)
		default:
			return fmt.Errorf("Invalid op %v", rv)
		}
	default:
		return fmt.Errorf("unknown type %v", rv)
	}
	return nil
}

func (t *TileFlag) OverrideBits(v tile.Tile) {
	t.ClearByTileFlag(overrideList[v])
	t.SetByTile(v)
}

const ( // bit groups
	layer_1 = SeaFlag | MagmaFlag | SoilFlag | SandFlag | StoneFlag | SwampFlag | IceFlag
	layer_2 = TreeFlag
	layer_3 = FogFlag | SmokeFlag
	layer_4 = RoadFlag | RoomFlag | WallFlag | WindowFlag | DoorFlag
)

var overrideList = []TileFlag{
	tile.Swamp: layer_1,
	tile.Soil:  layer_1,
	tile.Sand:  layer_1,
	tile.Stone: layer_1,
	tile.Sea:   layer_1,
	tile.Magma: layer_1,
	tile.Ice:   layer_1,
	tile.Grass: layer_1,

	tile.Tree: layer_2,

	tile.Smoke: layer_3,
	tile.Fog:   layer_3,

	tile.Road: layer_1 | layer_2 | layer_4,
	tile.Room: layer_1 | layer_2 | layer_4,
	// Wall:   layer_1 | layer_2 | layer_4,
	// Window: layer_1 | layer_2 | layer_4,
	// Door:   layer_1 | layer_2 | layer_4,
	tile.Window: WallFlag,
	tile.Door:   WallFlag | WindowFlag,
}
