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

package wasmclientgl

import (
	"math"
	"time"

	"github.com/kasworld/goguelike/lib/canvastext"
)

var tcsRebirth = canvastext.TextTemplate{
	Color:    "yellow",
	SizeRate: 1.0,
	Dur:      1 * time.Second,
}

var tcsInfo = canvastext.TextTemplate{
	Color:    "yellow",
	SizeRate: 1.0,
	Dur:      5 * time.Second,
}

var tcsWarn = canvastext.TextTemplate{
	Color:    "OrangeRed",
	SizeRate: 1.0,
	Dur:      5 * time.Second,
}

func calcSizeRateByValue(v int) float64 {
	return 0.9 + math.Log2(float64(v))/10
}
