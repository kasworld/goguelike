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

package drawtile

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/kasworld/htmlcolors"

	"github.com/kasworld/walk2d"
)

type DrawTile struct {
	// input data
	GroupName string
	Name      string
	Rect      image.Rectangle

	// output
	Img draw.Image
}

func New(GroupName, Name string, w, h int) *DrawTile {
	rtn := &DrawTile{
		GroupName: GroupName,
		Name:      Name,
	}
	rtn.Rect = image.Rect(0, 0, w, h)
	rtn.Img = image.NewRGBA(rtn.Rect)
	return rtn
}

func (dt *DrawTile) DrawRect(co htmlcolors.Color24) {
	cc := color.RGBA{
		co.R(),
		co.G(),
		co.B(),
		255,
	}
	walk2d.Rect(
		dt.Rect.Min.X,
		dt.Rect.Min.Y,
		dt.Rect.Max.X,
		dt.Rect.Max.Y,
		func(x, y int) bool {
			dt.Img.Set(x, y, cc)
			return false
		},
	)
}

func (dt *DrawTile) DrawX(co htmlcolors.Color24) {
	cc := color.RGBA{
		co.R(),
		co.G(),
		co.B(),
		255,
	}
	walk2d.Line(
		dt.Rect.Min.X,
		dt.Rect.Min.Y,
		dt.Rect.Max.X,
		dt.Rect.Max.Y,
		func(x, y int) bool {
			dt.Img.Set(x, y, cc)
			return false
		},
	)
	walk2d.Line(
		dt.Rect.Max.X,
		dt.Rect.Min.Y,
		dt.Rect.Min.X,
		dt.Rect.Max.Y,
		func(x, y int) bool {
			dt.Img.Set(x, y, cc)
			return false
		},
	)
}
