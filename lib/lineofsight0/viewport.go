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
	"math"

	"github.com/kasworld/findnear"
)

func CalcXYLenListLine(x1, y1, x2, y2 int) (findnear.XYLenList, error) {
	srcPos := [2]float64{float64(x1) + 0.5, float64(y1) + 0.5}
	dstPos := [2]float64{float64(x2) + 0.5, float64(y2) + 0.5}

	if srcPos == dstPos {
		return nil, nil
	}
	srcXf64, srcYf64 := srcPos[0], srcPos[1]
	deltaX, deltaY := dstPos[0]-srcPos[0], dstPos[1]-srcPos[1]
	deltaXi, deltaYi := x2-x1, y2-y1
	xSign, xLen := SignAbsf(deltaX)
	ySign, yLen := SignAbsf(deltaY)

	// prealloc for speed
	xyLenList := make(findnear.XYLenList, 0, int(xLen+yLen+3))

	currentXf64, currentYf64 := srcXf64, srcYf64
	var advXX, advXY, advYX, advYY, advXLen, advYLen float64
	var curXi, curYi int

	for {
		// get next x by y
		advYY = calcNextAxisCrossPos(currentYf64, ySign)
		advYX = srcXf64 + (advYY-srcYf64)*deltaX/deltaY
		advYLen = Lenf(currentXf64, currentYf64, advYX, advYY)

		// get next y by x
		advXX = calcNextAxisCrossPos(currentXf64, xSign)
		advXY = srcYf64 + (advXX-srcXf64)*deltaY/deltaX
		advXLen = Lenf(currentXf64, currentYf64, advXX, advXY)

		var nextXf64, nextYf64 float64
		var xInc, yInc int
		if (deltaXi != 0 && advXLen < advYLen) || deltaYi == 0 {
			// step X
			nextXf64, nextYf64 = advXX, advXY
			xInc = int(xSign)
		} else if (deltaYi != 0 && advXLen > advYLen) || deltaXi == 0 {
			// step Y
			nextXf64, nextYf64 = advYX, advYY
			yInc = int(ySign)
		} else {
			nextXf64, nextYf64 = advYX, advYY
			xInc = int(xSign)
			yInc = int(ySign)
		}

		if (deltaXi != 0 && Absf(nextXf64-srcXf64) >= xLen) ||
			(deltaYi != 0 && Absf(nextYf64-srcYf64) >= yLen) {
			break
		}
		lenInTile := Lenf(currentXf64, currentYf64, nextXf64, nextYf64)
		xyLenList = append(xyLenList,
			findnear.XYLen{X: curXi, Y: curYi, L: lenInTile},
		)
		curXi += xInc
		curYi += yInc
		currentXf64, currentYf64 = nextXf64, nextYf64
		if len(xyLenList) > int(xLen+yLen+2) {
			return nil, fmt.Errorf(
				"line too long %v %v %v %v",
				srcPos, dstPos,
				xLen+yLen,
				xyLenList)
		}
	}
	lenInTile := Lenf(currentXf64, currentYf64, dstPos[0], dstPos[1])
	xyLenList = append(xyLenList,
		findnear.XYLen{X: curXi, Y: curYi, L: lenInTile},
	)
	return xyLenList, nil
}

func calcNextAxisCrossPos(currentPos, moveDirection float64) float64 {
	intCurrentPos := math.Trunc(currentPos)
	if currentPos > 0 {
		if moveDirection > 0 {
			return intCurrentPos + 1
		} else if moveDirection < 0 {
			return intCurrentPos
		} else {
			return currentPos
		}
	} else if currentPos < 0 {
		if moveDirection > 0 {
			return intCurrentPos
		} else if moveDirection < 0 {
			return intCurrentPos - 1
		} else {
			return currentPos
		}
	} else {
		if moveDirection > 0 {
			return 1
		} else if moveDirection < 0 {
			return -1
		} else {
			return 0
		}
	}
}

func calcXYLen(x1, y1, x2, y2 float64) (int, int, float64) {
	l := Lenf(x1, y1, x2, y2)
	x := int((x1 + x2) / 2)
	y := int((y1 + y2) / 2)
	return x, y, l
}

func Lenf(x1, y1, x2, y2 float64) float64 {
	xl := x1 - x2
	yl := y1 - y2
	r := math.Sqrt(xl*xl + yl*yl)
	return r
}

func FloatEqual(a, b float64) bool {
	absGap := math.Abs(a - b)
	pow := math.Pow10(-5)

	return pow > absGap
}

func SignAbsf(n float64) (float64, float64) {
	if n > 0 {
		return 1, n
	} else if n < 0 {
		return -1, -n
	} else {
		return 0, 0
	}
}

func Absf(n float64) float64 {
	if n > 0 {
		return n
	} else if n < 0 {
		return -n
	} else {
		return 0
	}
}
