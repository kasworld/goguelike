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

package condition_flag

import (
	"github.com/kasworld/goguelike/enum/condition"
)

func HideSelf() ConditionFlag {
	return hideSelf
}

func HideOther() ConditionFlag {
	return hideOther
}

var hideSelf ConditionFlag
var hideOther ConditionFlag

func init() {
	for i := 0; i < condition.Condition_Count; i++ {
		cn := condition.Condition(i)
		if cn.HideSelf() {
			hideSelf.SetByCondition(cn)
		}
		if cn.HideOther() {
			hideOther.SetByCondition(cn)
		}
	}
}
