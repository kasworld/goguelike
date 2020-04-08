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

package configutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	ini "gopkg.in/ini.v1"
)

func LoadData(urlpath string) ([]byte, error) {
	var fd io.Reader
	u, err := url.Parse(urlpath)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		resp, err := http.Get(urlpath)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		fd = resp.Body
	} else {
		ffd, err := os.Open(urlpath)
		if err != nil {
			return nil, err
		}
		defer ffd.Close()
		fd = ffd
	}
	return ioutil.ReadAll(fd)
}

func LoadIni(urlpath string, config interface{}) error {
	datas, err := LoadData(urlpath)
	if err != nil {
		return err
	}
	f, err := ini.Load(datas)
	if err != nil {
		return err
	}
	if err := f.MapTo(config); err != nil {
		return err
	}
	return nil
}

func SaveIni(filename string, config interface{}) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return GenerateToIni(config, fd)
}

func GenerateToIni(config interface{}, w io.Writer) error {
	cfg := ini.Empty()
	ini.ReflectFrom(cfg, config)
	_, err := cfg.WriteToIndent(w, "\t")
	return err
}

func LoadJSON(urlpath string, config interface{}) error {
	fd, err := os.Open(urlpath)
	if err != nil {
		return err
	}
	defer fd.Close()

	dec := json.NewDecoder(fd)
	err = dec.Decode(config)
	if err != nil {
		return err
	}
	return nil
}

func SaveJSON(filename string, config interface{}) error {
	j, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	if _, err := fd.Write(j); err != nil {
		return err
	}
	return nil
}
