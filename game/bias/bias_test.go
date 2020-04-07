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

package bias

import (
	"testing"

	"github.com/kasworld/g2rand"
)

func TestSelectRandSkill(t *testing.T) {
	rnd := g2rand.New()
	b := Bias{-0.2, -0.9, -3}
	t.Log(b)

	for i := 0; i < 10; i++ {
		skilli := rnd.Intn(3)
		skillv := b.SelectSkill(skilli)
		t.Log(skilli, skillv)
		if skillv < 0 {
			t.Fail()
		}
	}
}

func TestBias_IntString(t *testing.T) {
	b := Bias{-0.2, -0.9, -3}

	t.Logf("%v", b.String())
	t.Logf("%v", b.IntString())
	t.Logf("%v", b.F42String())
}
