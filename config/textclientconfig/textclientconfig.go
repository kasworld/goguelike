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

package textclientconfig

import (
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/prettystring"
)

type TextClientConfig struct {
	LogDir        string        `default:"/tmp/textclient.logfile" argname:""`
	LogLevel      g2log.LL_Type `argname:""`
	SplitLogLevel g2log.LL_Type `argname:""`
	PidFilename   string        `default:"/tmp/textclient.pid" argname:""`

	ConnectToTower string `default:"localhost:14101" argname:""`
	PlayerName     string `default:"Player" argname:""`
}

func (config *TextClientConfig) StringForm() string {
	return prettystring.PrettyString(config, 4)
}
