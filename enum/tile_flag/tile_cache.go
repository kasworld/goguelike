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
	"bytes"
	"fmt"
	"math"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/htmlcolors"
)

var calcedTileValuesCache [1 << uint(tile.Tile_Count)]*tile.AttribEle

func (t TileFlag) calcTileValues() *tile.AttribEle {
	rtn := tile.AttribEle{
		AtkMod:    1,
		DefMod:    1,
		Color:     htmlcolors.Black.UInt8Array(),
		ObjOnTile: tile.ObjOnTile_1,
	}
	var buf bytes.Buffer
	for i := 0; i < tile.Tile_Count; i++ {
		tlt := tile.Tile(i)
		if t.TestByTile(tlt) {
			v := tlt.Attrib()
			rtn.BlockSight += v.BlockSight
			rtn.NeedHPOther += v.NeedHPOther
			rtn.NeedHPMove += v.NeedHPMove
			rtn.NeedHPAttack += v.NeedHPAttack
			rtn.NeedSPOther += v.NeedSPOther
			rtn.NeedSPMove += v.NeedSPMove
			rtn.NeedSPAttack += v.NeedSPAttack
			rtn.AtkMod *= v.AtkMod
			rtn.DefMod *= v.DefMod
			for j, _ := range rtn.Color {
				rtn.Color[j] = rtn.Color[j]/2 + v.Color[j]/2
			}
			rtn.Slippery = rtn.Slippery || v.Slippery
			if rtn.ObjOnTile == tile.ObjOnTile_1 && v.ObjOnTile == tile.ObjOnTile_N {
				rtn.ObjOnTile = tile.ObjOnTile_N
			}
			if v.ObjOnTile == tile.ObjOnTile_0 {
				rtn.ObjOnTile = tile.ObjOnTile_0
			}
			if buf.Len() > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(v.Name)
		}
	}
	rtn.Name = buf.String()
	if t == 0 {
		rtn.BlockSight = math.Inf(1)
		rtn.ObjOnTile = tile.ObjOnTile_0
	}
	return &rtn
}

func (t TileFlag) getTileValuesByCached() *tile.AttribEle {
	rtn := calcedTileValuesCache[t]
	if rtn == nil {
		rtn = t.calcTileValues()
		calcedTileValuesCache[t] = rtn
	}
	return rtn
}

func (t TileFlag) ToRGB() (uint8, uint8, uint8) {
	rtn := t.getTileValuesByCached().Color
	return rtn[0], rtn[1], rtn[2]
}

func (t TileFlag) BlockSight() float64 {
	return t.getTileValuesByCached().BlockSight
}

func (t TileFlag) Slippery() bool {
	return t.getTileValuesByCached().Slippery
}

func (t TileFlag) ActHPSP(act c2t_idcmd.CommandID) (float64, float64) {
	rtn := t.getTileValuesByCached()
	switch act {
	default:
		return rtn.NeedHPOther, rtn.NeedSPOther
	case c2t_idcmd.Move:
		return rtn.NeedHPMove, rtn.NeedSPMove
	case c2t_idcmd.Attack:
		return rtn.NeedHPAttack, rtn.NeedSPAttack
	}
}

func (t TileFlag) ActHPSPCalced(act c2t_idcmd.CommandID, dir way9type.Way9Type) (float64, float64) {
	rtn := t.getTileValuesByCached()
	switch act {
	default:
		return rtn.NeedHPOther, rtn.NeedSPOther
	case c2t_idcmd.Move:
		if dir == way9type.Center {
			return rtn.NeedHPMove, rtn.NeedSPMove
		}
		return rtn.NeedHPMove * dir.Len(), rtn.NeedSPMove * dir.Len()
	case c2t_idcmd.Attack:
		return rtn.NeedHPAttack, rtn.NeedSPAttack
	}
}

func (t TileFlag) AtkMod() float64 {
	return t.getTileValuesByCached().AtkMod
}

func (t TileFlag) DefMod() float64 {
	return t.getTileValuesByCached().DefMod
}

func (t TileFlag) CannotPlaceObj() bool {
	return t.getTileValuesByCached().ObjOnTile == tile.ObjOnTile_0
}

func (t TileFlag) CharPlaceable() bool {
	return t.getTileValuesByCached().ObjOnTile != tile.ObjOnTile_0
}

func (t TileFlag) CanPlaceMultiObj() bool {
	return t.getTileValuesByCached().ObjOnTile == tile.ObjOnTile_N
}

func (t TileFlag) NoBattle() bool {
	return t.getTileValuesByCached().ObjOnTile == tile.ObjOnTile_N
}

func (t TileFlag) CanBattle() bool {
	return t.getTileValuesByCached().ObjOnTile == tile.ObjOnTile_1
}

func (t TileFlag) Name() string {
	return t.getTileValuesByCached().Name
}

func (t TileFlag) String() string {
	return fmt.Sprintf("TileFlag[%s]", t.Name())
}

func TileCacheCount() int {
	rtn := 0
	for _, v := range calcedTileValuesCache {
		if v != nil {
			rtn++
		}
	}
	return rtn
}
