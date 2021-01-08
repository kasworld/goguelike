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

package tile_flag

import (
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

func (t TileFlag) Meditateable() bool {
	hp, sp := t.ActHPSP(c2t_idcmd.Meditate)
	return hp == 0 && sp == 0
}

func (t TileFlag) Empty() bool {
	return t == 0
}

func (t *TileFlag) OverrideBits(v tile.Tile) {
	t.ClearByTileFlag(overrideList[v])
	t.SetByTile(v)
}

const ( // bit groups
	Layer_1 = SeaFlag | MagmaFlag | SoilFlag | SandFlag | StoneFlag | SwampFlag | IceFlag
	Layer_2 = TreeFlag
	Layer_3 = FogFlag | SmokeFlag
	Layer_4 = RoadFlag | RoomFlag | WallFlag | WindowFlag | DoorFlag
)

var overrideList = []TileFlag{
	tile.Swamp: Layer_1,
	tile.Soil:  Layer_1,
	tile.Sand:  Layer_1,
	tile.Stone: Layer_1,
	tile.Sea:   Layer_1,
	tile.Magma: Layer_1,
	tile.Ice:   Layer_1,
	tile.Grass: Layer_1,

	tile.Tree: Layer_2,

	tile.Smoke: Layer_3,
	tile.Fog:   Layer_3,

	// tile.Road: Layer_1 | Layer_2 | Layer_4,
	// tile.Room: Layer_1 | Layer_2 | Layer_4,
	// Wall:   Layer_1 | Layer_2 | Layer_4,
	// Window: Layer_1 | Layer_2 | Layer_4,
	// Door:   Layer_1 | Layer_2 | Layer_4,
	tile.Window: WallFlag,
	tile.Door:   WallFlag | WindowFlag,
}
