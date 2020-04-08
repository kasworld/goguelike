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

package weblib

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// type IMuxHandle interface {
// 	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
// }
type IAuthMuxHandle interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	HandleFuncAuth(pattern string, handler func(http.ResponseWriter, *http.Request))
}

func (am AuthMux) String() string {
	return fmt.Sprintf("AuthMux[AuthAct:%v NoAuthAct:%v %v]",
		len(am.AuthAction),
		len(am.NoAuthAction),
		am.authData)
}

type AuthMux struct {
	log loggerI `webformhide:"" stringformhide:""`

	*http.ServeMux
	authData     *AuthData
	AuthAction   []string
	NoAuthAction []string
}

func NewAuthMux(authdata *AuthData, l loggerI) *AuthMux {
	return &AuthMux{
		log:          l,
		ServeMux:     http.NewServeMux(),
		authData:     authdata,
		AuthAction:   make([]string, 0),
		NoAuthAction: make([]string, 0),
	}
}

func (am *AuthMux) HandleFuncAuth(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if err := am.authData.AddAction(pattern); err != nil {
		am.log.Warn("%v", err)
	}
	am.AuthAction = append(am.AuthAction, pattern)
	am.ServeMux.HandleFunc(pattern, am.Auth(pattern, handler))
}

// func (am *AuthMux) HandleFuncNoAuth(pattern string, handler func(http.ResponseWriter, *http.Request)) {
// 	am.NoAuthAction = append(am.NoAuthAction, pattern)
// 	am.ServeMux.HandleFunc(pattern, handler)
// }

func (am *AuthMux) Auth(pattern string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", am.authData.GetBasicAuthRealmString())
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			am.log.Error("Invalid Authenticate %v %v", pattern, s)
			http.Error(w, "Not authorized", 401)
			return
		}
		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			am.log.Error("Invalid Authenticate %v %v %v", pattern, b, err)
			http.Error(w, err.Error(), 401)
			return
		}
		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			am.log.Error("Invalid Authenticate %v %v", pattern, pair)
			http.Error(w, "Not authorized", 401)
			return
		}
		if err := am.authData.CheckLogin(pair[0], pair[1]); err != nil {
			am.log.Error("Check login fail %v %v", pattern, err)
			http.Error(w, "Not authorized", 401)
			return
		}
		if err := am.authData.CheckAction(pattern, pair[0]); err != nil {
			am.log.Error("Check action fail %v %v", pattern, err)
			http.Error(w, "Not authorized", 401)
			return
		}
		am.log.Debug("Do %v %v", pair[0], pattern)
		handler(w, r)
	}
}
