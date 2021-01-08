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

// Package aoscore score data for ground server
package aoscore

import (
	"html/template"
	"net/http"

	"github.com/kasworld/weblib"
)

func (aol ActiveObjScoreList) ToWeb(w http.ResponseWriter, r *http.Request) error {
	weblib.WebFormBegin("activeobject list", w, r)
	if err := aol.ToWebMid(w, r); err != nil {
		return err
	}
	weblib.WebFormEnd(w, r)
	return nil
}

func (aol ActiveObjScoreList) ToWebMid(w http.ResponseWriter, r *http.Request) error {
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
		return err
	}
	if err := tplIndex.Execute(w, aol); err != nil {
		return err
	}
	return nil
}

const (
	HTML_tableheader = `
	<tr>
	<td>Name</td>
	<td>Exp</td>
	<td>Wealth</td>
	<td>Born</td>
	<td>Bias</td>
	<td>Tower</td>
	<td>Date</td>
	</tr>	
`
	HTML_row = `
	<tr>
		<td>{{$v.NickName}}</td>
		<td>{{printf "%.2f" $v.Exp}}</td>
		<td>{{printf "%.2f" $v.Wealth}}</td>
		<td>{{$v.BornFaction}}</td>
		<td>{{$v.CurrentBias}}</td>
		<td>{{$v.TowerName}}</td>
		<td>{{$v.RecordTime.Format "2006-01-02T15:04:05Z07:00"}}</td>
	</tr>
`
)
