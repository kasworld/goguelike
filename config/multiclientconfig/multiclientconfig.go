// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multiclientconfig

import (
	"fmt"
	"path/filepath"

	"github.com/kasworld/goguelike/lib/g2log"
)

type MultiClientConfig struct {
	BaseLogDir string `default:"/tmp/"  argname:""`

	LogLevel      g2log.LL_Type `argname:""`
	SplitLogLevel g2log.LL_Type `argname:""`

	Net               string `default:"web" argname:""`
	ConnectToTower    string `default:"localhost:14101" argname:""`
	ListenWebInfoPort string `default:":14011" argname:""`
	PlayerNameBase    string `default:"MC_" argname:""`

	Concurrent        int  `default:"1000" argname:""`
	AccountPool       int  `default:"0" argname:""`
	AccountOverlap    int  `default:"0" argname:""`
	LimitStartCount   int  `default:"0" argname:""`
	LimitEndCount     int  `default:"0" argname:""`
	RetryDelayTimeOut int  `default:"-1" argname:""`
	DisconnectOnDeath bool `default:"false" argname:""`
}

func (config *MultiClientConfig) MakeLogDir() string {
	return filepath.Join(config.BaseLogDir,
		fmt.Sprintf("goguelike_multiclient_%v.logfiles", config.PlayerNameBase),
	)
}
