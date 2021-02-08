/*
Package KaguyaKernel : The kernel for Kaguya

    Copyright 2021 Star Inc.(https://starinc.xyz)

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
package talk

import (
	"fmt"
	"github.com/star-inc/kaguya_kernel"
	"gopkg.in/olahol/melody.v1"
	"reflect"
	"strings"
)

const (
	ErrorEmptyContent = "Content is empty"
	TargetNotExists   = "Target not exist"
)

type Service struct {
	data      *DataInterface
	authorize *KaguyaKernel.Authorize
	session   *KaguyaKernel.ResponseHandler
}

func NewServiceInterface(session *melody.Session) ServiceInterface {
	Handler := new(Service)
	Handler.data = NewDataInterface()
	Handler.authorize = KaguyaKernel.NewAuthorizeHandler()
	Handler.session = KaguyaKernel.NewResponseHandler(session)
	return Handler
}

func (handler *Service) Run(request KaguyaKernel.Request) {
	fooType := reflect.TypeOf(Service{})
	for i := 0; i < fooType.NumMethod(); i++ {
		method := fooType.Method(i)
		fmt.Println(method.Name)
	}
}

func (handler *Service) fetchMessage() {
	messages := handler.data.FetchMessage(handler.authorize.User.Identity)
	handler.session.Response(messages.([]*Message))
}

func (handler *Service) syncMessageBox() {
	messages := handler.data.SyncMessageBox(handler.authorize.User.Identity)
	handler.session.Response(messages.([]*Message))
}

func (handler *Service) getMessageBox(request KaguyaKernel.Request) {
	requestData := (request.Data).(map[string]interface{})
	messages := handler.data.GetMessageBox(
		handler.authorize.User.Identity,
		requestData["target"].(string),
	)
	handler.session.Response(messages.([]*Message))
}

func (handler *Service) getMessage(request KaguyaKernel.Request) {
	handler.session.Response((request.Data).(Message))
}

func (handler *Service) sendMessage(request KaguyaKernel.Request) {
	requestData := (request.Data).(Message)
	if len(strings.Trim(string(requestData.Content), " ")) == 0 {
		handler.session.ErrorRaise(ErrorEmptyContent)
		return
	}
	if !handler.authorize.CheckUserExisted(requestData.Target) {
		handler.session.ErrorRaise(TargetNotExists)
		return
	}
	handler.data.SaveMessage(requestData)
	handler.session.Response(requestData)
}
