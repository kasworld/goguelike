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

package clientai

import (
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/clientfloor"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

func (cai *ClientAI) GetRunResult() error {
	return cai.runResult
}

func (cai *ClientAI) String() string {
	return cai.config.Nickname
}

func (cai *ClientAI) GetArg() interface{} {
	return cai.config
}

func (cai *ClientAI) currentFloor() *clientfloor.ClientFloor {
	if fi := cai.FloorInfo; fi == nil {
		return nil
	} else {
		return cai.Name2ClientFloor[fi.Name]
	}
}

func (cai *ClientAI) TowerBias() bias.Bias {
	if cai.OLNotiData == nil {
		return bias.Bias{}
	}
	ft := cai.TowerInfo.Factor
	dur := cai.OLNotiData.Time.Sub(cai.TowerInfo.StartTime)
	return bias.MakeBiasByProgress(ft, dur.Seconds(), gameconst.TowerBaseBiasLen)
}

func (cai *ClientAI) GetPlayerXY() (int, int) {
	ao := cai.playerActiveObjClient
	if ao != nil {
		return ao.X, ao.Y
	}
	return 0, 0
}

func (cai *ClientAI) CanUseCmd(cmd c2t_idcmd.CommandID) bool {
	if acinfo := cai.AccountInfo; acinfo != nil {
		return acinfo.CmdList[cmd]
	}
	return false
}
