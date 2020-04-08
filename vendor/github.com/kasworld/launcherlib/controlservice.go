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
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func RunAndWaitPidFile(
	exe string, args []string, outfile string, pidfile string) error {

	if IsFileExist(pidfile) {
		return fmt.Errorf("already run? file exist %v", pidfile)
	}
	var err error
	go func() {
		err = RunExeSync(exe, args, outfile)
	}()
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		if IsFileExist(pidfile) {
			return nil
		}
	}
	return fmt.Errorf("run fail? not found %v %v", pidfile, err)
}

func RunExeSync(exe string, args []string, outfile string) error {
	fd, err := os.OpenFile(outfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, args...)
	cmd.Stderr = fd
	cmd.Stdout = fd
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Setpgid: true,
	// 	Pgid:    0,
	// }
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func IsFileExist(filename string) bool {
	if _, err := os.Stat(filename); err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func SignalByPidFile(pidfilename string, sig os.Signal) error {
	pid, err := GetPidIntFromPidFile(pidfilename)
	if err != nil {
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	// os.Interrupt, syscall.SIGTERM ,syscall.SIGINT os.Signal
	if err := p.Signal(sig); err != nil {
		return err
	}
	return nil
}

func GetPidIntFromPidFile(pidfilename string) (int, error) {
	pidstr, err := getPidStringFromPidFile(pidfilename)
	if err != nil {
		return 0, err
	}
	int64v, err := strconv.ParseInt(pidstr, 0, 64)
	if err != nil {
		return 0, err
	}
	pid := int(int64v)
	return pid, nil
}

func ProcessAlive(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	} else {
		err := process.Signal(syscall.Signal(0))
		return err
	}
}
