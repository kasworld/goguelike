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
	"syscall/js"
	"time"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/potiontype"
	"github.com/kasworld/goguelike/enum/scrolltype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/lib/jsobj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/wrapspan"
	"github.com/kasworld/ifop"
)

func (app *WasmClient) makeBaseInfoHTML() string {
	if gInitData.TowerInfo == nil {
		return "Wait to enter tower"
	}
	var buf bytes.Buffer
	cf := app.currentFloor()

	towerBias := app.TowerBias()
	fmt.Fprintf(&buf,
		`Tower <span class="tooltip" style="color:%s">%s %v
			<span class="tooltiptext-bottom">%v</span>
		</span><br/>`,
		towerBias.ToHTMLColorString(),
		gInitData.TowerInfo.Name, towerBias.NearFaction(),
		towerBias.IntString(),
	)
	floorBias := app.FloorInfo.Bias
	fmt.Fprintf(&buf,
		`Floor <span class="tooltip" style="color:%s">%s %v
			<span class="tooltiptext-bottom">%v</span>
		</span><br/>`,
		floorBias.ToHTMLColorString(),
		app.FloorInfo.Name, floorBias.NearFaction(),
		floorBias.IntString(),
	)

	pao := app.olNotiData.ActiveObj
	envBias := app.GetEnvBias()
	lv := leveldata.CalcLevelFromExp(float64(pao.Exp))
	app.level = int(lv)

	effBias := pao.Bias.Add(envBias)
	if effBias.AbsSum() < pao.Bias.AbsSum() {
		buf.WriteString(wrapspan.ColorText("OrangeRed", "Env bad "))
	} else {
		buf.WriteString(wrapspan.ColorText("LimeGreen", "Env good "))
	}
	if effBias.AbsSum() < app.lastEffBias.AbsSum() {
		buf.WriteString(wrapspan.ColorText("OrangeRed", "to bad "))
	} else {
		buf.WriteString(wrapspan.ColorText("LimeGreen", "to good "))
	}

	fmt.Fprintf(&buf,
		`<span class="tooltip" style="color:%s">%s
			<span class="tooltiptext-bottom">%v</span>
		</span><br/>`,
		envBias.ToHTMLColorString(), envBias.NearFaction(), envBias.IntString(),
	)
	app.lastEffBias = effBias

	playerX, playerY := app.GetPlayerXY()
	atkStr := ""
	defStr := ""
	tlHp := 0.0
	tlAp := 0.0
	if cf != nil && cf.IsValidPos(playerX, playerY) {
		tl := cf.Tiles[playerX][playerY]
		co := "white"
		if tl.Safe() {
			co = "LightGreen"
		}
		act := c2t_idcmd.Meditate
		dir := way9type.Center
		if pao.Act != nil && pao.Act.Acted() {
			act = pao.Act.Done.Act
			dir = pao.Act.Done.Dir
		}
		tlHp, tlAp = tl.ActHPSPCalced(act, dir)
		if tl.CanBattle() {
			atk := tl.AtkMod() * 100
			def := tl.DefMod() * 100
			co = "red"
			atkStr = wrapspan.ColorTextf(co, " x %.0f%%", atk)
			defStr = wrapspan.ColorTextf(co, " x %.0f%%", def)
		}
		buf.WriteString(wrapspan.ColorTextf(co,
			"[%v %v] %v<br/>",
			playerX, playerY, tl.Name()))
	} else {
		buf.WriteString(wrapspan.ColorTextf("red", "[%v %v]<br/>", playerX, playerY))
		jslog.Error("ao pos out of floor %v %v %v", cf, playerX, playerY)
	}

	fmt.Fprintf(&buf,
		`%s <span class="tooltip" style="color:%s">%s
			<span class="tooltiptext-bottom">%v</span>
		</span><br/>`,
		gInitData.AccountInfo.NickName,
		pao.Bias.ToHTMLColorString(),
		pao.Bias.NearFaction(),
		pao.Bias.IntString(),
	)

	fmt.Fprintf(&buf, "Level %4.2f<br/>", lv)
	fmt.Fprintf(&buf, "Exp %v<br/>", pao.Exp)
	fmt.Fprintf(&buf, "Rank %v<br/>", pao.Ranking+1)

	fmt.Fprintf(&buf, "%s %v/%v", wrapspan.ColorText("Red", "HP"), pao.HP, pao.HPMax)
	if tlHp != 0 {
		buf.WriteString(wrapspan.ColorTextf("red", " -%.0f", tlHp))
	}
	buf.WriteString("<br/>")

	fmt.Fprintf(&buf, "%s  %v/%v", wrapspan.ColorTextf("Yellow", "SP"), pao.SP, pao.SPMax)
	if tlAp != 0 {
		buf.WriteString(wrapspan.ColorTextf("red", " -%.0f", tlAp))
	}
	buf.WriteString("<br/>")

	fmt.Fprintf(&buf, "%s  %.2f", wrapspan.ColorTextf("Lime", "Delay"), pao.RemainTurn2Act)
	buf.WriteString("<br/>")

	if autoActs.GetByIDBase("AutoPlay").State == 0 {
		buf.WriteString(wrapspan.ColorTextf(
			pao.AIPlan.Color24().ToHTMLColorString(), "AI %v", pao.AIPlan))
		buf.WriteString("<br/>")
	} else {
		buf.WriteString(wrapspan.ColorTextf(
			app.ClientColtrolMode.Color24().ToHTMLColorString(), "Control %v", app.ClientColtrolMode))
		buf.WriteString("<br/>")
	}

	atk := app.getAttackBias().Add(envBias)
	fmt.Fprintf(&buf,
		`<span class="tooltip" style="color:%s">Atk %4.1f
			<span class="tooltiptext-bottom">%v</span>
		</span>`,
		"white", atk.AbsSum()/3+lv, atk.IntString(),
	)
	fmt.Fprintf(&buf, "%v<br/>", atkStr)
	def := app.getDefenceBias().Add(envBias)
	fmt.Fprintf(&buf,
		`<span class="tooltip" style="color:%s">Def %4.1f
		<span class="tooltiptext-bottom">%v</span>
		</span>`,
		"white", def.AbsSum()/3+lv, def.IntString(),
	)
	fmt.Fprintf(&buf, "%v<br/>", defStr)

	if pao.Sight == leveldata.Sight(int(lv)) {
		fmt.Fprintf(&buf, "Sight %v", pao.Sight)
	} else if pao.Sight >= gameconst.SightXray {
		buf.WriteString(wrapspan.ColorText("Yellow", "X-ray sight"))
	} else if pao.Sight >= leveldata.Sight(int(lv)) {
		buf.WriteString(wrapspan.ColorTextf("Lime", "Sight %v", pao.Sight))
	} else {
		buf.WriteString(wrapspan.ColorTextf("Red", "Sight %v", pao.Sight))
	}
	buf.WriteString("<br/>")
	fmt.Fprintf(&buf, "Kill %v Death %v<br/>", pao.Kill, pao.Death)

	ldco := ifop.TrueString(app.OverLoadRate > 1, "Red", "white")
	wLimit := leveldata.WeightLimit(int(lv))
	if pao.Conditions.TestByCondition(condition.Burden) {
		wLimit /= 2
	}
	buf.WriteString(wrapspan.ColorTextf(ldco,
		"Carry %.0f/%v<br/>", pao.CalcWeight(), wLimit))
	fmt.Fprintf(&buf, "Wealth %v<br/>", makeMoneyColor(pao.Wealth))
	fmt.Fprintf(&buf, "Equip %v Bag %v<br/>", len(pao.EquippedPo), len(pao.EquipBag))
	fmt.Fprintf(&buf, "Potion %v Scroll %v<br/>", len(pao.PotionBag), len(pao.ScrollBag))
	fmt.Fprintf(&buf, "Wallet %v<br/>", makeMoneyColor(pao.Wallet))
	return buf.String()
}

