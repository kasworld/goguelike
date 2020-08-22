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

package signalhandlewin

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kasworld/launcherlib"
)

var (
	signalhandle_ArgAdded bool
	signalhandle_service  *string
)

func AddArgs() {
	signalhandle_ArgAdded = true
	signalhandle_service = flag.String("service", "start", "start,stop,restart,forcestart")
}

func StartByArgs(svr ServiceI) error {
	if !signalhandle_ArgAdded {
		return fmt.Errorf("AddArgs Not called")
	}
	switch *signalhandle_service {
	default:
		return fmt.Errorf("unknown service arg %v", *signalhandle_service)
	case "start":
		svr.GetLogger().(LoggerI).Debug("Service Start %v", svr)
		RunWithSignalHandle(svr)
	case "restart":
		svr.GetLogger().(LoggerI).Debug("Service ReStart %v", svr)
		if err := SignalToStopAndWaitServiceEnd(svr); err != nil {
			svr.GetLogger().(LoggerI).Debug("%v", err)
		}
		if !launcherlib.IsFileExist(svr.GetServiceLockFilename()) {
			RunWithSignalHandle(svr)
			return nil
		}
		return fmt.Errorf("fail to stop service %v", svr)
	case "forcestart":
		svr.GetLogger().(LoggerI).Debug("Service Force Start %v", svr)
		if err := SignalToStopAndWaitServiceEnd(svr); err != nil {
			svr.GetLogger().(LoggerI).Debug("%v", err)
		}
		if !launcherlib.IsFileExist(svr.GetServiceLockFilename()) {
			RunWithSignalHandle(svr)
			return nil
		}
		// force start
		os.Remove(svr.GetServiceLockFilename())
		RunWithSignalHandle(svr)
		return nil
	case "stop":
		svr.GetLogger().(LoggerI).Debug("Service Stop %v", svr)
		return SignalToStopAndWaitServiceEnd(svr)
		// case "logreopen":
		// 	svr.GetLogger().(LoggerI).Debug("Service log reopen %v", svr)
		// 	launcherlib.SignalByPidFile(svr.GetServiceLockFilename(), syscall.SIGUSR1)
		// 	return nil
	}
	return nil
}

func SignalToStopAndWaitServiceEnd(svr ServiceI) error {
	pid, err := launcherlib.GetPidIntFromPidFile(svr.GetServiceLockFilename())
	if err != nil {
		svr.GetLogger().(LoggerI).Debug("%v", err)
		return nil
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		svr.GetLogger().(LoggerI).Debug("%v", err)
		return nil
	}
	if err := p.Signal(os.Interrupt); err != nil {
		svr.GetLogger().(LoggerI).Debug("%v", err)
		return nil
	}

	// wait end
	for i := 0; i < 60; i++ {
		if !launcherlib.IsFileExist(svr.GetServiceLockFilename()) {
			return nil
		}
		svr.GetLogger().(LoggerI).Debug("wait service end %v", i)
		time.Sleep(1000 * time.Millisecond)
	}
	return fmt.Errorf("fail to stop service")
}
