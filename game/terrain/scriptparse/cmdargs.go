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
	"strings"
)

func (ca CmdArgs) String() string {
	return fmt.Sprintf("CmdArgs[%v %v]", ca.Cmd, ca.Name2Value)
}

type CmdArgs struct {
	Cmd        string
	Name2Value map[string]string
	NameList   []string
	Name2Type  map[string]string
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
	convFn, exist := Type2ConvFn[ca.Name2Type[argName]]
	if !exist {
		return fmt.Errorf("not supported type %v %v", ca.Name2Type[argName], argValue)
	}
	err := convFn(argValue, v)
	if err != nil {
		return err
	}
	return nil
}

func SplitCmdArgstr(cmdline string, sep string) (string, string) {
	cmdline = strings.TrimSpace(cmdline)
	splited := strings.SplitN(cmdline, sep, 2)
	switch len(splited) {
	case 1:
		return strings.TrimSpace(splited[0]), ""
	case 2:
		return strings.TrimSpace(splited[0]), strings.TrimSpace(splited[1])
	default:
		return "", ""
	}
}

func SetArgsByFormat(
	argStr string, argSep1, argSep2 string,
	format string, fmtSep1, fmtSep2 string,
	varList ...interface{}) error {

	_, name2Value, err := Split2ListMap(argStr, argSep1, argSep2)
	if err != nil {
		return err
	}

	nameList2, name2Type, err := Split2ListMap(format, fmtSep1, fmtSep2)
	if err != nil {
		return err
	}

	for argPos, argName := range nameList2 {
		argValue, exist := name2Value[argName]
		if !exist {
			return fmt.Errorf("arg %v not found %v", argName, name2Value)
		}
		convFn, exist := Type2ConvFn[name2Type[argName]]
		if !exist {
			return fmt.Errorf("not supported type %v %v", name2Type[argName], argValue)
		}
		err := convFn(argValue, varList[argPos])
		if err != nil {
			return err
		}
	}
	return nil
}