func (app *WasmClient) makeBuffListHTML() string {
	var buf bytes.Buffer
	pao := app.olNotiData.ActiveObj
	fmt.Fprintf(&buf, "Active Buff %v<br/>", len(pao.ActiveBuff))
	for i, v := range pao.ActiveBuff {
		if i > DisplayLineLimit-20 { // limit buff display
			break
		}
		fmt.Fprintf(&buf, "%v %v<br/>",
			v.RemainCount,
			v.Name,
		)
	}
	return buf.String()
}

func (app *WasmClient) getAttackBias() bias.Bias {
	rtn := app.olNotiData.ActiveObj.Bias
	for _, po := range app.olNotiData.ActiveObj.EquippedPo {
		if po.EquipType.Attack() {
			rtn = rtn.Add(po.GetBias())
		}
	}
	return rtn
}

func (app *WasmClient) getDefenceBias() bias.Bias {
	rtn := app.olNotiData.ActiveObj.Bias
	for _, po := range app.olNotiData.ActiveObj.EquippedPo {
		if po.EquipType.Defence() {
			rtn = rtn.Add(po.GetBias())
		}
	}
	return rtn
}

var makeRecycleButton = `<button style="font-size: %vpx" onclick="recycle('%s')">Recycle</button> `
var makeUnequipButton = `<button style="font-size: %vpx" onclick="unequip('%s')">Unequip</button> `
var makeEquipButton = `<button style="font-size: %vpx" onclick="equip('%s')">Equip</button> `
var makeDropButton = `<button style="font-size: %vpx" onclick="drop('%s')">Drop</button> `
var makeDrinkPotionButton = `<button style="font-size: %vpx" onclick="drinkpotion('%s')">DrinkPotion</button> `
var makeReadScrollButton = `<button style="font-size: %vpx" onclick="readscroll('%s')">ReadScroll</button> `

