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

package aoactreqrsp

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/condition_flag"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

type Act struct {
	Act  c2t_idcmd.CommandID
	Dir  way9type.Way9Type
	UUID string
}

func (act Act) CalcAPByActAndCondition(cndflag condition_flag.ConditionFlag) float64 {
	turn2need := act.Act.NeedTurn()
	if cndflag.TestByCondition(condition.Slow) {
		turn2need *= 2
	}
	if cndflag.TestByCondition(condition.Haste) {
		turn2need /= 2
	}
	if act.Act == c2t_idcmd.Move && act.Dir != way9type.Center {
		turn2need *= act.Dir.Len()
	}
	return turn2need
}

type ActReqRsp struct {
	Req  Act // act requested
	Done Act // act done

	Error c2t_error.ErrorCode
}

func (arr ActReqRsp) Acted() bool {
	return arr.Done.Act != c2t_idcmd.Invalid
}

func (arr ActReqRsp) IsSuccess() bool {
	return arr.Acted() && arr.Error == c2t_error.None
}

func (arr *ActReqRsp) SetDone(done Act, Error c2t_error.ErrorCode) {
	if arr.Acted() {
		fmt.Printf("already Acted %v %+v", Error, arr)
	}
	arr.Done = done
	arr.Error = Error
}
