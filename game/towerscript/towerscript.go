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

package towerscript

import (
	"encoding/json"
	"fmt"
	"os"
)

// type TerrainScript []string

type TowerScript [][]string

func (m TowerScript) SaveJSON(filename string) error {
	j, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return fmt.Errorf("err in make json %v", err)
	}
	fd, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("err in create %v", err)
	}
	defer fd.Close()
	n, err := fd.Write(j)
	if err != nil {
		return fmt.Errorf("err in write %v %v", n, err)
	}
	return nil
}
func LoadJSON(filename string) (TowerScript, error) {
	m := TowerScript{}
	fd, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Fail to %v", err)
	}
	defer fd.Close()

	dec := json.NewDecoder(fd)
	err = dec.Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("Error in decode %v ", err)
	}
	return m, nil
}
