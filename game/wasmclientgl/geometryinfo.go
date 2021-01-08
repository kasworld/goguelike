// Copyright 2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
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
	"syscall/js"
)

type GeoInfo struct {
	Min [3]float64
	Max [3]float64
	Len [3]float64
}

func GetGeoInfo(geo js.Value) GeoInfo {
	var rtn GeoInfo
	geo.Call("computeBoundingBox")
	minbox := geo.Get("boundingBox").Get("min")
	maxbox := geo.Get("boundingBox").Get("max")
	rtn.Min = [3]float64{
		minbox.Get("x").Float(),
		minbox.Get("y").Float(),
		minbox.Get("z").Float(),
	}
	rtn.Max = [3]float64{
		maxbox.Get("x").Float(),
		maxbox.Get("y").Float(),
		maxbox.Get("z").Float(),
	}
	rtn.Len = [3]float64{
		maxbox.Get("x").Float() - minbox.Get("x").Float(),
		maxbox.Get("y").Float() - minbox.Get("y").Float(),
		maxbox.Get("z").Float() - minbox.Get("z").Float(),
	}
	return rtn
}
