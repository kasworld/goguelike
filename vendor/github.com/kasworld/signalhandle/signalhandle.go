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

package signalhandle

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
	GetLogger() LoggerI
}

func makeLockfile(svr ServiceI, filename string) (lockfile.Lockfile, error) {
	lock, err := lockfile.New(filename)
	if err != nil {
		svr.GetLogger().Error("Cannot init lock. reason: %v\n", err)
		return lock, err
	}
	// Error handling is essential, as we only try to get the lock.
	if err := lock.TryLock(); err != nil {
		svr.GetLogger().Error("Cannot lock \"%v\", reason: %v\n", lock, err)
		return lock, err
	}
	return lock, nil
}

func RunWithSignalHandle(svr ServiceI) error {
	svr.GetLogger().Debug("Start RunWithSignalHandle %v", svr)
	defer func() { svr.GetLogger().Debug("End RunWithSignalHandle %v", svr) }()

	if lockpfilename := svr.GetServiceLockFilename(); lockpfilename != "" {
		lockfileobj, err := makeLockfile(svr, svr.GetServiceLockFilename())
		if err != nil {
			svr.GetLogger().Fatal("Make Lockfile fail, %v", err)
			return err
		}
		defer lockfileobj.Unlock()
	}

	if err := svr.ServiceInit(); err != nil {
		svr.GetLogger().Fatal("Server Initialize Failed, %v", err)
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
	sigUsr1Ch := make(chan os.Signal, 10)
	signal.Notify(sigUsr1Ch, syscall.SIGUSR1) // for logrotate

	signalTermIntCh := make(chan os.Signal, 1)
	signal.Notify(signalTermIntCh, os.Interrupt, syscall.SIGTERM)
loop:
	for {
		select {
		case <-ctxMain.Done():
			wg.Wait()
			svr.ServiceCleanup()
			svr.GetLogger().Reload()
			break loop

		case s := <-signalTermIntCh:
			svr.GetLogger().Debug("Got signal: %v", s)
			signal.Stop(signalTermIntCh)
			endService()

		case <-sigUsr1Ch:
			fmt.Println("Catch sigusr1, reopen logfile")
			svr.GetLogger().Reload()
		}
	}
	return nil
}
