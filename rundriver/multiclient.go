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

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/kasworld/argdefault"
	"github.com/kasworld/configutil"
	profile "github.com/kasworld/go-profile"
	"github.com/kasworld/goguelike/config/multiclientconfig"
	"github.com/kasworld/goguelike/game/clientai"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/multirun"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/rangestat"
)

func main() {
	profile.AddArgs()
	configurl := flag.String("i", "", "client config file or url")
	ads := argdefault.New(&multiclientconfig.MultiClientConfig{})
	ads.RegisterFlag()
	flag.Parse()
	config := &multiclientconfig.MultiClientConfig{
		LogLevel:      g2log.LL_Fatal | g2log.LL_Error | g2log.LL_Warn | g2log.LL_TraceService,
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
	if profile.IsCpu() {
		fn := profile.StartCPUProfile()
		defer fn()
	}
	twlog, err := g2log.NewWithDstDir(
		config.PlayerNameBase,
		config.MakeLogDir(),
		logflags.DefaultValue(false).BitClear(logflags.LF_functionname),
		config.LogLevel,
		config.SplitLogLevel,
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	g2log.GlobalLogger = twlog
	// mc := NewMultiClient(*config, g2log.GlobalLogger)
	// mc.Run()

	chErr := make(chan error)
	go func() {
		for err := range chErr {
			fmt.Printf("%v\n", err)
		}
	}()
	multirun.Run(
		context.Background(),
		config.Concurrent,
		config.AccountPool,
		config.AccountOverlap,
		config.LimitStartCount,
		config.LimitEndCount,
		func(config interface{}) multirun.ClientI {
			return clientai.New(
				config.(clientai.ClientAIConfig),
				g2log.GlobalLogger,
			)
		},
		func(i int) interface{} {
			return clientai.ClientAIConfig{
				ConnectToTower:    config.ConnectToTower,
				Nickname:          fmt.Sprintf("%s%d", config.PlayerNameBase, i),
				SessionG2ID:       "",
				DisconnectOnDeath: config.DisconnectOnDeath,
				Auth:              "6e9456cf-ab29-99b2-f223-1459e00cfcd5",
			}
		},
		chErr,
		rangestat.New("", 0, config.Concurrent),
	)

	if profile.IsMem() {
		profile.WriteHeapProfile()
	}
}
