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
	"net/http"
	"strings"
)

func Use(
	h http.HandlerFunc,
	middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {

	for _, m := range middleware {
		h = m(h)
	}
	return h
}

type NoAuth struct {
}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (na *NoAuth) Auth(h http.HandlerFunc) http.HandlerFunc {
	return h
	// return h.ServeHTTP
}

type BasicAuth struct {
	realm    string
	name     string
	password string
	realmstr string
}

func NewBasicAuth(realm string, name string, password string) *BasicAuth {
	return &BasicAuth{
		realm:    realm,
		name:     name,
		password: password,
		realmstr: `Basic realm="` + realm + `"`,
	}
}

// Leverages nemo's answer in http://stackoverflow.com/a/21937924/556573
func (ba *BasicAuth) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", ba.realmstr)
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}
		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}
		if pair[0] != ba.name || pair[1] != ba.password {
			http.Error(w, "Not authorized", 401)
			return
		}
		h(w, r)
		// h.ServeHTTP(w, r)
	}
}
