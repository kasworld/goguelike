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

package ifop

import "fmt"

func F64ZeroString(f float64, zerostr, notzerostr string) string {
	if f == 0 {
		return zerostr
	} else {
		return notzerostr
	}
}

func F64ZeroFormat(f float64, zerostr, notzerostr string) string {
	if f == 0 {
		if zerostr == "" {
			return ""
		}
		return fmt.Sprintf(zerostr, f)
	} else {
		if notzerostr == "" {
			return ""
		}
		return fmt.Sprintf(notzerostr, f)
	}
}

func IntZeroString(i int, zerostr, notzerostr string) string {
	if i == 0 {
		return zerostr
	} else {
		return notzerostr
	}
}
func IntZeroFormat(i int, zerostr, notzerostr string) string {
	if i == 0 {
		if zerostr == "" {
			return ""
		}
		return fmt.Sprintf(zerostr, i)
	} else {
		if notzerostr == "" {
			return ""
		}
		return fmt.Sprintf(notzerostr, i)
	}
}
