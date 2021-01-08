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
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/config/authdata"
	"github.com/kasworld/goguelike/config/dataversion"
	"github.com/kasworld/goguelike/config/gamedata"
	"github.com/kasworld/goguelike/config/towerconfig"
	"github.com/kasworld/goguelike/enum/towerachieve_vector"
	"github.com/kasworld/goguelike/game/activeobject"
	"github.com/kasworld/goguelike/game/aoexpsort"
	"github.com/kasworld/goguelike/game/aoid2activeobject"
	"github.com/kasworld/goguelike/game/aoid2floor"
	"github.com/kasworld/goguelike/game/floormanager"
	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/game/towerscript"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/lib/loadlines"
	"github.com/kasworld/goguelike/lib/sessionmanager"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_connbytemanager"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_statapierror"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_statnoti"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_statserveapi"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_version"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_idcmd"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_obj"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_packet"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/rangestat"
	"github.com/kasworld/recordduration"
	"github.com/kasworld/uuidstr"
	"github.com/kasworld/version"
	"github.com/kasworld/weblib/retrylistenandserve"
)

var _ gamei.TowerI = &Tower{}

func (tw *Tower) String() string {
	return fmt.Sprintf("Tower[%v Seed:%v %v]",
		tw.sconfig.TowerName,
		tw.seed,
		tw.uuid,
	)
}

type Tower struct {
	mutex   sync.RWMutex   `prettystring:"hide"`
	doClose func()         `prettystring:"hide"`
	rnd     *g2rand.G2Rand `prettystring:"hide"`
	log     *g2log.LogBase `prettystring:"hide"`

	recvRequestCh chan interface{}

	sconfig    *towerconfig.TowerConfig
	seed       int64
	uuid       string
	biasFactor [3]int64  `prettystring:"simple"`
	startTime  time.Time `prettystring:"simple"`

	floorMan              *floormanager.FloorManager                  `prettystring:"simple"`
	ao2Floor              *aoid2floor.ActiveObjID2Floor               `prettystring:"simple"`
	id2ao                 *aoid2activeobject.ActiveObjID2ActiveObject `prettystring:"simple"`
	id2aoSuspend          *aoid2activeobject.ActiveObjID2ActiveObject `prettystring:"simple"`
	aoExpRankingSuspended aoexpsort.ByExp                             `prettystring:"simple"`
	aoExpRanking          aoexpsort.ByExp                             `prettystring:"simple"`

	serviceInfo *c2t_obj.ServiceInfo
	towerInfo   *c2t_obj.TowerInfo
	conn2ground *Conn2Ground `prettystring:"simple"`
	registered  bool

	// for server
	// limit client connection
	clientConnLimitStat *rangestat.RangeStat `prettystring:"simple"`

	// manage session (playing and played)
	sessionManager *sessionmanager.SessionManager
	// manage connection
	connManager *c2t_connbytemanager.Manager

	towerAchieveStat       *towerachieve_vector.TowerAchieveVector `prettystring:"simple"`
	sendStat               *actpersec.ActPerSec                    `prettystring:"simple"`
	recvStat               *actpersec.ActPerSec                    `prettystring:"simple"`
	protocolStat           *c2t_statserveapi.StatServeAPI          `prettystring:"simple"`
	notiStat               *c2t_statnoti.StatNotification          `prettystring:"simple"`
	errorStat              *c2t_statapierror.StatAPIError          `prettystring:"simple"`
	listenClientPaused     bool
	demuxReq2BytesAPIFnMap [c2t_idcmd.CommandID_Count]func(
		me interface{}, hd c2t_packet.Header, rbody []byte) (
		c2t_packet.Header, interface{}, error) `prettystring:"hide"`

	// tower cmd stats
	towerCmdActStat *actpersec.ActPerSec `prettystring:"simple"`

	adminWeb  *http.Server `prettystring:"simple"`
	clientWeb *http.Server `prettystring:"simple"`
}

