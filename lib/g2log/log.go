// Copyright 2014,2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package g2log

import (
	"fmt"
	"os"
	"strings"

	"github.com/kasworld/log/logdestination_file"
	"github.com/kasworld/log/logflagi"
	"github.com/kasworld/log/logflags"
)

func NewWithDstDir(prefix string, logdir string, lf logflagi.LogFlagI,
	loglevel LL_Type, splitLogLevel LL_Type) (*LogBase, error) {
	logdir = strings.TrimSpace(logdir)
	if logdir == "" {
		return nil, fmt.Errorf("logdir empty")
	}
	if err := os.MkdirAll(logdir, os.ModePerm); err != nil {
		return nil, err
	}
	newDestForOther, err := logdestination_file.New(
		makeLogFilename(logdir, "Other"))
	if err != nil {
		return nil, err
	}
	newlg := New(prefix, lf, loglevel)
	newlg.AddDestination(LL_All^splitLogLevel, newDestForOther)
	for ll := loglevel.StartLevel(); !ll.IsLastLevel(); ll = ll.NextLevel(1) {
		if splitLogLevel.IsLevel(ll) {
			newDestForLL, serr := logdestination_file.New(
				makeLogFilename(logdir, ll.String()))
			if serr != nil {
				return nil, serr
			}
			newlg.AddDestination(ll, newDestForLL)
		}
	}
	newlg.AddDestination(LL_Fatal, OutputStderr)
	return newlg, nil
}

func init() {
	GlobalLogger = New("", logflags.DefaultValue(false), LL_All)
	GlobalLogger.AddDestination(LL_All, OutputStderr)
}

const (
	LL_ServiceDefault = LL_Fatal | LL_Error | LL_Warn | LL_TraceService | LL_TraceUser | LL_AdminAudit
	LL_TestDefault    = LL_Fatal | LL_Error | LL_Warn | LL_TraceService | LL_TraceUser | LL_Debug | LL_AdminAudit
)

//////////////////////////////////////////////////////////////////
