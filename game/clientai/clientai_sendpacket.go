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
	"time"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_gob"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_version"
	"github.com/kasworld/version"
)

func (cai *ClientAI) reqLogin(sessionUUID string, nickname string, authkey string) error {
	return cai.ReqWithRspFn(
		c2t_idcmd.Login,
		&c2t_obj.ReqLogin_data{
			SessionUUID: sessionUUID,
			NickName:    nickname,
			AuthKey:     authkey,
		},
		func(hd c2t_packet.Header, rsp interface{}) error {
			robj, err := c2t_gob.UnmarshalPacket(hd, rsp.([]byte))
			if err != nil {
				cai.log.Fatal("%v %v %v %v", cai, hd, rsp, err)
				return err
			}

			rpk := robj.(*c2t_obj.RspLogin_data)

			cai.config.SessionUUID = rpk.AccountInfo.SessionUUID
			cai.ServiceInfo = rpk.ServiceInfo
			cai.AccountInfo = rpk.AccountInfo

			if !version.IsSame(cai.ServiceInfo.Version) {
				cai.log.Error("Version mismatch client %v server %v",
					version.GetVersion(), cai.ServiceInfo.Version)
			}
			if gameconst.DataVersion != cai.ServiceInfo.DataVersion {
				cai.log.Error("DataVersion mismatch client %v server %v",
					gameconst.DataVersion, cai.ServiceInfo.DataVersion)
			}
			if c2t_version.ProtocolVersion != cai.ServiceInfo.ProtocolVersion {
				cai.log.Error("ProtocolVersion mismatch client %v server %v",
					c2t_version.ProtocolVersion, cai.ServiceInfo.ProtocolVersion)
			}
			return cai.reqAIPlay(true)
		},
	)
}

func (cai *ClientAI) reqAIPlay(onoff bool) error {
	return cai.ReqWithRspFnWithAuth(
		c2t_idcmd.AIPlay,
		&c2t_obj.ReqAIPlay_data{onoff},
		func(hd c2t_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}

func (cai *ClientAI) reqAchieveInfo() error {
	return cai.ReqWithRspFnWithAuth(
		c2t_idcmd.AchieveInfo,
		&c2t_obj.ReqAchieveInfo_data{},
		func(hd c2t_packet.Header, rsp interface{}) error {
			robj, err := c2t_gob.UnmarshalPacket(hd, rsp.([]byte))
			if err != nil {
				cai.log.Fatal("%v %v %v %v", cai, hd, rsp, err)
				return err
			}
			rpk := robj.(*c2t_obj.RspAchieveInfo_data)
			cai.log.Debug("=================================")
			for i, v := range rpk.Achieve {
				cai.log.Debug("%v : %v", achievetype.AchieveType(i).String(), v)
			}
			cai.log.Debug("Achieve list ====================")
			return nil
		},
	)
}

func (cai *ClientAI) reqHeartbeat() error {
	return cai.ReqWithRspFnWithAuth(
		c2t_idcmd.Heartbeat,
		&c2t_obj.ReqHeartbeat_data{
			Time: time.Now(),
		},
		func(hd c2t_packet.Header, rsp interface{}) error {
			robj, err := c2t_gob.UnmarshalPacket(hd, rsp.([]byte))
			if err != nil {
				cai.log.Fatal("%v %v %v %v", cai, hd, rsp, err)
				return err
			}
			rpk := robj.(*c2t_obj.RspHeartbeat_data)
			PingDur := time.Now().Sub(rpk.Time)
			cai.log.Monitor("Ping %v", PingDur)
			return nil
		},
	)
}
