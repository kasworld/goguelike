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

package launcherlib

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func PS(pidlist []string) ([]string, error) {
	pidstr := strings.Join(pidlist, " ")
	cmd := exec.Command(
		"ps",
		"-p",
		pidstr,
		"-o",
		"pid,comm,cputime,pcpu,rss,pmem",
	)
	outbytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	rtn := splitAndTrimString(string(outbytes), "\n")
	return rtn, nil
}

func ListPidFiles(prefix string, pidDir string) ([]string, error) {
	rtn := make([]string, 0)
	files, err := ioutil.ReadDir(pidDir)
	if err != nil {
		return rtn, err
	}

	for _, file := range files {
		// g2log.Debug("%v", file)
		if !file.IsDir() &&
			strings.HasSuffix(file.Name(), ".pid") &&
			strings.HasPrefix(file.Name(), prefix) {
			rtn = append(rtn, file.Name())
		}
	}
	return rtn, nil
}

func ListExeFiles(exedir string) ([]string, error) {
	rtn := make([]string, 0)
	files, err := ioutil.ReadDir(exedir)
	if err != nil {
		return rtn, err
	}
	for _, file := range files {
		if !file.IsDir() {
			rtn = append(rtn, file.Name())
		}
	}
	return rtn, nil
}

func splitAndTrimString(src string, sep string) []string {
	splitedSrc := strings.Split(strings.TrimSpace(src), sep)
	rtn := make([]string, 0, len(splitedSrc))
	for _, v := range splitedSrc {
		ss := strings.TrimSpace(v)
		if len(ss) == 0 {
			continue
		}
		rtn = append(rtn, ss)
	}
	return rtn
}

func getPidStringFromPidFile(pidfilename string) (string, error) {
	fd, err := os.Open(pidfilename)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	pidbytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return "", err
	}
	pidstr := strings.TrimSpace(string(pidbytes))
	return pidstr, nil
}
