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

package c2t_serveconnbyte

import "github.com/kasworld/goguelike/lib/conndata"

const (
	HTML_tableheader = `<tr>
<th>UUID</th>
<th>Remote Address</th>
<th>Session</th>
<th>Author</th>
<th>Command</th>
</tr>`
	HTML_row = `<tr>
	<td>{{$v.WebConnData.UUID}}</td>
	<td>{{$v.WebConnData.RemoteAddr}}</td>
<td>{{$v.WebConnData.Session}}</td>
<td>{{$v.GetAuthorCmdList}}</td>
<td><a href="/KickConnection?id={{$v.WebConnData.Session.SessionUUID}}" target="_blank">[Kick]</a></td>
</tr>
`
)

func (scb *ServeConnByte) WebConnData() *conndata.ConnData {
	return scb.connData.(*conndata.ConnData)
}
