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

package fieldobjdisplaytype

import "github.com/kasworld/htmlcolors"

func (v FieldObjDisplayType) Color24() htmlcolors.Color24 {
	return attrib[v].Color24
}

func (v FieldObjDisplayType) Rune() string {
	return attrib[v].Rune
}

var attrib = [FieldObjDisplayType_Count]struct {
	Rune    string
	Color24 htmlcolors.Color24
}{
	None:         {"?", htmlcolors.Black},
	StairUp:      {"/\\", htmlcolors.Black},
	StairDn:      {"\\/", htmlcolors.Black},
	PortalIn:     {"[+]", htmlcolors.Black},
	PortalAutoIn: {"{+}", htmlcolors.Black},
	PortalOut:    {"[-]", htmlcolors.Black},
	Recycler:     {"*", htmlcolors.Black},
	LightHouse:   {"X", htmlcolors.Black},
	GateKeeper:   {"-|-", htmlcolors.Black},
}
