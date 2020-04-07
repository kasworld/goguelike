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

package soundmap

import (
	"syscall/js"

	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/gowasmlib/jslog"
)

var SoundOn bool = true

func Play(name string) {
	if !SoundOn {
		return
	}
	snd := js.Global().Get("document").Call("getElementById", name)
	if snd.Truthy() {
		snd.Call("play")
	} else {
		jslog.Errorf("no sound %v", name)
	}
}

var act2name = map[c2t_idcmd.CommandID]string{
	// c2t_idcmd.KillSelf:    "",
	c2t_idcmd.Move:        "stepsound",
	c2t_idcmd.Attack:      "attacksound",
	c2t_idcmd.Pickup:      "pickupsound",
	c2t_idcmd.Drop:        "dropsound",
	c2t_idcmd.Equip:       "equipsound",
	c2t_idcmd.UnEquip:     "unequipsound",
	c2t_idcmd.DrinkPotion: "usesound",
	c2t_idcmd.ReadScroll:  "usesound",
	c2t_idcmd.Recycle:     "recyclesound",
	// c2t_idcmd.EnterPortal: "",
}

func PlayByAct(act c2t_idcmd.CommandID) {
	if !SoundOn {
		return
	}
	name := act2name[act]
	if name != "" {
		Play(name)
	}
}