func New(config *towerconfig.TowerConfig) *Tower {
	fmt.Printf("%v\n", config.StringForm())

	if config.BaseLogDir != "" {
		log, err := g2log.NewWithDstDir(
			config.TowerName,
			config.MakeLogDir(),
			logflags.DefaultValue(false).BitClear(logflags.LF_functionname),
			config.LogLevel,
			config.SplitLogLevel,
		)
		if err == nil {
			g2log.GlobalLogger = log
		} else {
			fmt.Printf("%v\n", err)
			g2log.GlobalLogger.SetFlags(
				g2log.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
			g2log.GlobalLogger.SetLevel(
				config.LogLevel)
		}
	} else {
		g2log.GlobalLogger.SetFlags(
			g2log.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
		g2log.GlobalLogger.SetLevel(
			config.LogLevel)
	}

	tw := &Tower{
		uuid:         uuidstr.New(),
		id2ao:        aoid2activeobject.New("ActiveObject working"),
		id2aoSuspend: aoid2activeobject.New("ActiveObject suspended"),
		recvRequestCh: make(chan interface{},
			int(float64(config.ConcurrentConnections*2)*config.TurnPerSec)),

		sconfig: config,
		log:     g2log.GlobalLogger,

		clientConnLimitStat: rangestat.New("", 0, config.ConcurrentConnections),
		sendStat:            actpersec.New(),
		recvStat:            actpersec.New(),
		protocolStat:        c2t_statserveapi.New(),
		notiStat:            c2t_statnoti.New(),
		errorStat:           c2t_statapierror.New(),
		towerCmdActStat:     actpersec.New(),
		towerAchieveStat:    new(towerachieve_vector.TowerAchieveVector),
	}

	tw.seed = int64(config.Seed)
	if tw.seed <= 0 {
		tw.seed = time.Now().UnixNano()
	}
	tw.rnd = g2rand.NewWithSeed(int64(tw.seed))

	tw.connManager = c2t_connbytemanager.New()

	tw.sessionManager = sessionmanager.New("",
		tw.sconfig.ConcurrentConnections*10, tw.log)

	tw.doClose = func() {
		tw.log.Fatal("Too early doClose call %v", tw)
	}
	tw.serviceInfo = &c2t_obj.ServiceInfo{
		Version:         version.GetVersion(),
		ProtocolVersion: c2t_version.ProtocolVersion,
		DataVersion:     dataversion.DataVersion,
	}
	authdata.AddAdminKey(config.AdminAuthKey)
	return tw
}

// return implement signalhandle.LoggerI
func (tw *Tower) GetLogger() interface{} {
	return tw.log
}

func (tw *Tower) GetServiceLockFilename() string {
	return tw.sconfig.MakePIDFileFullpath()
}

func (tw *Tower) ServiceInit() error {
	rd := recordduration.New(tw.String())

	tw.log.TraceService("Start ServiceInit %v %v", tw, rd)
	defer func() {
		tw.log.TraceService("End ServiceInit %v %v", tw, rd)
		fmt.Println(rd)
	}()

	tw.log.TraceService("%v", tw.serviceInfo.StringForm())
	tw.log.TraceService("%v", tw.sconfig.StringForm())

	var err error

	gamedata.ActiveObjNameList, err = loadlines.LoadLineList(
		filepath.Join(tw.Config().DataFolder, "ainames.txt"),
	)
	if err != nil {
		tw.log.Fatal("load ainame fail %v", err)
		return err
	}

	gamedata.ChatData, err = loadlines.LoadLineList(
		filepath.Join(tw.Config().DataFolder, "chatdata.txt"),
	)
	if err != nil {
		tw.log.Fatal("load chatdata fail %v", err)
		return err
	}

	tScript, err := towerscript.LoadJSON(
		tw.sconfig.MakeTowerFileFullpath(),
	)
	if err != nil {
		return err
	}

	tw.ao2Floor = aoid2floor.New(tw)
	tw.biasFactor = tw.NewRandFactor()

	tw.floorMan = floormanager.New(tScript, tw)
	if err := tw.floorMan.Init(tw.rnd); err != nil {
		return err
	}
	tw.startTime = time.Now()
	tw.towerInfo = &c2t_obj.TowerInfo{
		StartTime:     tw.startTime,
		UUID:          tw.uuid,
		Name:          tw.sconfig.TowerName,
		Factor:        tw.biasFactor,
		TotalFloorNum: tw.floorMan.GetFloorCount(),
		TurnPerSec:    tw.sconfig.TurnPerSec,
	}

	tw.conn2ground = NewConn2Ground(tw.sconfig.GroundRPC)

	tw.log.TraceService("%v", tw.towerInfo.StringForm())
	fmt.Printf("%v\n", tw.towerInfo.StringForm())
	if tw.sconfig.StandAlone {
		fmt.Printf("WebAdmin  : %v:%v id:%v pass:%v\n",
			tw.sconfig.ServiceHostBase, tw.sconfig.AdminPort, tw.sconfig.WebAdminID, tw.sconfig.WebAdminPass)
		// fmt.Printf("WebClient : %v:%v\n", tw.sconfig.ServiceHostBase, tw.sconfig.ServicePort)
		// fmt.Printf("WebClient with authkey : %v:%v/?authkey=%v\n", tw.sconfig.ServiceHostBase, tw.sconfig.ServicePort, tw.sconfig.AdminAuthKey)
		fmt.Printf("WebClient : %v:%v/\n", tw.sconfig.ServiceHostBase, tw.sconfig.ServicePort)
		fmt.Printf("WebClient with authkey : %v:%v/?authkey=%v\n", tw.sconfig.ServiceHostBase, tw.sconfig.ServicePort, tw.sconfig.AdminAuthKey)
	}

	return nil
}

func (tw *Tower) ServiceCleanup() {
	tw.log.TraceService("Start ServiceCleanup %v", tw)
	defer func() { tw.log.TraceService("End ServiceCleanup %v", tw) }()

	tw.id2ao.Cleanup()
	tw.ao2Floor.Cleanup()
	for _, f := range tw.floorMan.GetFloorList() {
		f.Cleanup()
	}
	tw.floorMan.Cleanup()
	tw.sessionManager.Cleanup()
}

func (tw *Tower) ServiceMain(mainctx context.Context) {
	tw.log.TraceService("Start ServiceMain %v", tw)
	defer func() { tw.log.TraceService("End ServiceMain %v", tw) }()
	ctx, closeCtx := context.WithCancel(mainctx)
	tw.doClose = closeCtx

	defer closeCtx()

	go tw.runTower(ctx)
	for _, f := range tw.floorMan.GetFloorList() {
		go func(f gamei.FloorI) {
			f.Run(ctx)
			closeCtx()
		}(f)
	}

	totalaocount := 0
	for _, f := range tw.floorMan.GetFloorList() {
		for i := 0; i < f.GetTerrain().GetActiveObjCount(); i++ {
			ao := activeobject.NewSystemActiveObj(tw.rnd.Int63(), f, tw.log, tw.towerAchieveStat)
			if err := tw.ao2Floor.ActiveObjEnterTower(f, ao); err != nil {
				tw.log.Error("%v", err)
				continue
			}
			if err := tw.id2ao.Add(ao); err != nil {
				tw.log.Error("%v", err)
			}
		}
		totalaocount += f.GetTerrain().GetActiveObjCount()
	}
	tw.log.Monitor("Total system ActiveObj in tower %v", totalaocount)

	tw.initAdminWeb()
	tw.initServiceWeb(ctx)

	go retrylistenandserve.RetryListenAndServe(tw.adminWeb, tw.log, "serveAdminWeb")
	go retrylistenandserve.RetryListenAndServe(tw.clientWeb, tw.log, "serveServiceWeb")
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		}
	}
}

