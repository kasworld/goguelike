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

package htmlcolors

import "fmt"

type Color24 uint32

func (i Color24) String() string {
	if str, exist := Color24ToName[i]; exist {
		return str
	}
	return fmt.Sprintf("Color24[0x%08x]", uint32(i))
}

func (c Color24) R() uint8 {
	return uint8((c >> 16) & 0xff)
}
func (c Color24) G() uint8 {
	return uint8((c >> 8) & 0xff)
}
func (c Color24) B() uint8 {
	return uint8((c >> 0) & 0xff)
}

func (c Color24) SetR(r uint8) Color24 {
	return (c & 0x00ffff) | Color24(r)<<16
}
func (c Color24) SetG(r uint8) Color24 {
	return (c & 0xff00ff) | Color24(r)<<8
}
func (c Color24) SetB(r uint8) Color24 {
	return (c & 0xffff00) | Color24(r)<<0
}

func (c Color24) UInt8Array() [3]uint8 {
	return [3]uint8{
		c.R(),
		c.G(),
		c.B(),
	}
}

func (c Color24) UInt8ArrayAlpha(a uint8) [4]uint8 {
	return [4]uint8{
		c.R(),
		c.G(),
		c.B(),
		a,
	}
}

func (c Color24) Rev() Color24 {
	return NewColor24(c.R()+128, c.G()+128, c.B()+128)
}

func (c Color24) Neg() Color24 {
	return 0xffffff - c
}

func (c Color24) ShiftRight(n uint) Color24 {
	return NewColor24(c.R()>>n, c.G()>>n, c.B()>>n)
}

func (c Color24) MulF(v float64) Color24 {
	return NewColor24(uint8(float64(c.R())*v), uint8(float64(c.G())*v), uint8(float64(c.B())*v))
}

func (c Color24) ToHTMLColorString() string {
	return fmt.Sprintf("#%06x", uint32(c))
}

func NewColor24(r, g, b uint8) Color24 {
	return Color24(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

func (src Color24) NearNamedColor24() Color24 {
	minlen := 4 * 256 * 256
	idx := 0
	for i, v := range Color24List {
		if src.LenTo(v) < minlen {
			idx = i
			minlen = src.LenTo(v)
		}
	}
	return Color24List[idx]
}

func (c Color24) LenTo(dst Color24) int {
	srcRGB := c.UInt8Array()
	dstRGB := dst.UInt8Array()
	sum := 0
	for i := range srcRGB {
		diff := int(srcRGB[i]) - int(dstRGB[i])
		sum += diff * diff
	}
	return sum
}
