// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package statusoptype

import "github.com/kasworld/htmlcolors"

var statusOp2Color = [StatusOpType_Count]htmlcolors.Color24{
	None:      htmlcolors.Black,
	AddHP:     htmlcolors.PaleVioletRed,
	AddSP:     htmlcolors.LightGreen,
	AddHPRate: htmlcolors.OrangeRed,
	AddSPRate: htmlcolors.LimeGreen,

	ModSight: htmlcolors.LightYellow,

	RndFaction:   htmlcolors.MidnightBlue,
	IncFaction:   htmlcolors.MidnightBlue,
	SetFaction:   htmlcolors.MidnightBlue,
	ResetFaction: htmlcolors.MidnightBlue,

	NegBias:         htmlcolors.BlueViolet,
	RotateBiasRight: htmlcolors.BlueViolet,
	RotateBiasLeft:  htmlcolors.BlueViolet,

	SetCondition: htmlcolors.DeepPink,

	ForgetFloor:    htmlcolors.PaleVioletRed,
	ForgetOneFloor: htmlcolors.PaleVioletRed,
}

func (sot StatusOpType) Color24() htmlcolors.Color24 {
	return statusOp2Color[sot]
}