func (app *WasmClient) makeInvenInfoHTML() string {

	var buf bytes.Buffer
	pao := app.olNotiData.ActiveObj
	win := js.Global().Get("window")
	winH := win.Get("innerHeight").Float()
	ftSize := winH / 100
	canRecycle := false
	if app.onFieldObj != nil && app.onFieldObj.ActType == fieldobjacttype.RecycleCarryObj {
		canRecycle = true
	}
	displayedLine := 4 // text not in loop

	potionType2info := make([]struct {
		UUID  string
		Count int
	}, potiontype.PotionType_Count)
	for _, v := range pao.PotionBag {
		potionType2info[v.PotionType].UUID = v.UUID
		potionType2info[v.PotionType].Count++
	}
	fmt.Fprintf(&buf, "Potion %v<br/>", len(pao.PotionBag))
	for i, v := range potionType2info {
		if v.Count == 0 {
			continue
		}
		if displayedLine > DisplayLineLimit {
			break
		}
		pt := potiontype.PotionType(i)
		displayedLine++
		poStr := wrapspan.THCSTextf(pt.Color24(),
			"%v %v(%v)", pt.String(), pt.Rune(), v.Count)
		buf.WriteString(poStr)
		if canRecycle {
			fmt.Fprintf(&buf, makeRecycleButton, ftSize, v.UUID)
		}
		fmt.Fprintf(&buf, makeDrinkPotionButton, ftSize, v.UUID)
		fmt.Fprintf(&buf, makeDropButton, ftSize, v.UUID)
		buf.WriteString("<br/>")
	}
	scrollType2info := make([]struct {
		UUID  string
		Count int
	}, scrolltype.ScrollType_Count)
	for _, v := range pao.ScrollBag {
		scrollType2info[v.ScrollType].UUID = v.UUID
		scrollType2info[v.ScrollType].Count++
	}

	fmt.Fprintf(&buf, "Scroll %v<br/>", len(pao.ScrollBag))
	for i, v := range scrollType2info {
		if v.Count == 0 {
			continue
		}
		if displayedLine > DisplayLineLimit {
			break
		}
		st := scrolltype.ScrollType(i)
		displayedLine++
		poStr := wrapspan.THCSTextf(st.Color24(),
			"%v %v(%v)", st.String(), st.Rune(), v.Count)
		buf.WriteString(poStr)
		if canRecycle {
			fmt.Fprintf(&buf, makeRecycleButton, ftSize, v.UUID)
		}
		fmt.Fprintf(&buf, makeReadScrollButton, ftSize, v.UUID)
		fmt.Fprintf(&buf, makeDropButton, ftSize, v.UUID)
		buf.WriteString("<br/>")
	}

	fmt.Fprintf(&buf, "Equip %v<br/>", len(pao.EquippedPo))
	for _, v := range pao.EquippedPo {
		if displayedLine > DisplayLineLimit {
			break
		}
		displayedLine++
		fmt.Fprintf(&buf, "%s ", v.Name)
		poStr := wrapspan.THCSTextf(v.GetBias(), "%v%v%.0f",
			v.EquipType.Rune(), v.Faction.Rune(), v.BiasLen)
		buf.WriteString(poStr)
		if canRecycle {
			fmt.Fprintf(&buf, makeRecycleButton, ftSize, v.UUID)
		}
		fmt.Fprintf(&buf, makeUnequipButton, ftSize, v.UUID)
		fmt.Fprintf(&buf, makeDropButton, ftSize, v.UUID)
		buf.WriteString("<br/>")
	}

	fmt.Fprintf(&buf, "Equip Bag %v<br/>", len(pao.EquipBag))
	for _, v := range pao.EquipBag {
		if displayedLine > DisplayLineLimit {
			break
		}
		displayedLine++
		fmt.Fprintf(&buf, "%s ", v.Name)
		poStr := wrapspan.THCSTextf(v.GetBias(), "%v%v%.0f",
			v.EquipType.Rune(), v.Faction.Rune(), v.BiasLen)
		buf.WriteString(poStr)
		if canRecycle {
			fmt.Fprintf(&buf, makeRecycleButton, ftSize, v.UUID)
		}
		fmt.Fprintf(&buf, makeEquipButton, ftSize, v.UUID)
		fmt.Fprintf(&buf, makeDropButton, ftSize, v.UUID)
		buf.WriteString("<br/>")
	}
	return buf.String()
}

