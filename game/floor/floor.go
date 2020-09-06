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
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/game/terrain"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/lib/uuidposman"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/intervalduration"
)

var _ gamei.FloorI = &Floor{}

func (f *Floor) String() string {
	return fmt.Sprintf("Floor[%v %v/%v]",
		f.GetName(),
		len(f.recvRequestCh), cap(f.recvRequestCh),
	)
}

func (f *Floor) ReqState() string {
	return fmt.Sprintf("%v/%v",
		len(f.recvRequestCh), cap(f.recvRequestCh),
	)
}

type Floor struct {
	rnd *g2rand.G2Rand `prettystring:"hide"`
	log *g2log.LogBase `prettystring:"hide"`

	tower       gamei.TowerI
	w           int
	h           int
	bias        bias.Bias        `prettystring:"simple"`
	terrain     *terrain.Terrain `prettystring:"simple"`
	initialized bool

	// async ageing one at a time
	inAgeing int32

	aoPosMan *uuidposman.UUIDPosMan `prettystring:"simple"`
	poPosMan *uuidposman.UUIDPosMan `prettystring:"simple"`
	foPosMan *uuidposman.UUIDPosMan `prettystring:"simple"`
	doPosMan *uuidposman.UUIDPosMan `prettystring:"simple"`

	interDur          *intervalduration.IntervalDuration `prettystring:"simple"`
	statPacketObjOver *actpersec.ActPerSec               `prettystring:"simple"`
	floorCmdActStat   *actpersec.ActPerSec               `prettystring:"simple"`

	// for actturn data
	recvRequestCh chan interface{}

	aiWG sync.WaitGroup // for ai run
}

func New(seed int64, ts []string, tw gamei.TowerI) *Floor {
	queuesize := int(float64(tw.Config().ConcurrentConnections) * tw.Config().TurnPerSec)
	f := &Floor{
		log:               tw.Log(),
		tower:             tw,
		rnd:               g2rand.NewWithSeed(seed),
		interDur:          intervalduration.New(""),
		statPacketObjOver: actpersec.New(),
		floorCmdActStat:   actpersec.New(),
		recvRequestCh:     make(chan interface{}, queuesize),
	}
	f.terrain = terrain.New(f.rnd.Int63(), ts, f.tower.Config().DataFolder, f.log)
	return f
}

func (f *Floor) Cleanup() {
	f.log.TraceService("Start Cleanup Floor %v", f.GetName())
	defer func() { f.log.TraceService("End Cleanup Floor %v", f.GetName()) }()

	f.aoPosMan.Cleanup()
	f.poPosMan.Cleanup()
	f.terrain.Cleanup()
	f.doPosMan.Cleanup()
}

// Init bi need for randomness
func (f *Floor) Init() error {
	f.log.TraceService("Start Init %v", f)
	defer func() { f.log.TraceService("End Init %v", f) }()

	if err := f.terrain.Init(); err != nil {
		return fmt.Errorf("fail to make terrain %v", err)
	}
	if f.terrain.GetName() == "" {
		return nil // skip no name floor terrain
	}
	f.w, f.h = f.GetTerrain().GetXYLen()
	f.aoPosMan = uuidposman.New(f.w, f.h)
	f.poPosMan = uuidposman.New(f.w, f.h)
	f.foPosMan = f.terrain.GetFieldObjPosMan()
	f.doPosMan = uuidposman.New(f.w, f.h)
	f.bias = bias.Bias{
		f.rnd.Float64() - 0.5,
		f.rnd.Float64() - 0.5,
		f.rnd.Float64() - 0.5,
	}.MakeAbsSumTo(gameconst.FloorBaseBiasLen)

	f.initialized = true
	return nil
}

func (f *Floor) Run(ctx context.Context) {
	f.log.TraceService("start Run %v", f)
	defer func() { f.log.TraceService("End Run %v", f) }()

	ageingSec := time.Duration(f.terrain.GetMSPerAgeing()) * time.Millisecond
	if ageingSec == 0 {
		ageingSec = time.Hour * 24 * 365
	}
	ageingTk := time.NewTicker(ageingSec)
	if f.terrain.GetMSPerAgeing() == 0 {
		ageingTk.Stop()
	} else {
		defer ageingTk.Stop()
	}

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

	turnDur := time.Duration(
		float64(time.Second) / f.tower.Config().TurnPerSec / f.terrain.ActTurnBoost)
	timerTurnTk := time.NewTicker(turnDur)
	defer timerTurnTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerInfoTk.C:
			f.statPacketObjOver.UpdateLap()
			f.floorCmdActStat.UpdateLap()
			if len(f.recvRequestCh) > cap(f.recvRequestCh)/2 {
				f.log.Fatal("Floor %v %v reqch overloaded %v/%v",
					f.terrain.Name, f.GetName(),
					len(f.recvRequestCh), cap(f.recvRequestCh))
			}
			if len(f.recvRequestCh) >= cap(f.recvRequestCh) {
				break loop
			}

		case data := <-f.recvRequestCh:
			f.floorCmdActStat.Inc()
			f.processCmd2Floor(data)

		case <-ageingTk.C:
			go f.processAgeing()

		case turnTime := <-timerTurnTk.C:
			f.processTurn(turnTime)
			lastTurnDur := f.interDur.GetDuration().GetLastDuration()
			if lastTurnDur > turnDur {
				// 	f.log.Monitor("%v %v slow %v > %v", f.GetName(), f, lastTurnDur, turnDur)
			}
		}
	}
}

func (f *Floor) processAgeing() {
	if atomic.CompareAndSwapInt32(&f.inAgeing, 0, 1) {
		defer atomic.AddInt32(&f.inAgeing, -1)

		var err error
		err = f.terrain.Ageing()
		if err != nil {
			f.log.Fatal("%v %v", f, err)
			return
		}
		NotiAgeing := f.ToPacket_NotiAgeing()
		for _, v := range f.aoPosMan.GetAllList() {
			ao := v.(gamei.ActiveObjectI)
			ao.SetNeedTANoti()
			// send ageing noti
			if aoconn := ao.GetClientConn(); aoconn != nil {
				if err := aoconn.SendNotiPacket(
					c2t_idnoti.Ageing,
					NotiAgeing,
				); err != nil {
					f.log.Error("%v %v %v", f, ao, err)
				}
			}
		}
	} else {
		f.log.Fatal("processAgeing skipped %v", f)
	}
}
