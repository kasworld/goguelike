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

package aiplan

import "github.com/kasworld/htmlcolors"

func (ct AIPlan) Color24() htmlcolors.Color24 {
	return attrib[ct].Color
}

var attrib = [AIPlan_Count]struct {
	Color htmlcolors.Color24
}{
	None:           {htmlcolors.Yellow},
	Chat:           {htmlcolors.Yellow},
	StrollAround:   {htmlcolors.Yellow},
	Move2Dest:      {htmlcolors.Yellow},
	Revenge:        {htmlcolors.Yellow},
	UsePortal:      {htmlcolors.Yellow},
	MoveToRecycler: {htmlcolors.Yellow},
	RechargeSafe:   {htmlcolors.Yellow},
	RechargeCan:    {htmlcolors.Yellow},
	PickupCarryObj: {htmlcolors.Yellow},
	Equip:          {htmlcolors.Yellow},
	UsePotion:      {htmlcolors.Yellow},
	Attack:         {htmlcolors.Yellow},
	MoveStraight3:  {htmlcolors.Yellow},
	MoveStraight5:  {htmlcolors.Yellow},
}
