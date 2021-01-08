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

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/kasworld/argdefault"
	"github.com/kasworld/configutil"
	"github.com/kasworld/goguelike/config/textclientconfig"
	"github.com/kasworld/goguelike/game/clientai"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/version"
)

var Ver = ""

func init() {
	version.Set(Ver)
}

func main() {
	configurl := flag.String("i", "", "client config file or url")

	ads := argdefault.New(&textclientconfig.TextClientConfig{})
	ads.RegisterFlag()
	flag.Parse()
	config := &textclientconfig.TextClientConfig{
		LogLevel:      g2log.LL_All,
		SplitLogLevel: 0,
	}
	ads.SetDefaultToNonZeroField(config)
	if *configurl != "" {
		if err := configutil.LoadIni(*configurl, config); err != nil {
			g2log.Error("%v", err)
		}
	}
	ads.ApplyFlagTo(config)
	fmt.Println(prettystring.PrettyString(config, 4))

	g2log.GlobalLogger.SetFlags(
		g2log.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
	g2log.GlobalLogger.SetLevel(
		config.LogLevel)

	aiconfig := clientai.ClientAIConfig{
		ConnectToTower:    config.ConnectToTower,
		Nickname:          config.PlayerName,
		SessionUUID:       "",
		DisconnectOnDeath: false,
		Auth:              "6e9456cf-ab29-99b2-f223-1459e00cfcd5",
	}
	app := clientai.New(
		aiconfig,
		g2log.GlobalLogger,
	)
	app.Run(context.Background())
	g2log.Error("%v", app.GetRunResult())
}
