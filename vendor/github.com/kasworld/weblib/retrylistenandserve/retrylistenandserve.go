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

package retrylistenandserve

import (
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

type loggerI interface {
	Warn(format string, v ...interface{})
	TraceService(format string, v ...interface{})
}

func RetryListenAndServe(httpServer *http.Server, l loggerI, name string) {
	l.TraceService("Start RetryListenAndServe %v", name)
	for {
		err := httpServer.ListenAndServe()
		if err != nil {
			operr, ok := err.(*net.OpError)
			if ok {
				syscallerr, ok := operr.Err.(*os.SyscallError)
				if ok && syscallerr.Err == syscall.EADDRINUSE {
					l.Warn("Retry %v %v", name, err) // bind: address already in use
					time.Sleep(time.Second)
					continue
				}
			}
			l.TraceService("End RetryListenAndServe %v %v", name, err)
			break
		}
	}
}
