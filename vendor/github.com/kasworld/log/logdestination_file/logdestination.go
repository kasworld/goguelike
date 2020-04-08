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

// file log destination
package logdestination_file

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type LogDestinationFile struct {
	mutex sync.Mutex

	name  string
	fpath string
	out   io.WriteCloser
}

func New(fpath string) (*LogDestinationFile, error) {
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	rtn := &LogDestinationFile{
		name: fpath,
		out:  f,
	}
	rtn.fpath = fpath
	return rtn, nil
}

func (lo LogDestinationFile) String() string {
	return fmt.Sprintf(
		"LogDestinationFile[Name:%v Path:%v]",
		lo.name,
		lo.fpath)
}

func (lo *LogDestinationFile) Name() string {
	return lo.name
}

func (lo *LogDestinationFile) Reload() error {
	lo.mutex.Lock()
	defer lo.mutex.Unlock()

	lo.out.Close()

	out, err := os.OpenFile(lo.fpath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("%v failed to reload, %v", lo, err)
	}
	lo.out = out
	return nil
}

func (lo *LogDestinationFile) Write(b []byte) error {
	lo.mutex.Lock()
	defer lo.mutex.Unlock()

	_, err := lo.out.Write(b)
	return err
}
