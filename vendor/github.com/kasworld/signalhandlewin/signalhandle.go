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
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nightlyone/lockfile"
)

type LoggerI interface {
	Reload() error
	Fatal(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type ServiceI interface {
	GetServiceLockFilename() string
	ServiceInit() error
	ServiceMain(ctx context.Context)
	ServiceCleanup()
	GetLogger() interface{} // cannot use LoggerI
}

func makeLockfile(svr ServiceI, filename string) (lockfile.Lockfile, error) {
	lock, err := lockfile.New(filename)
	if err != nil {
		svr.GetLogger().(LoggerI).Error("Cannot init lock. reason: %v\n", err)
		return lock, err
	}
	// Error handling is essential, as we only try to get the lock.
	if err := lock.TryLock(); err != nil {
		svr.GetLogger().(LoggerI).Error("Cannot lock \"%v\", reason: %v\n", lock, err)
		return lock, err
	}
	return lock, nil
}

func RunWithSignalHandle(svr ServiceI) error {
	svr.GetLogger().(LoggerI).Debug("Start RunWithSignalHandle %v", svr)
	defer func() { svr.GetLogger().(LoggerI).Debug("End RunWithSignalHandle %v", svr) }()

	_, ok := svr.GetLogger().(LoggerI)
	if !ok {
		return fmt.Errorf("GetLogger return logger implement LoggerI")
	}

	if lockpfilename := svr.GetServiceLockFilename(); lockpfilename != "" {
		lockfileobj, err := makeLockfile(svr, svr.GetServiceLockFilename())
		if err != nil {
			svr.GetLogger().(LoggerI).Fatal("Make Lockfile fail, %v", err)
			return err
		}
		defer lockfileobj.Unlock()
	}

	if err := svr.ServiceInit(); err != nil {
		svr.GetLogger().(LoggerI).Fatal("Server Initialize Failed, %v", err)
		return err
	}
	var wg sync.WaitGroup
	ctxMain, endService := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		svr.ServiceMain(ctxMain)
		wg.Done()
		endService()
	}()

	// signal handle
	// kill -s SIGUSR1 `cat pidfilename`
	// sigUsr1Ch := make(chan os.Signal, 10)
	// signal.Notify(sigUsr1Ch, syscall.SIGUSR1) // for logrotate

	signalTermIntCh := make(chan os.Signal, 1)
	signal.Notify(signalTermIntCh, os.Interrupt, syscall.SIGTERM)
loop:
	for {
		select {
		case <-ctxMain.Done():
			wg.Wait()
			svr.ServiceCleanup()
			svr.GetLogger().(LoggerI).Reload()
			break loop

		case s := <-signalTermIntCh:
			svr.GetLogger().(LoggerI).Debug("Got signal: %v", s)
			signal.Stop(signalTermIntCh)
			endService()

			// case <-sigUsr1Ch:
			// 	fmt.Println("Catch sigusr1, reopen logfile")
			// 	svr.GetLogger().(LoggerI).Reload()
		}
	}
	return nil
}
