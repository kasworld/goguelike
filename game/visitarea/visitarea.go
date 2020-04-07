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

package visitarea

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"sync"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/wrapper"
)

type BitContainder uint8

const FullBit = 0xff
const BitLen = 8
const MaskBit = BitLen - 1
const PosShift = 3

type floorI interface {
	GetName() string
	GetUUID() string
	GetWidth() int
	GetHeight() int
	VisitableCount() int
}

type VisitArea struct {
	mutex sync.RWMutex `prettystring:"hide"`
	name  string
	uuid  string
	w     int
	h     int
	xWrap func(i int) int
	yWrap func(i int) int

	discoverExp         int
	discoveredTileCount int
	completeTileCount   int
	bitsList            []BitContainder
}

func NewVisitArea(fi floorI) *VisitArea {
	va := &VisitArea{
		name:              fi.GetName(),
		uuid:              fi.GetUUID(),
		w:                 fi.GetWidth(),
		h:                 fi.GetHeight(),
		completeTileCount: fi.VisitableCount(),
	}
	va.xWrap = wrapper.New(va.w).GetWrapFn()
	va.yWrap = wrapper.New(va.h).GetWrapFn()
	va.bitsList = make([]BitContainder, va.w*va.h/BitLen+1)
	return va
}

func (va *VisitArea) Cleanup() {
	va.mutex.Lock()
	defer va.mutex.Unlock()

}

func (va *VisitArea) Dup() *VisitArea {
	va.mutex.RLock()
	defer va.mutex.RUnlock()
	rtn := &VisitArea{
		name:                va.name,
		w:                   va.w,
		h:                   va.h,
		xWrap:               va.xWrap,
		yWrap:               va.yWrap,
		discoveredTileCount: va.discoveredTileCount,
		completeTileCount:   va.completeTileCount,
		bitsList:            make([]BitContainder, len(va.bitsList)),
	}
	copy(rtn.bitsList, va.bitsList)
	return rtn
}

func (va *VisitArea) String() string {
	return fmt.Sprintf("%v[%v/%v]",
		va.name, va.discoveredTileCount, va.completeTileCount)
}

func (va *VisitArea) GetName() string {
	return va.name
}

func (va *VisitArea) GetUUID() string {
	return va.uuid
}

func (va *VisitArea) GetDiscoveredTileCount() int {
	va.mutex.RLock()
	defer va.mutex.RUnlock()
	return va.discoveredTileCount
}
func (va *VisitArea) GetCompleteTileCount() int {
	va.mutex.RLock()
	defer va.mutex.RUnlock()
	return va.completeTileCount
}

func (va *VisitArea) CalcCompleteRate() float64 {
	va.mutex.RLock()
	defer va.mutex.RUnlock()
	return float64(va.discoveredTileCount) / float64(va.completeTileCount)
}

func (va *VisitArea) IsComplete() bool {
	va.mutex.RLock()
	defer va.mutex.RUnlock()
	return va.IsCompleteNoLock()
}

func (va *VisitArea) IsCompleteNoLock() bool {
	return va.discoveredTileCount >= va.completeTileCount
}

func (va *VisitArea) CheckAndSetNolock(x, y int) {
	index, bit := va.calcIndexAndBitForXY(x, y)
	if va.bitsList[index]&bit == 0 { // if not set
		va.bitsList[index] |= bit
		va.discoveredTileCount++
	}
}

func (va *VisitArea) MakeComplete() int {
	va.mutex.Lock()
	defer va.mutex.Unlock()
	inc := va.completeTileCount - va.discoveredTileCount
	if inc <= 0 {
		return 0
	}
	for i := range va.bitsList {
		va.bitsList[i] = FullBit
	}
	va.discoveredTileCount = va.completeTileCount
	va.updateExp()
	return inc
}

