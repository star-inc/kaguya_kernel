/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

import (
	"log"

	"github.com/gorilla/websocket"
)

func HandleRequest(wsHandle *websocket.Conn) {
	for {
		mtype, msg, err := wsHandle.ReadMessage()
		DeBug("WS Read", err)
		log.Printf("Received: %s\n", msg)
		data := NewDataInterface()
		data.LogMessage(msg)
		err = wsHandle.WriteMessage(mtype, msg)
		DeBug("WS Write", err)
	}
}