func (tw *Tower) runTower(ctx context.Context) {
	tw.log.TraceService("Start Run %v", tw)
	defer func() { tw.log.TraceService("End Run %v", tw) }()

	defer func() {
		tw.log.TraceService("remain connection %v", tw.clientConnLimitStat)
	}()

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

	rankMakeTk := time.NewTicker(1 * time.Second)
	defer rankMakeTk.Stop()

	groundHeartbeat := time.NewTicker(1 * time.Second)
	defer groundHeartbeat.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			tw.log.TraceService("AdminWeb Close %v", tw.adminWeb.Close())
			tw.log.TraceService("ClientWeb Close %v", tw.clientWeb.Close())
			break loop

		case <-timerInfoTk.C:
			tw.towerCmdActStat.UpdateLap()
			tw.sendStat.UpdateLap()
			tw.recvStat.UpdateLap()
			tw.clientConnLimitStat.UpdateLap()
			if len(tw.recvRequestCh) > cap(tw.recvRequestCh)/2 {
				tw.log.Fatal("Tower reqch overloaded %v/%v",
					len(tw.recvRequestCh), cap(tw.recvRequestCh))
			}
			if len(tw.recvRequestCh) >= cap(tw.recvRequestCh) {
				break loop
			}

		case data := <-tw.recvRequestCh:
			tw.towerCmdActStat.Inc()
			go tw.processCmd2Tower(data)

		case <-rankMakeTk.C:
			go tw.makeActiveObjExpRank()
			go tw.makeActiveObjExpRankSuspended()

		case <-groundHeartbeat.C:
			if !tw.sconfig.StandAlone {
				if !tw.conn2ground.IsConnected() {
					go tw.conn2ground.Run(ctx)
					tw.registered = false
				} else if !tw.registered {
					go tw.Ground_Register()
				} else {
					go tw.Ground_Heartbeat([]string{
						fmt.Sprintf("Connection: %v", tw.connManager.Len()),
						fmt.Sprintf("Pause: %v", tw.listenClientPaused),
						fmt.Sprintf("Session: %v", tw.sessionManager.Count()),
					})
				}
			}
		}
	}
}

