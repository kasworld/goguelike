// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"

	"github.com/kasworld/fortune"
)

func main() {
	chatlines, err := fortune.LoadDir("fortunes")

	if err != nil {
		log.Fatal("load fortune fail %v", err)
		return
	}
	fd, err := os.Create("chatdata.txt")
	if err != nil {
		log.Fatal("Create chatdata fail %v", err)
		return
	}
	defer fd.Close()
	for _, v := range chatlines {
		if _, err := fd.WriteString(v); err != nil {
			log.Fatal("Write chatlines fail %v", err)
			return
		}
	}
}
