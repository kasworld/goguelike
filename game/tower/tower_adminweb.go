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

package tower

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/towerachieve_vector"
	"github.com/kasworld/goguelike/game/activeobject"
	"github.com/kasworld/goguelike/game/aoid2activeobject"
	"github.com/kasworld/goguelike/game/carryingobject"
	"github.com/kasworld/goguelike/game/dangerobject"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/goguelike/lib/sessionmanager"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_connbytemanager"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_statapierror"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_statnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_statserveapi"
	"github.com/kasworld/version"
	"github.com/kasworld/weblib"
	"github.com/kasworld/weblib/webprofile"
	"github.com/kasworld/wrapper"
)

func (tw *Tower) web_FaviconIco(w http.ResponseWriter, r *http.Request) {
}

func (tw *Tower) initAdminWeb() {
	authdata := weblib.NewAuthData("tower")
	authdata.ReLoadUserData([][2]string{
		{tw.sconfig.WebAdminID, tw.sconfig.WebAdminPass},
	})
	webMux := weblib.NewAuthMux(authdata, tw.log)

	if !version.IsRelease() {
		webprofile.AddWebProfile(webMux)
	}

	webMux.HandleFunc("/favicon.ico", tw.web_FaviconIco)

	webMux.HandleFunc("/service", tw.web_Service)

	webMux.HandleFuncAuth("/", tw.web_TowerInfo)
	webMux.HandleFuncAuth("/floor", tw.web_FloorInfo)
	webMux.HandleFuncAuth("/floorimagezoom", tw.web_FloorImageZoom)
	webMux.HandleFuncAuth("/floorimageautozoom", tw.web_FloorImageAutoZoom)
	webMux.HandleFuncAuth("/floortile", tw.web_FloorTile)

	webMux.HandleFuncAuth("/ActiveObj", tw.web_ActiveObjInfo)
	webMux.HandleFuncAuth("/ActiveObjVisitImgae", tw.web_ActiveObjVisitFloorImage)
	webMux.HandleFuncAuth("/ActiveObjRankingList", tw.web_ActiveObjRankingList)
	webMux.HandleFuncAuth("/ActiveObjSuspendedList", tw.web_ActiveObjSuspendedList)
	webMux.HandleFuncAuth("/StatServeAPI", tw.web_ProtocolStat)
	webMux.HandleFuncAuth("/StatNotification", tw.web_NotiStat)
	webMux.HandleFuncAuth("/StatAPIError", tw.web_ErrorStat)
	webMux.HandleFuncAuth("/towerStat", tw.web_towerStat)

	webMux.HandleFuncAuth("/terrain", tw.web_TerrainInfo)
	webMux.HandleFuncAuth("/terrainimagezoom", tw.web_TerrainImageZoom)
	webMux.HandleFuncAuth("/terrainimageautozoom", tw.web_TerrainImageAutoZoom)
	webMux.HandleFuncAuth("/terraintile", tw.web_TerrainTile)

	webMux.HandleFuncAuth("/ConnectionList", tw.web_ConnectionList)
	webMux.HandleFuncAuth("/KickConnection", tw.web_KickConnection)
	webMux.HandleFuncAuth("/KickActiveObj", tw.web_KickActiveObj)
	webMux.HandleFuncAuth("/SessionList", tw.web_SessionList)
	webMux.HandleFuncAuth("/DelSession", tw.web_DelSession)

	webMux.HandleFuncAuth("/Broadcast", tw.web_Broadcast)
	webMux.HandleFuncAuth("/ListenClientPause", tw.web_ListenClientPause)
	webMux.HandleFuncAuth("/ListenClientResume", tw.web_ListenClientResume)
	webMux.HandleFuncAuth("/SetSoftMax_Connection", tw.web_SetSoftMax_Connection)

	webMux.HandleFunc("/Config", tw.json_Config)

	authdata.AddAllActionName(tw.sconfig.WebAdminID)
	tw.log.TraceService("%v", webMux)

	tw.adminWeb = &http.Server{
		Handler: webMux,
		Addr:    fmt.Sprintf(":%v", tw.sconfig.AdminPort),
	}
}

func (tw *Tower) BuildDate() time.Time {
	return version.GetBuildDate()
}

func (tw *Tower) NumGoroutine() int {
	return runtime.NumGoroutine()
}

func (tw *Tower) TileCacheCount() int {
	return tile_flag.TileCacheCount()
}

