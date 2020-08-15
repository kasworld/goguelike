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
	"testing"
)

func TestMakePosLenList(t *testing.T) {
	pll := MakePosLenList(0.5, -0.5, 13.5, -11).DelDup().ToInLen()
	for i, v := range pll {
		cellx, celly := v.CellXY()
		fmt.Printf("%v %v (%v %v)\n", i, v, cellx, celly)
	}
}
