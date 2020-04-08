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

package weblib

import (
	"fmt"
	"strings"
)

func (ad AuthData) String() string {
	return fmt.Sprintf("AuthData[%s Name:%v Action:%v %v]",
		ad.Realm,
		len(ad.Name2Password),
		len(ad.ActionMap),
		len(ad.ActionNameMap),
	)
}

type AuthData struct {
	Realm         string
	Name2Password map[string]string
	ActionNameMap map[[2]string]bool
	ActionMap     map[string]bool
}

func NewAuthData(realm string) *AuthData {
	return &AuthData{
		Realm: strings.TrimSpace(realm),
	}
}

func (ad *AuthData) GetBasicAuthRealmString() string {
	return fmt.Sprintf("Basic realm=\"%s\"",
		ad.Realm)
}

// name password list
func (ad *AuthData) ReLoadUserData(data [][2]string) error {
	ad.Name2Password = make(map[string]string)
	for _, v := range data {
		name := strings.TrimSpace(v[0])
		pass := strings.TrimSpace(v[1])
		if _, exist := ad.Name2Password[name]; exist {
			return fmt.Errorf("Name exist %v", name)
		}
		ad.Name2Password[name] = pass
	}
	return nil
}

func (ad *AuthData) CheckLogin(name, password string) error {
	name = strings.TrimSpace(name)
	password = strings.TrimSpace(password)
	pass, exist := ad.Name2Password[name]
	if !exist {
		return fmt.Errorf("Name not exist %v", name)
	}
	if pass != password {
		return fmt.Errorf("Password not match %v", name)
	}
	return nil
}

// action name list
func (ad *AuthData) ReLoadActionName(data [][2]string) error {
	ad.ActionNameMap = make(map[[2]string]bool)
	ad.ActionMap = make(map[string]bool)
	for _, v := range data {
		action := strings.TrimSpace(v[0])
		name := strings.TrimSpace(v[1])
		key := [2]string{action, name}
		if _, exist := ad.Name2Password[name]; !exist {
			return fmt.Errorf("Name not exist %v", name)
		}
		if _, exist := ad.ActionNameMap[key]; exist {
			return fmt.Errorf("Action,Name Exist %v", key)
		}
		ad.ActionNameMap[key] = true
		ad.ActionMap[action] = true
	}
	return nil
}

func (ad *AuthData) CheckAction(action, name string) error {
	if _, exist := ad.Name2Password[name]; !exist {
		return fmt.Errorf("Name not exist %v", name)
	}
	val, exist := ad.ActionNameMap[[2]string{action, name}]
	if !exist {
		return fmt.Errorf("Invalid Action,Name %v %v", action, name)
	}
	if !val {
		return fmt.Errorf("Action,Name disabled %v %v", action, name)
	}
	return nil
}

func (ad *AuthData) AddAction(action string) error {
	if ad.ActionMap == nil {
		ad.ActionMap = make(map[string]bool)
	}
	action = strings.TrimSpace(action)
	_, exist := ad.ActionMap[action]
	if exist {
		return fmt.Errorf("Action exist %v", action)
	}
	ad.ActionMap[action] = true
	return nil
}

func (ad *AuthData) AddAllActionName(name string) error {
	if ad.ActionNameMap == nil {
		ad.ActionNameMap = make(map[[2]string]bool)
	}
	name = strings.TrimSpace(name)
	if _, exist := ad.Name2Password[name]; !exist {
		return fmt.Errorf("Name not exist %v", name)
	}
	for action, v := range ad.ActionMap {
		if !v {
			// disabled action
		}
		key := [2]string{action, name}
		ad.ActionNameMap[key] = true
	}
	return nil
}
