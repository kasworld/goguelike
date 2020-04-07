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

package runetile

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/kasworld/htmlcolors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func LoadFontFile(fullpath string) ([]byte, error) {
	fd, err := os.Open(fullpath)
	if err != nil {
		return nil, fmt.Errorf("%v %v", fullpath, err)
	}
	defer fd.Close()
	return ioutil.ReadAll(fd)
}

func MakeFace(fontdata []byte, size float64) (font.Face, error) {
	ft, err := truetype.Parse(fontdata)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(ft,
		&truetype.Options{
			Size: size,
		})
	return face, nil
}

func NewRuneTileListByString(
	groupname string, runelist string, co htmlcolors.Color24, ft []byte, ftsize float64,
	dx, dy, shift, w, h int,
) []RuneTile {
	rtn := make([]RuneTile, 0)
	for _, v := range runelist {
		rtn = append(rtn, RuneTile{
			GroupName:   groupname,
			Rune:        string(v),
			Color1:      co,
			FontData:    ft,
			FontSize:    ftsize,
			FontDx:      dx,
			FontDy:      dy,
			ShadowShift: shift,
			TileW:       w,
			TileH:       h,
		})
	}
	return rtn
}

type RuneTile struct {
	// input data
	GroupName string
	Name      string
	Rune      string
	Color1    htmlcolors.Color24
	FontData  []byte
	FontSize  float64

	FontDx      int
	FontDy      int
	ShadowShift int
	TileW       int
	TileH       int

	// output
	Img draw.Image
}

func (tmd *RuneTile) MakeImg() error {
	face, err := MakeFace(tmd.FontData, tmd.FontSize)
	if err != nil {
		return err
	}
	Img := image.NewRGBA(image.Rect(0, 0, tmd.TileW, tmd.TileH))
	Color2 := tmd.Color1.Rev()
	bd, ad := font.BoundString(face, tmd.Rune)
	_ = ad
	runeW := bd.Max.X.Ceil() - bd.Min.X.Ceil()
	runeH := bd.Max.Y.Ceil() - bd.Min.Y.Ceil()
	_ = runeW
	_ = runeH
	x, yb := tmd.TileW/2-runeW/2-tmd.ShadowShift, tmd.TileH-runeH/6-tmd.ShadowShift
	x += tmd.FontDx
	yb += tmd.FontDy
	d := &font.Drawer{
		Dst: Img,
		Src: image.NewUniform(
			color.RGBA{
				Color2.R(),
				Color2.G(),
				Color2.B(),
				255,
			}),
		Face: face,
	}
	d.Dot = fixed.P(x+tmd.ShadowShift/2, yb+tmd.ShadowShift/2)
	d.DrawString(tmd.Rune)

	d.Src = image.NewUniform(color.RGBA{
		tmd.Color1.R(),
		tmd.Color1.G(),
		tmd.Color1.B(),
		255,
	})
	d.Dot = fixed.P(x-tmd.ShadowShift/2, yb-tmd.ShadowShift/2)
	d.DrawString(tmd.Rune)
	tmd.Img = Img
	return nil
}
