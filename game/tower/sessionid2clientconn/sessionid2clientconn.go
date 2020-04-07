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

// client connection manager
package sessionid2clientconn

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/rangestat"
	"github.com/kasworld/weblib"
)

func (cman SessionID2ClientConnection) String() string {
	return fmt.Sprintf(
		"%s[Count:%v %v]",
		cman.name,
		cman.Count(),
		cman.rstat,
	)
}

type SessionID2ClientConnection struct {
	mutex sync.RWMutex   `prettystring:"hide"`
	log   *g2log.LogBase `prettystring:"hide"`

	name                 string
	sessionID2Connection map[string]gamei.ServeClientConnI
	rstat                *rangestat.RangeStat
}

func New(name string, size int, l *g2log.LogBase) *SessionID2ClientConnection {
	if name == "" {
		name = "SessionID2ClientConnection"
	}
	rtn := &SessionID2ClientConnection{
		log:                  l,
		name:                 name,
		sessionID2Connection: make(map[string]gamei.ServeClientConnI),
		rstat:                rangestat.New("", 0, size),
	}
	return rtn
}

func (cman *SessionID2ClientConnection) Cleanup() {

}

func (cman *SessionID2ClientConnection) UpdateLap() {
	cman.rstat.UpdateLap()
}

func (cman *SessionID2ClientConnection) Count() int {
	return len(cman.sessionID2Connection)
}

func (cman *SessionID2ClientConnection) GetAllConnectionList() []gamei.ServeClientConnI {
	cman.mutex.RLock()
	defer cman.mutex.RUnlock()
	rtn := make([]gamei.ServeClientConnI, len(cman.sessionID2Connection))
	i := 0
	for _, v := range cman.sessionID2Connection {
		rtn[i] = v
		i++
	}
	return rtn
}

func (cman *SessionID2ClientConnection) GetRangeStat() *rangestat.RangeStat {
	return cman.rstat
}

////////////////////////////////////////////////////////////////////////////////

func (cman *SessionID2ClientConnection) GetBySessionID(key string) gamei.ServeClientConnI {

	cman.mutex.Lock()
	defer cman.mutex.Unlock()

	if oldc2cc, exist := cman.sessionID2Connection[key]; exist {
		return oldc2cc
	}
	return nil
}

func (cman *SessionID2ClientConnection) AddOrSwap(c2sc gamei.ServeClientConnI) (
	gamei.ServeClientConnI, error) {

	cman.mutex.Lock()
	defer cman.mutex.Unlock()

	key := c2sc.GetSession().GetUUID()

	oldc2cc, exist := cman.sessionID2Connection[key]
	if oldc2cc == c2sc {
		cman.log.Warn("add same connection. ignore, continue %v", c2sc)
		return nil, nil
	}
	if !exist {
		if !cman.rstat.Inc() {
			return nil, fmt.Errorf("connection Full %v", cman)
		}
	}
	cman.sessionID2Connection[key] = c2sc
	return oldc2cc, nil
}

func (cman *SessionID2ClientConnection) DelBySessionID(key string) (gamei.ServeClientConnI, error) {

	cman.mutex.Lock()
	defer cman.mutex.Unlock()

	oldc2cc, exist := cman.sessionID2Connection[key]
	if !exist {
		return nil, fmt.Errorf("connection not exist %v", key)
	}
	delete(cman.sessionID2Connection, key)

	if !cman.rstat.Dec() {
		cman.log.Fatal("Fail to Dec stat %v", cman)
	}
	return oldc2cc, nil
}

////////////////////////////////////////////////////////////////////////////////

func (cman *SessionID2ClientConnection) ToWeb(w http.ResponseWriter, r *http.Request) {
	weblib.WebFormBegin("client connection list", w, r)
	cman.ToWebMid(w, r)
	weblib.WebFormEnd(w, r)
}

func (cman *SessionID2ClientConnection) ToWebMid(w http.ResponseWriter, r *http.Request) {

	connList := cman.GetAllConnectionList()
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
	cman.log.Debug("Page %v,  %v-%v", page, st, ed)

	rtn := connList[st:ed]

	tplIndex, err := template.New("index").Parse(`
	<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table>
	<br/>
	`)
	if err != nil {
		cman.log.Error("%v", err)
	}
	if err := tplIndex.Execute(w, rtn); err != nil {
		cman.log.Error("%v", err)
	}
}

const (
	HTML_tableheader = `<tr>
<th>Remote Address</th>
<th>Session</th>
<th>Author</th>
<th>ActiveObj</th>
<th>Command</th>
</tr>`
	HTML_row = `<tr>
<td>{{$v.GetSession.RemoteAddr}}</td>
<td>{{$v.GetSession}}</td>
<td>{{$v.GetAuthorCmdList}}</td>
<td>{{$v.GetActiveObj}}</td>
<td><a href="/KickConnection?id={{$v.GetSession.SessionUUID}}" target="_blank">[Kick]</a></td>
</tr>
`
)
