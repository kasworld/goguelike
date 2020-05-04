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

package c2t_connbytemanager

import (
	"net/http"
	"text/template"

	"github.com/kasworld/goguelike/protocol_c2t/c2t_serveconnbyte"
	"github.com/kasworld/weblib"
)

func (cman *Manager) ToWeb(w http.ResponseWriter, r *http.Request) {
	weblib.WebFormBegin("client connection list", w, r)
	cman.ToWebMid(w, r)
	weblib.WebFormEnd(w, r)
}

func (cman *Manager) ToWebMid(w http.ResponseWriter, r *http.Request) {

	connList := cman.GetList()
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

	rtn := connList[st:ed]

	tplIndex, err := template.New("index").Parse(`
	<table border=1 style="border-collapse:collapse;">` +
		c2t_serveconnbyte.HTML_tableheader +
		`{{range $i, $v := .}}` +
		c2t_serveconnbyte.HTML_row +
		`{{end}}` +
		c2t_serveconnbyte.HTML_tableheader +
		`</table>
	<br/>
	`)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	if err := tplIndex.Execute(w, rtn); err != nil {
		http.Error(w, err.Error(), 404)
	}
}
