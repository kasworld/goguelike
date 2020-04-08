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

package primenum

import (
	"math"
)

func isprime(n int64, foundPrimes []int64) bool {
	to := int64(math.Sqrt(float64(n)))
	for _, v := range foundPrimes {
		if n%v == 0 {
			return false
		}
		if v > to {
			break
		}
	}
	return true
}

func MakePrimes(n int64) []int64 {
	foundPrimes := append(make([]int64, 0, n/10), 2)

	for i := int64(3); i < n; i += 2 {
		if isprime(i, foundPrimes) {
			foundPrimes = append(foundPrimes, i)
		}
	}
	return foundPrimes
}
