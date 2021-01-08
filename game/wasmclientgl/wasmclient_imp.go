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

package wasmclientgl

import (
	"fmt"
	"time"

	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (app *WasmClient) ServerTime() time.Time {
	return time.Now().Add(app.ServerClientTimeDiff + app.PingDur/2)
}

func (app *WasmClient) GetEnvBias() bias.Bias {
	if gInitData.TowerInfo == nil {
		return bias.Bias{}
	}
	var envBias bias.Bias
	if fd := app.olNotiData; fd != nil {
		envBias = app.TowerBias().Add(app.CurrentFloor.FloorInfo.Bias)
	}
	return envBias
}

func (app *WasmClient) TowerBias() bias.Bias {
	if gInitData.TowerInfo == nil {
		return bias.Bias{}
	}
	if app.olNotiData == nil {
		return bias.Bias{}
	}
	ft := gInitData.TowerInfo.Factor
	dur := app.ServerTime().Sub(gInitData.TowerInfo.StartTime)
	// dur := app.olNotiData.Time.Sub(gInitData.TowerInfo.StartTime)
	return bias.MakeBiasByProgress(ft, dur.Seconds(), gameconst.TowerBaseBiasLen)
}

func (app *WasmClient) ActionResult2String(ar *aoactreqrsp.ActReqRsp) string {
	if ar == nil {
		return ""
	}

	postr := ar.Done.UUID
	po, exist := app.CaObjUUID2CaObjClient[ar.Done.UUID]
	if exist {
		postr = obj2ColorStr(po)
	}
	successStr := ""
	if !ar.IsSuccess() {
		successStr = fmt.Sprintf(", %v fail", ar.Req.Act)
	}
	if ar.IsSuccess() && ar.Done.Act == c2t_idcmd.Move {
		return ""
	}
	if ar.IsSuccess() && ar.Done.Act == c2t_idcmd.Attack {
		return ""
	}

	dirStr := ar.Done.Dir.String()
	if ar.Done.Dir == way9type.Center {
		dirStr = ""
	}

	return fmt.Sprintf("%v %v %v %v",
		ar.Done.Act,
		dirStr,
		postr,
		successStr,
	)
}

func (app *WasmClient) GetPlayerActiveObjClient() *c2t_obj.ActiveObjClient {
	if gInitData.AccountInfo == nil {
		return nil
	}
	ao, _ := app.AOUUID2AOClient[gInitData.AccountInfo.ActiveObjUUID]
	return ao
}

func (app *WasmClient) GetPlayerXY() (int, int) {
	if gInitData.AccountInfo == nil {
		return 0, 0
	}
	ao, _ := app.AOUUID2AOClient[gInitData.AccountInfo.ActiveObjUUID]
	if ao != nil {
		return ao.X, ao.Y
	}
	return 0, 0
}
