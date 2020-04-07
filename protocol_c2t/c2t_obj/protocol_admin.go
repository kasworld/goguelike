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

package c2t_obj

type ReqAdminFloorMove_data struct {
	Floor string
}
type RspAdminFloorMove_data struct {
	Dummy uint8
}

type ReqAdminTeleport_data struct {
	X int
	Y int
}
type RspAdminTeleport_data struct {
	Dummy uint8
}

type ReqAdminTowerCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminTowerCmd_data struct {
	Dummy uint8
}

type ReqAdminFloorCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminFloorCmd_data struct {
	Dummy uint8
}

type ReqAdminActiveObjCmd_data struct {
	Cmd string
	Arg string
}
type RspAdminActiveObjCmd_data struct {
	Dummy uint8
}

type ReqAIPlay_data struct {
	On bool
}
type RspAIPlay_data struct {
	Dummy uint8
}
