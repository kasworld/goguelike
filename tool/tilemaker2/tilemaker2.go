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

package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/lib/webtilegroup"
	"github.com/kasworld/goguelike/tool/tilemaker2/drawtile"
	"github.com/kasworld/goguelike/tool/tilemaker2/runetile"
	"github.com/kasworld/goguelike/tool/tilemaker2/tile4merge"
	"github.com/kasworld/htmlcolors"
)

func main() {

	var tiList []*tile4merge.Tile4Merge
	tiList = append(tiList, tile4merge.LoadImageDir("bitmaptiles/")...)
	rList, err := MakeAllRuneTile(64)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	for _, rti := range rList {
		imgSize := rti.Img.Bounds().Size()
		tiList = append(tiList, &tile4merge.Tile4Merge{
			GroupName: rti.GroupName,
			TileName:  rti.Name,
			Image:     rti.Img,
			Rect: webtilegroup.Rect{
				W: imgSize.X,
				H: imgSize.Y,
			},
		})
	}

	dtList := []*drawtile.DrawTile{
		drawtile.New("Cursor", "cursor", 64, 64),
		drawtile.New("Cursor", "cursorred", 64, 64),
		drawtile.New("Cursor", "cursorx", 64, 64),
	}
	dtList[0].DrawRect(htmlcolors.White)
	dtList[1].DrawRect(htmlcolors.Red)
	dtList[2].DrawRect(htmlcolors.Red)
	dtList[2].DrawX(htmlcolors.Red)
	for _, dt := range dtList {
		imgSize := dt.Img.Bounds().Size()
		tiList = append(tiList, &tile4merge.Tile4Merge{
			GroupName: dt.GroupName,
			TileName:  dt.Name,
			Image:     dt.Img,
			Rect: webtilegroup.Rect{
				W: imgSize.X,
				H: imgSize.Y,
			},
		})
	}

	if err := tile4merge.SaveMergedPNG(tiList, "merged.png"); err != nil {
		fmt.Printf("%v\n", err)
	}

	if err := tile4merge.SaveWebJSON(tiList, "merged.json"); err != nil {
		fmt.Printf("%v\n", err)
	}
}

/*
Dir : •↑↗→↘↓↙←↖

∘⇑⇗⇒⇘⇓⇙⇐⇖


Marker : 1234567890
PortalAutoIn : ⦾⦿
PortalIn : ♠♥♣♦
PortalOut : ♤♡♧♢
Recycler : ♲♳♴♵♶♷♸♹♺♻♼♽
Trap : ☠💀☢☣

portal in/out ⎆⎒ ⍚⎐

wall ═║╔╗╚╝╠╣╦╩╬
window ▦
door ◫ ▪▫



▓▀▐╚▄║╔╠▌╝═╩╗╣╦╬
▓╹╺╚╻║╔╠╸╝═╩╗╣╦╬
▓╹╺┗╻┃┏┣╸┛━┻┓┫┳╋



▓╸╹╺╻━┃┏┓┗┛┣┫┳┻╋
╍╏┅┇┉┋

◯╳

↧↥

up/down stair  ⥮⮁⇅ ⥯⮃⇵   ⭿⇕⬍↕ ⍗⍐ ⇱⇲ ⍏⍖ ⍍⍔

wall : ▓◆♦⯁✦❖

money 💰

gem : 💎

conditions
high volt ⚡
Sparkles ✨
Dizzy 💫
Anger Symbol 💢
Collision 💥
Hole 🕳️
Right Anger Bubble 	🗯️
Dashing Away 💨
Sweat Droplets 💦
Droplet 💧
Fountain ⛲
Syringe 💉
Fog 🌫️
Thought Balloon 💭
Hammer and Pick ⚒️
Pick ⛏️
Dagger 🗡️
Bomb 💣
Fire 🔥
Key 🔑
Locked 🔒
Unlocked 🔓
Old Key 🗝️
Door 🚪
Bow and Arrow 🏹

food
Meat on Bone 🍖
Hamburger 🍔
Poultry Leg 🍗
Cut of Meat 🥩
Candy 🍬
Bread 🍞
Grapes 🍇
Cherries 🍒
Steaming Bowl 🍜
Soft Ice Cream 🍦
Beer Mug 🍺
*/

