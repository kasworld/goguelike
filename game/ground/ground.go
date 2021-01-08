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

package ground

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/configutil"
	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/config/dataversion"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/groundconfig"
	"github.com/kasworld/goguelike/config/groundconst"
	"github.com/kasworld/goguelike/config/towerdata"
	"github.com/kasworld/goguelike/enum/factiontype"
	"github.com/kasworld/goguelike/game/aoscore"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_version"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_connbytemanager"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_idcmd"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_packet"
	"github.com/kasworld/launcherlib"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/recordduration"
	"github.com/kasworld/version"
	"github.com/kasworld/weblib/retrylistenandserve"
)

func (grd *Ground) String() string {
	return fmt.Sprintf("Ground[]")
}

func (grd *Ground) StringForm() string {
	return prettystring.PrettyString(grd, 4)
}

type Ground struct {
	mutex   sync.RWMutex   `prettystring:"hide"`
	DoClose func()         `prettystring:"hide"`
	rnd     *g2rand.G2Rand `prettystring:"hide"`
	log     *g2log.LogBase `prettystring:"hide"`

	sconfig     *groundconfig.GroundConfig
	serviceInfo c2t_obj.ServiceInfo
	twMan       *TowerManager

	mutexHighScore sync.RWMutex               `prettystring:"hide"`
	highScore      aoscore.ActiveObjScoreList `prettystring:"simple"`

	RecvStat *actpersec.ActPerSec `prettystring:"simple"`
	SendStat *actpersec.ActPerSec `prettystring:"simple"`

	adminWeb  *http.Server `prettystring:"simple"`
	clientWeb *http.Server `prettystring:"simple"`

	connManager            *t2g_connbytemanager.Manager
	demuxReq2BytesAPIFnMap [t2g_idcmd.CommandID_Count]func(
		me interface{}, hd t2g_packet.Header, rbody []byte) (
		t2g_packet.Header, interface{}, error) `prettystring:"hide"`
}

