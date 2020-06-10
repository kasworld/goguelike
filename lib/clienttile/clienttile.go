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

package clienttile

import (
	"net/url"
	"syscall/js"

	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/lib/webtilegroup"
	"github.com/kasworld/goguelike/lib/imagecanvas"
	"github.com/kasworld/gowasmlib/jslog"
)

type ClientTile struct {
	TilePNG *imagecanvas.ImageCanvas

	// premade from map tilegroup
	FloorTiles    [tile.Tile_Count][]webtilegroup.TileInfo
	FieldObjTiles [fieldobjdisplaytype.FieldObjDisplayType_Count][]webtilegroup.TileInfo
	CharTiles     [26][]webtilegroup.TileInfo
	EquipTiles    [equipslottype.EquipSlotType_Count][]webtilegroup.TileInfo
	GoldTiles     []webtilegroup.TileInfo
	PotionTiles   []webtilegroup.TileInfo
	ScrollTiles   []webtilegroup.TileInfo
	DirTiles      []webtilegroup.TileInfo
	Dir2Tiles     []webtilegroup.TileInfo
	CursorTiles   []webtilegroup.TileInfo
}

func New() *ClientTile {
	ct := ClientTile{}
	ct.TilePNG = imagecanvas.NewByID("tilepng")
	if ct.TilePNG == nil {
		jslog.Error("fail to get tilepng")
	}
	tileGroup := refreshTileGroup()
	for i := 0; i < tile.Tile_Count; i++ {
		ct.FloorTiles[i] = tileGroup[tile.Tile(i).String()]
	}

	ct.FieldObjTiles[fieldobjdisplaytype.None] = tileGroup["Trap"]
	ct.FieldObjTiles[fieldobjdisplaytype.StairUp] = tileGroup["StairUp"]
	ct.FieldObjTiles[fieldobjdisplaytype.StairDn] = tileGroup["StairDn"]
	ct.FieldObjTiles[fieldobjdisplaytype.PortalIn] = tileGroup["PortalIn"]
	ct.FieldObjTiles[fieldobjdisplaytype.PortalAutoIn] = tileGroup["PortalAutoIn"]
	ct.FieldObjTiles[fieldobjdisplaytype.PortalOut] = tileGroup["PortalOut"]
	ct.FieldObjTiles[fieldobjdisplaytype.Recycler] = tileGroup["Recycler"]

	for i := 0; i < factiontype.FactionType_Count; i++ {
		charTileName := "Char_" + factiontype.FactionType(i).Rune()
		ct.CharTiles[i] = tileGroup[charTileName]
	}

	for i := 0; i < equipslottype.EquipSlotType_Count; i++ {
		n := equipslottype.EquipSlotType(i).String()
		ct.EquipTiles[i] = tileGroup[n]
	}
	ct.GoldTiles = tileGroup["Gold"]
	ct.PotionTiles = tileGroup["Potion"]
	ct.ScrollTiles = tileGroup["Scroll"]
	ct.DirTiles = tileGroup["Dir"]
	ct.Dir2Tiles = tileGroup["Dir2"]
	ct.CursorTiles = tileGroup["Cursor"]
	return &ct
}

func refreshTileGroup() webtilegroup.TileGroup {
	tg, err := webtilegroup.GetTileGroup(replacePathFromHref("merged.json"))
	if err != nil {
		jslog.Errorf("tilegroup load fail %v\n", err)
	}
	return tg
}

func replacePathFromHref(s string) string {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	if err != nil {
		jslog.Errorf("%v", err)
		return ""
	}
	u.Path = s
	return u.String()
}
