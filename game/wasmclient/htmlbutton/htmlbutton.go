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

// Package htmlbutton handle html button(group)
package htmlbutton

import (
	"bytes"
	"fmt"
	"net/url"
	"syscall/js"
)

type HTMLButton struct {
	// new args
	KeyCode    string
	IDBase     string
	ButtonText []string // state count
	ToolTip    string
	ClickFn    func(obj interface{}, v *HTMLButton)
	State      int

	// made obj
	JSID     string
	JSFnName string
	JSFn     func(this js.Value, args []js.Value) interface{}
}

func New(
	KeyCode string,
	IDBase string,
	ButtonText []string,
	ToolTip string,
	ClickFn func(obj interface{}, v *HTMLButton),
	State int,
) *HTMLButton {
	return &HTMLButton{
		KeyCode:    KeyCode,
		IDBase:     IDBase,
		ButtonText: ButtonText,
		ToolTip:    ToolTip,
		ClickFn:    ClickFn,
		State:      State,

		JSID:     "jsid_" + IDBase,
		JSFnName: "jsfn_" + IDBase,
	}
}

func (v *HTMLButton) RegisterJSFn(obj interface{}) {
	v.JSFn = func(this js.Value, args []js.Value) interface{} {
		v.State++
		v.State %= len(v.ButtonText)
		v.UpdateButtonText()
		if v.ClickFn != nil {
			v.ClickFn(obj, v)
		}
		return nil
	}
	js.Global().Set(v.JSFnName, js.FuncOf(v.JSFn))
}

func (v HTMLButton) JSButton() js.Value {
	return js.Global().Get("document").Call("getElementById", v.JSID)
}

func (v HTMLButton) UpdateButtonText() {
	btnStr := fmt.Sprintf("%s(%s)", v.ButtonText[v.State], v.KeyCode)
	v.JSButton().Set("innerHTML", btnStr)
}

func (v HTMLButton) Enable() {
	v.JSButton().Call("removeAttribute", "disabled")
}

func (v HTMLButton) Disable() {
	v.JSButton().Set("disabled", true)
}

func (v HTMLButton) Show() {
	v.JSButton().Get("style").Set("display", "initial")
}

func (v HTMLButton) Hide() {
	v.JSButton().Get("style").Set("display", "none")
}

type HTMLButtonGroup struct {
	Name           string
	ButtonList     []*HTMLButton
	ID2Button      map[string]*HTMLButton
	KeyCode2Button map[string]*HTMLButton
}

func NewButtonGroup(name string, buttons []*HTMLButton) *HTMLButtonGroup {
	rtn := &HTMLButtonGroup{
		Name:           name,
		ButtonList:     buttons,
		ID2Button:      make(map[string]*HTMLButton),
		KeyCode2Button: make(map[string]*HTMLButton),
	}
	for _, v := range buttons {
		rtn.ID2Button[v.IDBase] = v
		rtn.KeyCode2Button[v.KeyCode] = v
	}
	return rtn
}

func (hbl HTMLButtonGroup) GetByIDBase(idb string) *HTMLButton {
	return hbl.ID2Button[idb]
}

func (hbl HTMLButtonGroup) GetByKeyCode(kcode string) *HTMLButton {
	return hbl.KeyCode2Button[kcode]
}

func (hbl *HTMLButtonGroup) RegisterJSFn(obj interface{}) {
	for _, v := range hbl.ButtonList {
		v.RegisterJSFn(obj)
	}
}

func (hbl HTMLButtonGroup) MakeHTML(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%v", hbl.Name)
	for _, v := range hbl.ButtonList {
		fmt.Fprintf(buf,
			`<button id="%v" title="%v" onclick="%v()">%s(%s)</button>`,
			v.JSID, v.ToolTip, v.JSFnName, v.ButtonText[v.State], v.KeyCode,
		)
	}
}

func (hbl HTMLButtonGroup) MakeButtonToolTipTop(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%v", hbl.Name)
	for _, v := range hbl.ButtonList {
		fmt.Fprintf(buf,
			`<button class="tooltip" id="%v" onclick="%v()">%s(%s)
			<span class="tooltiptext-top">%v</span>
			</button>`,
			v.JSID, v.JSFnName, v.ButtonText[v.State], v.KeyCode, v.ToolTip,
		)
	}
	buf.WriteString("<br/>")
}

// SetFromURLArg set button state from url name=value
func (hbl *HTMLButtonGroup) SetFromURLArg() error {
	var err error
	qr, err := GetQuery()
	if err != nil {
		return err
	}
loop:
	for _, v := range hbl.ButtonList {
		optV := qr.Get(v.IDBase)
		if optV == "" {
			continue
		}
		for j, w := range v.ButtonText {
			if optV == w {
				v.State = j
				continue loop
			}
		}
		err = fmt.Errorf("invalid option %v %v", v.IDBase, optV)
	}
	return err
}

func GetQuery() (url.Values, error) {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	return u.Query(), err
}
