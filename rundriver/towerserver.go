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
	"flag"
	"fmt"

	"github.com/kasworld/argdefault"
	"github.com/kasworld/configutil"
	profile "github.com/kasworld/go-profile"
	"github.com/kasworld/goguelike/config/dataversion"
	"github.com/kasworld/goguelike/config/towerconfig"
	"github.com/kasworld/goguelike/game/tower"
	"github.com/kasworld/goguelike/lib/g2log"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_version"
	"github.com/kasworld/signalhandle"
	"github.com/kasworld/version"
)

var Ver = ""

func init() {
	version.Set(Ver)
}

func printVersion() {
	fmt.Println("Goguelike service tower")
	fmt.Println("Build     ", version.GetVersion())
	fmt.Println("Data      ", dataversion.DataVersion)
	fmt.Println("Protocol  ", c2t_version.ProtocolVersion)
	fmt.Println()
}

func main() {
	printVersion()

	configurl := flag.String("i", "", "server config file or url")
	signalhandle.AddArgs()
	profile.AddArgs()

	ads := argdefault.New(&towerconfig.TowerConfig{})
	ads.RegisterFlag()
	flag.Parse()
	config := &towerconfig.TowerConfig{}
	ads.SetDefaultToNonZeroField(config)
	if *configurl != "" {
		// fmt.Printf("load ini %v\n", *configurl)
		if err := configutil.LoadIni(*configurl, &config); err != nil {
			g2log.Fatal("%v", err)
		}
	}
	ads.ApplyFlagTo(config)
	// fmt.Println("after ApplyArgsTo", prettystring.PrettyString(config, 4))

	if profile.IsCpu() {
		fn := profile.StartCPUProfile()
		defer fn()
	}

	tw := tower.New(config)
	if err := signalhandle.StartByArgs(tw); err != nil {
		fmt.Printf("%v\n", err)
	}

	if profile.IsMem() {
		profile.WriteHeapProfile()
	}

}