func New(config *groundconfig.GroundConfig) *Ground {
	fmt.Printf("%v", config.StringForm())

	if config.BaseLogDir != "" {
		log, err := g2log.NewWithDstDir(
			"ground",
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

	grd := Ground{
		rnd:         g2rand.New(),
		sconfig:     config,
		log:         g2log.GlobalLogger,
		SendStat:    actpersec.New(),
		RecvStat:    actpersec.New(),
		connManager: t2g_connbytemanager.New(),
	}
	grd.demuxReq2BytesAPIFnMap = [t2g_idcmd.CommandID_Count]func(
		me interface{}, hd t2g_packet.Header, rbody []byte) (
		t2g_packet.Header, interface{}, error){
		t2g_idcmd.Invalid:   grd.bytesAPIFn_ReqInvalid,
		t2g_idcmd.Register:  grd.bytesAPIFn_ReqRegister,
		t2g_idcmd.Heartbeat: grd.bytesAPIFn_ReqHeartbeat,
		t2g_idcmd.HighScore: grd.bytesAPIFn_ReqHighScore,
	} // DemuxReq2BytesAPIFnMap

	// grd.log = g2log.GlobalLogger
	grd.DoClose = func() {
		grd.log.Fatal("Too early DoClose call %v", grd)
	}
	grd.serviceInfo = c2t_obj.ServiceInfo{
		Version:         version.GetVersion(),
		ProtocolVersion: c2t_version.ProtocolVersion,
		DataVersion:     dataversion.DataVersion,
	}
	return &grd
}

func (grd *Ground) GetLogger() interface{} {
	return grd.log
}

func (grd *Ground) GetServiceLockFilename() string {
	return grd.sconfig.MakePIDFileFullpath()
}

func (grd *Ground) ServiceInit() error {
	rd := recordduration.New("GroundInit")

	grdlog, err := g2log.NewWithDstDir("ground",
		grd.sconfig.MakeLogDir(),
		logflags.DefaultValue(false).BitClear(logflags.LF_functionname),
		grd.sconfig.LogLevel,
		grd.sconfig.SplitLogLevel,
	)
	if err != nil {
		return err
	}
	grd.log = grdlog

	grd.log.TraceService("Start ServiceInit %v %v", grd, rd)
	defer func() { grd.log.TraceService("End ServiceInit %v %v", grd, rd) }()

	grd.log.TraceService("%v", grd.serviceInfo.StringForm())
	grd.log.TraceService("%v", grd.sconfig.StringForm())

	var towerdatalist []towerdata.TowerData
	if err := configutil.LoadJSON(grd.sconfig.MakeTowerDataFileFullpath(), &towerdatalist); err == nil {
		grd.log.Debug("load towerdatalist %v", grd.sconfig.MakeTowerDataFileFullpath())
		fmt.Printf("load towerdatalist %v\n", grd.sconfig.MakeTowerDataFileFullpath())
		grd.twMan = NewTowerManager(grd.sconfig, grd.log, towerdatalist)
	} else {
		grd.log.Debug("save towerdatalist %v", grd.sconfig.MakeTowerDataFileFullpath())
		fmt.Printf("save towerdatalist %v\n", grd.sconfig.MakeTowerDataFileFullpath())
		configutil.SaveJSON(grd.sconfig.MakeTowerDataFileFullpath(), towerdata.Default)
		grd.twMan = NewTowerManager(grd.sconfig, grd.log, towerdata.Default)
	}

	if err := grd.highScore.LoadJSON(grd.sconfig.MakeHighScoreFileFullpath()); err != nil {
		grd.log.Warn("fail to load high score %v, make new %v",
			err, grd.sconfig.MakeHighScoreFileFullpath())

		for lv := 1; lv <= 10; lv++ {
			bornFaction := factiontype.FactionType(grd.rnd.Intn(factiontype.FactionType_Count))
			CurrentBias := bias.Bias{
				grd.rnd.Float64() - 0.5,
				grd.rnd.Float64() - 0.5,
				grd.rnd.Float64() - 0.5,
			}.MakeAbsSumTo(gameconst.ActiveObjBaseBiasLen)

			grd.highScore = append(grd.highScore,
				aoscore.NewActiveObjScoreByLevel(lv, bornFaction, CurrentBias))
		}

		grd.highScore.SortByExp()
		grd.highScore.TrimLen(groundconst.HighScoreLen)

		if err := grd.highScore.SaveJSON(grd.sconfig.MakeHighScoreFileFullpath()); err != nil {
			grd.log.Error("fail to save high score %v %v",
				grd.sconfig.MakeHighScoreFileFullpath(), err)
		}
	}

	grd.initAdminWeb()
	grd.initServiceWeb()

	fmt.Printf("WebAdmin  : http://%v:%v\n", grd.sconfig.GroundHost, grd.sconfig.GroundAdminWebPort)
	fmt.Printf("WebClient : http://%v:%v\n", grd.sconfig.GroundHost, grd.sconfig.GroundServiceWebPort)

	return nil
}

func (grd *Ground) ServiceMain(mainctx context.Context) {
	grd.log.TraceService("Start ServiceMain %v", grd)
	defer func() { grd.log.TraceService("End ServiceMain %v", grd) }()

	ctx, closeCtx := context.WithCancel(mainctx)
	grd.DoClose = closeCtx

	go retrylistenandserve.RetryListenAndServe(grd.adminWeb, grd.log, "serveAdminWeb")
	go retrylistenandserve.RetryListenAndServe(grd.clientWeb, grd.log, "serveServiceWeb")
	go grd.listenTowerTCP(ctx)

	tw2autostart := grd.twMan.GetAutoStartTowerList()
	for _, v := range tw2autostart {
		err := grd.ControlTower(v, "start")
		if err != nil {
			grd.log.Fatal("fail to start tower %v", v.TowerConfigMade.TowerName)
		} else {
			grd.log.Debug("start tower %v", v.TowerConfigMade.TowerName)
			fmt.Printf("start tower %v\n", v.TowerConfigMade.TowerName)
		}
	}

loop:
	for {
		select {
		case <-ctx.Done():
			grd.log.TraceService("AdminWeb Close %v", grd.adminWeb.Close())
			grd.log.TraceService("ClientWeb Close %v", grd.clientWeb.Close())
			break loop
		}
	}
}

func (grd *Ground) ServiceCleanup() {
	grd.log.TraceService("Start ServiceCleanup %v", grd)
	defer func() { grd.log.TraceService("End ServiceCleanup %v", grd) }()
}

func (grd *Ground) GetTowerManager() *TowerManager {
	return grd.twMan
}

func (grd *Ground) GetConfig() *groundconfig.GroundConfig {
	return grd.sconfig
}

func (grd *Ground) GetServiceInfo() c2t_obj.ServiceInfo {
	return grd.serviceInfo
}

func (grd *Ground) AddActiveObj2HighScoreAndSort(aos *aoscore.ActiveObjScore) {
	scoreFilename := grd.sconfig.MakeHighScoreFileFullpath()

	grd.mutexHighScore.Lock()
	grd.highScore.AddOrUpdate(aos)
	grd.highScore.SortByExp()
	grd.highScore.TrimLen(groundconst.HighScoreLen)
	err := grd.highScore.SaveJSON(scoreFilename)
	grd.mutexHighScore.Unlock()

	if err != nil {
		grd.log.Error("fail to save high score %v %v",
			scoreFilename, err)
	}
}

// ControlTower control tower process
// cmd : start,stop,restart,forcestart,logreopen (default "start")
func (grd *Ground) ControlTower(te *TowerRunning, cmd string) error {
	exe := grd.sconfig.MakeTowerExeFullpath()
	args := []string{
		fmt.Sprintf("-service=%s", cmd),
		"-i",
		te.MakeTowerConfigDownURL(),
	}
	outfile := te.TowerConfigMade.MakeOutfileFullpath()
	grd.log.Debug("%v %v", exe, args)
	go func() {
		err := launcherlib.RunExeSync(
			exe,
			args,
			outfile,
		)
		if err != nil {
			grd.log.Error("%v", err)
		}
	}()
	return nil
}

func (grd *Ground) BuildDate() time.Time {
	return version.GetBuildDate()
}

func (grd *Ground) NumGoroutine() int {
	return runtime.NumGoroutine()
}
