/*
Package kaguya : The library for kaguya

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package kaguya

func TalkService_SendMessage(handle *Handle) {
	msg := (handle.request.Data).(map[string]interface{})
	output := new(Message)
	output.ContentType = int(msg["contentType"].(float64))
	output.TargetType = int(msg["targetType"].(float64))
	output.Target = msg["target"].(string)
	output.Content = []byte(msg["content"].(string))
	output.Origin = handle.identity
	handle.dataInterface.LogMessage(output)
	handle.Response(false, handle.request.ActionType, handle.request.Action, output)
}

func TalkService_ReceiveMessage(handle *Handle) {
	msg := (handle.request.Data).(map[string]interface{})
	output := new(Message)
	output.ContentType = int(msg["contentType"].(float64))
	output.TargetType = int(msg["targetType"].(float64))
	output.Target = msg["target"].(string)
	output.Content = []byte(msg["content"].(string))
	handle.Response(true, handle.request.ActionType, handle.request.Action, output)
}
