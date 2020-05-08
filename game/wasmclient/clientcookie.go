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
	"fmt"
	"net/http"
	"syscall/js"
	"time"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/gowasmlib/wasmcookie"
)

func sessionKeyName(towerindex int) string {
	return fmt.Sprintf("sessionkey%d", towerindex)
}

func ClearSession(towerindex int) {
	wasmcookie.Set(&http.Cookie{
		Name:    sessionKeyName(towerindex),
		Value:   "",
		Path:    "/",
		Expires: time.Now().AddDate(1, 0, 0),
	})
}

func SetSession(towerindex int, sessionkey string, nick string) {
	wasmcookie.Set(&http.Cookie{
		Name:    sessionKeyName(towerindex),
		Value:   sessionkey,
		Path:    "/",
		Expires: time.Now().AddDate(1, 0, 0),
	})
	wasmcookie.Set(&http.Cookie{
		Name:    "nickname",
		Value:   nick,
		Path:    "/",
		Expires: time.Now().AddDate(1, 0, 0),
	})
}

func InitNickname() {
	ck := wasmcookie.GetMap()
	var nickname string
	if oldnick, exist := ck["nickname"]; exist {
		nickname = oldnick
	} else {
		nickname = fmt.Sprintf("unnamed_%x", g2rand.New().Uint32())
	}
	js.Global().Get("document").Call("getElementById", "nickname").Set("value", nickname)
}
