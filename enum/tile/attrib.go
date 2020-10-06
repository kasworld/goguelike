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

package tile

import (
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/htmlcolors"
)

type ObjOnTile_Type uint8

const (
	ObjOnTile_0 ObjOnTile_Type = iota // forbid place obj
	ObjOnTile_1                       // can place a obj, can battle
	ObjOnTile_N                       // can place multiful obj, cannot battle
)

type AttribEle struct {
	NeedHPOther  float64
	NeedHPMove   float64
	NeedHPAttack float64
	NeedSPOther  float64
	NeedSPMove   float64
	NeedSPAttack float64
	AtkMod       float64
	DefMod       float64
	BlockSight   float64
	Slippery     bool
	Color        [3]uint8
	ObjOnTile    ObjOnTile_Type
	Name         string
}

func (tt Tile) Attrib() AttribEle {
	return attrib[tt]
}

var attrib = [Tile_Count]AttribEle{
	Swamp: {1.0, 1.0, 1.0, 1.0, 2.0, 3.0, 0.8, 1.2, 1.0, false,
		htmlcolors.LightSeaGreen.UInt8Array(), ObjOnTile_1, "Swamp"},
	Soil: {0.0, 0.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0, false,
		htmlcolors.SaddleBrown.UInt8Array(), ObjOnTile_1, "Soil"},
	Sand: {0.0, 0.0, 0.0, 1.0, 1.5, 1.5, 0.9, 0.9, 1.0, false,
		htmlcolors.SandyBrown.UInt8Array(), ObjOnTile_1, "Sand"},
	Stone: {0.0, 0.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0, false,
		htmlcolors.LightSteelBlue.UInt8Array(), ObjOnTile_1, "Stone"},
	Sea: {1.0, 1.0, 1.0, 5.0, 7.0, 9.0, 0.5, 2.0, 1.0, false,
		htmlcolors.Navy.UInt8Array(), ObjOnTile_1, "Sea"},
	Magma: {5.0, 5.0, 5.0, 2.0, 2.0, 2.0, 1.1, 0.9, 1.1, false,
		htmlcolors.Crimson.UInt8Array(), ObjOnTile_1, "Magma"},
	Road: {0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 1.0, 1.0, false,
		htmlcolors.Gray.UInt8Array(), ObjOnTile_N, "Road"},
	Room: {0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 1.0, 1.0, false,
		htmlcolors.Yellow.UInt8Array(), ObjOnTile_N, "Room"},
	Ice: {1.0, 1.0, 1.0, 1.0, 2.0, 3.0, 0.8, 0.8, 1.0, true,
		htmlcolors.LightBlue.UInt8Array(), ObjOnTile_1, "Ice"},
	Grass: {0.0, 0.0, 0.0, 0.0, 1.0, 2.0, 0.9, 1.1, 1.1, false,
		htmlcolors.LimeGreen.UInt8Array(), ObjOnTile_1, "Grass"},
	Tree: {0.0, 0.0, 0.0, 0.0, 2.0, 3.0, 0.5, 2.0, 2.0, false,
		htmlcolors.ForestGreen.UInt8Array(), ObjOnTile_1, "Tree"},
	Wall: {9.9, 9.9, 9.9, 9.9, 9.9, 9.9, 1.0, 1.0, gameconst.SightBlock, false,
		htmlcolors.Silver.UInt8Array(), ObjOnTile_0, "Wall"},
	Window: {9.9, 9.9, 9.9, 9.9, 9.9, 9.9, 1.0, 1.0, 1.1, false,
		htmlcolors.YellowGreen.UInt8Array(), ObjOnTile_0, "Window"},
	Door: {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, false,
		htmlcolors.GreenYellow.UInt8Array(), ObjOnTile_N, "Door"},
	Smoke: {1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.9, 1.2, 1.3, false,
		htmlcolors.DimGray.UInt8Array(), ObjOnTile_1, "Smoke"},
	Fog: {0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.8, 1.5, 1.2, false,
		htmlcolors.SlateGray.UInt8Array(), ObjOnTile_1, "Fog"},
}

var TileScrollAttrib = [Tile_Count]struct {
	Texture      bool
	SrcXBiasAxis int // 0-2 in bias
	AniXSpeed    float64
	SrcYBiasAxis int // 0-2 in bias
	AniYSpeed    float64
}{
	Fog:   {true, 0, 64, 1, 64},
	Smoke: {true, 0, 48, 1, 48},

	Magma: {true, 2, 64, 0, 64},

	Swamp: {true, 0, 48, 2, 48},
	Sea:   {true, 0, 24, 2, 24},
	Ice:   {true, 0, 8, 2, 8},

	Grass: {true, 1, 8, 2, 8},
	Tree:  {true, 1, 4, 2, 4},

	Sand:  {true, 0, 8, 1, 8},
	Soil:  {true, 0, 2, 1, 2},
	Stone: {true, 0, 1, 1, 1},

	Road:   {true, 0, 0, 0, 0},
	Room:   {true, 0, 0, 0, 0},
	Wall:   {false, 0, 0, 0, 0},
	Window: {false, 0, 0, 0, 0},
	Door:   {false, 0, 0, 0, 0},
}

var TileShadowAttrib = [Tile_Count]struct {
	Cast    bool
	Receive bool
}{
	Swamp:  {false, true},
	Soil:   {false, true},
	Stone:  {false, true},
	Sand:   {false, true},
	Sea:    {false, true},
	Magma:  {false, true},
	Ice:    {false, true},
	Grass:  {false, true},
	Tree:   {true, false},
	Road:   {false, true},
	Room:   {false, true},
	Wall:   {true, false},
	Window: {false, false},
	Door:   {false, false},
	Fog:    {false, false},
	Smoke:  {false, false},
}
