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
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
)

type Maker struct {
	prefix     []byte
	buf16      []byte
	currentu64 uint64
}

func New(prefix string) *Maker {
	mk := &Maker{
		prefix: []byte(prefix),
	}
	mk.buf16 = make([]byte, 16+len(mk.prefix))
	copy(mk.buf16, mk.prefix)
	return mk
}

const hextable = "0123456789abcdef"

func (mk *Maker) New() string {
	newValue := atomic.AddUint64(&mk.currentu64, 1)
	l := len(mk.prefix)
	for i := 15; i >= 0; i-- {
		mk.buf16[i+l] = hextable[newValue&0xf]
		newValue >>= 4
	}
	return string(mk.buf16)
}

var G_Maker = New("gmaker")

func (mk *Maker) New1() string {
	newValue := atomic.AddUint64(&mk.currentu64, 1)
	l := len(mk.prefix)
	rtn8 := make([]byte, 8)
	binary.BigEndian.PutUint64(rtn8, newValue)
	rtn16 := make([]byte, 16+l)
	copy(rtn16, mk.prefix)
	hex.Encode(rtn16[l:], rtn8)
	return string(rtn16)
}
