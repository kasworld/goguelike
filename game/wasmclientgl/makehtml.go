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

package wasmclientgl

import (
	"bytes"
	"fmt"

	"github.com/kasworld/htmlcolors"

	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/tile"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_version"
	"github.com/kasworld/gowasmlib/wrapspan"
	"github.com/kasworld/ifop"
	"github.com/kasworld/version"
)

func makeHelpInfoHTML() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, `Environment = Tower + Floor<br/>
	Atk = Env + Bias + attack equip + lv<br/>
	Def = Env + Bias + defence equip + lv<br/>
	OverLoaded, Penalty at HP/SP gain per Turn<br/>
	Item weight and value are proportional to power.<br/>
	When you die, item or money of same value dropped.<br/>
	You can teleport in floor or move to floor which discovered completely.<br/>
	When level up HPMax %v SPMax %v Sight %v Load %vg increased.<br/>`,
		gameconst.HPPerLevel, gameconst.SPPerLevel, gameconst.SightPerLevel, gameconst.LvGram,
	)

	help := []string{
		"",
		"Keypad : 8 way move",
		"Mouse left : move to pathfind pos",
		"Mouse right : follow mouse on/off",
		"",
		"Keyboard shortcut",
	}
	for _, v := range help {
		buf.WriteString(v)
		buf.WriteString("<br/>")
	}
	fmt.Fprintf(&buf, "%v<br/>", commandButtons.Name)
	for _, v := range commandButtons.ButtonList {
		fmt.Fprintf(&buf, "%v : %v<br/>", v.KeyCode, v.ToolTip)
	}

	fmt.Fprintf(&buf, "%v<br/>", autoActs.Name)
	for _, v := range autoActs.ButtonList {
		fmt.Fprintf(&buf, "%v : %v<br/>", v.KeyCode, v.ToolTip)
	}

	fmt.Fprintf(&buf, "%v<br/>", gameOptions.Name)
	for _, v := range gameOptions.ButtonList {
		fmt.Fprintf(&buf, "%v : %v<br/>", v.KeyCode, v.ToolTip)
	}

	fmt.Fprintf(&buf, "%v can set by url arg<br/>", gameOptions.Name)
	for _, v := range gameOptions.ButtonList {
		fmt.Fprintf(&buf, "%v=", v.IDBase)
		for _, w := range v.ButtonText {
			fmt.Fprintf(&buf, "%v ", w)
		}
		buf.WriteString("<br/>")
	}

	fmt.Fprintf(&buf, "%v<br/>", adminCommandButtons.Name)
	for i, v := range adminCommandButtons.ButtonList {
		if gInitData.CanUseCmd(adminCmds[i]) {
			buf.WriteString(wrapspan.ColorTextf("red",
				"%v : %v", v.KeyCode, v.ToolTip),
			)
			buf.WriteString("<br/>")
		}
	}
	return buf.String()
}

func bias2NameIntStr(b bias.Bias) string {
	return fmt.Sprintf("%v %v", b.NearFaction().String(), b.IntString())
}

func obj2ColorStr(v interface{}) string {
	if v == nil {
		return ""
	}
	switch o := v.(type) {
	default:
		return fmt.Sprintf("unknonw type %v", v)

	case *c2t_obj.FieldObjClient:
		return wrapspan.THCSTextf(o.ActType.Color24(),
			"%v", o.ActType.String())
	case *c2t_obj.ActiveObjClient:
		return wrapspan.THCSTextf(o.Faction.Color24(),
			"%v", o.NickName)
	case *c2t_obj.PlayerActiveObjInfo:
		return wrapspan.THCSText(o.Bias, bias2NameIntStr(o.Bias))

	case *c2t_obj.CarryObjClientOnFloor:
		switch o.CarryingObjectType {
		default:
			return fmt.Sprintf("unknonw po %#v", o)
		case carryingobjecttype.Equip:
			return wrapspan.THCSTextf(o.Faction.Color24(),
				"%v%v", o.EquipType.Rune(), o.Faction.Rune())
		case carryingobjecttype.Money:
			return makeMoneyColor(o.Value)
		case carryingobjecttype.Potion:
			return wrapspan.THCSTextf(o.PotionType.Color24(),
				"%v%v", o.PotionType.Rune(), o.PotionType.String())
		case carryingobjecttype.Scroll:
			return wrapspan.THCSTextf(o.ScrollType.Color24(),
				"%v%v", o.ScrollType.Rune(), o.ScrollType.String())
		}

	case *c2t_obj.EquipClient:
		return wrapspan.THCSTextf(o.GetBias(),
			"%v%v%.0f", o.EquipType.Rune(), o.Faction.Rune(), o.BiasLen)
	case *c2t_obj.PotionClient:
		return wrapspan.THCSTextf(o.PotionType.Color24(),
			"%v%v", o.PotionType.Rune(), o.PotionType.String())
	case *c2t_obj.ScrollClient:
		return wrapspan.THCSTextf(o.ScrollType.Color24(),
			"%v%v", o.ScrollType.Rune(), o.ScrollType.String())
	}
}

func makeMoneyColor(mo int) string {
	for _, v := range moneycolor.Attrib {
		if mo < v.UpLimit {
			return wrapspan.ColorTextf(v.Color.String(), "$%v", mo)
		}
	}
	v := moneycolor.Attrib[len(moneycolor.Attrib)-1]
	return wrapspan.ColorTextf(v.Color.String(), "$%v", mo)
}

func makeClientInfoHTML() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Goguelike webclient<br/>")
	fmt.Fprintf(&buf, "Protocol %v<br/>", c2t_version.ProtocolVersion)
	fmt.Fprintf(&buf, "Data %v<br/>", gameconst.DataVersion)
	fmt.Fprintf(&buf, "Build %v<br/>", version.GetVersion())
	fmt.Fprintf(&buf, "Build %v<br/>", version.GetBuildDate())
	fmt.Fprintf(&buf, "%v<br/>", msgCopyright)
	return buf.String()
}

func makeHelpFactionHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Faction Info<br/>`)
	buf.WriteString(`<table border=2>`)
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

func makeHelpCarryObjectHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`CarryObject Info<br/>`)
	buf.WriteString(`<table border=2>`)
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

func makeHelpPotionHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Potion Info<br/>`)
	buf.WriteString(`<table border=2>`)
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

func makeHelpScrollHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Scroll Info<br/>`)
	buf.WriteString(`<table border=2>`)
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

func makeHelpMoneyColorHTML() string {
	var buf bytes.Buffer
	buf.WriteString(`Money color Info<br/>`)
	buf.WriteString(`<table border=2>`)
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

func makeHelpTileHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Floor Tile Info<br/>`)
	buf.WriteString(`<table border=2>`)
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

func makeHelpConditionHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`Condition Info<br/>`)
	buf.WriteString(`<table border=2>`)
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

func makeHelpFieldObjHTML() string {
	var buf bytes.Buffer

	buf.WriteString(`FieldObject Info<br/>`)
	buf.WriteString(`<table border=2>`)
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