func (va *VisitArea) calcIndexAndBitForXY(x, y int) (uint, BitContainder) {
	pos := uint(x + y*va.w)
	indexPos := pos >> PosShift
	bitPos := pos & MaskBit
	return indexPos, 1 << bitPos
}
func (va *VisitArea) SetXYNolock(x, y int) {
	index, bit := va.calcIndexAndBitForXY(x, y)
	va.bitsList[index] |= bit
}

func (va *VisitArea) GetXYNolock(x, y int) bool {
	index, bit := va.calcIndexAndBitForXY(x, y)
	return va.bitsList[index]&bit != 0
}

func (va *VisitArea) UpdateByViewport2(
	vpCenterX, vpCenterY int,
	vpTiles *viewportdata.ViewportTileArea2) int {
	va.mutex.Lock()
	defer va.mutex.Unlock()

	xwrapper := va.xWrap
	ywrapper := va.yWrap

	old := va.discoveredTileCount
	for i, v := range viewportdata.ViewportXYLenList {
		fx := xwrapper(v.X + vpCenterX)
		fy := ywrapper(v.Y + vpCenterY)
		if vpTiles[i] != 0 {
			index, bit := va.calcIndexAndBitForXY(fx, fy)
			if va.bitsList[index]&bit == 0 { // if not set
				va.bitsList[index] |= bit
				va.discoveredTileCount++
			}
		}
	}
	va.updateExp()
	return va.discoveredTileCount - old
}

func (va *VisitArea) UpdateBySightMat2(
	floorTiles tilearea.TileArea,
	vpCenterX, vpCenterY int,
	sightMat *viewportdata.ViewportSight2,
	sight float32) int {
	va.mutex.Lock()
	defer va.mutex.Unlock()

	xwrapper := va.xWrap
	ywrapper := va.yWrap

	old := va.discoveredTileCount
	for i, v := range viewportdata.ViewportXYLenList {
		fx := xwrapper(v.X + vpCenterX)
		fy := ywrapper(v.Y + vpCenterY)
		if floorTiles[fx][fy] != 0 && sightMat[i] <= sight {
			index, bit := va.calcIndexAndBitForXY(fx, fy)
			if va.bitsList[index]&bit == 0 { // if not set
				va.bitsList[index] |= bit
				va.discoveredTileCount++
			}
		}
	}
	va.updateExp()
	return va.discoveredTileCount - old
}

func (va *VisitArea) GetDiscoverExp() int {
	return va.discoverExp
}
func (va *VisitArea) updateExp() {
	discoverExp := va.discoveredTileCount * gameconst.ActiveObjExp_DiscoverTile
	findExp := gameconst.ActiveObjExp_VisitFloor
	completeExp := 0
	if va.IsCompleteNoLock() {
		completeExp = va.completeTileCount * gameconst.ActiveObjExp_CompleteFloorPerTile
	}
	va.discoverExp = discoverExp + findExp + completeExp
}

// Forget clear discover data and return old state
func (va *VisitArea) Forget() int {
	va.mutex.Lock()
	defer va.mutex.Unlock()
	forgetTileCount := va.discoveredTileCount
	va.discoveredTileCount = 0
	va.bitsList = make([]BitContainder, va.w*va.h/BitLen+1)
	va.updateExp()
	return forgetTileCount
}

func (va *VisitArea) ToImage(zoom int) *image.RGBA {
	Xlen, Ylen := va.w, va.h
	img := image.NewRGBA(image.Rect(0, 0, Xlen*zoom, Ylen*zoom))
	coSet := color.RGBA{0xcf, 0xcf, 0xcf, 0xff}
	coClear := color.RGBA{0x3f, 0x3f, 0x3f, 0xff}
	va.mutex.RLock()
	defer va.mutex.RUnlock()
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			srcX := x / zoom
			srcY := y / zoom
			if va.GetXYNolock(srcX, srcY) {
				img.Set(x, y, coSet)
			} else {
				img.Set(x, y, coClear)
			}
		}
	}
	return img
}

func (va *VisitArea) Web_Image(w http.ResponseWriter, r *http.Request) error {
	img := va.ToImage(2)
	return png.Encode(w, img)
}
