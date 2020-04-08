// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package argdefault

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
)

type ArgStatue struct {
	obj           interface{}
	intDefault    map[string]int
	uintDefault   map[string]uint
	floatDefault  map[string]float64
	boolDefault   map[string]bool
	stringDefault map[string]string

	intArg    map[string]*int
	uintArg   map[string]*uint
	floatArg  map[string]*float64
	boolArg   map[string]*bool
	stringArg map[string]*string
}

// New prepare default and arg flag
func New(inObj interface{}) *ArgStatue {
	as := &ArgStatue{
		obj:           inObj,
		intDefault:    make(map[string]int),
		uintDefault:   make(map[string]uint),
		floatDefault:  make(map[string]float64),
		boolDefault:   make(map[string]bool),
		stringDefault: make(map[string]string),

		intArg:    make(map[string]*int),
		uintArg:   make(map[string]*uint),
		floatArg:  make(map[string]*float64),
		boolArg:   make(map[string]*bool),
		stringArg: make(map[string]*string),
	}
	as.processDefaultTag()
	return as
}

func (as *ArgStatue) processDefaultTag() {
	objValue := reflect.ValueOf(as.obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		structField := objValue.Field(i)
		fieldName := objValue.Type().Field(i).Name
		fieldTag := objValue.Type().Field(i).Tag
		defaultVal, defaultExist := fieldTag.Lookup("default")
		if defaultExist {
			switch structField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, err := strconv.Atoi(defaultVal)
				if err != nil {
					fmt.Printf("invalid int field %v %v\n", defaultVal, err)
				} else {
					as.intDefault[fieldName] = v
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v, err := strconv.Atoi(defaultVal)
				if err != nil {
					fmt.Printf("invalid int field %v %v\n", defaultVal, err)
				} else {
					as.uintDefault[fieldName] = uint(v)
				}
			case reflect.Float64, reflect.Float32:
				v, err := strconv.ParseFloat(defaultVal, 64)
				if err != nil {
					fmt.Printf("invalid int field %v %v\n", defaultVal, err)
				} else {
					as.floatDefault[fieldName] = v
				}
			case reflect.Bool:
				v, err := strconv.ParseBool(defaultVal)
				if err != nil {
					fmt.Printf("invalid int field %v %v\n", defaultVal, err)
				} else {
					as.boolDefault[fieldName] = v
				}
			case reflect.String:
				as.stringDefault[fieldName] = defaultVal
			default:
				fmt.Printf("unprocessed %v\n", structField)
			}
		}
	}
}

// SetDefaultToNonZeroField set inObj default value if field is zero value
func (as *ArgStatue) SetDefaultToNonZeroField(inObj interface{}) {
	destObjValue := reflect.ValueOf(inObj).Elem()
	for i := 0; i < destObjValue.NumField(); i++ {
		fieldName := destObjValue.Type().Field(i).Name
		fieldTag := destObjValue.Type().Field(i).Tag
		destObjField := destObjValue.Field(i)
		// skip non zero field
		if !destObjField.IsZero() {
			continue
		}
		_, defexist := fieldTag.Lookup("default")
		if defexist {
			switch destObjField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				def, _ := as.intDefault[fieldName]
				destObjField.SetInt(int64(def))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				def, _ := as.uintDefault[fieldName]
				destObjField.SetUint(uint64(def))
			case reflect.Float64, reflect.Float32:
				def, _ := as.floatDefault[fieldName]
				destObjField.SetFloat(def)
			case reflect.Bool:
				def, _ := as.boolDefault[fieldName]
				destObjField.SetBool(def)
			case reflect.String:
				def, _ := as.stringDefault[fieldName]
				destObjField.SetString(def)
			default:
				fmt.Printf("unprocessed %v\n", fieldName)
			}
		}
	}
}

