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

package groundconfig

import (
	"fmt"
	"path/filepath"

	"github.com/kasworld/argdefault"
	"github.com/kasworld/goguelike/config/towerconfig"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/prettystring"
)

type GroundConfig struct {
	LogLevel      g2log.LL_Type `default:"7" argname:""`
	SplitLogLevel g2log.LL_Type `default:"0" argname:""`
	BaseLogDir    string        `default:"/tmp/"  argname:""`

	DataFolder       string `default:"./serverdata" argname:""`
	ExeFolder        string `default:"./" argname:""`
	ClientDataFolder string `default:"./clientdata" argname:""`

	GroundHost           string `default:"localhost" argname:""`
	GroundAdminWebPort   int    `default:"14001" argname:""`
	GroundRPCPort        int    `default:"14002" argname:""`
	GroundServiceWebPort int    `default:"14003" argname:""`

	// TowerServicePortBase + towernum
	TowerServicePortBase int `default:"14100" argname:""`
	// TowerAdminWebPortBase + towernum
	TowerAdminWebPortBase int `default:"14200" argname:""`

	WebAdminID   string `default:"root" argname:""`
	WebAdminPass string `default:"password" argname:"" prettystring:"hidevalue"`
	AdminAuthKey string `default:"6e9456cf-ab29-99b2-f223-1459e00cfcd5" argname:"" prettystring:"hidevalue"`

	TowerDataFile        string `default:"towerdata.json" argname:""`
	HighScoreFile        string `default:"highscore.json" argname:""`
	TowerBin             string `default:"towerserver" argname:""`
	TowerAdminHostBase   string `default:"http://localhost" argname:""`
	TowerServiceHostBase string `default:"http://localhost" argname:""`
}

func (config *GroundConfig) MakeTowerConfig(
	portIndex int,
	displayName string,
	towerFilename string,
	turnPerSec float64,
) *towerconfig.TowerConfig {

	ads := argdefault.New(&towerconfig.TowerConfig{})
	tconfig := towerconfig.TowerConfig{}
	ads.SetDefaultToNonZeroField(&tconfig)

	tconfig.TowerName = displayName
	tconfig.ScriptFilename = towerFilename
	tconfig.TurnPerSec = turnPerSec

	tconfig.LogLevel = config.LogLevel
	tconfig.SplitLogLevel = config.SplitLogLevel
	tconfig.BaseLogDir = config.BaseLogDir

	tconfig.DataFolder = config.DataFolder
	tconfig.ClientDataFolder = config.ClientDataFolder

	tconfig.GroundRPC = fmt.Sprintf("%s:%d",
		config.GroundHost,
		config.GroundRPCPort)

	tconfig.ServicePort = config.TowerServicePortBase + portIndex
	tconfig.AdminPort = config.TowerAdminWebPortBase + portIndex

	tconfig.StandAlone = false
	tconfig.WebAdminID = config.WebAdminID
	tconfig.WebAdminPass = config.WebAdminPass
	tconfig.AdminAuthKey = config.AdminAuthKey
	return &tconfig
}

func (config *GroundConfig) MakeLogDir() string {
	rstr := filepath.Join(config.BaseLogDir,
		fmt.Sprintf("goguelike_groundserver.logfiles"),
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *GroundConfig) MakePIDFileFullpath() string {
	rstr := filepath.Join(config.BaseLogDir,
		fmt.Sprintf("goguelike_groundserver.pid"),
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *GroundConfig) MakeHighScoreFileFullpath() string {
	rstr := filepath.Join(config.ClientDataFolder,
		config.HighScoreFile,
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *GroundConfig) MakeTowerDataFileFullpath() string {
	rstr := filepath.Join(config.DataFolder,
		config.TowerDataFile,
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *GroundConfig) MakeTowerExeFullpath() string {
	rstr := filepath.Join(config.ExeFolder,
		config.TowerBin,
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *GroundConfig) MakeOutfileFullpath() string {
	rstr := fmt.Sprintf("goguelike_groundserver.out")
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config *GroundConfig) StringForm() string {
	return prettystring.PrettyString(config, 4)
}
