/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

func (handle *Handle) QueryServices() {
	switch handle.request.ActionType {
	case "authService":
		handle.authService()
		break
	case "talkService":
		handle.talkService()
		break
	}
}

func (handle *Handle) authService() {
	switch handle.request.Action {
	case "registerUser":
		account := handle.request.Data.(map[string]interface{})
		data := handle.RegisterUser(account["displayName"].(string), account["username"].(string), account["password"].(string))
		handle.Response(false, handle.request.ActionType, handle.request.Action, data)
		break
	case "getAccess":
		id := handle.request.Data.(map[string]interface{})
		data := handle.GetAccess(id["username"].(string), id["password"].(string))
		handle.Response(false, handle.request.ActionType, handle.request.Action, data)
		break
	}
}

func (handle *Handle) talkService() {
	switch handle.request.Action {
	case "sendMessage":
		TalkService_SendMessage(handle)
		break
	}
}