// /usr/share/fonts/truetype/noto/NotoMono-Regular.ttf
// /usr/share/fonts/truetype/nanum/NanumBarunGothic.ttf
// /usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf

func MakeAllRuneTile(size int) ([]runetile.RuneTile, error) {
	shift := size / 16
	letterFont, err := runetile.LoadFontFile("/usr/share/fonts/truetype/nanum/NanumBarunGothic.ttf")
	if err != nil {
		return nil, err
	}
	symbola, err := runetile.LoadFontFile("/usr/share/fonts/truetype/ancient-scripts/Symbola_hint.ttf")
	if err != nil {
		return nil, err
	}

	rtn := make([]runetile.RuneTile, 0)

	// rtn = append(rtn,
	//		runetile.NewRuneTileListByString("WallCenter", "◆♦⯁✦❖", htmlcolors.Gray, symbola, 64, 0,0,shift,64,64)...)

	rtn = append(rtn, runetile.NewRuneTileListByString("Wall", "▓╹╺╚╻║╔╠╸╝═╩╗╣╦╬",
		htmlcolors.Gray, symbola, 64, 6, 4, shift, 64, 64)...)
	rtn = append(rtn, runetile.NewRuneTileListByString("Window", "┉┋",
		htmlcolors.Gray, symbola, 64, 6, 4, shift, 64, 64)...)
	rtn = append(rtn, runetile.NewRuneTileListByString("Door", "⛶",
		htmlcolors.Gray, symbola, 64, 0, 0, shift, 64, 64)...)

	// rtn = append(rtn, []runetile.RuneTile{
	// 	{GroupName: "Cursor", Rune: "◯", Name: "cursor", Color1: htmlcolors.White, FontData: symbola, FontSize: 56},
	// 	{GroupName: "Cursor", Rune: "◯", Name: "cursorred", Color1: htmlcolors.Red, FontData: symbola, FontSize: 56},
	// 	{GroupName: "Cursor", Rune: "⛝", Name: "cursorx", Color1: htmlcolors.Red, FontData: symbola, FontSize: 64},
	// }...)

	rtn = append(rtn, runetile.NewRuneTileListByString("StairDn", "⍗⍔",
		htmlcolors.Gray, symbola, 64, 0, 0, shift, 64, 64)...)
	rtn = append(rtn, runetile.NewRuneTileListByString("StairUp", "⍐⍍",
		htmlcolors.Gray, symbola, 64, 0, 0, shift, 64, 64)...)

	rtn = append(rtn,
		runetile.NewRuneTileListByString("Dir", "•↑↗→↘↓↙←↖", htmlcolors.Gray, symbola, 56, 0, 0, shift, 64, 64)...)
	rtn = append(rtn,
		runetile.NewRuneTileListByString("Dir2", "∘⇑⇗⇒⇘⇓⇙⇐⇖", htmlcolors.Gray, symbola, 56, 0, 0, shift, 64, 64)...)
	rtn = append(rtn,
		runetile.NewRuneTileListByString("PortalAutoIn", "⎆", htmlcolors.Gray, symbola, 56, 0, 0, shift, 64, 64)...)
	rtn = append(rtn,
		runetile.NewRuneTileListByString("PortalIn", "⍚", htmlcolors.Gray, symbola, 56, 0, 0, shift, 64, 64)...)
	rtn = append(rtn,
		runetile.NewRuneTileListByString("PortalOut", "⎒", htmlcolors.Gray, symbola, 56, 0, 0, shift, 64, 64)...)
	rtn = append(rtn,
		runetile.NewRuneTileListByString("Recycler", "♲♻♼♽", htmlcolors.Gray, symbola, 56, 0, 0, shift, 64, 64)...)
	rtn = append(rtn,
		runetile.NewRuneTileListByString("Trap", "☠💀☢☣", htmlcolors.Gray, symbola, 56, 0, 0, shift, 64, 64)...)

	// gold
	for _, v := range moneycolor.Attrib {
		mo := runetile.RuneTile{
			GroupName: "Gold", Rune: "$",
			Name:     fmt.Sprintf("%d", v.UpLimit),
			Color1:   v.Color,
			FontData: letterFont, FontSize: 56,
			FontDx: 0, FontDy: 0, ShadowShift: shift, TileW: 64, TileH: 64,
		}
		rtn = append(rtn, mo)
	}

	// potion
	for i := 0; i < potiontype.PotionType_Count; i++ {
		pot := potiontype.PotionType(i)
		rtn = append(rtn, runetile.RuneTile{
			GroupName:   "Potion",
			Rune:        pot.Rune(),
			Name:        pot.String(),
			Color1:      pot.Color24(),
			FontData:    symbola,
			FontSize:    56,
			FontDx:      0,
			FontDy:      0,
			ShadowShift: shift,
			TileW:       64,
			TileH:       64,
		})
	}
	// scroll
	for i := 0; i < scrolltype.ScrollType_Count; i++ {
		sct := scrolltype.ScrollType(i)
		rtn = append(rtn, runetile.RuneTile{
			GroupName:   "Scroll",
			Rune:        sct.Rune(),
			Name:        sct.String(),
			Color1:      sct.Color24(),
			FontData:    symbola,
			FontSize:    56,
			FontDx:      0,
			FontDy:      0,
			ShadowShift: shift,
			TileW:       64,
			TileH:       64,
		})
	}
	// faction char
	for i := 0; i < factiontype.FactionType_Count; i++ {
		ft := factiontype.FactionType(i)
		rtn = append(rtn, []runetile.RuneTile{
			{GroupName: "Char_" + ft.Rune(), Rune: ft.Rune(), Name: "Char", Color1: ft.Color24(),
				FontData: letterFont, FontSize: 56,
				FontDx: 0, FontDy: 0, ShadowShift: shift, TileW: 64, TileH: 64,
			},
			{GroupName: "Char_" + ft.Rune(), Rune: ft.Rune(), Name: "Dead", Color1: htmlcolors.Gray,
				FontData: letterFont, FontSize: 56,
				FontDx: 0, FontDy: 0, ShadowShift: shift, TileW: 64, TileH: 64,
			},
			{GroupName: "Char_" + ft.Rune(), Rune: ft.Rune(), Name: "Rev", Color1: ft.Color24().Neg(),
				FontData: letterFont, FontSize: 56,
				FontDx: 0, FontDy: 0, ShadowShift: shift, TileW: 64, TileH: 64,
			},
		}...)
	}
	// equip items
	for eqi := 0; eqi < equipslottype.EquipSlotType_Count; eqi++ {
		eqt := equipslottype.EquipSlotType(eqi)
		dx := 0
		dy := -2
		if eqt == equipslottype.Gauntlet {
			dx = -12
		} else if eqt == equipslottype.Ring {
			dy = 0
			dx = -10
		}
		for fti := 0; fti < factiontype.FactionType_Count; fti++ {
			ft := factiontype.FactionType(fti)

			rtn = append(rtn, runetile.RuneTile{
				GroupName:   equipslottype.EquipSlotType(eqi).String(),
				Rune:        eqt.Rune(),
				Name:        factiontype.FactionType(fti).String(),
				Color1:      ft.Color24(),
				FontData:    symbola,
				FontSize:    64,
				FontDx:      dx,
				FontDy:      dy,
				ShadowShift: shift,
				TileW:       64,
				TileH:       64,
			})
		}
	}
	for i := range rtn {
		err := rtn[i].MakeImg()
		if err != nil {
			return nil, err
		}
	}
	return rtn, nil
}
