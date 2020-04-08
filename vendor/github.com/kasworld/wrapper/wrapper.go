// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wrapper

import (
	"fmt"
	"sync"
)

// reuse same (witdh) wrapper
var g_wrapper map[int]*Wrapper
var g_mutex sync.Mutex

func init() {
	g_wrapper = make(map[int]*Wrapper)
}

func G_WrapperInfo() string {
	return fmt.Sprintf("%v", g_wrapper)
}

func New(w int) *Wrapper {
	g_mutex.Lock()
	defer g_mutex.Unlock()

	wr, exist := g_wrapper[w]
	if exist {
		return wr
	}
	wr = makeNew(w)
	g_wrapper[w] = wr
	return wr
}

type Wrapper struct {
	wrapTable []int
	width     int
	widthMask int
	Wrap      func(i int) int
	WrapSafe  func(i int) int
}

func (wr Wrapper) String() string {
	return fmt.Sprintf("Wrapper[%v]", wr.width)
}

func makeNew(w int) *Wrapper {
	if w < 2 {
		return nil
	}
	if isPowerOfTwo(w) {
		wr := &Wrapper{
			width:     w,
			widthMask: w - 1,
		}
		wr.Wrap = wr.wrapPowerOfTwo
		wr.WrapSafe = wr.wrapPowerOfTwo
		return wr
	} else {
		wr := &Wrapper{
			wrapTable: make([]int, w*3),
			width:     w,
		}
		for i, _ := range wr.wrapTable {
			wr.wrapTable[i] = i % w
		}
		wr.Wrap = wr.wrapByTable
		wr.WrapSafe = wr.wrapByCalc
		return wr
	}
}

func isPowerOfTwo(i int) bool {
	return (i & (i - 1)) == 0
}

func (wr *Wrapper) GetWrapFn() func(i int) int {
	return wr.Wrap
}
func (wr *Wrapper) GetWrapSafeFn() func(i int) int {
	return wr.WrapSafe
}

func (wr *Wrapper) GetWidth() int {
	return wr.width
}

func (wr *Wrapper) IsIn(i int) bool {
	return i >= 0 && i < wr.width
}

func (wr *Wrapper) wrapByCalc(i int) int {
	return (i%wr.width + wr.width) % wr.width
}

func (wr *Wrapper) wrapPowerOfTwo(i int) int {
	return i & wr.widthMask
}

func (wr *Wrapper) wrapByTable(i int) int {
	return wr.wrapTable[i+wr.width]
}

func (wr *Wrapper) wrapSimple(i int) int {
	return (i + wr.width) % wr.width
}

func (wr *Wrapper) wrapIf(i int) int {
	if i < 0 {
		return i%wr.width + wr.width
	} else if i >= wr.width {
		return i % wr.width
	} else {
		return i
	}
}

func (wr *Wrapper) wrapIfSimple(i int) int {
	if i < 0 {
		return i + wr.width
	} else if i >= wr.width {
		return i - wr.width
	} else {
		return i
	}
}
