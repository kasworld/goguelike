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

package wasmclient

import (
	"syscall/js"
)

type UITextObj struct {
	leftinfo   js.Value
	rightinfo  js.Value
	centerinfo js.Value
}

func NewUITextObj() *UITextObj {
	uio := &UITextObj{}

	uio.leftinfo = js.Global().Get("document").Call("getElementById", "leftinfo")
	uio.rightinfo = js.Global().Get("document").Call("getElementById", "rightinfo")
	JSObjSetColor(uio.rightinfo, "white")

	uio.centerinfo = js.Global().Get("document").Call("getElementById", "centerinfo")
	JSObjSetColor(uio.centerinfo, "white")

	return uio
}

func JSObjShow(jsobj js.Value) {
	jsobj.Get("style").Set("display", "initial")
}

func JSObjHide(jsobj js.Value) {
	jsobj.Get("style").Set("display", "none")
}

func JSObjSetColor(jsobj js.Value, color string) {
	jsobj.Get("style").Set("color", color)
}

func JSObjSetBGColor(jsobj js.Value, color string) {
	jsobj.Get("style").Set("backgroundColor", color)
}

func GetTextValueFromInputText(jsobjid string) string {
	jsobj := js.Global().Get("document").Call("getElementById", jsobjid).Get("value")
	return jsobj.String()
}
