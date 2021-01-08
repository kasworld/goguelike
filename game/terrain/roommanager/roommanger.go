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

package roommanager

import (
	"fmt"
	"sync"

	"github.com/kasworld/goguelike/game/terrain/room"
	"github.com/kasworld/wrapper"
)

type RoomManager struct {
	mutex sync.RWMutex `prettystring:"hide"`

	roomList []*room.Room     `prettystring:"simple"`
	roomArea [][]*room.Room   `prettystring:"simple"`
	XWrapper *wrapper.Wrapper `prettystring:"simple"`
	YWrapper *wrapper.Wrapper `prettystring:"simple"`
}

func New(w, h int) *RoomManager {
	rm := &RoomManager{
		roomList: make([]*room.Room, 0),
		roomArea: make([][]*room.Room, w),
		XWrapper: wrapper.New(w),
		YWrapper: wrapper.New(h),
	}
	for x := range rm.roomArea {
		rm.roomArea[x] = make([]*room.Room, h)
	}
	return rm
}

func (rm *RoomManager) GetCount() int {
	return len(rm.roomList)
}

func (rm *RoomManager) GetRoomList() []*room.Room {
	return rm.roomList
}

func (rm *RoomManager) GetRoomArea() [][]*room.Room {
	return rm.roomArea
}

func (rm *RoomManager) AddRoom(r *room.Room) error {
	// check space to add
	for x := r.Area.X; x < r.Area.X2(); x++ {
		for y := r.Area.Y; y < r.Area.Y2(); y++ {
			if rm.roomArea[rm.XWrapper.WrapSafe(x)][rm.YWrapper.WrapSafe(y)] != nil {
				return fmt.Errorf("No space to add room")
			}
		}
	}
	for x := r.Area.X; x < r.Area.X2(); x++ {
		for y := r.Area.Y; y < r.Area.Y2(); y++ {
			rm.roomArea[rm.XWrapper.WrapSafe(x)][rm.YWrapper.WrapSafe(y)] = r
		}
	}
	rm.roomList = append(rm.roomList, r)
	return nil
}

func (rm *RoomManager) GetRoomByPos(x, y int) *room.Room {
	return rm.roomArea[rm.XWrapper.WrapSafe(x)][rm.YWrapper.WrapSafe(y)]
}
