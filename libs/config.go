/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type configStruct struct {
	DBhost string `json:"dbhost"`
	DBname   string `json:"dbname"`
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
