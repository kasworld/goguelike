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

package aoexpsort

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/weblib"
)

type ByExp []gamei.ActiveObjectI

func (aol ByExp) Len() int { return len(aol) }
func (aol ByExp) Swap(i, j int) {
	aol[i], aol[j] = aol[j], aol[i]
}
func (aol ByExp) Less(i, j int) bool {
	ao1 := aol[i]
	ao2 := aol[j]
	return ao1.GetExpCopy() > ao2.GetExpCopy()
}
func (aol ByExp) Sort() {
	sort.Stable(aol)
}

func (aol ByExp) FindExist(ao gamei.ActiveObjectI) (int, bool) {
	findPos := sort.Search(
		len(aol),
		func(i int) bool {
			return aol[i].GetExpCopy() <= ao.GetExpCopy()
		},
	)

	if findPos < len(aol) && aol[findPos].GetUUID() == ao.GetUUID() {
		return findPos, true
	}
	return findPos, false
}

// if not found return len
func (aol ByExp) FindPosByKey(exp float64) int {
	findPos := sort.Search(
		len(aol),
		func(i int) bool {
			return aol[i].GetExpCopy() <= exp
		},
	)
	return findPos
}

func (aol *ByExp) Add(ao gamei.ActiveObjectI) error {
	// assume sorted
	_, exist := aol.FindExist(ao)
	if exist {
		return fmt.Errorf("exist %v", ao)
	}
	*aol = append(*aol, ao)
	return nil
}

func (aol *ByExp) Del(ao gamei.ActiveObjectI) error {
	// assume sorted
	pos, exist := aol.FindExist(ao)
	if !exist {
		return fmt.Errorf("not exist %v", ao)
	}
	*aol = append((*aol)[:pos], (*aol)[pos+1:]...)
	return nil
}

func (aol ByExp) GetPage(page int, pagesize int) ByExp {
	st := page * pagesize
	if st < 0 || st >= len(aol) {
		return nil
	}

	ed := st + pagesize
	if ed > len(aol) {
		ed = len(aol)
	}

	rtn := aol[st:ed]
	return rtn
}

func (aol ByExp) ToWeb(w http.ResponseWriter, r *http.Request) error {
	weblib.WebFormBegin("activeobject list", w, r)
	if err := aol.ToWebMid(w, r); err != nil {
		return err
	}
	weblib.WebFormEnd(w, r)
	return nil
}

func (aol ByExp) ToWebMid(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Parse(`
	<table border=1 style="border-collapse:collapse;">` +
		gamei.ActiveObjectI_HTML_tableheader +
		`{{range $i, $v := .}}` +
		gamei.ActiveObjectI_HTML_row +
		`{{end}}` +
		gamei.ActiveObjectI_HTML_tableheader +
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
