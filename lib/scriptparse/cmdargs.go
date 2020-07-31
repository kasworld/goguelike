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

package scriptparse

import (
	"fmt"
)

func (ca CmdArgs) String() string {
	return fmt.Sprintf("CmdArgs[%v %v]", ca.Cmd, ca.Name2Value)
}

type CmdArgs struct {
	Type2ConvFn map[string]func(valStr string, dstValue interface{}) error
	Cmd         string
	Name2Value  map[string]string
	NameList    []string
	Name2Type   map[string]string
}

func (ca *CmdArgs) GetArgs(varList ...interface{}) error {
	if len(ca.Name2Value) != len(ca.NameList) {
		return fmt.Errorf("invalid arg count %v, %v", ca.Name2Value, ca.NameList)
	}
	for argPos, argName := range ca.NameList {
		err := ca.SetArgByName(argName, varList[argPos])
		if err != nil {
			return err
		}
	}
	return nil
}

func (ca *CmdArgs) SetArgByName(argName string, v interface{}) error {
	argValue, exist := ca.Name2Value[argName]
	if !exist {
		return fmt.Errorf("arg %v not found %v", argName, ca.Name2Value)
	}
	convFn, exist := ca.Type2ConvFn[ca.Name2Type[argName]]
	if !exist {
		return fmt.Errorf("not supported type %v %v", ca.Name2Type[argName], argValue)
	}
	err := convFn(argValue, v)
	if err != nil {
		return err
	}
	return nil
}
