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

package g2rand

import (
	"math/rand"
	"sync"
	"time"
)

type G2Rand struct {
	mutex sync.Mutex
	mr    *rand.Rand
	src   rand.Source64
}

func New() *G2Rand {
	return NewWithSeed(time.Now().UnixNano())
}

func NewWithSeed(seed int64) *G2Rand {
	rndSource := rand.NewSource(seed)
	return &G2Rand{
		mr:  rand.New(rndSource),
		src: rndSource.(rand.Source64),
	}
}

// include s, not include e
func (rnd *G2Rand) IntRange(s, e int) int {
	r := e - s
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return rnd.mr.Intn(r) + s
}

func (rnd *G2Rand) Read(dst []byte) (int, error) {
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return rnd.mr.Read(dst)
}

func (rnd *G2Rand) NormIntRange(desiredMean, desiredStdDev int) int {
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return int(rnd.mr.NormFloat64()*float64(desiredStdDev)) + desiredMean
}

func (rnd *G2Rand) NormFloat64Range(desiredMean, desiredStdDev float64) float64 {
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return rnd.mr.NormFloat64()*desiredStdDev + desiredMean
}

func (rnd *G2Rand) Float64() float64 {
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return rnd.mr.Float64()
}

func (rnd *G2Rand) Perm(n int) []int {
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return rnd.mr.Perm(n)
}

func (rnd *G2Rand) Intn(n int) int {
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return rnd.mr.Intn(n)
}

func (rnd *G2Rand) Int63() int64 {
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return rnd.mr.Int63()
}

func (rnd *G2Rand) Uint32() uint32 {
	rnd.mutex.Lock()
	defer rnd.mutex.Unlock()
	return rnd.mr.Uint32()
}

func (r *G2Rand) Shuffle(n int, swap func(i, j int)) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.mr.Shuffle(n, swap)
}
