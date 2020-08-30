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

package floor

import (
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/slippperydata"
	"github.com/kasworld/goguelike/enum/achievetype"
	"github.com/kasworld/goguelike/enum/condition"
	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/turnresulttype"
	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/activeobject/turnresult"
	"github.com/kasworld/goguelike/game/aoactreqrsp"
	"github.com/kasworld/goguelike/game/dangerobject"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
)

func (f *Floor) addBasicAttack(ao gamei.ActiveObjectI, arr *aoactreqrsp.ActReqRsp) {
	aox, aoy, exist := f.aoPosMan.GetXYByUUID(ao.GetUUID())
	if !exist {
		f.log.Error("ao not in currentfloor %v %v", f, ao)
		arr.SetDone(
			aoactreqrsp.Act{Act: c2t_idcmd.Meditate},
			c2t_error.ActionProhibited)
		return
	}
	atkdir := arr.Req.Dir
	if ao.GetTurnData().Condition.TestByCondition(condition.Drunken) {
		turnmod := slippperydata.Drunken[f.rnd.Intn(len(slippperydata.Drunken))]
		atkdir = atkdir.TurnDir(turnmod)
	}
	// add dopoaman near attack

	// check valid attack
	if !atkdir.IsValid() || atkdir == way9type.Center {
		arr.SetDone(aoactreqrsp.Act{Act: c2t_idcmd.Attack, Dir: atkdir},
			c2t_error.InvalidDirection)
		return
	}
	srcTile := f.terrain.GetTiles()[aox][aoy]
	if srcTile.NoBattle() {
		arr.SetDone(aoactreqrsp.Act{Act: c2t_idcmd.Attack, Dir: atkdir},
			c2t_error.ActionProhibited)
		return
	}
	dstX, dstY := aox+atkdir.Dx(), aoy+atkdir.Dy()
	dstX, dstY = f.terrain.WrapXY(dstX, dstY)
	dstTile := f.terrain.GetTiles()[dstX][dstY]
	if dstTile.NoBattle() {
		arr.SetDone(aoactreqrsp.Act{Act: c2t_idcmd.Attack, Dir: atkdir},
			c2t_error.ActionProhibited)
		return
	}
	if err := f.doPosMan.AddToXY(
		dangerobject.NewBasicAttact(ao, aox, aoy),
		dstX, dstY); err != nil {
		f.log.Fatal("fail to AddToXY %v", err)
		arr.SetDone(aoactreqrsp.Act{Act: c2t_idcmd.Attack, Dir: atkdir},
			c2t_error.ActionCanceled)
		return
	}
	arr.SetDone(
		aoactreqrsp.Act{Act: c2t_idcmd.Attack, Dir: atkdir},
		c2t_error.None)
}

func (f *Floor) aoAttackActiveObj(src, dst gamei.ActiveObjectI, srcTile, dstTile tile_flag.TileFlag) {

	// attack to invisible ao miss 50%
	if dst.GetTurnData().Condition.TestByCondition(condition.Invisible) && f.rnd.Intn(2) == 0 {
		src.GetAchieveStat().Inc(achievetype.AttackMiss)
		return
	}

	// blind ao attack miss 50%
	if src.GetTurnData().Condition.TestByCondition(condition.Blind) && f.rnd.Intn(2) == 0 {
		src.GetAchieveStat().Inc(achievetype.AttackMiss)
		return
	}

	envbias := f.GetEnvBias()
	srcbias := src.GetTurnData().AttackBias.Add(envbias)
	dstbias := dst.GetTurnData().DefenceBias.Add(envbias)

	atkMod := srcTile.AtkMod()
	defMod := dstTile.DefMod()
	atkValue := srcbias.SelectSkill(f.rnd.Intn(3))
	defValue := dstbias.SelectSkill(f.rnd.Intn(3))
	diffValue := atkValue*atkMod - defValue*defMod +
		src.GetTurnData().Level - dst.GetTurnData().Level +
		f.rnd.NormFloat64Range(gameconst.ActiveObjBaseBiasLen, 0)
	atkSuccess := diffValue > 0
	atkCritical := false
	rndValue := f.rnd.Intn(20)
	if rndValue == 0 {
		atkSuccess = false
	} else if rndValue == 19 {
		atkSuccess = true
	}
	if atkSuccess && f.rnd.Intn(20) == 19 {
		atkCritical = true
	}

	src.GetAchieveStat().Inc(achievetype.Attack)
	dst.GetAchieveStat().Inc(achievetype.Attacked)

	if !atkSuccess {
		src.GetAchieveStat().Inc(achievetype.AttackMiss)
		return
	}

	src.GetAchieveStat().Inc(achievetype.AttackHit)
	if diffValue < 0 {
		diffValue = -diffValue
	}

	damage := diffValue

	if atkCritical {
		damage *= 2
		src.GetAchieveStat().Inc(achievetype.AttackCritical)
	}

	src.AppendTurnResult(turnresult.New(turnresulttype.AttackTo, dst, damage))
	src.GetAchieveStat().Add(achievetype.DamageTotalGive, damage)
	src.GetAchieveStat().SetIfGt(achievetype.DamageMaxGive, damage)
	dst.AppendTurnResult(turnresult.New(turnresulttype.AttackedFrom, src, damage))
	dst.GetAchieveStat().Add(achievetype.DamageTotalRecv, damage)
	dst.GetAchieveStat().SetIfGt(achievetype.DamageMaxRecv, damage)

	src.AddBattleExp(damage * gameconst.ActiveObjExp_Damage)
}
