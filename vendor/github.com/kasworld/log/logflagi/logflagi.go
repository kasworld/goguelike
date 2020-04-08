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

// Package logflagi define log flag interface
package logflagi

import "time"

type LogFlagI interface {
	BitAnd(lf2 LogFlagI) LogFlagI
	BitOr(lf2 LogFlagI) LogFlagI
	BitNeg(lf2 LogFlagI) LogFlagI
	BitClear(lf2 LogFlagI) LogFlagI
	BitTest(lf2 LogFlagI) bool

	FlagString() string
	FormatHeader(
		buf *[]byte,
		calldepth int,
		now time.Time,
		prefix string, llinfo string)

	ParseHeader(buf []byte) (
		prefix string, llinfo string,
		datestr string, timestr string,
		filestr string,
		remainbuf []byte)
}
