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

func (bs Bias) SelectSkill(i int) float64 {
	v := bs[i]
	if v < 0 {
		v = -v
	}
	return v
}

func (bs Bias) Add(other Bias) Bias {
	return Bias{bs[0] + other[0], bs[1] + other[1], bs[2] + other[2]}
}

func (bs Bias) Sub(other Bias) Bias {
	return Bias{bs[0] - other[0], bs[1] - other[1], bs[2] - other[2]}
}

func (bs Bias) Unsigned(other Bias) Bias {
	rtn := bs
	for i, v := range rtn {
		if v < 0 {
			rtn[i] = -v
		}
	}
	return rtn
}

func (bs Bias) Neg() Bias {
	return Bias{-bs[0], -bs[1], -bs[2]}
}

func (bs Bias) RotateRight() Bias {
	return Bias{bs[2], bs[0], bs[1]}
}

func (bs Bias) RotateLeft() Bias {
	return Bias{bs[1], bs[2], bs[0]}
}

func (bs Bias) Idiv(other float64) Bias {
	return Bias{bs[0] / other, bs[1] / other, bs[2] / other}
}

func (bs Bias) AbsSum() float64 {
	sum := 0.0
	for _, v := range bs {
		if v < 0 {
			sum -= v
		} else {
			sum += v
		}
	}
	return sum
}

func (bs Bias) MakeAbsSumTo(l float64) Bias {
	d := bs.AbsSum() / l
	if d != 0 {
		return bs.Idiv(d)
	}
	return bs
}
