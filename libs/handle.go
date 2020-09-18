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
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	identity      string
	userIdentity  string
	request       KaguyaRequest
	wsHandle      *websocket.Conn
	dataInterface *DataInterface
	startedPoll   bool
}

func NewHandleInterface(wsHandle *websocket.Conn) *Handler {
	Handler := new(Handler)
	Handler.wsHandle = wsHandle
	Handler.dataInterface = NewDataInterface()
	Handler.startedPoll = false
	return Handler
}

// Start :
func (Handler *Handler) Start() {
	for {
		err := Handler.wsHandle.ReadJSON(&Handler.request)
		DeBug("WS Read", err)
		if Handler.request.Version < 1 {
			Handler.ErrorRaise(false, "core", "version", "End of Support")
			return
		}
		if Handler.request.ActionType == "authService" {
			go Handler.QueryServices()
		} else if Handler.request.AuthToken != "" {
			Handler.VerfiyAccess(Handler.request.AuthToken)
			if Handler.identity != "" {
				go Handler.QueryServices()
				if !Handler.startedPoll {
					go Handler.PollServices()
					Handler.startedPoll = true
				}
			} else {
				go Handler.ErrorRaise(false, "core", "verify", "Unauthorized")
			}
		} else {
			go Handler.ErrorRaise(false, "core", "verify", "Unauthorized")
		}
	}
}

// Response :
func (Handler *Handler) Response(initiative bool, serviceCode string, actionCode string, data interface{}) {
	var actionID string
	now := time.Now().Unix()
	if initiative {
		actionID = uuid.New().String()
	} else {
		actionID = Handler.request.ActionID
	}
	Handler.wsHandle.WriteJSON(
		&KaguyaResponse{
			Time:       now,
			ActionID:   actionID,
			ActionType: serviceCode,
			Action:     actionCode,
			Data:       data,
		},
	)
}

func (Handler *Handler) ErrorRaise(initiative bool, serviceCode string, actionCode string, message string) {
	Handler.Response(initiative, serviceCode, actionCode, &KaguyaErrorRaise{Error: message})
}
