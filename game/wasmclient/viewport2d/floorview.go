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

package viewport2d

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/way9type"
	"github.com/kasworld/goguelike/game/bias"
	"github.com/kasworld/goguelike/game/clientfloor"
)

func (vp *Viewport2d) DrawFloorVP(
	envBias bias.Bias,
	cf *clientfloor.ClientFloor,
	floorVPPosX int,
	floorVPPosY int,
	actMoveDir way9type.Way9Type,
) error {

	defer vp.drawNoti()

	vp.context2d.Set("fillStyle", "black")
	vp.context2d.Call("fillRect", 0, 0, vp.ViewWidth, vp.ViewHeight)

	vp.drawFloorTiles(envBias, cf, 0, 0, floorVPPosX, floorVPPosY)
	vp.drawFloorFieldObjs(cf, 0, 0, floorVPPosX, floorVPPosY)

	dstRect := vp.calcDstRect(vp.CellXCount/2, vp.CellYCount/2, 0, 0)
	vp.drawCenterArrowDir(dstRect, actMoveDir)

	fontH := vp.CellSize
	vp.context2d.Set("font", fmt.Sprintf("%dpx sans-serif", fontH))
	posx := vp.ViewWidth / 2
	posy := vp.ViewHeight/2 - vp.CellSize
	vp.context2d.Set("fillStyle", "white")
	txt := fmt.Sprintf("%v %v", floorVPPosX, floorVPPosY)
	vp.context2d.Call("fillText", txt, posx, posy)
	return nil
}
