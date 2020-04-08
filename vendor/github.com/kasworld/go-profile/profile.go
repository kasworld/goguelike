// Copyright 2015 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package profile

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

func WriteHeapProfile() {
	f, err := os.Create(*memprofilename)
	if err != nil {
		log.Fatalf("mem profile %v", err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
}

func StartCPUProfile() func() {
	f, err := os.Create(*cpuprofilename)
	if err != nil {
		log.Fatalf("cpu profile %v", err)
	}
	pprof.StartCPUProfile(f)
	return pprof.StopCPUProfile
}

var cpuprofilename *string
var memprofilename *string

func AddArgs() {
	cpuprofilename = flag.String("cpuprofilename", "", "cpu profile filename")
	memprofilename = flag.String("memprofilename", "", "memory profile filename")
}

func IsCpu() bool {
	return *cpuprofilename != ""
}

func IsMem() bool {
	return *memprofilename != ""
}
