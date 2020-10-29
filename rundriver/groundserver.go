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
	"github.com/kasworld/goguelike/config/dataversion"
	"github.com/kasworld/goguelike/config/groundconfig"
	"github.com/kasworld/goguelike/game/ground"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_version"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/signalhandle"
	"github.com/kasworld/version"
)

var Ver = ""

func init() {
	version.Set(Ver)
}

func printVersion() {
	fmt.Println("Goguelike ground server")
	fmt.Println("Build     ", version.GetVersion())
	fmt.Println("Data      ", dataversion.DataVersion)
	fmt.Println("Protocol  ", c2t_version.ProtocolVersion)
	fmt.Println()
}

func main() {
	printVersion()

	configurl := flag.String("i", "", "server config file or url")
	signalhandle.AddArgs()

	ads := argdefault.New(&groundconfig.GroundConfig{})
	ads.RegisterFlag()
	flag.Parse()
	config := &groundconfig.GroundConfig{}
	ads.SetDefaultToNonZeroField(config)
	if *configurl != "" {
		if err := configutil.LoadIni(*configurl, config); err != nil {
			fmt.Printf("%v\n", err)
		}
	}
	ads.ApplyFlagTo(config)
	fmt.Println(prettystring.PrettyString(config, 4))

	tw := ground.New(config)
	if err := signalhandle.StartByArgs(tw); err != nil {
		fmt.Printf("%v\n", err)
	}
}
