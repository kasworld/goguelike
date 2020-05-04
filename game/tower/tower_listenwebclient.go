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

package tower

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/groundconst"
	"github.com/kasworld/goguelike/config/towerwsurl"
	"github.com/kasworld/goguelike/enum/aotype"
	"github.com/kasworld/goguelike/game/aoexpsort"
	"github.com/kasworld/goguelike/game/aoscore"
	"github.com/kasworld/goguelike/game/cmd2tower"
	"github.com/kasworld/goguelike/game/conndata"
	"github.com/kasworld/goguelike/game/towerlist4client"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_authorize"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_msgp"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_serveconnbyte"
	"github.com/kasworld/uuidstr"
	"github.com/kasworld/weblib"
)

func (tw *Tower) initServiceWeb(ctx context.Context) {
	webMux := http.NewServeMux()
	if tw.sconfig.StandAlone {
		webMux.Handle("/",
			http.FileServer(http.Dir(tw.Config().ClientDataFolder)),
		)
		webMux.HandleFunc("/towerlist.json", tw.json_TowerList)
		webMux.HandleFunc("/highscore.json", tw.json_HighScore)
	}
	webMux.HandleFunc("/TowerInfo", tw.json_TowerInfo)
	webMux.HandleFunc("/ServiceInfo", tw.json_ServiceInfo)
	webMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		tw.serveWebSocketClient(ctx, w, r)
	})
	tw.log.TraceService("%v", webMux)
	tw.clientWeb = &http.Server{
		Handler: webMux,
		Addr:    fmt.Sprintf(":%v", tw.sconfig.ServicePort),
	}
	tw.demuxReq2BytesAPIFnMap = [c2t_idcmd.CommandID_Count]func(
		me interface{}, hd c2t_packet.Header, rbody []byte) (
		c2t_packet.Header, interface{}, error){
		c2t_idcmd.Invalid:           tw.bytesAPIFn_ReqInvalid,
		c2t_idcmd.Login:             tw.bytesAPIFn_ReqLogin,
		c2t_idcmd.Heartbeat:         tw.bytesAPIFn_ReqHeartbeat,
		c2t_idcmd.Chat:              tw.bytesAPIFn_ReqChat,
		c2t_idcmd.AchieveInfo:       tw.bytesAPIFn_ReqAchieveInfo,
		c2t_idcmd.Rebirth:           tw.bytesAPIFn_ReqRebirth,
		c2t_idcmd.Meditate:          tw.bytesAPIFn_ReqMeditate,
		c2t_idcmd.KillSelf:          tw.bytesAPIFn_ReqKillSelf,
		c2t_idcmd.Move:              tw.bytesAPIFn_ReqMove,
		c2t_idcmd.Attack:            tw.bytesAPIFn_ReqAttack,
		c2t_idcmd.Pickup:            tw.bytesAPIFn_ReqPickup,
		c2t_idcmd.Drop:              tw.bytesAPIFn_ReqDrop,
		c2t_idcmd.Equip:             tw.bytesAPIFn_ReqEquip,
		c2t_idcmd.UnEquip:           tw.bytesAPIFn_ReqUnEquip,
		c2t_idcmd.DrinkPotion:       tw.bytesAPIFn_ReqDrinkPotion,
		c2t_idcmd.ReadScroll:        tw.bytesAPIFn_ReqReadScroll,
		c2t_idcmd.Recycle:           tw.bytesAPIFn_ReqRecycle,
		c2t_idcmd.EnterPortal:       tw.bytesAPIFn_ReqEnterPortal,
		c2t_idcmd.MoveFloor:         tw.bytesAPIFn_ReqMoveFloor,
		c2t_idcmd.ActTeleport:       tw.bytesAPIFn_ReqActTeleport,
		c2t_idcmd.AdminTowerCmd:     tw.bytesAPIFn_ReqAdminTowerCmd,
		c2t_idcmd.AdminFloorCmd:     tw.bytesAPIFn_ReqAdminFloorCmd,
		c2t_idcmd.AdminActiveObjCmd: tw.bytesAPIFn_ReqAdminActiveObjCmd,
		c2t_idcmd.AdminFloorMove:    tw.bytesAPIFn_ReqAdminFloorMove,
		c2t_idcmd.AdminTeleport:     tw.bytesAPIFn_ReqAdminTeleport,
		c2t_idcmd.AIPlay:            tw.bytesAPIFn_ReqAIPlay,
	}
}

