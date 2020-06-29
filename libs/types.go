/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

const (
	ContentType_Text   = 0
	TargetType_Contact = 0
)

type KaguyaRequest struct {
	Version    int         `json:"version"`
	ActionID   string      `json:"actionID"`
	AuthToken  string      `json:"authToken"`
	ActionType string      `json:"actionType"`
	Action     string      `json:"action"`
	Data       interface{} `json:"data"`
}

type KaguyaResponse struct {
	ActionID   string      `json:"actionID"`
	ActionType string      `json:"actionType"`
	Action     string      `json:"action"`
	Data       interface{} `json:"data"`
}

type Message struct {
	ContentType int    `json:"contentType"`
	TargetType  int    `json:"targetType"`
	Origin      string `json:"origin"`
	Target      string `json:"target"`
	Content     []byte `json:"content"`
}

type User struct {
	DisplayName string `json:"displayName"`
	Identity    string `json:"identity"`
	Password    string `json:"password"`
}
