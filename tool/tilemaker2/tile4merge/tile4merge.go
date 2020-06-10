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

package tile4merge

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kasworld/goguelike/lib/webtilegroup"
)

type Tile4Merge struct {
	GroupName string
	TileName  string
	Image     image.Image

	Rect webtilegroup.Rect // not image rect, rect in merged image
}

func SaveMergedPNG(tiList []*Tile4Merge, filename string) error {
	tileXcount := int(math.Sqrt(float64(len(tiList)))) + 1
	curXcount := 0
	curXpos := 0
	curYpos := 0
	lineHight := 0
	surfaceW := 0
	surfaceH := 0
	for _, ti := range tiList {
		ti.Rect.X = curXpos
		ti.Rect.Y = curYpos
		if curXpos+ti.Rect.W > surfaceW {
			surfaceW = curXpos + ti.Rect.W
		}
		if curYpos+ti.Rect.H > surfaceH {
			surfaceH = curYpos + ti.Rect.H
		}
		curXpos += ti.Rect.W
		if lineHight < ti.Rect.H {
			lineHight = ti.Rect.H
		}
		curXcount++
		if curXcount >= tileXcount {
			curXcount = 0
			curXpos = 0
			curYpos += lineHight
			lineHight = 0
		}
	}
	Img := image.NewRGBA(image.Rect(0, 0, surfaceW, surfaceH))
	for _, ti := range tiList {
		sr := ti.Image.Bounds()
		dp := image.Point{ti.Rect.X, ti.Rect.Y}
		r := image.Rectangle{dp, dp.Add(sr.Size())}
		draw.Draw(Img, r, ti.Image, sr.Min, draw.Src)
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	if err := png.Encode(out, Img); err != nil {
		return err
	}

	return nil
}

func SaveWebJSON(tiList []*Tile4Merge, filename string) error {
	rtn := make(webtilegroup.TileGroup)
	for _, v := range tiList {
		rtn[v.GroupName] = append(rtn[v.GroupName],
			webtilegroup.TileInfo{
				Name: v.TileName,
				Rect: v.Rect,
			})
	}
	j, err := json.Marshal(&rtn)
	if err != nil {
		return err
	}
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	if _, err := fd.Write(j); err != nil {
		return err
	}
	return nil
}

func LoadImageDir(baseDir string) []*Tile4Merge {
	fileSuffix := ".png"
	filenameList := make([]string, 0)
	rcsWalk := func(ppath string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() &&
			strings.HasSuffix(info.Name(), fileSuffix) {
			filenameList = append(filenameList, ppath)
		}
		return nil //errors.New(text)
	}
	filepath.Walk(baseDir, rcsWalk)
	sort.StringSlice(filenameList).Sort()

	tiList := make([]*Tile4Merge, 0)
	for _, v := range filenameList {
		ti, err := NewTileInfoByImageFile(v)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		tiList = append(tiList, ti)
	}
	return tiList
}

func NewTileInfoByImageFile(fullfilepath string) (*Tile4Merge, error) {
	fd, err := os.Open(fullfilepath)
	if err != nil {
		return nil, fmt.Errorf("Fail to %v", err)
	}
	defer fd.Close()
	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, fmt.Errorf("Fail to %v", err)
	}
	dirname, filename := filepath.Split(fullfilepath)
	groupname := filepath.Base(dirname)
	tilename := filename[:len(filename)-len(filepath.Ext(filename))]
	imgSize := img.Bounds().Size()
	return &Tile4Merge{
		GroupName: groupname,
		TileName:  tilename,
		Image:     img,
		Rect: webtilegroup.Rect{
			W: imgSize.X,
			H: imgSize.Y,
		},
	}, nil
}
