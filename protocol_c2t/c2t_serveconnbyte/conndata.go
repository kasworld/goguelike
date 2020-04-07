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

package c2t_serveconnbyte

import (
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/game/session"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

var _ gamei.ServeClientConnI = &ServeConnByte{}

type ConnData struct {
	RemoteAddr string
	Session    *session.Session
	ActiveObj  gamei.ActiveObjectI
}

func (scb *ServeConnByte) ToPacket_AccountInfo() *c2t_obj.AccountInfo {
	connData := scb.connData.(*ConnData)
	return &c2t_obj.AccountInfo{
		SessionG2ID:   g2id.NewFromString(connData.Session.GetUUID()),
		ActiveObjG2ID: g2id.NewFromString(connData.ActiveObj.GetUUID()),
		NickName:      connData.Session.NickName,
		CmdList:       *scb.authorCmdList,
	}
}

func (scb *ServeConnByte) GetActiveObj() gamei.ActiveObjectI {
	connData := scb.connData.(*ConnData)
	return connData.ActiveObj
}

func (scb *ServeConnByte) GetSession() *session.Session {
	connData := scb.connData.(*ConnData)
	return connData.Session
}
