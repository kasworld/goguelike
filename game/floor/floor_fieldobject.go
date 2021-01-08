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

package floor

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/game/fieldobject"
)

func (f *Floor) FindUsablePortalPairAt(
	x, y int) (*fieldobject.FieldObject, *fieldobject.FieldObject, error) {

	srcPortal, ok := f.foPosMan.Get1stObjAt(x, y).(*fieldobject.FieldObject)
	if !ok {
		return nil, nil, fmt.Errorf("not found src %v %v", x, y)
	}
	if srcPortal.ActType != fieldobjacttype.PortalInOut &&
		srcPortal.ActType != fieldobjacttype.PortalIn &&
		srcPortal.ActType != fieldobjacttype.PortalAutoIn {
		return nil, nil, fmt.Errorf("no in able portal %v", srcPortal)
	}
	fm := f.tower.GetFloorManager()
	if fm == nil {
		return nil, nil, fmt.Errorf("floor manager not ready")
	}
	dstPortal := fm.FindPortalByID(srcPortal.DstPortalID)
	if dstPortal == nil {
		return nil, nil, fmt.Errorf("not found dest %v", dstPortal)
	}
	if dstPortal.ActType != fieldobjacttype.PortalOut &&
		dstPortal.ActType != fieldobjacttype.PortalInOut {
		return nil, nil, fmt.Errorf("not out-able %v", dstPortal)
	}
	return srcPortal, dstPortal, nil
}