func (app *WasmClient) makeFieldObjListHTML() string {
	var buf bytes.Buffer
	cf := app.currentFloor()
	objList := cf.FieldObjPosMan.GetAllList()
	foList := make([]*c2t_obj.FieldObjClient, len(objList))
	for i, v := range objList {
		foList[i] = v.(*c2t_obj.FieldObjClient)
	}
	c2t_obj.FieldObjByType(foList).Sort()
	fmt.Fprintf(&buf, "Field object found %v<br/>", len(foList))
	for i, fo := range foList {
		fmt.Fprintf(&buf, "%v [%v %v]<br/>", fo.Message, fo.X, fo.Y)
		if i > DisplayLineLimit {
			break // limit display line
		}
	}
	return buf.String()
}

func (app *WasmClient) makeFloorListHTML() string {
	var buf bytes.Buffer
	if len(app.UUID2ClientFloor) == gInitData.TowerInfo.TotalFloorNum {
		fmt.Fprintf(&buf, "Found All Floor %v<br/>",
			len(app.UUID2ClientFloor))
	} else {
		fmt.Fprintf(&buf, "Floor Found %v<br/>",
			len(app.UUID2ClientFloor))
	}
	cfList := make(ClientFloorGLList, 0)
	for _, v := range app.UUID2ClientFloor {
		cfList = append(cfList, v)
	}
	cfList.Sort()

	win := js.Global().Get("window")
	winH := win.Get("innerHeight").Float()
	ftSize := winH / 100
	for i, cf := range cfList {
		floorStr := wrapspan.THCSTextf(cf.FloorInfo.Bias, "%v", cf.FloorInfo.Bias.NearFaction().String())
		fmt.Fprintf(&buf, "%v %v ", floorStr, cf.Visited)
		if cf.Visited.CalcCompleteRate() >= 1.0 {
			fmt.Fprintf(&buf,
				`<button style="font-size: %vpx" onclick="moveFloor('%s')">Teleport</button>`,
				ftSize, cf.FloorInfo.UUID,
			)
		} else {
			fmt.Fprintf(&buf,
				`<button style="font-size: %vpx" onclick="moveFloor('%s')" disabled>Teleport</button>`,
				ftSize, cf.FloorInfo.UUID,
			)
		}
		buf.WriteString("<br/>")
		if i > DisplayLineLimit {
			break // limit display line
		}
	}
	return buf.String()
}

