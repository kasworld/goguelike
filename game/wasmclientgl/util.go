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

package wasmclientgl

import (
	"fmt"
	"syscall/js"
	"time"
)

func CalcCurrentFrame(difftick int64, fps float64) int {
	diffsec := float64(difftick) / float64(time.Second)
	frame := fps * diffsec
	return int(frame)
}

func getImgWH(srcImageID string) (js.Value, float64, float64) {
	img := js.Global().Get("document").Call("getElementById", srcImageID)
	if !img.Truthy() {
		fmt.Printf("fail to get %v", srcImageID)
		return js.Null(), 0, 0
	}
	srcw := img.Get("naturalWidth").Float()
	srch := img.Get("naturalHeight").Float()
	return img, srcw, srch
}

func SetPosition(jso js.Value, pos ...interface{}) {
	po := jso.Get("position")
	if len(pos) >= 1 {
		po.Set("x", pos[0])
	}
	if len(pos) >= 2 {
		po.Set("y", pos[1])
	}
	if len(pos) >= 3 {
		po.Set("z", pos[2])
	}
}
