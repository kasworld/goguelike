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

// Package mazearea bool matrix of maze wall, road
package mazearea

import (
	"bytes"
)

type MazeArea [][]bool

func New(w, h int) MazeArea {
	ma := make([][]bool, w)
	for x := range ma {
		ma[x] = make([]bool, h)
	}
	return ma
}

func (ma MazeArea) String() string {
	return ma.ToString('0', '1')
}

func (ma MazeArea) ToString(r1, r2 rune) string {
	var buf bytes.Buffer
	for _, v := range ma {
		for _, w := range v {
			if w {
				buf.WriteRune(r2)
			} else {
				buf.WriteRune(r1)
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func (ma MazeArea) HLine(x1, x2, y int, v bool) {
	w := len(ma)
	h := len(ma[0])
	if y >= h {
		y = h - 1
	}
	if x1 >= w {
		x1 = w - 1
	}
	if x2 >= w {
		x2 = w - 1
	}
	for i := x1; i < x2; i++ {
		ma[i][y] = v
	}
}

func (ma MazeArea) VLine(x, y1, y2 int, v bool) {
	w := len(ma)
	h := len(ma[0])
	if y1 >= h {
		y1 = h - 1
	}
	if y2 >= h {
		y2 = h - 1
	}
	if x >= w {
		x = w - 1
	}
	for i := y1; i < y2; i++ {
		ma[x][i] = v
	}
}
