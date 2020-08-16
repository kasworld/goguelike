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

package lineofsight0

import (
	"fmt"
	"testing"

	"github.com/kasworld/goguelike/config/viewportdata"
)

// func Test_calcNextAxisCrossPos(t *testing.T) {
// 	type args struct {
// 		currentPos    float64
// 		moveDirection float64
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want float64
// 	}{
// 		{"++", args{3.5, 1}, 4},
// 		{"+-", args{3.5, -1}, 3},
// 		{"-+", args{-3.5, 1}, -3},
// 		{"--", args{-3.5, -1}, -4},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := calcNextAxisCrossPos(tt.args.currentPos, tt.args.moveDirection); got != tt.want {
// 				t.Errorf("calcNextAxisCrossPos() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestMakePosLenList(t *testing.T) {
	pll, _ := CalcXYLenListLine(0, -1, 13, -11)
	for i, v := range pll {
		fmt.Printf("%v %v\n", i, v)
	}
}

func TestMakeSightlinesByXYLenList(t *testing.T) {
	sightlinesByXYLenList := MakeSightlinesByXYLenList(viewportdata.ViewportXYLenList)
	for _, v := range sightlinesByXYLenList {
		fmt.Printf("%v\n", v)
	}
}
