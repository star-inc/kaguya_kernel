/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
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
