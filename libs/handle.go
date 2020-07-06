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

type Handle struct {
	identity      string
	request       KaguyaRequest
	wsHandle      *websocket.Conn
	dataInterface *DataInterface
	startedPoll   bool
}

func NewHandleInterface(wsHandle *websocket.Conn) *Handle {
	handle := new(Handle)
	handle.wsHandle = wsHandle
	handle.dataInterface = NewDataInterface()
	handle.startedPoll = false
	return handle
}

// Start :
func (handle *Handle) Start() {
	for {
		err := handle.wsHandle.ReadJSON(&handle.request)
		DeBug("WS Read", err)
		if handle.request.Version < 1 {
			handle.ErrorRaise(false, "core", "version", "End of Support")
			return
		}
		if handle.request.ActionType == "authService" {
			go handle.QueryServices()
		} else if handle.request.AuthToken != "" {
			handle.VerfiyAccess(handle.request.AuthToken)
			if handle.identity != "" {
				go handle.QueryServices()
				if !handle.startedPoll {
					go handle.PollServices()
					handle.startedPoll = true
				}
			} else {
				go handle.ErrorRaise(false, "core", "verify", "Unauthorized")
			}
		} else {
			go handle.ErrorRaise(false, "core", "verify", "Unauthorized")
		}
	}
}

// Response :
func (handle *Handle) Response(initiative bool, serviceCode string, actionCode string, data interface{}) {
	var actionID string
	now := time.Now().Unix()
	if initiative {
		actionID = uuid.New().String()
	} else {
		actionID = handle.request.ActionID
	}
	handle.wsHandle.WriteJSON(
		&KaguyaResponse{
			Time:       now,
			ActionID:   actionID,
			ActionType: serviceCode,
			Action:     actionCode,
			Data:       data,
		},
	)
}

func (handle *Handle) ErrorRaise(initiative bool, serviceCode string, actionCode string, message string) {
	handle.Response(initiative, serviceCode, actionCode, &KaguyaErrorRaise{Error: message})
}
