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
	prefix     string
	currentu64 uint64
}

func New(prefix string) *Maker {
	mk := &Maker{
		prefix: prefix,
	}
	return mk
}

func (mk *Maker) New() string {
	newValue := atomic.AddUint64(&mk.currentu64, 1)
	// return fmt.Sprintf("%s%X", mk.prefix, newValue)
	return mk.prefix + makeStr2(newValue)
}

func makeStr1(src uint64) string {
	rtn8 := make([]byte, 8)
	binary.BigEndian.PutUint64(rtn8, src)
	rtn16 := make([]byte, 16)
	hex.Encode(rtn16, rtn8)
	return string(rtn16)
}

const hextable = "0123456789abcdef"

func makeStr2(src uint64) string {
	rtn16 := make([]byte, 16)
	for i := 15; i >= 0; i-- {
		remain := src & 0xf
		src >>= 4
		rtn16[i] = hextable[remain]
	}
	return string(rtn16)
}

var G_Maker = New("gmaker")
