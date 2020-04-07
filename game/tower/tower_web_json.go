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

package tower

import (
	"net/http"

	"github.com/kasworld/weblib"
)

func (tw *Tower) json_TowerInfo(w http.ResponseWriter, r *http.Request) {
	//Access-Control-Allow-Origin: *
	w.Header().Set("Access-Control-Allow-Origin", "*")
	weblib.ServeJSON2HTTP(tw.towerInfo, w)
}

func (tw *Tower) json_ServiceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	weblib.ServeJSON2HTTP(tw.serviceInfo, w)
}

func (tw *Tower) json_Config(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	weblib.ServeJSON2HTTP(tw.sconfig, w)
}
