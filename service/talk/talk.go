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
)

type Service struct {
	Kernel.Service
	data *Data
}

func NewServiceInterface() ServiceInterface {
	service := new(Service)
	service.data = NewData()
	return service
}

func (service *Service) fetchMessage() {
	messages := service.data.FetchMessage(service.GetGuard().User.Identity)
	service.GetSession().Response(messages.([]*Message))
}

func (service *Service) syncMessageBox() {
	messages := service.data.SyncMessageBox(service.GetGuard().User.Identity)
	service.GetSession().Response(messages.([]*Message))
}

func (service *Service) getMessageBox(request Kernel.Request) {
	message := (request.Data).(Message)
	messages := service.data.GetMessageBox(
		service.GetGuard().User.Identity,
		message.Target,
	)
	service.GetSession().Response(messages.([]*Message))
}

func (service *Service) getMessage(request Kernel.Request) {
	service.GetSession().Response((request.Data).(Message))
}

func (service *Service) sendMessage(request Kernel.Request) {
	message := (request.Data).(Message)
	if len(strings.Trim(string(message.Content), " ")) == 0 {
		service.GetSession().RaiseError(ErrorEmptyContent)
		return
	}
	if !service.GetGuard().CheckUserExisted(message.Target) {
		service.GetSession().RaiseError(TargetNotExists)
		return
	}
	service.data.SaveMessage(message)
	service.GetSession().Response(message)
}
