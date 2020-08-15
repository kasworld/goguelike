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

package condition

func (cn Condition) HideSelf() bool {
	return attrib[cn].hideSelf
}

func (cn Condition) HideOther() bool {
	return attrib[cn].hideOther
}

var attrib = [Condition_Count]struct {
	hideSelf  bool // hide to client self ao
	hideOther bool // hide to client other ao
}{
	Blind:     {false, false},
	Invisible: {false, false},
	Burden:    {false, false},
	Float:     {false, false},
	Greasy:    {false, false},
	Drunken:   {false, false},
	Sleep:     {false, false},
	Contagion: {false, false},
}
