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

// StrList2OrderedMap make string(name sep value) list to namelist and map
func StrList2OrderedMap(strList []string, sep string) ([]string, map[string]string, error) {
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

// Split2ListMap make list , map
func Split2ListMap(srcText string, sep1, sep2 string) ([]string, map[string]string, error) {
	list := SplitTrim(srcText, sep1)
	nameList, name2value, err := StrList2OrderedMap(list, sep2)
	return nameList, name2value, err
}
