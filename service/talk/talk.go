/*
Package Kernel : The kernel for Kaguya

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
	Kernel "github.com/star-inc/kaguya_kernel"
	"strings"
)

const (
	ErrorEmptyContent = "Content is empty"
	TargetNotExists   = "Target not exist"
	Forbidden         = "Forbidden"
)

type Service struct {
	Kernel.Service
	data *Data
}

func NewServiceInterface(dbConfig Kernel.RethinkConfig, dbTable string) ServiceInterface {
	service := new(Service)
	service.data = newData(dbConfig, dbTable)
	service.SetMap(service)
	return service
}

func (service *Service) Fetch() {
	service.data.fetchMessage(service)
}

func (service *Service) syncMessageBox() {
	messages := service.data.syncMessageBox(service.GetGuard().Me().Identity)
	service.GetSession().Response(messages)
}

func (service *Service) getMessageBox(request Kernel.Request) {
	messages := service.data.getMessageBox(
		service.GetGuard().Me().Identity,
		(request.Data).(string),
	)
	service.GetSession().Response(messages)
}

func (service *Service) getMessage(request Kernel.Request) {
	dbMessage := service.data.getMessage((request.Data).(string))
	message := dbMessage.Message
	identity := service.GetGuard().Me().Identity
	if message.Target != identity && message.Origin != identity {
		service.GetSession().RaiseError(Forbidden)
	}
	service.GetSession().Response(dbMessage)
}

func (service *Service) sendMessage(request Kernel.Request) {
	message := (request.Data).(*Message)
	if len(strings.Trim(string(message.Content), " ")) == 0 {
		service.GetSession().RaiseError(ErrorEmptyContent)
		return
	}
	if !service.GetGuard().CheckUserExists(message.Target) {
		service.GetSession().RaiseError(TargetNotExists)
		return
	}
	service.data.saveMessage(message)
	service.GetSession().Response(message)
}
