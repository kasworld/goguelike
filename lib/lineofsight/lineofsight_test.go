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

package lineofsight

import (
	"fmt"
	"math"
	"testing"

	"github.com/kasworld/goguelike/config/viewportdata"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
)

func TestMakePosLenList(t *testing.T) {
	for deg := 0.0; deg < 360; deg += 10 {
		rad := deg / 180 * math.Pi
		dx := fieldobjacttype.LightHouseRadius * math.Cos(rad)
		dy := fieldobjacttype.LightHouseRadius * math.Sin(rad)
		xylenline := MakePosLenList(0.5, 0.5, dx+0.5, dy+0.5).ToCellLenList()
		fmt.Printf("%v %v\n", deg, xylenline)
	}
}

func TestMakeSightlinesByXYLenList(t *testing.T) {
	vp := viewportdata.ViewportXYLenList[:9]
	sightlinesByXYLenList := MakeSightlinesByXYLenList(vp)
	for i, v := range sightlinesByXYLenList {
		fmt.Printf("%v %v\n", vp[i], v)
	}
}