func (tw *Tower) makeActiveObjExpRank() {
	rtn := tw.id2ao.GetAllList()
	for _, v := range rtn {
		v.UpdateExpCopy()
	}
	aoexpsort.ByExp(rtn).Sort()
	tw.aoExpRanking = rtn
}
func (tw *Tower) makeActiveObjExpRankSuspended() {
	rtn := tw.id2aoSuspend.GetAllList()
	aoexpsort.ByExp(rtn).Sort()
	tw.aoExpRankingSuspended = rtn
}

func (tw *Tower) Ground_Register() {
	if !tw.conn2ground.IsConnected() {
		return
	}
	tw.conn2ground.ReqWithRspFn(
		t2g_idcmd.Register,
		&t2g_obj.ReqRegister_data{
			TI: tw.towerInfo,
			SI: tw.serviceInfo,
			TC: tw.sconfig,
		},
		func(hd t2g_packet.Header, rsp interface{}) error {
			tw.registered = true
			return nil
		},
	)
}

func (tw *Tower) Ground_Heartbeat(StatusInfo []string) {
	if !tw.conn2ground.IsConnected() {
		return
	}
	tw.conn2ground.ReqWithRspFn(
		t2g_idcmd.Heartbeat,
		&t2g_obj.ReqHeartbeat_data{
			TowerUUID:  tw.uuid,
			StatusInfo: StatusInfo,
		},
		func(hd t2g_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}

func (tw *Tower) Ground_HighScore(ao gamei.ActiveObjectI) {
	if !tw.conn2ground.IsConnected() {
		return
	}
	aos := ao.To_ActiveObjScore()
	aos.TowerName = tw.towerInfo.Name
	aos.TowerUUID = tw.uuid
	aos.RecordTime = time.Now()

	tw.conn2ground.ReqWithRspFn(
		t2g_idcmd.HighScore,
		&t2g_obj.ReqHighScore_data{
			ActiveObjScore: aos,
		},
		func(hd t2g_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}

func (tw *Tower) NewRandFactor() [3]int64 {
	// st := 0
	lenPrimes := len(gamedata.Primes)
	rtn := [3]int64{}
loop:
	for i := 0; i < 3; {
		rtn[i] = gamedata.Primes[lenPrimes/2+tw.rnd.Intn(lenPrimes/2)]
		for j := 0; j < i; j++ {
			if rtn[j] == rtn[i] {
				continue loop
			}
		}
		i++
	}
	for i, _ := range rtn {
		if tw.rnd.Intn(2) == 0 {
			rtn[i] = -rtn[i]
		}
	}
	return rtn
}
