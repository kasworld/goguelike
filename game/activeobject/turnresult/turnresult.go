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

package turnresult

import (
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/lib/g2id"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

type uuidObjI interface {
	GetUUID() string
}

type TurnResult struct {
	ResultType turnresulttype.TurnResultType
	DstObj     uuidObjI
	Arg        float64
}

func New(ResultType turnresulttype.TurnResultType,
	DstObj uuidObjI, Arg float64) TurnResult {
	return TurnResult{
		ResultType: ResultType,
		DstObj:     DstObj,
		Arg:        Arg,
	}
}

func (ar TurnResult) GetTurnResultType() turnresulttype.TurnResultType {
	return ar.ResultType
}

func (ar TurnResult) GetDstObj() uuidObjI {
	return ar.DstObj
}

func (ar TurnResult) GetDamage() float64 {
	return ar.Arg
}

func (ar TurnResult) ToPacket_TurnResultClient() c2t_obj.TurnResultClient {
	if dstao := ar.DstObj; dstao != nil {
		return c2t_obj.TurnResultClient{
			ResultType: ar.ResultType,
			DstG2ID:    g2id.NewFromString(dstao.GetUUID()),
			Arg:        ar.Arg,
		}
	}
	return c2t_obj.TurnResultClient{
		ResultType: ar.ResultType,
		DstG2ID:    "",
		Arg:        ar.Arg,
	}
}
