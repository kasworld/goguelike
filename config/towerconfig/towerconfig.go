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

package towerconfig

import (
	"fmt"
	"path/filepath"

	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/prettystring"
)

type TowerConfig struct {
	// common to all tower
	LogLevel         g2log.LL_Type `default:"7" argname:""`
	SplitLogLevel    g2log.LL_Type `default:"0" argname:""`
	BaseLogDir       string        `default:"/tmp/"  argname:""`
	DataFolder       string        `default:"./serverdata" argname:""`
	ClientDataFolder string        `default:"./clientdata" argname:""`
	GroundRPC        string        `default:"localhost:14002"  argname:""`
	WebAdminID       string        `default:"root" argname:""`
	WebAdminPass     string        `default:"password" argname:"" prettystring:"hidevalue"`
	AdminAuthKey     string        `default:"6e9456cf-ab29-99b2-f223-1459e00cfcd5" argname:"" prettystring:"hidevalue"`

	// config for each tower
	ServicePort           int     `default:"14101"  argname:""`
	AdminPort             int     `default:"14201"  argname:""`
	TowerFilename         string  `default:"start" argname:""`
	TowerNumber           int     `default:"1" argname:""`
	DisplayName           string  `default:"Default" argname:""`
	ConcurrentConnections int     `default:"10000" argname:""`
	TurnPerSec            float64 `default:"2.0" argname:""`
	StandAlone            bool    `default:"true" argname:""`
	ServiceHostBase       string  `default:"http://localhost" argname:""` // for StandAlone mode
}

func (config *TowerConfig) MakeLogDir() string {
	rstr := filepath.Join(config.BaseLogDir,
		fmt.Sprintf("goguelike_towerserver_%v.logfiles",
			config.TowerNumber),
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *TowerConfig) MakePIDFileFullpath() string {
	rstr := filepath.Join(config.BaseLogDir,
		fmt.Sprintf("goguelike_towerserver_%v.pid",
			config.TowerNumber),
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}
func (config *TowerConfig) MakeTowerFileFullpath() string {
	rstr := filepath.Join(config.DataFolder,
		fmt.Sprintf("%v.tower", config.TowerFilename),
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *TowerConfig) MakeOutfileFullpath() string {
	rstr := fmt.Sprintf("goguelike_towerserver_%v.out",
		config.TowerNumber)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *TowerConfig) StringForm() string {
	return prettystring.PrettyString(config, 4)
}
