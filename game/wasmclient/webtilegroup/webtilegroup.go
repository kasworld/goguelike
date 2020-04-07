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

// Package webtilegroup tile group from json
package webtilegroup

import (
	"encoding/json"
	"net/http"
)

type TileInfo struct {
	Name string
	Rect Rect
}

type TileGroup map[string][]TileInfo

func GetTileGroup(u string) (TileGroup, error) {
	tlgp := make(TileGroup)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tlgp)
	if err != nil {
		return nil, err
	}
	return tlgp, nil
}
