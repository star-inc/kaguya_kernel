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