func (tw *Tower) WrapInfo() string {
	return wrapper.G_WrapperInfo()
}

func (tw *Tower) GetTowerAchieveStat() *towerachieve_vector.TowerAchieveVector {
	return tw.towerAchieveStat
}

func (tw *Tower) GetConnManager() *c2t_connbytemanager.Manager {
	return tw.connManager
}

func (tw *Tower) GetSessionManager() *sessionmanager.SessionManager {
	return tw.sessionManager
}

func (tw *Tower) GetStartTime() time.Time {
	return tw.startTime
}

func (tw *Tower) GetID2ActiveObj() *aoid2activeobject.ActiveObjID2ActiveObject {
	return tw.id2ao
}

func (tw *Tower) GetID2ActiveObjSuspend() *aoid2activeobject.ActiveObjID2ActiveObject {
	return tw.id2aoSuspend
}

func (tw *Tower) GetSendStat() *actpersec.ActPerSec {
	return tw.sendStat
}
func (tw *Tower) GetRecvStat() *actpersec.ActPerSec {
	return tw.recvStat
}
func (tw *Tower) GetProtocolStat() *c2t_statserveapi.StatServeAPI {
	return tw.protocolStat
}
func (tw *Tower) GetNotiStat() *c2t_statnoti.StatNotification {
	return tw.notiStat
}
func (tw *Tower) GetErrorStat() *c2t_statapierror.StatAPIError {
	return tw.errorStat
}
func (tw *Tower) GetServiceInfo() *c2t_obj.ServiceInfo {
	return tw.serviceInfo
}

func (tw *Tower) GetTowerInfo() *c2t_obj.TowerInfo {
	return tw.towerInfo
}

func (tw *Tower) GetTowerCmdActStat() *actpersec.ActPerSec {
	return tw.towerCmdActStat
}

func (tw *Tower) SysAOID() string {
	return activeobject.SysAOIDMaker.String()
}
func (tw *Tower) EquipID() string {
	return carryingobject.EquipIDMaker.String()
}
func (tw *Tower) MoneyID() string {
	return carryingobject.MoneyIDMaker.String()
}
func (tw *Tower) PotionID() string {
	return carryingobject.PotionIDMaker.String()
}
func (tw *Tower) ScrollID() string {
	return carryingobject.ScrollIDMaker.String()
}
func (tw *Tower) DOID() string {
	return dangerobject.DOIDMaker.String()
}
func (tw *Tower) FOID() string {
	return fieldobject.FOIDMaker.String()
}
func (tw *Tower) RoomID() string {
	return room.RoomIDMaker.String()
}

