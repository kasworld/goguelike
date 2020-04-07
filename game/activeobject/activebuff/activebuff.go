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

// Package activebuff timed buff in activeobject
package activebuff

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/statusoptype"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

func (as ActiveBuff) String() string {
	return fmt.Sprintf("ActiveBuff[%v %v %v]",
		as.Name,
		as.Buff,
		as.AppliedCount,
	)
}

type ActiveBuff struct {
	Name         string
	ClearOnDeath bool
	Buff         []statusoptype.OpArg
	AppliedCount int
}

// GetOpArgToApply return current oparg and exist, inc appliedcount
func (as *ActiveBuff) GetOpArgToApply() (statusoptype.OpArg, bool) {
	if as.AppliedCount >= len(as.Buff) {
		return statusoptype.OpArg{}, false // ended
	}
	oa := as.Buff[as.AppliedCount]
	as.AppliedCount++
	return oa, true
}

func (as *ActiveBuff) IsEnded() bool {
	return as.AppliedCount >= len(as.Buff)
}

func (as *ActiveBuff) ToPacket_ActiveObjBuff() *c2t_obj.ActiveObjBuff {
	return &c2t_obj.ActiveObjBuff{
		Name:        as.Name,
		RemainCount: len(as.Buff) - as.AppliedCount,
	}
}
