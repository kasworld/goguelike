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

package viewportcache

import (
	"fmt"
	"sync"

	"github.com/kasworld/findnear"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/game/tilearea"
	"github.com/kasworld/hitrate"
	"github.com/kasworld/wrapper"
)

type terrainI interface {
	GetTiles() tilearea.TileArea
	GetXYLen() (int, int)
	GetXWrapper() *wrapper.Wrapper
	GetYWrapper() *wrapper.Wrapper
}

func (vpc *ViewportCache) String() string {
	return fmt.Sprintf("ViewportCache%v",
		vpc.HitRate)
}

type ViewportCache struct {
	mutex              sync.Mutex `prettystring:"hide"`
	terrain            terrainI
	RequireSightFromXY [][]*viewportdata.ViewportSight2
	HitRate            hitrate.HitRate
}

func New(tr terrainI) *ViewportCache {
	return &ViewportCache{
		terrain: tr,
	}
}

func (vpc *ViewportCache) Reset() {
	vpc.mutex.Lock()
	defer vpc.mutex.Unlock()
	w, h := vpc.terrain.GetXYLen()
	vpc.RequireSightFromXY = make([][]*viewportdata.ViewportSight2, w)
	for i := range vpc.RequireSightFromXY {
		vpc.RequireSightFromXY[i] = make([]*viewportdata.ViewportSight2, h)
	}
}

// ClearAt clear sight cache around x,y by xyl
func (vpc *ViewportCache) ClearAt(x, y int, xyl findnear.XYLenList) {
	vpc.mutex.Lock()
	defer vpc.mutex.Unlock()
	xWrap := vpc.terrain.GetXWrapper().GetWrapFn()
	yWrap := vpc.terrain.GetYWrapper().GetWrapFn()
	for _, v := range xyl {
		vpc.RequireSightFromXY[xWrap(x+v.X)][yWrap(y+v.Y)] = nil
	}
}

func (vpc *ViewportCache) GetByCache(x, y int) *viewportdata.ViewportSight2 {
	vpc.mutex.Lock()
	defer vpc.mutex.Unlock()
	if vpc.RequireSightFromXY[x][y] == nil {
		vpc.HitRate.Miss()
		vpc.RequireSightFromXY[x][y] = vpc.makeAt3(
			x, y)
	} else {
		vpc.HitRate.Hit()
	}
	return vpc.RequireSightFromXY[x][y]
}

// make visible map, lineofsight
func (vpc *ViewportCache) makeAt3(centerX, centerY int) *viewportdata.ViewportSight2 {
	vpSightMat := viewportdata.ViewportSight2{}

	tiles := vpc.terrain.GetTiles()
	xWrap := vpc.terrain.GetXWrapper().GetWrapFn()
	yWrap := vpc.terrain.GetYWrapper().GetWrapFn()

	for i, sightLine := range viewportdata.SightlinesByXYLenList {
		needSight := 0.0
		if len(sightLine) > 0 {
			for _, w := range sightLine {
				tx := xWrap(centerX + w.X)
				ty := yWrap(centerY + w.Y)
				needSight += tiles[tx][ty].BlockSight() * w.L
			}
		}
		vpSightMat[i] = float32(needSight)
	}
	return &vpSightMat
}
