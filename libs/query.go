/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

func (Handler *Handler) QueryServices() {
	switch Handler.request.ActionType {
	case "authService":
		Handler.authService()
		break
	case "talkService":
		Handler.talkService()
		break
	}
}

func (Handler *Handler) authService() {
	switch Handler.request.Action {
	case "registerUser":
		account := Handler.request.Data.(map[string]interface{})
		data := Handler.RegisterUser(account["displayName"].(string), account["username"].(string), account["password"].(string))
		Handler.Response(false, Handler.request.ActionType, Handler.request.Action, data)
		break
	case "getAccess":
		id := Handler.request.Data.(map[string]interface{})
		data := Handler.GetAccess(id["username"].(string), id["password"].(string))
		Handler.Response(false, Handler.request.ActionType, Handler.request.Action, data)
		break
	}
}

func (Handler *Handler) talkService() {
	switch Handler.request.Action {
	case "sendMessage":
		TalkService_SendMessage(Handler)
		break
	case "syncMessage":
		TalkService_LoadMessage(Handler)
		break
	}
}
