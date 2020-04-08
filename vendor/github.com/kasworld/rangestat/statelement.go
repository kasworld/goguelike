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

package rangestat

import (
	"fmt"
	"time"
)

type statElement struct {
	// for stat
	beginTime time.Time
	decCount  int
	incCount  int
}

func (se statElement) String() string {
	dur := time.Now().UTC().Sub(se.beginTime).Seconds()
	return fmt.Sprintf("statElement[Inc %v/%5.1f/s Dec %v/%5.1f/s]",
		se.incCount, float64(se.incCount)/dur,
		se.decCount, float64(se.decCount)/dur,
	)
}
