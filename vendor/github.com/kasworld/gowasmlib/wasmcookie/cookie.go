// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package wasmcookie manage browser cookie
package wasmcookie

import (
	"encoding/hex"
	"net/http"
	"strings"
	"syscall/js"

	"github.com/kasworld/gowasmlib/jslog"
)

func Set(ck *http.Cookie) {
	ck.Value = hex.EncodeToString([]byte(ck.Value))
	ckstr := ck.String()
	js.Global().Get("document").Set("cookie", ckstr)
}

func GetMap() map[string]string {
	rtn := make(map[string]string)

	fullstr := js.Global().Get("document").Get("cookie").String()
	slist := strings.Split(fullstr, ";")

	for _, v := range slist {
		v = strings.TrimSpace(v)
		nv := strings.Split(v, "=")
		if len(nv) != 2 {
			jslog.Errorf("invalid cookie %v", v)
			continue
		}
		name := strings.TrimSpace(nv[0])
		value, err := hex.DecodeString(strings.TrimSpace(nv[1]))
		if err != nil {
			jslog.Errorf("fail to decode cookie value %v %v", nv, err)
		}
		rtn[name] = string(value)
	}

	return rtn
}
