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

	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/conndata"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_serveconnbyte"
)

func (tw *Tower) api_me2ao(me interface{}) (gamei.ActiveObjectI, error) {
	c2sc, ok := me.(*c2t_serveconnbyte.ServeConnByte)
	if !ok {
		return nil, fmt.Errorf("invalid me not c2t_serveconnbyte.ServeConnByte %#v", me)
	}
	connData := c2sc.GetConnData().(*conndata.ConnData)
	if connData.Session == nil {
		return nil, fmt.Errorf("no session %#v", connData)
	}
	ao, exist := tw.id2ao.GetByUUID(connData.Session.ActiveObjUUID)
	if !exist {
		return nil, fmt.Errorf("ao not found %v", connData)
	}
	return ao, nil
}
