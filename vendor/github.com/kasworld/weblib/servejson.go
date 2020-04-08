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
	"encoding/json"
	"fmt"
	"net/http"
)

func ServeJSON2HTTP(obj interface{}, w http.ResponseWriter) error {
	data, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "fail to marshal json", 404)
		return fmt.Errorf("fail to marshal json %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "fail to write", 404)
		return fmt.Errorf("fail to write %v", err)
	}
	return nil
}
