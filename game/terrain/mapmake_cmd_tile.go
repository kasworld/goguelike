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

	"github.com/kasworld/goguelike/enum/tileoptype"

	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"

	"github.com/kasworld/goguelike/game/terrain/maze2"
	"github.com/kasworld/goguelike/lib/scriptparse"
	"github.com/kasworld/walk2d"
)

func cmdTileAt(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var x, y int
	if err := ca.GetArgs(&tl, &x, &y); err != nil {
		return err
	}
	tr.tileArea.OpXY(x, y, tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl})
	return nil
}

func cmdTileHLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var x, w, y int
	if err := ca.GetArgs(&tl, &x, &w, &y); err != nil {
		return err
	}

	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		tr.tileArea.OpXY(ax, ay, rv)
		return false
	}
	walk2d.HLine(x, x+w, y, fn)
	return nil
}

func cmdTileVLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var x, y, h int
	if err := ca.GetArgs(&tl, &x, &y, &h); err != nil {
		return err
	}

	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		tr.tileArea.OpXY(ax, ay, rv)
		return false
	}
	walk2d.VLine(y, y+h, x, fn)
	return nil
}

func cmdTileLine(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var x1, y1, x2, y2 int
	if err := ca.GetArgs(&tl, &x1, &y1, &x2, &y2); err != nil {
		return err
	}

	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		tr.tileArea.OpXY(ax, ay, rv)
		return false
	}
	walk2d.Line(x1, y1, x2, y2, fn)
	return nil
}

func cmdTileRect(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var x, w, y, h int
	if err := ca.GetArgs(&tl, &x, &w, &y, &h); err != nil {
		return err
	}

	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		tr.tileArea.OpXY(ax, ay, rv)
		return false
	}
	walk2d.Rect(x, y, x+w, y+h, fn)
	return nil
}

func cmdTileFillRect(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var x, w, y, h int
	if err := ca.GetArgs(&tl, &x, &w, &y, &h); err != nil {
		return err
	}
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		tr.tileArea.OpXY(ax, ay, rv)
		return false
	}
	walk2d.FillHV(x, y, x+w, y+h, fn)
	return nil
}

func cmdTileFillEllipses(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var x, w, y, h int
	if err := ca.GetArgs(&tl, &x, &w, &y, &h); err != nil {
		return err
	}
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	fn := func(ax, ay int) bool {
		tr.tileArea.OpXY(ax, ay, rv)
		return false
	}
	walk2d.Ellipses(x, y, x+w, y+h, fn)
	return nil
}

func cmdTileMazeWall(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var xn, yn int
	var conerFill bool
	var maX, maY, maW, maH int
	if err := ca.GetArgs(&tl, &maX, &maY, &maW, &maH, &xn, &yn, &conerFill); err != nil {
		return err
	}

	m := maze2.New(tr.rnd, xn, yn)
	ma, err := m.ToMazeArea(maW, maH, conerFill)
	if err != nil {
		return fmt.Errorf("tr %v %v", tr, err)
	}
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	for x, xv := range ma {
		for y, yv := range xv {
			if yv {
				tr.tileArea.OpXY(
					tr.XWrapper.WrapSafe(maX+x),
					tr.YWrapper.WrapSafe(maY+y),
					rv)
			}
		}
	}
	return nil
}

func cmdTileMazeWalk(tr *Terrain, ca *scriptparse.CmdArgs) error {
	var tl tile.Tile
	var xn, yn int
	var conerFill bool
	var maX, maY, maW, maH int
	if err := ca.GetArgs(&tl, &maX, &maY, &maW, &maH, &xn, &yn, &conerFill); err != nil {
		return err
	}

	m := maze2.New(tr.rnd, xn, yn)
	ma, err := m.ToMazeArea(tr.GetXLen(), tr.GetYLen(), conerFill)
	if err != nil {
		return fmt.Errorf("tr %v %v", tr, err)
	}
	rv := tile_flag.TileTypeValue{Op: tileoptype.OverrideBits, Arg: tl}
	for x, xv := range ma {
		for y, yv := range xv {
			if !yv {
				tr.tileArea.OpXY(
					tr.XWrapper.WrapSafe(maX+x),
					tr.YWrapper.WrapSafe(maY+y),
					rv)
			}
		}
	}
	return nil
}