// RegisterFlag set flag with default
func (as *ArgStatue) RegisterFlag() {
	objValue := reflect.ValueOf(as.obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		fieldName := objValue.Type().Field(i).Name
		structField := objValue.Field(i)
		fieldTag := objValue.Type().Field(i).Tag
		argname, argexist := fieldTag.Lookup("argname")
		if argexist {
			if argname == "" {
				argname = fieldName
			}
			switch structField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var argv int
				as.intArg[fieldName] = &argv
				if d, exist := as.intDefault[fieldName]; exist {
					flag.IntVar(&argv, argname, d, fieldName)
				} else {
					flag.IntVar(&argv, argname, 0, fieldName)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				var argv uint
				as.uintArg[fieldName] = &argv
				if d, exist := as.uintDefault[fieldName]; exist {
					flag.UintVar(&argv, argname, d, fieldName)
				} else {
					flag.UintVar(&argv, argname, 0, fieldName)
				}
			case reflect.Float64, reflect.Float32:
				var argv float64
				as.floatArg[fieldName] = &argv
				if d, exist := as.floatDefault[fieldName]; exist {
					flag.Float64Var(&argv, argname, d, fieldName)
				} else {
					flag.Float64Var(&argv, argname, 0, fieldName)
				}
			case reflect.Bool:
				var argv bool
				as.boolArg[fieldName] = &argv
				if d, exist := as.boolDefault[fieldName]; exist {
					flag.BoolVar(&argv, argname, d, fieldName)
				} else {
					flag.BoolVar(&argv, argname, false, fieldName)
				}
			case reflect.String:
				var argv string
				as.stringArg[fieldName] = &argv
				if d, exist := as.stringDefault[fieldName]; exist {
					flag.StringVar(&argv, argname, d, fieldName)
				} else {
					flag.StringVar(&argv, argname, "", fieldName)
				}
			default:
				fmt.Printf("unprocessed %v\n", structField)
			}
		}
	}
}

// ApplyFlagTo set inobj from flag and return
func (as *ArgStatue) ApplyFlagTo(inObj interface{}) interface{} {
	destObjValue := reflect.ValueOf(inObj).Elem()

	for i := 0; i < destObjValue.NumField(); i++ {
		fieldName := destObjValue.Type().Field(i).Name
		fieldTag := destObjValue.Type().Field(i).Tag
		destObjField := destObjValue.Field(i)
		_, argexist := fieldTag.Lookup("argname")
		if argexist {
			switch destObjField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				arg, existArg := as.intArg[fieldName]
				def, _ := as.intDefault[fieldName]
				if !existArg {
					continue
				}
				if *arg != def {
					// fmt.Printf("replace %v arg[%v] default[%v]\n", fieldName, *arg, def)
					destObjField.SetInt(int64(*arg))
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				arg, existArg := as.uintArg[fieldName]
				def, _ := as.uintDefault[fieldName]
				if !existArg {
					continue
				}
				if *arg != def {
					// fmt.Printf("replace %v arg[%v] default[%v]\n", fieldName, *arg, def)
					destObjField.SetUint(uint64(*arg))
				}
			case reflect.Float64, reflect.Float32:
				arg, existArg := as.floatArg[fieldName]
				def, _ := as.floatDefault[fieldName]
				if !existArg {
					continue
				}
				if *arg != def {
					// fmt.Printf("replace %v arg[%v] default[%v]\n", fieldName, *arg, def)
					destObjField.SetFloat(*arg)
				}
			case reflect.Bool:
				arg, existArg := as.boolArg[fieldName]
				def, _ := as.boolDefault[fieldName]
				if !existArg {
					continue
				}
				if *arg != def {
					// fmt.Printf("replace %v arg[%v] default[%v]\n", fieldName, *arg, def)
					destObjField.SetBool(*arg)
				}
			case reflect.String:
				arg, existArg := as.stringArg[fieldName]
				def, _ := as.stringDefault[fieldName]
				if !existArg {
					continue
				}
				if *arg != def {
					// fmt.Printf("replace %v arg[%v] default[%v]\n", fieldName, *arg, def)
					destObjField.SetString(*arg)
				}
			default:
				fmt.Printf("unprocessed %v\n", destObjField)
			}
		}
	}
	return inObj
}