func CheckOrigin(r *http.Request) bool {
	return true
}

func (tw *Tower) serveWebSocketClient(ctx context.Context,
	w http.ResponseWriter, r *http.Request) {

	if tw.IsListenClientPaused() {
		tw.log.Warn("ListenClientPaused %v %v", w, r)
		return
	}

	if !tw.clientConnLimitStat.Inc() {
		tw.log.Fatal(
			"Over limit connect made, cancel %v, continue %v",
			r,
			tw.clientConnLimitStat)
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: CheckOrigin,
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		tw.log.Error("upgrade %v", err)
		return
	}

	tw.log.TraceClient("Start serveWebSocketClient %v", r.RemoteAddr)
	defer func() {
		tw.log.TraceClient("End serveWebSocketClient %v", r.RemoteAddr)
	}()

	connData := &conndata.ConnData{
		UUID:       uuidstr.New(),
		RemoteAddr: r.RemoteAddr,
	}
	c2sc := c2t_serveconnbyte.NewWithStats(
		connData,
		tw.floorMan.GetFloorCount()*2,
		c2t_authorize.NewPreLoginAuthorCmdIDList(),
		tw.sendStat, tw.recvStat,
		tw.protocolStat,
		tw.notiStat,
		tw.errorStat,
		tw.demuxReq2BytesAPIFnMap,
	)
	// add to conn manager
	tw.connManager.Add(connData.UUID, c2sc)

	c2sc.StartServeWS(ctx, wsConn,
		gameconst.ServerPacketReadTimeOutSec*time.Second,
		gameconst.ServerPacketWriteTimeoutSec*time.Second,
		c2t_msgp.MarshalBodyFn,
	)

	// connected user play

	// end play

	// connData changed in user play
	ao, exist := tw.id2ao.GetByUUID(connData.Session.ActiveObjUUID)
	if !exist {
		panic(fmt.Sprintf("ao not found %v", connData))
	}
	if ao != nil && ao.GetActiveObjType() == aotype.User {
		go tw.Ground_HighScore(ao)
		ao.Suspend()
		rspCh := make(chan error, 1)
		tw.GetReqCh() <- &cmd2tower.ActiveObjSuspendFromTower{
			ActiveObj: ao,
			RspCh:     rspCh,
		}
		<-rspCh
	}
	wsConn.Close()

	if !tw.clientConnLimitStat.Dec() {
		tw.log.Fatal("Under limit connection delete, continue %v",
			tw.clientConnLimitStat)
	}
	// del from conn manager
	tw.connManager.Del(connData.UUID)
}

func (tw *Tower) json_TowerList(w http.ResponseWriter, r *http.Request) {
	// baseurl := "http://localhost"
	baseurl := tw.sconfig.ServiceHostBase
	tn := tw.sconfig.TowerNumber

	cu, err := towerwsurl.MakeWebsocketURL(baseurl, tw.sconfig.ServicePort, tn-1)
	if err != nil {
		tw.log.Warn("fail to MakeWebsocketURL %v", err)
		http.Error(w, "fail to MakeWebsocketURL", 404)
		return
	}

	tl := []towerlist4client.TowerInfo2Enter{
		towerlist4client.TowerInfo2Enter{
			Name:       tw.sconfig.DisplayName,
			ConnectURL: cu,
		},
	}
	weblib.ServeJSON2HTTP(tl, w)
}

func (tw *Tower) json_HighScore(w http.ResponseWriter, r *http.Request) {
	allActiveObj := aoexpsort.ByExp(tw.aoExpRanking)
	aoLen := len(allActiveObj)
	if aoLen >= groundconst.HighScoreLen {
		aoLen = groundconst.HighScoreLen
	}
	aol := make([]*aoscore.ActiveObjScore, aoLen)
	for i := 0; i < aoLen; i++ {
		aol[i] = allActiveObj[i].To_ActiveObjScore()
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	weblib.ServeJSON2HTTP(aol, w)
}
