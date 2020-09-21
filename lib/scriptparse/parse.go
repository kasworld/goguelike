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

// SplitCmdArgstr split smd and argstr from cmdline
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

// Split2ListMap make list , map
func Split2ListMap(srcText string, sep1, sep2 string) ([]string, map[string]string, error) {
	list := SplitTrim(srcText, sep1)
	nameList, name2value, err := strList2OrderedMap(list, sep2)
	return nameList, name2value, err
}

// SplitTrim split string with trimspace
func SplitTrim(str string, sep string) []string {
	str = strings.TrimSpace(str)
	strs := strings.Split(str, sep)
	var rtn []string
	for _, v := range strs {
		s := strings.TrimSpace(v)
		if len(s) > 0 {
			rtn = append(rtn, s)
		}
	}
	return rtn
}

// strList2OrderedMap make string(name sep value) list to namelist and map
func strList2OrderedMap(strList []string, sep string) ([]string, map[string]string, error) {
	nameList := make([]string, 0, len(strList))
	name2value := make(map[string]string)
	for _, v := range strList {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		nameValue := strings.SplitN(v, sep, 2)
		if len(nameValue) == 1 {
			name := strings.TrimSpace(nameValue[0])
			nameList = append(nameList, name)
			name2value[name] = ""
		} else if len(nameValue) == 2 {
			name := strings.TrimSpace(nameValue[0])
			nameList = append(nameList, name)
			name2value[name] = strings.TrimSpace(nameValue[1])
		} else { // 0 or more then 3
			panic(fmt.Sprintf("invalid splitN(%v,%v,2)", v, sep))
		}
	}
	return nameList, name2value, nil
}

// func SetArgsByFormat(
// 	argStr string, argSep1, argSep2 string,
// 	format string, fmtSep1, fmtSep2 string,
// 	varList ...interface{}) error {

// 	_, name2Value, err := Split2ListMap(argStr, argSep1, argSep2)
// 	if err != nil {
// 		return err
// 	}

// 	nameList2, name2Type, err := Split2ListMap(format, fmtSep1, fmtSep2)
// 	if err != nil {
// 		return err
// 	}

// 	for argPos, argName := range nameList2 {
// 		argValue, exist := name2Value[argName]
// 		if !exist {
// 			return fmt.Errorf("arg %v not found %v", argName, name2Value)
// 		}
// 		convFn, exist := Type2ConvFn[name2Type[argName]]
// 		if !exist {
// 			return fmt.Errorf("not supported type %v %v", name2Type[argName], argValue)
// 		}
// 		err := convFn(argValue, varList[argPos])
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