func (app *WasmClient) makeDebugInfoHTML() string {
	var buf bytes.Buffer
	fd := app.olNotiData
	fmt.Fprintf(&buf, "%v<br/>", app.taNotiHeader)
	fmt.Fprintf(&buf, "%v<br/>", app.olNotiHeader)
	fmt.Fprintf(&buf, "Tower Factor%v<br/>", gInitData.TowerInfo.Factor)
	fmt.Fprintf(&buf, "Server Run %v<br/>", fd.Time.Sub(gInitData.TowerInfo.StartTime))
	fmt.Fprintf(&buf, "Ping %v<br/>", app.PingDur)
	fmt.Fprintf(&buf, "TimeDiff %v<br/>", app.ServerClientTimeDiff)
	fmt.Fprintf(&buf, "ClientTime %v<br/>", time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"))
	fmt.Fprintf(&buf, "ServerTime %v<br/>", app.ServerTime().Format("2006-01-02T15:04:05.999999999Z07:00"))
	fmt.Fprintf(&buf, "%v<br/>", app.ClientJitter)
	fmt.Fprintf(&buf, "%v<br/>", app.ServerJitter)
	fmt.Fprintf(&buf, "%v<br/>", app.DispInterDur.GetInterval())
	fmt.Fprintf(&buf, "floor %.1f actturn/sec display %.1fFPS<br/>",
		app.FloorInfo.TurnPerSec,
		1.0/app.DispInterDur.GetInterval().GetLastDuration().Seconds())
	fmt.Fprintf(&buf, "%v<br/>", app.DispInterDur.GetDuration())
	fmt.Fprintf(&buf, "Display %d<br/>", app.DispInterDur.GetCount())
	fmt.Fprintf(&buf, "Known ActiveObj %d<br/>", len(app.AOUUID2AOClient))
	fmt.Fprintf(&buf, "Known CarryObj %d<br/>", len(app.CaObjUUID2CaObjClient))
	fmt.Fprintf(&buf, "Sent move packet lastTurn %d<br/>", app.movePacketPerTurn)
	fmt.Fprintf(&buf, "Sent act packet lastTurn %d<br/>", app.actPacketPerTurn)
	return buf.String()
}

func (app *WasmClient) DisplayTextInfo() {
	if app.onFieldObj != nil &&
		(app.onFieldObj.ActType == fieldobjacttype.PortalIn ||
			app.onFieldObj.ActType == fieldobjacttype.PortalInOut) {
		commandButtons.GetByIDBase("EnterPortal").Enable()
	} else {
		commandButtons.GetByIDBase("EnterPortal").Disable()
	}
	if cf := app.currentFloor(); cf.Visited.CalcCompleteRate() >= 1.0 {
		commandButtons.GetByIDBase("Teleport").Enable()
	} else {
		commandButtons.GetByIDBase("Teleport").Disable()
	}
	if dt := app.olNotiData; dt != nil && dt.ActiveObj.HP <= 0 {
		commandButtons.GetByIDBase("Rebirth").Enable()
		commandButtons.GetByIDBase("KillSelf").Disable()
	} else {
		commandButtons.GetByIDBase("Rebirth").Disable()
		commandButtons.GetByIDBase("KillSelf").Enable()
	}

	envColor := app.GetEnvBias().ToHTMLColorString()
	jsobj.SetBGColor(GetElementById("body"), envColor)

	app.updateLeftInfo()
	app.updateRightInfo()
	app.updateCenterInfo()
}

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
