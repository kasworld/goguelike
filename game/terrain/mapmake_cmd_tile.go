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

package terrain

import (
	"fmt"

	"github.com/kasworld/goguelike/enum/tile_flag"
	"github.com/kasworld/goguelike/enum/tileoptype"
	"github.com/kasworld/goguelike/lib/maze2"
	"github.com/kasworld/goguelike/lib/scriptparse"
)

// all function is operation to tileLayer

func cmdTileAt(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var x, y int
	if err := ca.GetArgs(&tl, &x, &y); err != nil {
		return err
	}
	tr.tileLayer.OpXY(x, y, tile_flag.TileTypeValue{Op: tileoptype.SetBits, Arg: tl})
	return nil
}

func cmdTileHLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var x, w, y int
	if err := ca.GetArgs(&tl, &x, &w, &y); err != nil {
		return err
	}

	tr.tileLayer.HLine(x, w, y, tl)
	return nil
}

func cmdTileVLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var x, y, h int
	if err := ca.GetArgs(&tl, &x, &y, &h); err != nil {
		return err
	}

	tr.tileLayer.VLine(x, y, h, tl)
	return nil
}

func cmdTileLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var x1, y1, x2, y2 int
	if err := ca.GetArgs(&tl, &x1, &y1, &x2, &y2); err != nil {
		return err
	}

	tr.tileLayer.Line(x1, y1, x2, y2, tl)
	return nil
}

func cmdTileRect(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var x, w, y, h int
	if err := ca.GetArgs(&tl, &x, &w, &y, &h); err != nil {
		return err
	}

	tr.tileLayer.Rect(x, w, y, h, tl)
	return nil
}

func cmdTileFillRect(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var x, w, y, h int
	if err := ca.GetArgs(&tl, &x, &w, &y, &h); err != nil {
		return err
	}
	tr.tileLayer.FillRect(x, w, y, h, tl)
	return nil
}

func cmdTileFillEllipses(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var x, w, y, h int
	if err := ca.GetArgs(&tl, &x, &w, &y, &h); err != nil {
		return err
	}
	tr.tileLayer.Ellipses(x, w, y, h, tl)
	return nil
}

func cmdTileMazeWall(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var xn, yn int
	var conerFill bool
	var maX, maY, maW, maH int
	if err := ca.GetArgs(&tl, &maX, &maY, &maW, &maH, &xn, &yn, &conerFill); err != nil {
		return err
	}

	m := maze2.New(tr.rnd, xn, yn)
	ma, err := m.ToBoolMatrix(maW, maH, conerFill)
	if err != nil {
		return fmt.Errorf("tr %v %v", tr, err)
	}

	tr.tileLayer.DrawBoolMapTrue(tr.XWrap, tr.YWrap, maX, maY, ma, tl)
	return nil
}

func cmdTileMazeWalk(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile_flag.TileFlag
	var xn, yn int
	var conerFill bool
	var maX, maY, maW, maH int
	if err := ca.GetArgs(&tl, &maX, &maY, &maW, &maH, &xn, &yn, &conerFill); err != nil {
		return err
	}

	m := maze2.New(tr.rnd, xn, yn)
	ma, err := m.ToBoolMatrix(tr.GetXLen(), tr.GetYLen(), conerFill)
	if err != nil {
		return fmt.Errorf("tr %v %v", tr, err)
	}
	tr.tileLayer.DrawBoolMapFalse(tr.XWrap, tr.YWrap, maX, maY, ma, tl)
	return nil
}
