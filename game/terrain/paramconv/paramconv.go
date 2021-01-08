// Copyright 2014,2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package paramconv

import (
	"fmt"
	"strconv"

	"github.com/kasworld/goguelike/lib/scriptparse"

	"github.com/kasworld/goguelike/enum/decaytype"
	"github.com/kasworld/goguelike/enum/fieldobjacttype"
	"github.com/kasworld/goguelike/enum/fieldobjdisplaytype"
	"github.com/kasworld/goguelike/enum/resourcetype"
	"github.com/kasworld/goguelike/enum/tile"
	"github.com/kasworld/goguelike/enum/tile_flag"
)

func SetFloat(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*float64)
	if !ok {
		return fmt.Errorf("fail to cast float64 %v", valStr)
	}
	r, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fmt.Errorf("ParseFloat fail %v , %v", valStr, err)
	}
	*iv = r
	return nil
}

func SetInt(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*int)
	if !ok {
		return fmt.Errorf("fail to cast int %v", valStr)
	}
	r, err := strconv.Atoi(valStr)
	if err != nil {
		return fmt.Errorf("atoi fail %v , %v", valStr, err)
	}
	*iv = r
	return nil
}

func SetBool(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*bool)
	if !ok {
		return fmt.Errorf("fail to cast bool %v", valStr)
	}
	r, err := strconv.ParseBool(valStr)
	if err != nil {
		return fmt.Errorf("ParseBool fail %v , %v", valStr, err)
	}
	*iv = r
	return nil
}

func SetString(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*string)
	if !ok {
		return fmt.Errorf("fail to cast string %v", valStr)
	}
	*iv = valStr
	return nil
}

func SetFieldObjActType(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*fieldobjacttype.FieldObjActType)
	if !ok {
		return fmt.Errorf("fail to cast FieldObjActType %v", valStr)
	}
	portalact, exist := fieldobjacttype.String2FieldObjActType(valStr)
	if !exist {
		return fmt.Errorf("unknown FieldObjActType %v", valStr)
	}
	*iv = portalact
	return nil
}

func SetFieldObjDisplayType(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*fieldobjdisplaytype.FieldObjDisplayType)
	if !ok {
		return fmt.Errorf("fail to cast FieldObjDisplayType %v", valStr)
	}
	portaldisp, exist := fieldobjdisplaytype.String2FieldObjDisplayType(valStr)
	if !exist {
		return fmt.Errorf("unknown FieldObjDisplayType %v", valStr)
	}
	*iv = portaldisp
	return nil
}

func SetTileType(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*tile.Tile)
	if !ok {
		return fmt.Errorf("fail to cast TileType %v", valStr)
	}
	tl, exist := tile.String2Tile(valStr)
	if !exist {
		return fmt.Errorf("unknown TileType %v", valStr)
	}
	*iv = tl
	return nil
}

func SetTileFlag(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*tile_flag.TileFlag)
	if !ok {
		return fmt.Errorf("fail to cast TileType %v", valStr)
	}
	tls := scriptparse.SplitTrim(valStr, ",")
	var tlf tile_flag.TileFlag
	for _, v := range tls {
		tl, exist := tile.String2Tile(v)
		if !exist {
			return fmt.Errorf("unknown TileType %v", valStr)
		}
		tlf.SetByTile(tl)
	}
	*iv = tlf
	return nil
}

func SetResourceType(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*resourcetype.ResourceType)
	if !ok {
		return fmt.Errorf("fail to cast ResourceType %v", valStr)
	}
	rsctl, exist := resourcetype.String2ResourceType(valStr)
	if !exist {
		return fmt.Errorf("unknown ResourceType %v", valStr)
	}
	*iv = rsctl
	return nil
}

func SetDecayType(valStr string, dstValue interface{}) error {
	iv, ok := dstValue.(*decaytype.DecayType)
	if !ok {
		return fmt.Errorf("fail to cast DecayType %v", valStr)
	}
	rsctl, exist := decaytype.String2DecayType(valStr)
	if !exist {
		return fmt.Errorf("unknown DecayType %v", valStr)
	}
	*iv = rsctl
	return nil
}

var Type2ConvFn = map[string]func(valStr string, dstValue interface{}) error{
	"float":               SetFloat,
	"int":                 SetInt,
	"bool":                SetBool,
	"string":              SetString,
	"FieldObjActType":     SetFieldObjActType,
	"FieldObjDisplayType": SetFieldObjDisplayType,
	"TileType":            SetTileType,
	"TileFlag":            SetTileFlag,
	"ResourceType":        SetResourceType,
	"DecayType":           SetDecayType,
}
