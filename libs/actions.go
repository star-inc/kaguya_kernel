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

func TalkService_SendMessage(Handler *Handler) {
	msg := (Handler.request.Data).(map[string]interface{})
	target := msg["target"].(string)
	message := []byte(msg["content"].(string))
	if !Handler.dataInterface.CheckUserExisted(target) {
		Handler.ErrorRaise(false, Handler.request.ActionType, Handler.request.Action, "Target not exist")
		return
	}
	if string(message) == "" {
		Handler.ErrorRaise(false, Handler.request.ActionType, Handler.request.Action, "Content is empty")
		return
	}
	output := new(Message)
	output.ContentType = int(msg["contentType"].(float64))
	output.TargetType = int(msg["targetType"].(float64))
	output.Target = target
	output.Content = message
	output.Origin = Handler.identity
	Handler.dataInterface.LogMessage(output)
	Handler.Response(false, Handler.request.ActionType, Handler.request.Action, output)
}

func TalkService_LoadMessage(Handler *Handler) {
	output := []*Message{}
	messages := Handler.dataInterface.GetMessageBox(Handler.identity)
	if messages != nil {
		output = messages.([]*Message)
	}
	Handler.Response(true, Handler.request.ActionType, Handler.request.Action, output)
}

func TalkService_ReceiveMessage(Handler *Handler) {
	msg := (Handler.request.Data).(map[string]interface{})
	output := new(Message)
	output.ContentType = int(msg["contentType"].(float64))
	output.TargetType = int(msg["targetType"].(float64))
	output.Target = msg["target"].(string)
	output.Content = []byte(msg["content"].(string))
	Handler.Response(true, Handler.request.ActionType, Handler.request.Action, output)
}
