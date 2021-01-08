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

package terrain

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/game/terrain/corridor"
	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/goguelike/lib/scriptparse"
)

func cmdConnectRooms(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var connectCount int
	var allconnect bool
	var diagonal bool
	if err := ca.GetArgs(&tl, &connectCount, &allconnect, &diagonal); err != nil {
		return err
	}
	tr.connectRooms(connectCount, allconnect, tl, diagonal)
	return nil
}

func (tr *Terrain) connectRooms(connectCount int, allConnect bool, ttile tile_flag.TileFlag, way8 bool) {
	roomGraph := make(map[string]map[string]bool)
	roomMap := make(map[string]*room.Room)

	for _, v := range tr.roomManager.GetRoomList() {
		roomMap[v.GetUUID()] = v
	}

	for i := 0; i < connectCount; i++ {
		roomList := tr.roomManager.GetRoomList()
		tr.rnd.Shuffle(len(roomList), func(i, j int) {
			roomList[i], roomList[j] = roomList[j], roomList[i]
		})

		for _, v := range roomList {
			if len(roomGraph[v.UUID]) > connectCount {
				continue
			}
			nearRoom := tr.findNearRoom(v, roomGraph[v.UUID])
			if nearRoom != nil {
				if err := tr.makeCorridor(v, nearRoom, ttile, way8); err == nil {
					addRoomGraph(roomGraph, v, nearRoom)
				} else {
					tr.log.Warn("%v %v", err, tr)
				}
			} else {
				tr.log.Warn("fail to find near room %v %v", v, tr)
			}
		}
	}

	if !allConnect {
		return
	}
	// do connect all room
	for try := 1; ; {
		_, notvisitedRoom := tr.findNotConnectedRooms(roomGraph)
		if len(notvisitedRoom) != 0 {
			tr.log.Warn("%v rooms not fully connected %v, try more",
				len(notvisitedRoom), tr)
		} else {
			break
		}
		try++
		for i := range notvisitedRoom {
			v := roomMap[i]
			nearRoom := tr.findNearRoom(v, notvisitedRoom)
			if nearRoom != nil {
				if err := tr.makeCorridor(v, nearRoom, ttile, way8); err == nil {
					addRoomGraph(roomGraph, v, nearRoom)
					break
				} else {
					tr.log.Warn("%v %v", err, tr)
				}
			} else {
				tr.log.Warn("fail to find near room %v %v", v, tr)
			}
		}
	}
}

func (tr *Terrain) findNearRoom(r1 *room.Room, skipList map[string]bool) *room.Room {
	cx := r1.Area.X + r1.Area.W/2
	cy := r1.Area.Y + r1.Area.H/2
	roomArea := tr.roomManager.GetRoomArea()
	for _, v := range tr.findList {
		x, y := tr.WrapXY(cx+v.X, cy+v.Y)
		r := roomArea[x][y]
		if r != nil && r.GetUUID() != r1.UUID && !skipList[r.GetUUID()] {
			return r
		}
	}
	return nil
}

func addRoomGraph(roomGraph map[string]map[string]bool, r1, r2 *room.Room) {
	if roomGraph[r1.UUID] == nil {
		roomGraph[r1.UUID] = make(map[string]bool)
	}
	if roomGraph[r2.UUID] == nil {
		roomGraph[r2.UUID] = make(map[string]bool)
	}
	roomGraph[r1.UUID][r2.UUID] = true
	roomGraph[r2.UUID][r1.UUID] = true
}

func (tr *Terrain) findNotConnectedRooms(roomGraph map[string]map[string]bool) (map[string]bool, map[string]bool) {
	notvisitedRoom := make(map[string]bool)

	visitedRoom := make(map[string]bool)
	startroom := tr.roomManager.GetRoomList()[0].UUID
	traverseRoom(startroom, roomGraph, visitedRoom)
	for _, v := range tr.roomManager.GetRoomList() {
		if !visitedRoom[v.GetUUID()] {
			notvisitedRoom[v.GetUUID()] = true
		}
	}
	return visitedRoom, notvisitedRoom
}

func traverseRoom(startroom string, roomGraph map[string]map[string]bool, visitedRoom map[string]bool) {
	if visitedRoom[startroom] {
		return
	}
	visitedRoom[startroom] = true
	for j := range roomGraph[startroom] {
		if visitedRoom[j] {
			continue
		}
		traverseRoom(j, roomGraph, visitedRoom)
	}
}

func (tr *Terrain) makeCorridor(r1, r2 *room.Room, ttile tile_flag.TileFlag, way8 bool) error {
	p1 := r1.RndDoorOuter(tr.rnd)
	p2 := r2.RndDoorOuter(tr.rnd)

	plist := tr.findCorridor(p1[0], p1[1], p2[0], p2[1], (tr.Xlen * tr.Ylen), way8)
	if plist == nil {
		return fmt.Errorf("make corridor fail %v %v", p1, p2)
	}
	tr.corridorList = append(tr.corridorList, &corridor.Corridor{
		Tile: ttile,
		P:    plist,
	})
	return nil
}
