/*
Package kaguya : The library for kaguya

    Copyright 2020 Star Inc.(https://starinc.xyz)

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
package kaguya

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type configStruct struct {
	Port   int    `json:"port"`
	DBhost string `json:"dbhost"`
	DBname string `json:"dbname"`
}

// Config : Global Settings for butterfly from config.json
var Config configStruct

// ConfigPath : Where the config file placed.
var ConfigPath string = "config.json"

// ReadConfig : Load configure file to Config
func ReadConfig() {
	jsonFile, err := os.Open(ConfigPath)
	DeBug("Get JSON config", err)
	defer jsonFile.Close()
	srcJSON, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(srcJSON, &Config)
	DeBug("Load JSON Initialization", err)
}
