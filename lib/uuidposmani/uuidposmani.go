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

// positioned object managment in 2d space
package uuidposmani

import "github.com/kasworld/findnear"

type UUIDPosI interface {
	GetUUID() string
}

type UUIDPosIList []UUIDPosI

// type DoFn func(fo UUIDPosI) bool

// VPIXYObj use for viewport by xylenlist,
type VPIXYObj struct {
	I    int // xylen pos
	X, Y int // pos in uuidposmani.
	O    UUIDPosI
}

type UUIDPosManI interface {
	Cleanup()
	Count() int
	Wrap(x, y int) (int, int)
	GetAllList() UUIDPosIList
	IterAll(iterfn func(o UUIDPosI, x, y int) bool)
	GetByUUID(id string) UUIDPosI
	GetByXYAndUUID(id string, x, y int) (UUIDPosI, error)
	GetXYByUUID(id string) (int, int, bool)
	Get1stObjAt(x, y int) UUIDPosI
	GetObjListAt(x, y int) UUIDPosIList
	Search1stByXYLenList(
		xylenlist findnear.XYLenList,
		sx, sy int,
		filterfn func(o UUIDPosI, x, y int, xylen findnear.XYLen) bool) (UUIDPosI, int, int)
	GetVPIXYObjByXYLenList(
		xylenlist findnear.XYLenList,
		sx, sy int, l int) []VPIXYObj
	IterByXYLenList(
		xylenlist findnear.XYLenList,
		sx, sy int, l int,
		stopFn func(o UUIDPosI, x, y int, i int, xylen findnear.XYLen) bool)
	AddOrUpdateToXY(o UUIDPosI, x, y int) error
	AddToXY(o UUIDPosI, x, y int) error
	DelByFilter(filter func(o UUIDPosI, x, y int) bool) error
	Del(o UUIDPosI) error
	UpdateToXY(o UUIDPosI, newx, newy int) error
}
