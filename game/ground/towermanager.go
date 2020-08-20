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
	"net/http"
	"sync"

	"github.com/kasworld/goguelike/config/groundconfig"
	"github.com/kasworld/goguelike/config/towerdata"
	"github.com/kasworld/goguelike/config/towerwsurl"
	"github.com/kasworld/goguelike/game/towerlist4client"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_obj"
	"github.com/kasworld/weblib"
)

type TowerManager struct {
	mutex           sync.RWMutex               `prettystring:"hide"`
	log             *g2log.LogBase             `prettystring:"hide"`
	towerList       []*TowerRunning            // don't use index
	sconfig         *groundconfig.GroundConfig `prettystring:"simple"`
	tower2AutoStart []*TowerRunning
}

func NewTowerManager(
	sconfig *groundconfig.GroundConfig,
	l *g2log.LogBase,
	towerdatalist []towerdata.TowerData,
) *TowerManager {
	tm := &TowerManager{
		towerList:       make([]*TowerRunning, 0),
		tower2AutoStart: make([]*TowerRunning, 0),
		sconfig:         sconfig,
		log:             l,
	}

	for i, v := range towerdatalist {
		towerNumber := i + 1
		cu, err := towerwsurl.MakeWebsocketURL(
			sconfig.TowerServiceHostBase,
			sconfig.TowerServicePortBase,
			towerNumber)
		if err != nil {
			tm.log.Fatal("fail to MakeWebsocketURL %v", err)
		}
		tr := &TowerRunning{
			log:        tm.log,
			sconfig:    tm.sconfig,
			ConnectURL: cu,
			TowerConfigMade: *tm.sconfig.MakeTowerConfig(
				towerNumber,
				v.TowerName,
				v.ScriptFilename,
				v.TurnPerSec),
		}
		tm.towerList = append(tm.towerList, tr)
		if v.AutoStart {
			tm.tower2AutoStart = append(tm.tower2AutoStart, tr)
		}
	}
	return tm
}

func (tm *TowerManager) GetAutoStartTowerList() []*TowerRunning {
	return tm.tower2AutoStart
}

func (tm *TowerManager) GetByTowerName(name string) *TowerRunning {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	for _, v := range tm.towerList {
		if v.TowerConfigMade.TowerName == name {
			return v
		}
	}
	return nil
}

func (tm *TowerManager) GetByTowerUUID(id string) *TowerRunning {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	for _, v := range tm.towerList {
		if v.TowerInfo.UUID == id {
			return v
		}
	}
	return nil
}

func (tm *TowerManager) Register(req *t2g_obj.ReqRegister_data) error {
	tw := tm.GetByTowerName(req.TI.Name)
	if tw == nil {
		return fmt.Errorf("not managed tower %v", req)
	}
	tw.Register(tm.sconfig, req)
	return nil
}

func (tm *TowerManager) Heartbeat(req *t2g_obj.ReqHeartbeat_data) error {
	tw := tm.GetByTowerUUID(req.TowerUUID)
	if tw == nil {
		return fmt.Errorf("tower not found %v", req)
	}
	tw.Heartbeat(req)
	return nil
}

func (tm *TowerManager) Web_TowerInfo(w http.ResponseWriter, r *http.Request) {
	towerName := weblib.GetStringByName("name", "", w, r)
	te := tm.GetByTowerName(towerName)
	if te == nil {
		tm.log.Error("invalid towerName")
		http.Error(w, "invalid towerName", http.StatusNotFound)
		return
	}
	te.Web_Info(w, r)
}

func (tm *TowerManager) GetTowerList() []*TowerRunning {
	return tm.towerList
}

func (tm *TowerManager) MakeTowerListJSON4Client() []towerlist4client.TowerInfo2Enter {
	tl := make([]towerlist4client.TowerInfo2Enter, 0)
	for _, v := range tm.towerList {
		te := towerlist4client.TowerInfo2Enter{
			Name: v.TowerConfigMade.TowerName,
		}
		if v.IsPing() {
			te.ConnectURL = v.ConnectURL
		}
		tl = append(tl, te)
	}
	return tl
}
