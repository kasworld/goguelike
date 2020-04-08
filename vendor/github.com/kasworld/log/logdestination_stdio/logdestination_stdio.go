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

// standard input/output/error log destination
package logdestination_stdio

import "os"

type StdErr struct {
}

func NewStdErr() *StdErr {
	return &StdErr{}
}

func (lo StdErr) String() string {
	return "StdErr"
}

func (lo *StdErr) Name() string {
	return "StdErr"
}

func (lo *StdErr) Reload() error {
	os.Stderr.Sync()
	return nil
}

func (lo *StdErr) Write(b []byte) error {
	_, err := os.Stderr.Write(b)
	return err
}

type StdOut struct {
}

func NewStdOut() *StdOut {
	return &StdOut{}
}

func (lo StdOut) String() string {
	return "StdOut"
}

func (lo *StdOut) Name() string {
	return "StdOut"
}

func (lo *StdOut) Reload() error {
	os.Stdout.Sync()
	return nil
}

func (lo *StdOut) Write(b []byte) error {
	_, err := os.Stdout.Write(b)
	return err
}
