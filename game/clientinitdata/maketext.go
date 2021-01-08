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

package clientinitdata

import (
	"bytes"
	"fmt"

	"github.com/kasworld/goguelike/config/dataversion"
	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_version"
	"github.com/kasworld/htmlcolors"
	"github.com/kasworld/ifop"
	"github.com/kasworld/version"
)

var MsgCopyright = `Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang 
		<a href="https://kasw.blogspot.com/" target="_blank">Goguelike</a>`

func MakeClientInfoHTML() string {
	var buf bytes.Buffer
	buf.WriteString(`<table style="text-align: left;" border=2><tr><td>`)
	fmt.Fprintf(&buf, "Goguelike webclient<br/>")
	fmt.Fprintf(&buf, "Protocol %v<br/>", c2t_version.ProtocolVersion)
	fmt.Fprintf(&buf, "Data %v<br/>", dataversion.DataVersion)
	fmt.Fprintf(&buf, "Build %v<br/>", version.GetVersion())
	fmt.Fprintf(&buf, "Build %v<br/>", version.GetBuildDate())
	fmt.Fprintf(&buf, "%v<br/>", MsgCopyright)
	buf.WriteString(`</td></tr></table>`)
	return buf.String()
}

func MakeHelpFactionHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Faction Info<br/>`)
	buf.WriteString(`<table style="text-align: left;" border=2>`)
	buf.WriteString(`<tr ><th>Faction</th><th>Rune</th><th>Value</th></tr>`)
	for i := 0; i < factiontype.FactionType_Count; i++ {
		ft := factiontype.FactionType(i)
		bs := bias.NewByFaction(ft, 100)
		fmt.Fprintf(&buf, `<tr style="color:%v; background-color:%v"> <td>%v</td> <td>%v</td> <td>%v</td> </tr>`,
			bs.ToColor24().Rev().ToHTMLColorString(), bs.ToColor24().ToHTMLColorString(),
			ft.String(),
			ft.Rune(),
			bs.IntString(),
		)
	}
	buf.WriteString(`<tr><th>Faction</th><th>Rune</th><th>Value</th></tr>`)
	buf.WriteString(`</table>`)
	return buf.String()
}

func MakeHelpCarryObjectHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`CarryObject Info<br/>`)
	buf.WriteString(`<table style="text-align: left;" border=2>`)
	buf.WriteString(`<tr><th>Rune</th><th>Name</th><th>Battle</th></tr>`)
	for i := 0; i < equipslottype.EquipSlotType_Count; i++ {
		eqt := equipslottype.EquipSlotType(i)
		eqbuf := ifop.TrueString(eqt.Attack(), "Attack", "")
		eqbuf += ifop.TrueString(eqt.Defence(), "Defence", "")
		fmt.Fprintf(&buf, "<tr> <td>%v</td> <td>%v</td> <td>%v</td> </tr>",
			eqt.Rune(), eqt.String(), eqbuf)
	}
	fmt.Fprintf(&buf, "<tr> <td> </td> <td>Potion</td> <td>DrinkPotion</td> </tr>")
	fmt.Fprintf(&buf, "<tr> <td> </td> <td>Scroll</td> <td>ReadScroll</td> </tr>")
	buf.WriteString(`<tr><th>Rune</th><th>Name</th><th>Battle</th></tr>`)
	buf.WriteString(`</table>`)

	return buf.String()
}

func MakeHelpPotionHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Potion Info<br/>`)
	buf.WriteString(`<table style="text-align: left;" border=2>`)
	buf.WriteString(`<tr ><th>Name</th><th>Rune</th><th>description</th></tr>`)
	for i := 0; i < potiontype.PotionType_Count; i++ {
		pt := potiontype.PotionType(i)
		fmt.Fprintf(&buf, `<tr style="color:%v; background-color:%v"><td>%v</td><td>%v</td><td>%v</td> </tr>`,
			pt.Color24().Rev().ToHTMLColorString(), pt.Color24().ToHTMLColorString(),
			pt.String(),
			pt.Rune(),
			pt.CommentString())
	}
	buf.WriteString(`<tr><th>Name</th><th>Rune</th><th>description</th></tr>`)
	buf.WriteString(`</table>`)

	return buf.String()
}

func MakeHelpScrollHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Scroll Info<br/>`)
	buf.WriteString(`<table style="text-align: left;" border=2>`)
	buf.WriteString(`<tr><th>Name</th><th>Rune</th><th>description</th></tr>`)
	for i := 0; i < scrolltype.ScrollType_Count; i++ {
		pt := scrolltype.ScrollType(i)
		fmt.Fprintf(&buf, `<tr style="color:%v; background-color:%v"><td>%v</td><td>%v</td><td>%v</td> </tr>`,
			pt.Color24().Rev().ToHTMLColorString(), pt.Color24().ToHTMLColorString(),
			pt.String(),
			pt.Rune(),
			pt.CommentString())
	}
	buf.WriteString(`<tr><th>Name</th><th>Rune</th><th>description</th></tr>`)
	buf.WriteString(`</table>`)

	return buf.String()
}

func MakeHelpMoneyColorHTML() string {
	var buf bytes.Buffer
	buf.WriteString(`Money color Info<br/>`)
	buf.WriteString(`<table style="text-align: left;" border=2>`)
	buf.WriteString(`<tr><th>Range</th><th>Color</th></tr>`)
	for _, v := range moneycolor.Attrib {
		fmt.Fprintf(&buf, `<tr><td>~$%v</td><td style="background-color:%v"></td> </tr>`,
			v.UpLimit,
			v.Color)
	}
	buf.WriteString(`<tr><th>Range</th><th>Color</th></tr>`)
	buf.WriteString(`</table>`)
	return buf.String()
}

func MakeHelpTileHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Floor Tile Info<br/>`)
	buf.WriteString(`<table style="text-align: left;" border=2>`)
	buf.WriteString(`<tr ><th>Tile</th><th>Attack mod %</th><th>Defence mod %</th></tr>`)
	for i := 0; i < tile.Tile_Count; i++ {
		ti := tile.Tile(i)
		co := ti.Attrib().Color
		co24 := htmlcolors.NewColor24(co[0], co[1], co[2])
		if ti.Attrib().ObjOnTile != tile.ObjOnTile_1 {
			fmt.Fprintf(&buf, `<tr style="color:%v; background-color:%v">
			<td>%v</td> <td>No Battle</td> <td>No Battle</td> </tr>`,
				co24.Rev().ToHTMLColorString(), co24.ToHTMLColorString(),
				ti.String())
		} else {
			fmt.Fprintf(&buf, `<tr style="color:%v; background-color:%v">
			<td>%v</td> <td>%3.0f%%</td> <td>%3.0f%%</td> </tr>`,
				co24.Rev().ToHTMLColorString(), co24.ToHTMLColorString(),
				ti.String(), ti.Attrib().AtkMod*100, ti.Attrib().DefMod*100)
		}
	}
	buf.WriteString(`<tr><th>Tile</th><th>Attack mod %</th><th>Defence mod %</th></tr>`)
	buf.WriteString(`</table>`)

	return buf.String()
}

func MakeHelpConditionHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Condition Info<br/>`)
	buf.WriteString(`<table style="text-align: left;" border=2>`)
	buf.WriteString(`<tr><th>Condition</th><th>description</th></tr>`)
	for i := 0; i < condition.Condition_Count; i++ {
		ti := condition.Condition(i)
		fmt.Fprintf(&buf, "<tr> <td>%v</td> <td>%v</td> </tr>",
			ti.String(), ti.CommentString())
	}
	buf.WriteString(`<tr><th>Condition</th><th>description</th></tr>`)
	buf.WriteString(`</table>`)

	return buf.String()
}

func MakeHelpFieldObjHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`FieldObject Info<br/>`)
	buf.WriteString(`<table style="text-align: left;" border=2>`)
	buf.WriteString(`<tr><th>Field object</th><th>Auto trigger</th><th>description</th></tr>`)
	for i := 0; i < fieldobjacttype.FieldObjActType_Count; i++ {
		if fieldobjacttype.ClientData[i].Text == "" {
			continue
		}
		ti := fieldobjacttype.FieldObjActType(i)
		fmt.Fprintf(&buf, `<tr style="color:%v; background-color:%v"> 
		<td>%v</td> <td>%v %v%%</td> <td>%v</td> </tr>`,
			ti.Color24().Rev().ToHTMLColorString(), ti.Color24().ToHTMLColorString(),
			ti.String(),
			ti.AutoTrigger(), ti.TriggerRate()*100,
			fieldobjacttype.ClientData[i].Text,
		)
	}
	buf.WriteString(`<tr><th>Field object</th><th>Auto trigger</th><th>description</th></tr>`)
	buf.WriteString(`</table>`)

	return buf.String()
}
