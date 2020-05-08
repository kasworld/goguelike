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

package session

import "time"

type Session struct {
	SessionUUID string
	ConnUUID    string
	RemoteAddr  string
	Create      time.Time `prettystring:"simple"`
	LastUse     time.Time `prettystring:"simple"`

	NickName      string
	ActiveObjUUID string
}

func (ac *Session) GetUUID() string {
	return ac.SessionUUID
}

const (
	HTML_tableheader = `<tr>
<th>Session</th>
<th>ConnUUID</th>
<th>RemoteAddr</th>
<th>Create</th>
<th>LastUse</th>
<th>NickName</th>
<th>ActiveObj</th>
<th>Command</th>
</tr>`
	HTML_row = `<tr>
<td>{{$v.SessionUUID}}</td>
<td>{{$v.ConnUUID}}</td>
<td>{{$v.RemoteAddr}}</td>
<td>{{$v.Create.Format "2006-01-02T15:04:05Z07:00"}}</td>
<td>{{$v.LastUse.Format "2006-01-02T15:04:05Z07:00"}}</td>
<td>{{$v.NickName}}</td>
<td>{{$v.ActiveObjUUID}}</td>
<td><a href="/DelSession?id={{$v.SessionUUID}}" target="_blank">[Del]</a></td>
</tr>
`
)
