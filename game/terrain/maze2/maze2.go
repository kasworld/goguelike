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

// Package maze2 make maze by Growing Tree algorithm
// from http://weblog.jamisbuck.org/2011/1/27/maze-generation-growing-tree-algorithm
package maze2

import (
	"bytes"
	"fmt"

	"github.com/kasworld/g2rand"
	"github.com/kasworld/goguelike/game/terrain/mazearea"
)

type Dir int

func (d Dir) String() string {
	str := []string{"N", "S", "E", "W"}
	var buf bytes.Buffer
	for i := uint(0); i < 4; i++ {
		if d&(1<<i) != 0 {
			buf.WriteString(str[i])
		} else {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

const (
	N Dir = 1
	S Dir = 2
	E Dir = 4
	W Dir = 8
)

var DX = map[Dir]int{E: 1, W: -1, N: 0, S: 0}
var DY = map[Dir]int{E: 0, W: 0, N: -1, S: 1}
var OPCarryObjSITE = map[Dir]Dir{E: W, W: E, N: S, S: N}

type Maze struct {
	W       int
	H       int
	PosList [][2]int
	Cells   [][]Dir
}

func New(rnd *g2rand.G2Rand, w, h int) *Maze {
	m := Maze{
		W:     w,
		H:     h,
		Cells: make([][]Dir, w),
	}
	for i := range m.Cells {
		m.Cells[i] = make([]Dir, h)
	}
	m.make(rnd)
	return &m
}

func (m *Maze) selectVisit(rnd *g2rand.G2Rand) int {
	if rnd.Intn(2) == 0 {
		return len(m.PosList) - 1
	}
	return rnd.Intn(len(m.PosList))
}

func (m *Maze) make(rnd *g2rand.G2Rand) {
	x, y := rnd.Intn(m.W), rnd.Intn(m.H)
	m.PosList = append(m.PosList, [2]int{x, y})
	for len(m.PosList) != 0 {
		curVisitI := m.selectVisit(rnd)
		x, y := m.PosList[curVisitI][0], m.PosList[curVisitI][1]
		delCur := true

		dirList := []Dir{N, S, E, W}
		for _, dirIndex := range rnd.Perm(4) {
			dir := dirList[dirIndex]
			nx, ny := x+DX[dir], y+DY[dir]
			if nx >= 0 && ny >= 0 && nx < m.W && ny < m.H && m.Cells[nx][ny] == 0 {
				m.Cells[x][y] |= dir
				m.Cells[nx][ny] |= OPCarryObjSITE[dir]
				m.PosList = append(m.PosList, [2]int{nx, ny})
				delCur = false
				break
			}
		}
		if delCur {
			m.PosList = append(m.PosList[:curVisitI], m.PosList[curVisitI+1:]...)
		}
	}
}

func (m *Maze) MazeString() string {
	var buf bytes.Buffer
	for _, v := range m.Cells {
		fmt.Fprintf(&buf, "%s\n", v)
	}
	return buf.String()
}

func (m *Maze) ToMazeArea(w, h int, conerFill bool) (mazearea.MazeArea, error) {
	ma := mazearea.New(w, h)
	cellx := w / m.W
	celly := h / m.H
	if cellx < 2 || celly < 2 {
		return ma, fmt.Errorf("area too small, cell %v %v", cellx, celly)
	}
	for x := 0; x < m.W; x++ {
		for y := 0; y < m.H; y++ {
			if m.Cells[x][y]&N == 0 {
				y1 := y * celly
				var x1, x2 int
				if conerFill {
					x1, x2 = x*cellx, (x+1)*cellx+1
				} else {
					x1, x2 = x*cellx+1, (x+1)*cellx
				}
				ma.HLine(x1, x2, y1, true)
			}
			if m.Cells[x][y]&W == 0 {
				x1 := x * cellx
				var y1, y2 int
				if conerFill {
					y1, y2 = y*celly, (y+1)*celly+1
				} else {
					y1, y2 = y*celly+1, (y+1)*celly
				}
				ma.VLine(x1, y1, y2, true)
			}
		}
	}
	return ma, nil
}
