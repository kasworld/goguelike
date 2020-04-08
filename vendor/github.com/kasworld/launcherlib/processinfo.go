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
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PSInfo struct {
	PID         string
	COMMAND     string
	TIME        string
	CPU         string
	RSS         string
	MEM         string
	PidFileName string
	WebURL      string
}

func PidFileList2PSInfo(
	pidfilelist []string, basedir string) ([]PSInfo, error) {

	pid2file := make(map[string]string)
	pidlist := make([]string, 0)
	for _, v := range pidfilelist {
		pidstr, err := getPidStringFromPidFile(basedir + v)
		if err != nil {
			// g2log.Error("%v", err)
			continue
		}
		pidlist = append(pidlist, pidstr)
		pid2file[pidstr] = v
	}
	if len(pidlist) < 1 {
		return nil, fmt.Errorf("No pid file found")
	}
	pslist, err := PS(pidlist)
	if err != nil {
		return nil, err
	}
	if len(pslist) <= 1 {
		return nil, fmt.Errorf("No pid file found")
	}
	PSInfoList := make([]PSInfo, 0)

	for _, v := range pslist[1:] {
		infos := splitAndTrimString(v, " ")
		pidfilename := pid2file[infos[0]]
		PSInfoList = append(PSInfoList, PSInfo{
			PID:         infos[0],
			COMMAND:     infos[1],
			TIME:        infos[2],
			CPU:         infos[3],
			RSS:         infos[4],
			MEM:         infos[5],
			PidFileName: pidfilename,
			// WebURL:      g2.GetListenPortFromPidFile(pidfilename),
		})
	}
	return PSInfoList, nil
}

////////////////////////////////

func (si PSInfo) String() string {
	return fmt.Sprintf("PID:%v %v %v CPU:%v RSS:%v MEM:%v %v %v",
		si.PID,
		si.COMMAND,
		si.TIME,
		si.CPU,
		si.RSS,
		si.MEM,
		si.PidFileName,
		si.WebURL,
	)
}

func (si *PSInfo) GetPid() (int, error) {
	int64v, err := strconv.ParseInt(si.PID, 0, 64)
	if err != nil {
		return 0, err
	}
	rtnInt := int(int64v)
	return rtnInt, nil
}

func (si *PSInfo) GetRss() (int, error) {
	int64v, err := strconv.ParseInt(si.RSS, 0, 64)
	if err != nil {
		return 0, err
	}
	rtnInt := int(int64v)
	return rtnInt, nil
}

func (si *PSInfo) GetCpu() (float64, error) {
	f64v, err := strconv.ParseFloat(si.CPU, 64)
	if err != nil {
		return 0, err
	}
	rtnF64 := float64(f64v)
	return rtnF64, nil
}

func (si *PSInfo) GetMem() (float64, error) {
	f64v, err := strconv.ParseFloat(si.MEM, 64)
	if err != nil {
		return 0, err
	}
	rtnF64 := float64(f64v)
	return rtnF64, nil
}

func (si *PSInfo) DownLoadHeap() error {
	hostbase := "http://localhost"
	heapurl := "/debug/pprof/heap"
	inurl := fmt.Sprintf("%v%v%v", hostbase, si.WebURL, heapurl)
	outname := fmt.Sprintf("%v.%v.heap", si.PidFileName, time.Now().UnixNano())

	inwebrsp, err := http.Get(inurl)
	if err != nil {
		return err
	}
	defer inwebrsp.Body.Close()

	outfd, err := os.Create(outname)
	if err != nil {
		return err
	}
	defer outfd.Close()

	_, err = io.Copy(outfd, inwebrsp.Body)
	return err
}
