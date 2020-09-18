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
	Time       int64       `json:"time"`
	ActionID   string      `json:"actionID"`
	ActionType string      `json:"actionType"`
	Action     string      `json:"action"`
	Data       interface{} `json:"data"`
}

type KaguyaErrorRaise struct {
	Error string `json:"error"`
}

type Message struct {
	ContentType int    `json:"contentType"`
	TargetType  int    `json:"targetType"`
	Origin      string `json:"origin"`
	Target      string `json:"target"`
	Content     []byte `json:"content"`
}

type User struct {
	Identity    string `json:"identity"`
	DisplayName string `json:"displayName"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}
