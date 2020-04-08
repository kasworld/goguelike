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

package multirun

import (
	"context"
	"time"

	"github.com/kasworld/rangestat"
)

type ClientI interface {
	Run(mainctx context.Context)
	GetArg() interface{}
	GetRunResult() error
}

func Run(
	ctx context.Context,
	Concurrent int,
	AccountPool int,
	AccountOverlap int,
	LimitStartCount int,
	LimitEndCount int,
	NewClientFn func(config interface{}) ClientI,
	NewClientArgFn func(i int) interface{},
	EndErrorCh chan<- error,
	RunStat *rangestat.RangeStat,
) {

	if AccountPool < Concurrent {
		AccountPool = Concurrent
	}

	var clientArgToRun chan interface{}
	waitStartCh := make(chan ClientI, Concurrent)
	endedCh := make(chan ClientI, Concurrent)

	ctx, endFn := context.WithCancel(ctx)
	defer endFn()

	// start worker
	for i := 0; i < Concurrent; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case ra, ok := <-waitStartCh:
					if !ok {
						return
					}
					if !RunStat.Inc() {
						panic("fail to inc")
					}
					ra.Run(ctx)
					if !RunStat.Dec() {
						panic("fail to dec")
					}
					endedCh <- ra
				}
			}
		}(ctx)
	}

	// make account pool to run
	if AccountOverlap == 1 {
		clientArgToRun = make(chan interface{}, AccountPool*2)
		for i := 0; i < AccountPool; i++ {
			clientArg := NewClientArgFn(i)
			clientArgToRun <- clientArg
			clientArgToRun <- clientArg
		}
	} else {
		clientArgToRun = make(chan interface{}, AccountPool)
		for i := 0; i < AccountPool; i++ {
			clientArg := NewClientArgFn(i)
			clientArgToRun <- clientArg
		}
	}

	// start new client to run
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case clinetArg := <-clientArgToRun:
				// check can start new
				if LimitStartCount == 0 || RunStat.GetTotalInc() < LimitStartCount {
					waitStartCh <- NewClientFn(clinetArg)
				}
				time.Sleep(time.Microsecond * 1)
			}
		}
	}()

	// handle ended client
	for {
		select {
		case <-ctx.Done():
			return

		case endedClient := <-endedCh:
			// check all end
			if LimitEndCount != 0 && RunStat.GetTotalDec() >= LimitEndCount {
				return
			}
			if err := endedClient.GetRunResult(); err != nil {
				EndErrorCh <- err
			}
			clientArgToRun <- endedClient.GetArg()
		}
	}
}
