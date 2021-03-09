/*
Package KaguyaKernel: The kernel for Kaguya

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
	"github.com/mitchellh/mapstructure"
	Kernel "github.com/star-inc/kaguya_kernel"
	"log"
	"strings"
	"time"
)

const (
	ErrorEmptyContent   = "Content_is_empty"
	ErrorInvalidContent = "Content_is_invalid"
	ErrorOriginNotEmpty = "Origin_is_not_empty"
)

type Service struct {
	Kernel.Service
	data             *Data
	contentValidator func(int, string) bool
	sendMessageHook  func(*DatabaseMessage)
}

func NewServiceInterface(
	dbConfig Kernel.RethinkConfig,
	chatRoomID string,
	contentValidator func(int, string) bool,
	sendMessageHook func(*DatabaseMessage),
) ServiceInterface {
	service := new(Service)
	service.data = newData(dbConfig, chatRoomID)
	service.contentValidator = contentValidator
	service.sendMessageHook = sendMessageHook
	return service
}

func (service *Service) CheckPermission() bool {
	if !service.GetGuard().Permission(service.data.chatRoomID) {
		return false
	}
	return true
}

func (service *Service) Fetch() {
	service.data.fetchMessage(service.GetSession())
}

func (service *Service) GetHistoryMessages(request *Kernel.Request) {
	data := request.Data.(map[string]interface{})
	messages := service.data.getHistoryMessages(
		int(data["timestamp"].(float64)),
		int(data["count"].(float64)),
	)
	service.GetSession().Response(messages)
}

func (service *Service) GetMessage(request *Kernel.Request) {
	dbMessage := service.data.getMessage((request.Data).(string))
	service.GetSession().Response(dbMessage)
}

func (service *Service) SendMessage(request *Kernel.Request) {
	message := new(Message)
	err := mapstructure.Decode(request.Data, message)
	if err != nil {
		log.Println(err)
		return
	}
	if len(strings.Trim(message.Content, " ")) == 0 {
		service.GetSession().RaiseError(ErrorEmptyContent)
		return
	}
	if !service.contentValidator(message.ContentType, message.Content) {
		service.GetSession().RaiseError(ErrorInvalidContent)
		return
	}
	if message.Origin != "" {
		service.GetSession().RaiseError(ErrorOriginNotEmpty)
		return
	}
	message.Origin = service.GetGuard().Me()
	savedMessage := service.data.insertMessage(message)
	service.sendMessageHook(savedMessage)
}

func (service *Service) CancelSentMessage(request *Kernel.Request) {
	message := service.data.getMessage((request.Data).(string))
	message.Message.Content = fmt.Sprint(time.Now().UnixNano())
	message.Canceled = true
	service.data.updateMessage(message)
}
