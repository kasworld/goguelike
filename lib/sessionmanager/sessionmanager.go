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

package sessionmanager

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/lib/session"
	"github.com/kasworld/rangestat"
	"github.com/kasworld/uuidstr"
	"github.com/kasworld/weblib"
)

func (sman SessionManager) String() string {
	return fmt.Sprintf(
		"%s[Count:%v %v]",
		sman.name,
		sman.Count(),
		sman.rstat,
	)
}

type SessionManager struct {
	log   *g2log.LogBase `prettystring:"hide"`
	rnd   *g2rand.G2Rand `prettystring:"hide"`
	name  string
	rstat *rangestat.RangeStat `prettystring:"simple"`

	mutex             sync.RWMutex                `prettystring:"hide"`
	sessionid2session map[string]*session.Session `prettystring:"simple"`
}

func New(name string, size int, l *g2log.LogBase) *SessionManager {
	if name == "" {
		name = "SessionManager"
	}
	sman := &SessionManager{
		log:               l,
		name:              name,
		rstat:             rangestat.New("", 0, size),
		rnd:               g2rand.New(),
		sessionid2session: make(map[string]*session.Session),
	}
	return sman
}

func (sman *SessionManager) Cleanup() {

}

func (sman *SessionManager) UpdateLap() {
	sman.rstat.UpdateLap()
}

func (sman *SessionManager) Count() int {
	return len(sman.sessionid2session)
}

func (sman *SessionManager) GetAllList() []*session.Session {
	sman.mutex.RLock()
	defer sman.mutex.RUnlock()
	rtn := make([]*session.Session, len(sman.sessionid2session))
	i := 0
	for _, v := range sman.sessionid2session {
		rtn[i] = v
		i++
	}
	return rtn
}

func (sman *SessionManager) GetRangeStat() *rangestat.RangeStat {
	return sman.rstat
}

////////////////////////////////////////////////////////////////////////////////

func (sman *SessionManager) GetBySessionID(id string) *session.Session {
	sman.mutex.RLock()
	defer sman.mutex.RUnlock()
	if ss, exist := sman.sessionid2session[id]; exist {
		ss.LastUse = time.Now()
		return ss
	}
	return nil
}

func (sman *SessionManager) DelBySessionID(id string) {
	sman.mutex.Lock()
	defer sman.mutex.Unlock()
	if _, exist := sman.sessionid2session[id]; exist {
		delete(sman.sessionid2session, id)
		sman.rstat.Dec()
	}
}

func (sman *SessionManager) UpdateOrNew(
	sessionuuid string,
	remoteaddr string,
	nickname string,
) *session.Session {

	sessionuuid = strings.TrimSpace(sessionuuid)
	now := time.Now()
	pad := fmt.Sprintf("_%08x", sman.rnd.Uint32())
	nickname = ValidatePlayername(nickname, pad)

	sman.mutex.Lock()
	defer sman.mutex.Unlock()

	if sessionuuid != "" {
		ss, exist := sman.sessionid2session[sessionuuid]
		if exist {
			ss.LastUse = now
			ss.RemoteAddr = remoteaddr
			ss.NickName = nickname
			sman.log.Debug("session reuse %v", ss)
			return ss
		}
	}
	ss := &session.Session{
		SessionUUID: uuidstr.New(),
		Create:      now,
		LastUse:     now,
		RemoteAddr:  remoteaddr,
		NickName:    nickname,
	}

	sman.sessionid2session[ss.SessionUUID] = ss
	if !sman.rstat.Inc() {
		sman.log.Fatal("Fail to Inc stat %v", sman)
	}
	return ss
}

////////////////////////////////////////////////////////////////////////////////

func (sman *SessionManager) ToWeb(w http.ResponseWriter, r *http.Request) {
	weblib.WebFormBegin("session list", w, r)
	sman.ToWebMid(w, r)
	weblib.WebFormEnd(w, r)
}

func (sman *SessionManager) ToWebMid(w http.ResponseWriter, r *http.Request) {

	connList := sman.GetAllList()
	page := weblib.GetIntByName("page", -1, w, r)
	if page < 0 {
		return
	}
	pagesize := 20
	st := page * pagesize
	if st < 0 || st >= len(connList) {
		st = 0
	}

	ed := st + pagesize
	if ed > len(connList) {
		ed = len(connList)
	}
	sman.log.Debug("Page %v,  %v-%v", page, st, ed)

	rtn := connList[st:ed]

	tplIndex, err := template.New("index").Parse(`
	<table border=1 style="border-collapse:collapse;">` +
		session.HTML_tableheader +
		`{{range $i, $v := .}}` +
		session.HTML_row +
		`{{end}}` +
		session.HTML_tableheader +
		`</table>
	<br/>
	`)
	if err != nil {
		sman.log.Error("%v", err)
	}
	if err := tplIndex.Execute(w, rtn); err != nil {
		sman.log.Error("%v", err)
	}
}
