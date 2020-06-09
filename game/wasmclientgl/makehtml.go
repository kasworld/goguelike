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

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/wrapspan"
)

func MakeHelpInfoHTML() string {
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
