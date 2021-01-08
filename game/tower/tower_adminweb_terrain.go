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

package tower

import (
	"net/http"

	"github.com/kasworld/goguelike/game/gamei"
	"github.com/kasworld/weblib"
)

func (tw *Tower) web_TerrainInfo(w http.ResponseWriter, r *http.Request) {
	floorname := weblib.GetStringByName("floorname", "", w, r)
	if floorname == "" {
		tw.log.Warn("Invalid floorname")
		http.Error(w, "Invalid floorname", 404)
		return
	}
	move := weblib.GetStringByName("move", "", w, r)
	var err error
	var i int
	var f gamei.FloorI
	switch move {
	case "Before":
		i, err = tw.floorMan.GetFloorIndexByName(floorname)
		f = tw.floorMan.GetFloorByIndexWrap(i - 1)
	case "Next":
		i, err = tw.floorMan.GetFloorIndexByName(floorname)
		f = tw.floorMan.GetFloorByIndexWrap(i + 1)
	default:
		f = tw.floorMan.GetFloorByName(floorname)
	}
	if err != nil {
		tw.log.Warn("floor not found %v", err)
		http.Error(w, "floor not found", 404)
		return
	}
	if f == nil {
		tw.log.Warn("floor not found %v", floorname)
		http.Error(w, "floor not found", 404)
		return
	}
	if err := weblib.SetFresh(w, r); err != nil {
		tw.log.Error("%v", err)
	}
	f.GetTerrain().Web_TerrainInfo(w, r)
}

func (tw *Tower) web_TerrainImageZoom(w http.ResponseWriter, r *http.Request) {
	if f := tw.getFloorFromHTTPArg(w, r); f != nil {
		f.GetTerrain().Web_TerrainImageZoom(w, r)
	}
}

func (tw *Tower) web_TerrainImageAutoZoom(w http.ResponseWriter, r *http.Request) {
	if f := tw.getFloorFromHTTPArg(w, r); f != nil {
		f.GetTerrain().Web_TerrainImageAutoZoom(w, r)
	}
}

func (tw *Tower) web_TerrainTile(w http.ResponseWriter, r *http.Request) {
	if f := tw.getFloorFromHTTPArg(w, r); f != nil {
		f.GetTerrain().Web_TileInfo(w, r)
	}
}
