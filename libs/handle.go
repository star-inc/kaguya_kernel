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

type kaguyaRequest struct {
	authToken string
	action    string
	data      interface{}
}

// HandleRequest :
func HandleRequest(wsHandle *websocket.Conn) {
	for {
		data := NewDataInterface()
		mtype, msg, err := wsHandle.ReadMessage()
		DeBug("WS Read", err)
		switch mtype {
		case 1:
			go func() {
				if msg != nil {
					go data.LogMessage(msg)
				}
				err = wsHandle.WriteMessage(mtype, msg)
				DeBug("WS Write", err)
			}()
			break
		}
	}
}
