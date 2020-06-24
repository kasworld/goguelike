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

package serverai2

import (
	"math/rand"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/enum/aiplan"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/intervalduration"
)

func (sai *ServerAI) ResetPlan() {
	sai.planRemainCount = 0
}

func (sai *ServerAI) GetAIDur() *intervalduration.IntervalDuration {
	return sai.interDur
}

func (sai *ServerAI) GetPlan() aiplan.AIPlan {
	return sai.runningPlanList.getCurrentPlan()
}

// for web
func (sai *ServerAI) GetPlanNameList() string {
	return sai.runningPlanList.String()
}

// ai util fns

// start pos == ao pos , end pos == dest pos
func (sai *ServerAI) makePath2Dest(dstx, dsty int) [][2]int {
	srcx, srcy := sai.aox, sai.aoy
	trylimit := gameconst.ViewPortWH
	rtn := sai.currentFloor.FindPath(dstx, dsty, srcx, srcy, trylimit)
	return rtn
}

func (sai *ServerAI) followPath2Dest() (way9type.Way9Type, bool) {
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

	w, h := sai.currentFloor.GetTerrain().GetXYLen()
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

func (sai *ServerAI) needUnEquipCarryObj(PoBias bias.Bias) bool {
	aoEnvBias := sai.ao.GetBias().Add(sai.currentFloor.GetEnvBias())

	currentBias := aoEnvBias.Add(PoBias)
	newBias := aoEnvBias
	return newBias.AbsSum() > currentBias.AbsSum()
}
func (sai *ServerAI) isBetterCarryObj2(PoEquipType equipslottype.EquipSlotType, PoBias bias.Bias) bool {
	aoEnvBias := sai.ao.GetBias().Add(sai.currentFloor.GetEnvBias())

	newBiasAbs := aoEnvBias.Add(PoBias).AbsSum()
	v := sai.ao.GetInven().GetEquipSlot()[PoEquipType]
	if v == nil {
		return newBiasAbs > aoEnvBias.AbsSum()
	} else {
		return newBiasAbs > aoEnvBias.Add(v.GetBias()).AbsSum()
	}
}

func (sai *ServerAI) canAttackTo(x1, y1, x2, y2 int) (way9type.Way9Type, bool) {
	ter := sai.currentFloor.GetTerrain()
	w, h := ter.GetXYLen()
	contact, dir := way9type.CalcContactDirWrappedXY(x1, y1, x2, y2, w, h)
	return dir, contact &&
		!ter.GetTiles()[x1][y1].Safe() && !ter.GetTiles()[x2][y2].Safe()
}

func (sai *ServerAI) findMovableDir5(x, y int, dir way9type.Way9Type) way9type.Way9Type {
	tiles := sai.currentFloor.GetTerrain().GetTiles()
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
		nextX, nextY := sai.posAddDir(x, y, dir)
		if tiles[nextX][nextY].CharPlaceable() {
			return dir
		}
	}
	return way9type.Center
}

func (sai *ServerAI) findMovableDir3(x, y int, dir way9type.Way9Type) way9type.Way9Type {
	tiles := sai.currentFloor.GetTerrain().GetTiles()
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
		nextX, nextY := sai.posAddDir(x, y, dir)
		if tiles[nextX][nextY].CharPlaceable() {
			return dir
		}
	}
	return way9type.Center
}

func (sai *ServerAI) posAddDir(x, y int, dir way9type.Way9Type) (int, int) {
	ter := sai.currentFloor.GetTerrain()
	nextX := x + dir.Dx()
	nextY := y + dir.Dy()
	nextX, nextY = ter.WrapXY(nextX, nextY)
	return nextX, nextY
}

func (sai *ServerAI) sendActNotiPacket2Floor(
	Act c2t_idcmd.CommandID,
	Dir way9type.Way9Type,
	UUID string,
) {
	pk := &aoactreqrsp.Act{
		Act:  Act,
		Dir:  Dir,
		UUID: UUID,
	}
	sai.ao.SetReq2Handle(pk)
}

func (sai *ServerAI) needRecharge() bool {
	return sai.ao.GetSPRate() < 0.3 || sai.ao.GetHPRate() < 0.3
}

func (sai *ServerAI) aoAttackLast() gamei.ActiveObjectI {
	for _, v := range sai.ao.GetTurnResultList() {
		if v.GetTurnResultType() == turnresulttype.AttackedFrom {
			dstObj := v.GetDstObj()
			if dstObj != nil {
				dstActiveObj := dstObj.(gamei.ActiveObjectI)
				if dstActiveObj.IsAlive() {
					return dstActiveObj
				}
			} else {
				sai.log.Fatal("dstao nil %v", v)
			}
		}
	}
	return nil
}

func (sai *ServerAI) overloadRate() float64 {
	return sai.ao.GetTurnData().LoadRate
}

func (sai *ServerAI) floorDiscoverRate() float64 {
	vf := sai.ao.GetVisitFloor(sai.currentFloor.GetUUID())
	return vf.CalcCompleteRate()
}
