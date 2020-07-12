// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wasmclientgl

import (
	"fmt"
	"sync"
	"syscall/js"

	"github.com/kasworld/goguelike/config/moneycolor"
	"github.com/kasworld/goguelike/enum/carryingobjecttype"
	"github.com/kasworld/goguelike/enum/equipslottype"
	"github.com/kasworld/goguelike/lib/webtilegroup"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
)

var gPoolCarryObj3D = NewPoolCarryObj3D(PoolSizeCarryObj3D)

type PoolCarryObj3D struct {
	mutex    sync.Mutex
	poolData []*CarryObj3D
}

func NewPoolCarryObj3D(initCap int) *PoolCarryObj3D {
	return &PoolCarryObj3D{
		poolData: make([]*CarryObj3D, 0, initCap),
	}
}

func (p *PoolCarryObj3D) String() string {
	return fmt.Sprintf("PacketPoolCarryObj3D[%v/%v]",
		len(p.poolData), cap(p.poolData),
	)
}

func (p *PoolCarryObj3D) Get() *CarryObj3D {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn *CarryObj3D
	if l := len(p.poolData); l > 0 {
		rtn = p.poolData[l-1]
		p.poolData = p.poolData[:l-1]
	} else {
		rtn = NewCarryObj3D()
	}
	return rtn
}

func (p *PoolCarryObj3D) Put(pb *CarryObj3D) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.poolData = append(p.poolData, pb)
}

type CarryObj3D struct {
	Cnv     js.Value
	Ctx     js.Value
	Tex     js.Value
	GeoInfo GeoInfo
	Mesh    js.Value
}

func NewCarryObj3D() *CarryObj3D {
	cnv := js.Global().Get("document").Call("createElement", "CANVAS")
	ctx := cnv.Call("getContext", "2d")
	ctx.Set("imageSmoothingEnabled", false)
	cnv.Set("width", DstCellSize)
	cnv.Set("height", DstCellSize)
	tex := ThreeJsNew("CanvasTexture", cnv)
	mat := ThreeJsNew("MeshStandardMaterial",
		map[string]interface{}{
			"map": tex,
		},
	)
	mat.Set("transparent", true)
	geo := ThreeJsNew("BoxGeometry",
		DstCellSize/3, DstCellSize/3, DstCellSize/3)
	mesh := ThreeJsNew("Mesh", geo, mat)
	return &CarryObj3D{
		Cnv:     cnv,
		Ctx:     ctx,
		Tex:     tex,
		GeoInfo: GetGeoInfo(geo),
		Mesh:    mesh,
	}
}

func (aog *CarryObj3D) ChangeTile(ti webtilegroup.TileInfo) {
	aog.Ctx.Call("clearRect", 0, 0, DstCellSize, DstCellSize)
	aog.Ctx.Call("drawImage", gClientTile.TilePNG.Cnv,
		ti.Rect.X, ti.Rect.Y, ti.Rect.W, ti.Rect.H,
		0, 0, DstCellSize, DstCellSize)
	aog.Tex.Set("needsUpdate", true)
}

func (aog *CarryObj3D) SetFieldPosition(fx, fy int, shInfo ShiftInfo) {
	SetPosition(
		aog.Mesh,
		float64(fx)*DstCellSize+aog.GeoInfo.Len[0]/2+shInfo.X,
		-float64(fy)*DstCellSize-aog.GeoInfo.Len[1]/2-shInfo.Y,
		aog.GeoInfo.Len[2]/2+shInfo.Z,
	)
}

func (aog *CarryObj3D) RotateX(rad float64) {
	aog.Mesh.Get("rotation").Set("x", rad)
}
func (aog *CarryObj3D) RotateY(rad float64) {
	aog.Mesh.Get("rotation").Set("y", rad)
}
func (aog *CarryObj3D) RotateZ(rad float64) {
	aog.Mesh.Get("rotation").Set("z", rad)
}

func (aog *CarryObj3D) Dispose() {
	// mesh do not need dispose
	aog.Mesh.Get("geometry").Call("dispose")
	aog.Mesh.Get("material").Call("dispose")
	aog.Tex.Call("dispose")

	aog.Cnv = js.Undefined()
	aog.Ctx = js.Undefined()
	aog.Mesh = js.Undefined()
	aog.Tex = js.Undefined()
	// no need createElement canvas dom obj
}

