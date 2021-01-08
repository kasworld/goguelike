// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bias

import (
	"fmt"
	"math"

	"github.com/kasworld/goguelike/enum/factiontype"
)

type Bias [3]float64

func NewByFaction(ft factiontype.FactionType, l float64) Bias {
	return Bias(ft.FactorBase()).MakeAbsSumTo(l)
}

func MakeBiasByProgress(ft [3]int64, sec float64, l float64) Bias {
	progress := sec * 2 * math.Pi
	shift := math.Pi / 9
	return Bias{
		math.Sin(progress/float64(ft[0])+shift) * l,
		math.Cos(progress/float64(ft[1])+shift*4) * l,
		-math.Sin(progress/float64(ft[2])+shift*7) * l,
	}
}

func NewByRGB(r, g, b uint8) Bias {
	rtn := Bias{
		float64(r) - 0x80,
		float64(g) - 0x80,
		float64(b) - 0x80,
	}
	return rtn
}

func (bs Bias) String() string {
	return fmt.Sprintf("Bias[%f %f %f]",
		bs[0], bs[1], bs[2])
}

func (bs Bias) FormatString(format string) string {
	return fmt.Sprintf(
		format,
		bs[0], bs[1], bs[2],
	)
}

func (bs Bias) IntString() string {
	return bs.FormatString("[%.0f %.0f %.0f]")
}

func (bs Bias) F42String() string {
	return bs.FormatString("[%4.2f %4.2f %4.2f]")
}

func (bs Bias) NearFaction() factiontype.FactionType {
	srcLen := bs.AbsSum()
	minLen := srcLen * 2
	rtn := 0
	for i := 0; i < factiontype.FactionType_Count; i++ {
		ft := factiontype.FactionType(i)
		ftBias := Bias(ft.FactorBase()).MakeAbsSumTo(srcLen)
		fl := bs.Sub(ftBias).AbsSum()
		if fl < minLen {
			rtn = i
			minLen = fl
		}
	}
	return factiontype.FactionType(rtn)
}