func (tw *Tower) web_TowerInfo(w http.ResponseWriter, r *http.Request) {
	tplIndex, err := template.New("index").Parse(`
	<html> <head>
	<title>Tower {{.GetTowerInfo.Name}} admin</title>
	</head>
	<body>

	service cmd <a href= "/service?cmd=logreopen" target="_blank">logreopen</a> 
	<br/>
	service cmd <a href= "/service?cmd=stop" target="_blank">stop</a> 
	<hr/>

	BuildDate : {{.BuildDate.Format "2006-01-02T15:04:05Z07:00"}}
	<br/>
	Version: {{.GetServiceInfo.Version}}
	<br/>
	ProtocolVersion : {{.GetServiceInfo.ProtocolVersion}}
	<br/>
	DataVersion : {{.GetServiceInfo.DataVersion}}
	<hr/>
	{{.}}
	<br/>
	Start : {{.GetStartTime}} / {{.GetRunDur}}
	<br/>
	Factor : {{.GetTowerInfo.Factor}}
	<br/>
	Current Bias : {{.GetBias}}
	<br/>
	TotalFloor : {{.GetFloorManager.GetFloorCount}}
	<br/>
	Tile2Discover : {{.GetFloorManager.CalcTiles2Discover}}
	<br/>
	Max Exp From Discover : {{.GetFloorManager.CalcFullDiscoverExp}} 
	<br/>
	Max Level From Discover :  {{.GetFloorManager.CalcFullDiscoverLevel}} 
	<br/>
	SysAOID : {{.SysAOID}}
	<br/>
	EquipID : {{.EquipID}}
	<br/>
	MoneyID : {{.MoneyID}}
	<br/>
	PotionID : {{.PotionID}}
	<br/>
	ScrollID : {{.ScrollID}}
	<br/>
	DOID : {{.DOID}}
	<br/>
	FOID : {{.FOID}}
	<br/>
	RoomID : {{.RoomID}}
	<br/>
	goroutine : {{.NumGoroutine}}	
	<br/>
	global wrapper : {{.WrapInfo}}	
	<br/>
	TileCache : {{.TileCacheCount}}
	<br/>
	SendStat : {{.GetSendStat}}
	<br/>
	RecvStat : {{.GetRecvStat}}
	<br/>
	<a href= "/towerStat" target="_blank">Tower Achieve</a>
    <br/>
	TowerCmd act : {{.ReqChState}} {{.GetTowerCmdActStat}}
    <br/>
    <a href="/StatServeAPI" target="_blank">{{.GetProtocolStat}}</a>
    <br/>
    <a href="/StatNotification" target="_blank">{{.GetNotiStat}}</a>
    <br/>
    <a href="/StatAPIError" target="_blank">{{.GetErrorStat}}</a>
    <br/>
    <a href="/ActiveObjRankingList?page=0" target="_blank">{{.GetID2ActiveObj}}</a>
    <br/>
    <a href="/ActiveObjSuspendedList?page=0" target="_blank">{{.GetID2ActiveObjSuspend}}</a>
    <br/>
    <a href="/ConnectionList?page=0" target="_blank">Connections:{{.GetConnManager}}</a>
    <br/>
    <a href="/SessionList?page=0" target="_blank">{{.GetSessionManager}}</a>
    <br/>
    Listen Client pasuse : {{.IsListenClientPaused}}
    <a href='/ListenClientPause' target="_blank">[Pause]</a>
    <a href='/ListenClientResume' target="_blank">[Resume]</a>
    <a href='/SetSoftMax_Connection?SoftMax=' target="_blank">[SetSoftMax]</a>
    <form action="/Broadcast" target="_blank">
		Broadcast Message: 
		<input type="text" name="Msg" value="" size="64">
		<input type="submit" value="Broadcast">
    </form>
	<table border=1 style="border-collapse:collapse;">
	` + floor_HTML_header + `
	{{range $i, $v := .GetFloorManager.GetFloorList}}
		{{if $v}}
	` + floor_HTML_row + `
		{{end}}
	{{end}}
	` + floor_HTML_header + `
	</table>
	<br/>
	<pre>{{.GetServiceInfo.StringForm}}</pre>
	<pre>{{.GetTowerInfo.StringForm}}</pre>
	<pre>{{.Config.StringForm}}</pre>
	<br/>
	</body> </html> 
	`)
	if err != nil {
		tw.log.Error("%v", err)
	}
	if err := weblib.SetFresh(w, r); err != nil {
		tw.log.Error("%v", err)
	}
	if err := tplIndex.Execute(w, tw); err != nil {
		tw.log.Error("%v", err)
	}
}

const (
	floor_HTML_header = `<tr>
	<td>Floor</td>
	<td>FloorCmd act</td>
	<td>Terrain Ageing</td>
	<td>Faction</td>
	<td>W/H</td>
	<td>ActiveObj/CarryObj</td>
	<td>ActTurnJitter</td>
	<td>ObjOver Packet</td>
	<td>Viewport Cache</td>
	</tr>`
	floor_HTML_row = `<tr>
	<td>
	<a href= "/floor?floorname={{$v.GetName}}" target="_blank">
		{{$v.GetName}}
	</a>
	</td>
	<td>{{$v.ReqChState}} {{$v.GetCmdFloorActStat}}
	</td>
	<td>
	<a href= "/terrain?floorname={{$v.GetName}}" target="_blank">
		{{$v.GetTerrain.AgeingCount}}/{{$v.GetTerrain.GetResetAfterNAgeing}}
	</a>
	</td>
	<td>{{$v.GetBias.NearFaction}}</td>
	<td>{{$v.GetWidth}}/{{$v.GetHeight}}</td>
	<td>{{$v.TotalActiveObjCount}} / {{$v.TotalCarryObjCount}}</td>
	<td>{{$v.GetInterDur.GetCount}} {{$v.GetInterDur.GetInterval}} {{$v.GetInterDur.GetDuration}}</td>
	<td>{{$v.GetStatPacketObjOver}}</td>
	<td>{{$v.GetTerrain.GetViewportCache.HitRate}}</td>
	</tr>`
)