func Equiped2TileInfo(o *c2t_obj.EquipClient) webtilegroup.TileInfo {
	return gClientTile.EquipTiles[o.EquipType][o.Faction]
}

func CarryObj2TileInfo(o *c2t_obj.CarryObjClientOnFloor) webtilegroup.TileInfo {
	var ti webtilegroup.TileInfo
	switch o.CarryingObjectType {
	case carryingobjecttype.Equip:
		ti = gClientTile.EquipTiles[o.EquipType][o.Faction]
	case carryingobjecttype.Money:
		var find bool
		for i, v := range moneycolor.Attrib {
			if o.Value < v.UpLimit {
				ti = gClientTile.GoldTiles[i]
				find = true
				break
			}
		}
		if !find {
			ti = gClientTile.GoldTiles[len(gClientTile.GoldTiles)-1]
		}
	case carryingobjecttype.Potion:
		ti = gClientTile.PotionTiles[o.PotionType]
	case carryingobjecttype.Scroll:
		ti = gClientTile.ScrollTiles[o.ScrollType]
	}
	return ti
}

type ShiftInfo struct {
	X float64
	Y float64
	Z float64
}

// equipped shift, around ao
var aoEqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	equipslottype.Helmet: {DstCellSize * -0.33, DstCellSize * 0.0, DstCellSize * 0.66},
	equipslottype.Amulet: {DstCellSize * 1.00, DstCellSize * 0.0, DstCellSize * 0.66},

	equipslottype.Weapon: {DstCellSize * -0.33, DstCellSize * 0.25, DstCellSize * 0.66},
	equipslottype.Shield: {DstCellSize * 1.00, DstCellSize * 0.25, DstCellSize * 0.66},

	equipslottype.Ring:     {DstCellSize * -0.33, DstCellSize * 0.50, DstCellSize * 0.66},
	equipslottype.Gauntlet: {DstCellSize * 1.00, DstCellSize * 0.50, DstCellSize * 0.66},

	equipslottype.Armor:    {DstCellSize * -0.33, DstCellSize * 0.75, DstCellSize * 0.66},
	equipslottype.Footwear: {DstCellSize * 1.00, DstCellSize * 0.75, DstCellSize * 0.66},
}

func CarryObjClientOnFloor2DrawInfo(
	co *c2t_obj.CarryObjClientOnFloor) ShiftInfo {
	switch co.CarryingObjectType {
	default:
		return otherCarryObjShift[co.CarryingObjectType]
	case carryingobjecttype.Equip:
		return eqPosShift[co.EquipType]
	}
}

// on floor in tile
var eqPosShift = [equipslottype.EquipSlotType_Count]ShiftInfo{
	equipslottype.Helmet: {DstCellSize * 0.0, DstCellSize * 0.0, DstCellSize * 0.33},
	equipslottype.Amulet: {DstCellSize * 0.75, DstCellSize * 0.0, DstCellSize * 0.33},

	equipslottype.Weapon: {DstCellSize * 0.0, DstCellSize * 0.25, DstCellSize * 0.33},
	equipslottype.Shield: {DstCellSize * 0.75, DstCellSize * 0.25, DstCellSize * 0.33},

	equipslottype.Ring:     {DstCellSize * 0.0, DstCellSize * 0.50, DstCellSize * 0.33},
	equipslottype.Gauntlet: {DstCellSize * 0.75, DstCellSize * 0.50, DstCellSize * 0.33},

	equipslottype.Armor:    {DstCellSize * 0.0, DstCellSize * 0.75, DstCellSize * 0.33},
	equipslottype.Footwear: {DstCellSize * 0.75, DstCellSize * 0.75, DstCellSize * 0.33},
}

var otherCarryObjShift = [carryingobjecttype.CarryingObjectType_Count]ShiftInfo{
	carryingobjecttype.Money:  {DstCellSize * 0.33, DstCellSize * 0.0, DstCellSize * 0.33},
	carryingobjecttype.Potion: {DstCellSize * 0.33, DstCellSize * 0.33, DstCellSize * 0.33},
	carryingobjecttype.Scroll: {DstCellSize * 0.33, DstCellSize * 0.66, DstCellSize * 0.33},
}
