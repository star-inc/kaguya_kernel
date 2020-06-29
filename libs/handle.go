/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

import (
	"github.com/gorilla/websocket"
)

type Handle struct {
	identify string
	wsHandle *websocket.Conn
}

func NewHandleInterface(wsHandle *websocket.Conn) *Handle {
	handle := new(Handle)
	handle.wsHandle = wsHandle
	return handle
}

// Start :
func (handle *Handle) Start() {
	for {
		mtype, msg, err := handle.wsHandle.ReadMessage()
		DeBug("WS Read", err)
		switch mtype {
		case 1:
			go func() {
				if handle.identify != "" {
					if msg != nil {
						go handle.HandleActions(msg)
					}
				} else {
					if string(msg) == "auth" {
						handle.identify = "msg"
						handle.wsHandle.WriteMessage(1, []byte("Authorizing"))
						handle.wsHandle.WriteMessage(1, []byte("Authorized"))
					} else {
						handle.wsHandle.WriteMessage(1, []byte("Not authorized"))
					}
				}
			}()
			break
		}
	}
}
