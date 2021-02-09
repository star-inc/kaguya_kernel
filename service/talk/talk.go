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
	"github.com/mitchellh/mapstructure"
	Kernel "github.com/star-inc/kaguya_kernel"
	"log"
	"strings"
)

const (
	ErrorEmptyContent   = "Content is empty"
	ErrorInvalidContent = "Content is invalid"
	ErrorOriginNotEmpty = "Origin is not empty"
)

type Service struct {
	Kernel.Service
	data             *Data
	contentValidator func(contentType int, content string) bool
}

func NewServiceInterface(
	dbConfig Kernel.RethinkConfig,
	tableName string,
	contentValidator func(int, string) bool,
) ServiceInterface {
	service := new(Service)
	service.data = newData(dbConfig, tableName)
	service.contentValidator = contentValidator
	return service
}

func (service *Service) CheckPermission() bool {
	if !service.GetGuard().Permission(service.data.tableName) {
		return false
	}
	return true
}

func (service *Service) Fetch() {
	service.data.fetchMessage(service.GetSession())
}

func (service *Service) GetHistoryMessages(request *Kernel.Request) {
	data := request.Data.(map[string]int)
	messages := service.data.getHistoryMessages(data["timestamp"], data["count"])
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
	service.data.saveMessage(message)
}
