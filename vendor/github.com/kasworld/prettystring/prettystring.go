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

package prettystring

import (
	"bytes"
	"fmt"
	"reflect"
)

func writeIndent(buf *bytes.Buffer, n int) {
	buf.WriteString("\n")
	for i := 0; i < n; i++ {
		buf.WriteString("  ")
		// buf.WriteString("-|")
	}
}

// PrettyString make struct to detailed string like python prettyprint
func PrettyString(inObj interface{}, depthLimit int) string {
	var buf bytes.Buffer
	formObj(&buf, reflect.ValueOf(inObj), 0, depthLimit)
	buf.WriteString("\n")
	return buf.String()
}

func formObj(buf *bytes.Buffer, objVal reflect.Value, depth int, depthLimit int) {
	if depth > depthLimit {
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
		return
	}
	switch objVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
	case reflect.Float64, reflect.Float32:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
	case reflect.Bool:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
	case reflect.String:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
	case reflect.Chan:
		fmt.Fprintf(buf, "%v %v/%v", objVal.Type(), objVal.Len(), objVal.Cap())
	case reflect.Func:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
	case reflect.Interface:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
		if !objVal.IsNil() {
			writeIndent(buf, depth+1)
			formObj(buf, objVal.Elem(), depth+1, depthLimit)
		}
	case reflect.Ptr:
		fmt.Fprintf(buf, "%v", objVal.Type())
		// fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
		if !objVal.IsNil() {
			writeIndent(buf, depth+1)
			formObj(buf, reflect.Indirect(objVal), depth+1, depthLimit)
		}
	case reflect.UnsafePointer:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
	case reflect.Array:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
		for i := 0; i < objVal.Len(); i++ {
			writeIndent(buf, depth+1)
			formObj(buf, objVal.Index(i), depth+1, depthLimit)
		}
	case reflect.Map:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
		iter := objVal.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			writeIndent(buf, depth+1)
			fmt.Fprintf(buf, "%v: ", k)
			formObj(buf, v, depth+1, depthLimit)
		}
	case reflect.Slice:
		fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
		for i := 0; i < objVal.Len(); i++ {
			writeIndent(buf, depth+1)
			formObj(buf, objVal.Index(i), depth+1, depthLimit)
		}
	case reflect.Struct:
		fmt.Fprintf(buf, "%v", objVal.Type())
		// fmt.Fprintf(buf, "%v %v", objVal.Type(), objVal)
		for i := 0; i < objVal.NumField(); i++ {
			structTypeField := objVal.Type().Field(i)
			if tag, exist := structTypeField.Tag.Lookup("prettystring"); exist {
				switch tag {
				default:
					writeIndent(buf, depth+1)
					fmt.Fprintf(buf, "%v %v %v unknowntag %v", structTypeField.Name, structTypeField.Type, objVal.Field(i), tag)
					continue
				case "hide":
					continue
				case "hidevalue":
					writeIndent(buf, depth+1)
					fmt.Fprintf(buf, "%v %v ****", structTypeField.Name, structTypeField.Type)
					continue
				case "simple":
					writeIndent(buf, depth+1)
					fmt.Fprintf(buf, "%v %v %v", structTypeField.Name, structTypeField.Type, objVal.Field(i))
					continue
				}
			}
			writeIndent(buf, depth+1)
			fmt.Fprintf(buf, "%v: ", structTypeField.Name)
			formObj(buf, objVal.Field(i), depth+1, depthLimit)
		}
	default:
		fmt.Fprintf(buf, "%v %v unprocessed", objVal.Type(), objVal)
	}
}
