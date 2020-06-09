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

package clientinitdata

import (
	"bytes"
	"fmt"
	"net/url"
	"syscall/js"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/config/leveldata"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/game/aoscore"
	"github.com/kasworld/goguelike/game/towerlist4client"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/wrapspan"
)

type InitData struct {
	TowerList  []towerlist4client.TowerInfo2Enter
	TowerIndex int

	// set by login
	ServiceInfo       *c2t_obj.ServiceInfo
	AccountInfo       *c2t_obj.AccountInfo
	TowerInfo         *c2t_obj.TowerInfo
	ViewportXYLenList findnear.XYLenList
}

func New() *InitData {
	rtn := &InitData{
		ViewportXYLenList: viewportdata.ViewportXYLenList,
	}
	rtn.TowerList = RefreshTowerListHTML("towerlist")
	return rtn
}
func (idata *InitData) String() string {
	return fmt.Sprintf("WasmClient[%v %v]", idata.GetNickName(), idata.GetConnectToTowerURL())
}

func (idata *InitData) GetNickName() string {
	var nickname string
	if idata.AccountInfo != nil {
		nickname = idata.AccountInfo.NickName
	}
	return nickname
}

func (idata *InitData) GetTowerName() string {
	return idata.TowerInfo.Name
}

func (idata *InitData) GetConnectToTowerURL() string {
	if idata.TowerList == nil {
		jslog.Error("TowerList not loaded")
		return ""
	}
	return idata.TowerList[idata.TowerIndex].ConnectURL
}

func (idata *InitData) MakeServiceInfoHTML() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Goguelike<br/>")
	fmt.Fprintf(&buf, "Protocol %v<br/>", idata.ServiceInfo.ProtocolVersion)
	fmt.Fprintf(&buf, "Data %v<br/>", idata.ServiceInfo.DataVersion)
	fmt.Fprintf(&buf, "Version %v<br/>", idata.ServiceInfo.Version)
	fmt.Fprintf(&buf, "SessionID %v<br/>", idata.AccountInfo.SessionG2ID)
	fmt.Fprintf(&buf, "ActiveObjID %v<br/>", idata.AccountInfo.ActiveObjG2ID)
	fmt.Fprintf(&buf, "%v<br/>", MsgCopyright)
	return buf.String()
}

func (idata *InitData) CanUseCmd(cmd c2t_idcmd.CommandID) bool {
	if idata.AccountInfo == nil {
		return false
	}
	return idata.AccountInfo.CmdList[cmd]
}

func RefreshTowerListHTML(objid string) []towerlist4client.TowerInfo2Enter {
	jsobj := js.Global().Get("document").Call("getElementById", objid)
	tlurl := ReplacePathFromHref("towerlist.json")
	tl, err := towerlist4client.LoadFromURL(tlurl)
	if err != nil {
		jslog.Errorf("towerlist load fail %v\n", err)
		jsobj.Set("innerHTML", "fail to load towerlist")
		return nil
	}

	var buf bytes.Buffer
	for i, v := range tl {
		if v.ConnectURL != "" {
			fmt.Fprintf(&buf, "Entrance found, Tower %v ",
				v.Name,
			)
			fmt.Fprintf(&buf,
				`<button type="button" style="font-size:20px;" onclick="enterTower(%d)">Enter</button>`,
				i,
			)
		} else {
			fmt.Fprintf(&buf, "Entrance not found, Tower %v ", v.Name)
		}
		fmt.Fprintf(&buf,
			` <button type="button" style="font-size:20px;" onclick="clearSession(%d)">Clear Memory</button>`,
			i,
		)
		buf.WriteString("<br/>")
	}
	jsobj.Set("innerHTML", buf.String())
	return tl
}

func LoadHighScoreHTML() string {
	tlurl := ReplacePathFromHref("highscore.json")
	aol, err := aoscore.LoadFromURL(tlurl)
	if err != nil {
		jslog.Errorf("highscore load fail %v", err)
		return "fail to load high score"
	}
	var buf bytes.Buffer
	buf.WriteString(`
		Top 10 Player in all tower
		<table border=2>
		<tr>
		<th>Ranking</th> <th>NickName</th>  <th>Level</th>  <th>Exp</th> <th>Wealth</th> <th>Born</th> <th>Current</th>
		</tr>	
		`)
	for i, ao := range aol {
		fmt.Fprintf(&buf, `
		<tr>
			<td>%v</td> <td>%v</td> <td>%.2f</td> <td>%.2f</td> <td>%.2f</td> <td>%v</td> <td>%v</td>
		</tr>`,
			i+1,
			ao.NickName,
			leveldata.CalcLevelFromExp(ao.Exp),
			ao.Exp,
			ao.Wealth,
			wrapspan.THCSText(ao.BornFaction.Color24(), ao.BornFaction.String()),
			wrapspan.THCSText(ao.CurrentBias, ao.CurrentBias.NearFaction().String()),
		)
	}
	buf.WriteString(`
		<tr>
		<th>Ranking</th> <th>NickName</th> <th>Level</th> <th>Exp</th> <th>Wealth</th> <th>Born</th> <th>Current</th>
		</tr>
		</table>
		`)
	return buf.String()
}

func ReplacePathFromHref(s string) string {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	if err != nil {
		jslog.Errorf("%v", err)
		return ""
	}
	u.Path = s
	return u.String()
}
