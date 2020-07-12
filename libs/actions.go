/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
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
