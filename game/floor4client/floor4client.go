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

package floor4client

import (
	"github.com/kasworld/goguelike/game/visitarea"
	"github.com/kasworld/goguelike/lib/uuidposman_map"
	"github.com/kasworld/goguelike/lib/uuidposmani"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

// Floor4Client has floor info for server ai, client
type Floor4Client struct {
	Visit         *visitarea.VisitArea
	FOPosMan      uuidposmani.UUIDPosManI          // must *c2t_obj.FieldObjClient
	ActiveObjList []*c2t_obj.ActiveObjClient       // changed every turn
	CarryObjList  []*c2t_obj.CarryObjClientOnFloor // changed every turn
	DangerObjList []*c2t_obj.DangerObjClient       // changed every turn
}

func New(fi visitarea.FloorI) *Floor4Client {
	f4c := &Floor4Client{
		Visit:    visitarea.New(fi),
		FOPosMan: uuidposman_map.New(fi.GetWidth(), fi.GetHeight()),
	}
	return f4c
}

func (f4c *Floor4Client) Forget() {
	f4c.Visit.Forget()
	f4c.FOPosMan.Cleanup()
}

func (f4c *Floor4Client) UpdateObjLists(
	ActiveObjList []*c2t_obj.ActiveObjClient,
	CarryObjList []*c2t_obj.CarryObjClientOnFloor,
	FieldObjList []*c2t_obj.FieldObjClient,
	DangerObjList []*c2t_obj.DangerObjClient,
) {
	f4c.ActiveObjList = ActiveObjList
	f4c.CarryObjList = CarryObjList
	f4c.DangerObjList = DangerObjList
	for _, fo := range FieldObjList {
		f4c.FOPosMan.AddOrUpdateToXY(fo, fo.X, fo.Y)
	}
}

// for web
func (f4c *Floor4Client) GetName() string {
	return f4c.Visit.GetName()
}
