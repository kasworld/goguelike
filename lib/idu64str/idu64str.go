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

package idu64str

import (
	"sync/atomic"
)

type Maker struct {
	buf16      []byte
	currentu64 uint64
}

func New(prefix string) *Maker {
	mk := &Maker{}
	prefixbyte := []byte(prefix)
	mk.buf16 = make([]byte, 16+len(prefixbyte))
	copy(mk.buf16, prefixbyte)
	return mk
}

const hextable = "0123456789abcdef"

func (mk *Maker) New() string {
	newValue := atomic.AddUint64(&mk.currentu64, 1)
	l := len(mk.buf16) - 1
	for i := 0; i < 16; i++ {
		mk.buf16[l-i] = hextable[newValue&0xf]
		newValue >>= 4
	}
	return string(mk.buf16)
}

var G_Maker = New("gmaker")
