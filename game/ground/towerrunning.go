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
	"time"

	"github.com/kasworld/goguelike/config/groundconfig"
	"github.com/kasworld/goguelike/config/towerconfig"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_obj"
	"github.com/kasworld/launcherlib"
	"github.com/kasworld/weblib"
)

type TowerRunning struct {
	log             *g2log.LogBase `prettystring:"hide"`
	sconfig         *groundconfig.GroundConfig
	ConnectURL      string
	TowerConfigMade towerconfig.TowerConfig

	StartTime       time.Time `prettystring:"simple"`
	TowerInfo       c2t_obj.TowerInfo
	ServiceInfo     c2t_obj.ServiceInfo
	TowerConfigRecv towerconfig.TowerConfig
	AdminURL        string
	ServiceURL      string

	LastPing   time.Time `prettystring:"simple"`
	StatusInfo []string
}

func (tw *TowerRunning) IsPing() bool {
	return time.Now().Sub(tw.LastPing) < time.Second*2
}

func (tw *TowerRunning) IsAlive() bool {
	pidfilename := tw.TowerConfigRecv.MakePIDFileFullpath()
	pid, err := launcherlib.GetPidIntFromPidFile(pidfilename)
	if err != nil {
		return false
	}
	return launcherlib.ProcessAlive(pid) == nil
}

func (tw *TowerRunning) IsConfigSame() bool {
	return tw.TowerConfigMade == tw.TowerConfigRecv
}

func (tw *TowerRunning) MakeTowerConfigDownURL() string {
	return fmt.Sprintf(
		"http://localhost:%d/TowerConfig?name=%s",
		tw.sconfig.GroundAdminWebPort,
		tw.TowerConfigMade.DisplayName,
	)
}

func (tw *TowerRunning) Register(
	sconfig *groundconfig.GroundConfig, req *t2g_obj.ReqRegister_data) {
	now := time.Now()
	tw.StartTime = now
	tw.TowerInfo = *req.TI
	tw.ServiceInfo = *req.SI
	tw.TowerConfigRecv = *req.TC
	tw.AdminURL = fmt.Sprintf("%s%s/",
		sconfig.TowerAdminHostBase, req.TC.AdminPort,
	)
	tw.ServiceURL = fmt.Sprintf("%s%s/",
		sconfig.TowerServiceHostBase, req.TC.ServicePort,
	)
}

func (tw *TowerRunning) Heartbeat(req *t2g_obj.ReqHeartbeat_data) {
	now := time.Now()
	tw.LastPing = now
	tw.StatusInfo = req.StatusInfo
}

const (
	TowerRunning_HTML_header = `<tr>
	<td>Name</td>
	<td>TowerFile</td>
	<td>Ping</td>
	<td>Process</td>
	<td>ConfigSame</td>
	<td>Link</td>
	<td>Cmd</td>
	<td>UUID</td>
	<td>Status</td>
	</tr>`
	TowerRunning_HTML_row = `<tr>
	<td>
	<a href="/TowerInfo?name={{$v.TowerConfigMade.DisplayName}}" target="_blank">{{$v.TowerConfigMade.DisplayName}}</a>
	</td>
	<td>{{$v.TowerConfigMade.TowerFilename}}</td>
	<td> {{$v.IsPing}} </td>
	<td> {{$v.IsAlive}} </td>
	<td> {{$v.IsConfigSame}} </td>
	<td>
    <a href="{{$v.AdminURL}}" target="_blank">[Admin]</a>
	<a href="/ViewTowerOutfile?name={{$v.TowerConfigMade.DisplayName}}" target="_blank">
	[Outfile]</a>
	</td>
	<td>
    <a href="/ControlTower?name={{$v.TowerConfigMade.DisplayName}}&cmd=start" target="_blank">[Start]</a>
    <a href="/ControlTower?name={{$v.TowerConfigMade.DisplayName}}&cmd=stop" target="_blank">[Stop]</a>
    <a href="/ControlTower?name={{$v.TowerConfigMade.DisplayName}}&cmd=forcestart" target="_blank">[ForceStart]</a>
	</td>
	<td>{{$v.TowerInfo.UUID}}</td>
	<td>	
	{{range $i, $v := .StatusInfo }}
		{{$v}} </br>
	{{end}}
	</td>
	</tr>`
)

func (tw *TowerRunning) Web_Info(w http.ResponseWriter, r *http.Request) {
	tplIndex, err := template.New("index").Parse(`
	<html> <head>
	<title>Tower {{.TowerConfigMade.DisplayName}} </title>
	</head>
	<body>
	Tower {{.TowerConfigMade.DisplayName}}
	<br/>
	<br/>
    <a href="{{.AdminURL}}" target="_blank">[Admin]</a>
    <a href="{{.ServiceURL}}" target="_blank">[Service]</a>
    <a href="/ViewTowerOutfile?name={{.TowerConfigMade.DisplayName}}" target="_blank">[ViewOutfile]</a>
	<a href="/ViewTowerOutfile?name={{.TowerConfigMade.DisplayName}}" target="_blank">
	[View {{.TowerConfigMade.MakeOutfileFullpath}}]</a>
    <a href="/ControlTower?name={{.TowerConfigMade.DisplayName}}&cmd=start" target="_blank">[Start]</a>
    <a href="/ControlTower?name={{.TowerConfigMade.DisplayName}}&cmd=stop" target="_blank">[Stop]</a>
    <a href="/ControlTower?name={{.TowerConfigMade.DisplayName}}&cmd=forcestart" target="_blank">[ForceStart]</a>
	<br/>
	<br/>
	Ping alive : {{.IsPing}}
	<br/>
	Process alive : {{.IsAlive}}
	<br/>
	Config same : {{.IsConfigSame}}
	<br/>
	<br/>
	Tower StatusInfo
	<br/>
	{{range $i, $v := .StatusInfo }}
		{{$v}} </br>
	{{end}}
	<br/>
	StartTime : {{.StartTime.Format "2006-01-02T15:04:05Z07:00"}}
	<br/>
	LastPing : {{.LastPing.Format "2006-01-02T15:04:05Z07:00"}}
	<br/>
	<pre>{{.TowerInfo.StringForm}}</pre>
	<pre>{{.ServiceInfo.StringForm}}</pre>
	<pre>{{.TowerConfigRecv.StringForm}}</pre>
	<pre>{{.TowerConfigMade.StringForm}}</pre>
	<br/>
	ConnectURL : {{.ConnectURL}}
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
