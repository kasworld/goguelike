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

// Package uuidstr make string form uuid
package g2id

import (
	"bytes"
	"crypto/rand"

	"github.com/kasworld/uuidstr"
)

type UUID string

func New() UUID {
	readFn := rand.Read
	var rtn [16]byte
	if _, err := readFn(rtn[:]); err != nil {
		panic(err.Error()) // rand should never fail
	}
	return UUID(rtn[:])
}

func (g2id UUID) String() string {
	buf := bytes.NewBufferString(string(g2id))
	return uuidstr.ToString(buf.Bytes())
}

func NewFromString(uuid string) UUID {
	return UUID(uuidstr.Parse(uuid))
}
