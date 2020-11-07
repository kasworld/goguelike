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

package floor

import (
	"fmt"

	"github.com/kasworld/goguelike/config/gameconst"
	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/lib/uuidposmani"
	"github.com/kasworld/hitrate"
)

// CacheVPIXYOList cache viewport index, pos(x,y), object for each turn
// to send viewport objectlist noti
// not inter-turn use
type CacheVPIXYOList struct {
	HitRate hitrate.HitRate
	olist   map[[2]int][4][]uuidposmani.VPIXYObj
	f       *Floor // for posman access
}

func (f *Floor) NewCacheVPIXYOList() *CacheVPIXYOList {
	return &CacheVPIXYOList{
		olist: make(map[[2]int][4][]uuidposmani.VPIXYObj),
		f:     f,
	}
}

func (cvpixyol *CacheVPIXYOList) String() string {
	return fmt.Sprintf("CacheVPIXYOList[len:%v %v]",
		len(cvpixyol.olist), cvpixyol.HitRate,
	)
}

func (cvpixyol *CacheVPIXYOList) GetAtByCache(
	x, y int) [4][]uuidposmani.VPIXYObj {
	rtn, exist := cvpixyol.olist[[2]int{x, y}]
	if !exist {
		rtn = [4][]uuidposmani.VPIXYObj{
			cvpixyol.f.aoPosMan.GetVPIXYObjByXYLenList(
				viewportdata.ViewportXYLenList,
				x, y,
				gameconst.ActiveObjCountInViewportLimit),
			cvpixyol.f.poPosMan.GetVPIXYObjByXYLenList(
				viewportdata.ViewportXYLenList,
				x, y,
				gameconst.CarryObjCountInViewportLimit),
			cvpixyol.f.foPosMan.GetVPIXYObjByXYLenList(
				viewportdata.ViewportXYLenList,
				x, y,
				gameconst.FieldObjCountInViewportLimit),
			cvpixyol.f.doPosMan.GetVPIXYObjByXYLenList(
				viewportdata.ViewportXYLenList,
				x, y,
				gameconst.DangerObjCountInViewportLimit),
		}
		cvpixyol.olist[[2]int{x, y}] = rtn
		cvpixyol.HitRate.Miss()
	} else {
		cvpixyol.HitRate.Hit()
	}
	return rtn
}
