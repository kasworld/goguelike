// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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
	GeoMin [3]float64
	GeoMax [3]float64
	GeoLen [3]float64
}

func GetGeoInfo(geo js.Value) GeoInfo {
	var rtn GeoInfo
	geo.Call("computeBoundingBox")
	minbox := geo.Get("boundingBox").Get("min")
	maxbox := geo.Get("boundingBox").Get("max")
	rtn.GeoMin = [3]float64{
		minbox.Get("x").Float(),
		minbox.Get("y").Float(),
		minbox.Get("z").Float(),
	}
	rtn.GeoMax = [3]float64{
		maxbox.Get("x").Float(),
		maxbox.Get("y").Float(),
		maxbox.Get("z").Float(),
	}
	rtn.GeoLen = [3]float64{
		maxbox.Get("x").Float() - minbox.Get("x").Float(),
		maxbox.Get("y").Float() - minbox.Get("y").Float(),
		maxbox.Get("z").Float() - minbox.Get("z").Float(),
	}
	return rtn
}

func CalcGeoMinMaxX(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("x").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("x").Float()
	return geoMin, geoMax
}

func CalcGeoMinMaxY(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("y").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("y").Float()
	return geoMin, geoMax
}

func CalcGeoMinMaxZ(geo js.Value) (float64, float64) {
	geo.Call("computeBoundingBox")
	geoMax := geo.Get("boundingBox").Get("max").Get("z").Float()
	geoMin := geo.Get("boundingBox").Get("min").Get("z").Float()
	return geoMin, geoMax
}
