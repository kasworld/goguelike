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

package ground

import (
	"fmt"
	"html/template"
	"net/http"
	"syscall"
	"time"

	"github.com/kasworld/configutil"
	"github.com/kasworld/goguelike/game/aoscore"
	"github.com/kasworld/launcherlib"
	"github.com/kasworld/version"
	"github.com/kasworld/weblib"
	"github.com/kasworld/weblib/webprofile"
)

func (grd *Ground) initAdminWeb() {
	authdata := weblib.NewAuthData("ground")
	authdata.ReLoadUserData([][2]string{
		{grd.sconfig.WebAdminID, grd.sconfig.WebAdminPass},
	})
	webMux := weblib.NewAuthMux(authdata, grd.log)

	if !version.IsRelease() {
		webprofile.AddWebProfile(webMux)
	}

	webMux.HandleFuncAuth("/", grd.web_GroundInfo)
	webMux.HandleFuncAuth("/ControlTower", grd.web_ControlTower)
	webMux.HandleFuncAuth("/KillSelf", grd.web_KillSelf)

	webMux.HandleFunc("/favicon.ico", grd.web_FaviconIco)
	webMux.HandleFunc("/ViewGroundOutfile", grd.web_ViewGroundOutFile)
	webMux.HandleFunc("/TowerConfig", grd.web_TowerConfig)
	webMux.HandleFunc("/ViewTowerOutfile", grd.web_ViewTowerOutFile)
	webMux.HandleFunc("/TowerInfo", grd.twMan.Web_TowerInfo)
	webMux.HandleFunc("/HighScore", grd.web_HighScore)

	authdata.AddAllActionName(grd.sconfig.WebAdminID)
	grd.log.TraceService("%v", webMux)
	grd.adminWeb = &http.Server{
		Handler: webMux,
		Addr:    fmt.Sprintf(":%v", grd.sconfig.GroundAdminWebPort),
	}
}

func (grd *Ground) initServiceWeb() {
	webMux := http.NewServeMux()
	webMux.Handle("/",
		http.FileServer(http.Dir(grd.sconfig.ClientDataFolder)),
	)
	webMux.HandleFunc("/towerlist.json", grd.json_TowerList)
	webMux.HandleFunc("/highscore.json", grd.json_HighScore)

	grd.clientWeb = &http.Server{
		Handler: webMux,
		Addr:    fmt.Sprintf(":%v", grd.sconfig.GroundServiceWebPort),
	}
}

func (grd *Ground) web_FaviconIco(w http.ResponseWriter, r *http.Request) {
}

func (grd *Ground) web_GroundInfo(w http.ResponseWriter, r *http.Request) {
	tplIndex, err := template.New("index").Parse(`
	<html> <head>
	<title>Ground  admin</title>
	</head>
	<body>
	<pre>{{.GetConfig.StringForm}}</pre>
	<hr/>
	<a href= "/KillSelf" target="_blank">[End ground]</a>
	<br/>
	<br/>
	<a href= "/ViewGroundOutfile" target="_blank">groundserver.out</a>
	<br/>
	BuildDate : {{.BuildDate.Format "2006-01-02T15:04:05Z07:00"}}
	<br/>
	Version: {{.GetServiceInfo.Version}}
	<br/>
	ProtocolVersion : {{.GetServiceInfo.ProtocolVersion}}
	<br/>
	DataVersion : {{.GetServiceInfo.DataVersion}}
	<br/>
	goroutine : {{.NumGoroutine}}	
	<br/>
    <a href="/HighScore?page=0" target="_blank">High score</a>
	<br/>
	<table border=1 style="border-collapse:collapse;">
	` + TowerRunning_HTML_header + `
	{{range $i, $v := .GetTowerManager.GetTowerList}}
		{{if $v}}
	` + TowerRunning_HTML_row + `
		{{end}}
	{{end}}
	` + TowerRunning_HTML_header + `
	</table>
	<hr/>
	<pre>{{.StringForm}}</pre>
	</body> </html> 
	`)
	if err != nil {
		grd.log.Error("%v", err)
	}
	if err := weblib.SetFresh(w, r); err != nil {
		grd.log.Error("%v", err)
	}
	if err := tplIndex.Execute(w, grd); err != nil {
		grd.log.Error("%v", err)
	}
}

func (grd *Ground) web_ViewGroundOutFile(w http.ResponseWriter, r *http.Request) {
	err := weblib.File2WebOut(w, grd.sconfig.MakeOutfileFullpath())
	if err != nil {
		grd.log.Error("%v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (grd *Ground) web_KillSelf(w http.ResponseWriter, r *http.Request) {
	err := launcherlib.SignalByPidFile(grd.sconfig.MakePIDFileFullpath(), syscall.SIGINT)
	if err != nil {
		grd.log.Error("%v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (grd *Ground) web_HighScore(w http.ResponseWriter, r *http.Request) {
	allActiveObj := aoscore.ActiveObjScoreList(grd.highScore)
	page := weblib.GetPage(w, r)
	listActiveObj := allActiveObj.GetPage(page, 40)
	aoscore.ActiveObjScoreList(listActiveObj).ToWeb(w, r)
}

func (grd *Ground) web_ControlTower(w http.ResponseWriter, r *http.Request) {
	towerName := weblib.GetStringByName("name", "", w, r)
	te := grd.twMan.GetByTowerName(towerName)
	if te == nil {
		grd.log.Error("invalid towerName")
		http.Error(w, "invalid towerName", http.StatusNotFound)
		return
	}
	cmd := weblib.GetStringByName("cmd", "", w, r)
	if cmd == "" {
		grd.log.Error("invalid cmd")
		http.Error(w, "invalid cmd", http.StatusNotFound)
		return
	}
	err := grd.ControlTower(te, cmd)
	if err != nil {
		grd.log.Error("%v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	time.Sleep(time.Second)
	err = weblib.File2WebOut(w, te.TowerConfigMade.MakeOutfileFullpath())
	if err != nil {
		grd.log.Error("%v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
func (grd *Ground) web_TowerConfig(w http.ResponseWriter, r *http.Request) {
	towerName := weblib.GetStringByName("name", "", w, r)
	te := grd.twMan.GetByTowerName(towerName)
	if te == nil {
		grd.log.Error("invalid towerName")
		http.Error(w, "invalid towerName", http.StatusNotFound)
		return
	}
	err := configutil.GenerateToIni(&te.TowerConfigMade, w)
	if err != nil {
		grd.log.Error("%v", err)
	}
}

func (grd *Ground) web_ViewTowerOutFile(w http.ResponseWriter, r *http.Request) {
	towerName := weblib.GetStringByName("name", "", w, r)
	te := grd.twMan.GetByTowerName(towerName)
	if te == nil {
		grd.log.Error("invalid towerName")
		http.Error(w, "invalid towerName", http.StatusNotFound)
		return
	}
	err := weblib.File2WebOut(w, te.TowerConfigMade.MakeOutfileFullpath())
	if err != nil {
		grd.log.Error("%v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (grd *Ground) json_TowerList(w http.ResponseWriter, r *http.Request) {
	tl := grd.twMan.MakeTowerListJSON4Client()
	weblib.ServeJSON2HTTP(tl, w)
}

func (grd *Ground) json_HighScore(w http.ResponseWriter, r *http.Request) {
	allActiveObj := aoscore.ActiveObjScoreList(grd.highScore)
	page := weblib.GetPage(w, r)
	listActiveObj := allActiveObj.GetPage(page, 40)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	weblib.ServeJSON2HTTP(listActiveObj, w)
}
