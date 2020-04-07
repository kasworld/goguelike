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
	"syscall/js"
)

type HTMLButton struct {
	KeyCode    string
	IDBase     string
	ButtonText []string // state count
	ToolTip    string

	ClickFn func(obj interface{}, v *HTMLButton)
	State   int
}

func (v HTMLButton) JSID() string {
	return "jsid_" + v.IDBase
}
func (v HTMLButton) JSFnName() string {
	return "jsfn_" + v.IDBase
}
func (v *HTMLButton) MakeJSFn(obj interface{}) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		v.State++
		v.State %= len(v.ButtonText)
		v.UpdateButtonText()
		if v.ClickFn != nil {
			v.ClickFn(obj, v)
		}
		return nil
	}
}

func (v HTMLButton) UpdateButtonText() {
	btnStr := fmt.Sprintf("%s(%s)", v.ButtonText[v.State], v.KeyCode)
	js.Global().Get("document").Call("getElementById", v.JSID()).Set("innerHTML", btnStr)
}

func (v HTMLButton) MakeHTML() string {
	btnStr := fmt.Sprintf("%s(%s)", v.ButtonText[v.State], v.KeyCode)
	return fmt.Sprintf(
		`<button id="%v" title="%v" onclick="%v()">%s</button>`,
		v.JSID(), v.ToolTip, v.JSFnName(), btnStr,
	)
}

func (v *HTMLButton) Enable() {
	btn := js.Global().Get("document").Call("getElementById", v.JSID())
	btn.Call("removeAttribute", "disabled")
}

func (v *HTMLButton) Disable() {
	btn := js.Global().Get("document").Call("getElementById", v.JSID())
	btn.Set("disabled", true)
}

func (v *HTMLButton) Show() {
	btn := js.Global().Get("document").Call("getElementById", v.JSID())
	btn.Get("style").Set("display", "initial")
}

func (v *HTMLButton) Hide() {
	btn := js.Global().Get("document").Call("getElementById", v.JSID())
	btn.Get("style").Set("display", "none")
}

func (v *HTMLButton) Register(obj interface{}) {
	js.Global().Set(v.JSFnName(), js.FuncOf(v.MakeJSFn(obj)))
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

func (hbl *HTMLButtonGroup) GetByIDBase(idb string) *HTMLButton {
	return hbl.ID2Button[idb]
}

func (hbl *HTMLButtonGroup) GetByKeyCode(kcode string) *HTMLButton {
	return hbl.KeyCode2Button[kcode]
}

func (hbl *HTMLButtonGroup) Register(obj interface{}) {
	for _, v := range hbl.ButtonList {
		v.Register(obj)
	}
}

func (hbl *HTMLButtonGroup) MakeHTML(obj interface{}) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%v", hbl.Name)
	for _, v := range hbl.ButtonList {
		buf.WriteString(v.MakeHTML())
	}
	return buf.String()
}

func (hbl *HTMLButtonGroup) MakeButtonToolTipTop(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%v", hbl.Name)
	for _, v := range hbl.ButtonList {
		fmt.Fprintf(buf,
			`<button class="tooltip" id="%v" onclick="%v()">%s(%s)
			<span class="tooltiptext-top">%v</span>
			</button>`,
			v.JSID(), v.JSFnName(), v.ButtonText[v.State], v.KeyCode, v.ToolTip,
		)
	}
	buf.WriteString("<br/>")
}
