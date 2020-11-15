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

package activeobject

import (
	"math/rand"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/fieldobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

func (ao *ActiveObject) ResetPlan(sai *ServerAIState) {
	sai.planRemainCount = 0
}

// func (ao *ActiveObject) GetAIDur(sai *ServerAIState) *intervalduration.IntervalDuration {
// 	return sai.interDur
// }

// func (ao *ActiveObject) GetPlan(sai *ServerAIState) aiplan.AIPlan {
// 	return sai.RunningPlanList.getCurrentPlan()
// }

// for web
// func (ao *ActiveObject) GetPlanNameList(sai *ServerAIState) string {
// 	return sai.RunningPlanList.String()
// }

// ai util fns

// start pos == ao pos , end pos == dest pos
func (ao *ActiveObject) makePath2Dest(sai *ServerAIState, dstx, dsty int) [][2]int {
	srcx, srcy := sai.aox, sai.aoy
	trylimit := gameconst.ViewPortWH
	rtn := ao.currentFloor.FindPath(dstx, dsty, srcx, srcy, trylimit)
	return rtn
}

func (ao *ActiveObject) followPath2Dest(sai *ServerAIState) (way9type.Way9Type, bool) {
	if sai.movePath2Dest == nil {
		return way9type.Center, false
	}
	if len(sai.movePath2Dest) < 2 {
		return way9type.Center, true
	}
	aopos := [2]int{sai.aox, sai.aoy}
	dstPos := sai.movePath2Dest[len(sai.movePath2Dest)-1]

	if aopos == dstPos {
		return way9type.Center, true
	}

	w, h := ao.currentFloor.GetTerrain().GetXYLen()
	if aopos == sai.movePath2Dest[0] {
		nextpos := sai.movePath2Dest[1]
		isContact, toMoveDir := way9type.CalcContactDirWrapped(
			aopos, nextpos, w, h)
		if isContact {
			sai.movePath2Dest = sai.movePath2Dest[1:]
		}
		return toMoveDir, isContact
	} else {
		nextpos := sai.movePath2Dest[1]
		isContact, toMoveDir := way9type.CalcContactDirWrapped(
			aopos, nextpos, w, h)
		if isContact {
			sai.movePath2Dest = sai.movePath2Dest[1:]
			return toMoveDir, isContact
		} else {
			nextpos := sai.movePath2Dest[0]
			isContact, toMoveDir := way9type.CalcContactDirWrapped(
				aopos, nextpos, w, h)
			return toMoveDir, isContact
		}
	}
}

func (ao *ActiveObject) needUnEquipCarryObj(sai *ServerAIState, PoBias bias.Bias) bool {
	aoEnvBias := ao.GetBias().Add(ao.currentFloor.GetEnvBias())

	currentBias := aoEnvBias.Add(PoBias)
	newBias := aoEnvBias
	return newBias.AbsSum() > currentBias.AbsSum()
}
func (ao *ActiveObject) isBetterCarryObj2(sai *ServerAIState, PoEquipType equipslottype.EquipSlotType, PoBias bias.Bias) bool {
	aoEnvBias := ao.GetBias().Add(ao.currentFloor.GetEnvBias())

	newBiasAbs := aoEnvBias.Add(PoBias).AbsSum()
	v := ao.GetInven().GetEquipSlot()[PoEquipType]
	if v == nil {
		return newBiasAbs > aoEnvBias.AbsSum()
	} else {
		return newBiasAbs > aoEnvBias.Add(v.GetBias()).AbsSum()
	}
}

func (ao *ActiveObject) findMovableDir5(sai *ServerAIState, x, y int, dir way9type.Way9Type) way9type.Way9Type {
	tiles := ao.currentFloor.GetTerrain().GetTiles()
	dirList := []way9type.Way9Type{
		dir,
		dir.TurnDir(1),
		dir.TurnDir(-1),
		dir.TurnDir(2),
		dir.TurnDir(-2),
	}
	if rand.Float64() >= 0.5 {
		dirList = []way9type.Way9Type{
			dir,
			dir.TurnDir(-1),
			dir.TurnDir(1),
			dir.TurnDir(-2),
			dir.TurnDir(2),
		}
	}
	for _, dir := range dirList {
		nextX, nextY := ao.posAddDir(sai, x, y, dir)
		if tiles[nextX][nextY].CharPlaceable() {
			return dir
		}
	}
	return way9type.Center
}

func (ao *ActiveObject) findMovableDir3(sai *ServerAIState, x, y int, dir way9type.Way9Type) way9type.Way9Type {
	tiles := ao.currentFloor.GetTerrain().GetTiles()
	dirList := []way9type.Way9Type{
		dir,
		dir.TurnDir(1),
		dir.TurnDir(-1),
	}
	if rand.Float64() >= 0.5 {
		dirList = []way9type.Way9Type{
			dir,
			dir.TurnDir(-1),
			dir.TurnDir(1),
		}
	}
	for _, dir := range dirList {
		nextX, nextY := ao.posAddDir(sai, x, y, dir)
		if tiles[nextX][nextY].CharPlaceable() {
			return dir
		}
	}
	return way9type.Center
}

func (ao *ActiveObject) posAddDir(sai *ServerAIState, x, y int, dir way9type.Way9Type) (int, int) {
	ter := ao.currentFloor.GetTerrain()
	nextX := x + dir.Dx()
	nextY := y + dir.Dy()
	nextX, nextY = ter.WrapXY(nextX, nextY)
	return nextX, nextY
}

func (ao *ActiveObject) sendActNotiPacket2Floor(sai *ServerAIState,
	Act c2t_idcmd.CommandID,
	Dir way9type.Way9Type,
	UUID string,
) {
	pk := &aoactreqrsp.Act{
		Act:  Act,
		Dir:  Dir,
		UUID: UUID,
	}
	ao.SetReq2Handle(pk)
}

func (ao *ActiveObject) needRecharge(sai *ServerAIState) bool {
	return ao.GetSPRate() < 0.3 || ao.GetHPRate() < 0.3
}

func (ao *ActiveObject) aoAttackLast(sai *ServerAIState) gamei.ActiveObjectI {
	for _, v := range ao.GetTurnResultList() {
		if v.GetTurnResultType() == turnresulttype.AttackedFrom {
			dstObj := v.GetDstObj()
			switch o := dstObj.(type) {
			default:
				ao.log.Fatal("unknown dstao %v", v)
			case gamei.ActiveObjectI:
				if o.IsAlive() {
					return o
				}
			case *fieldobject.FieldObject:
			}

		}
	}
	return nil
}

func (ao *ActiveObject) overloadRate(sai *ServerAIState) float64 {
	return ao.GetTurnData().LoadRate
}

func (ao *ActiveObject) floorDiscoverRate(sai *ServerAIState) float64 {
	vf := ao.GetFloor4Client(ao.currentFloor.GetName())
	return vf.Visit.CalcCompleteRate()
}
