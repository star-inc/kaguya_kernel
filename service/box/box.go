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
package box

import (
	Kernel "github.com/star-inc/kaguya_kernel"
	"sort"
)

type Service struct {
	Kernel.Service
	data              *Data
	syncExtradataAssigner func(SyncMessagebox, *interface{})
}

func NewServiceInterface(
	messageBoxConfig Kernel.RethinkConfig,
	syncExtradataAssigner func(SyncMessagebox, *interface{}),
	listenerID string,
) ServiceInterface {
	service := new(Service)
	service.data = newData(messageBoxConfig, listenerID)
	service.syncExtradataAssigner = syncExtradataAssigner
	return service
}

func (service *Service) CheckPermission() bool {
	if !service.GetGuard().Permission(service.data.listenerID) {
		return false
	}
	return true
}

func (service *Service) Fetch() {
	service.data.fetchMessagebox(service.GetSession())
}

func (service *Service) SyncMessagebox(request *Kernel.Request) {
	data := request.Data.(map[string]interface{})
	messages := service.data.getHistoryMessagebox(
		int(data["timestamp"].(float64)),
		int(data["count"].(float64)),
	)
	for _, message := range messages {
		service.syncExtradataAssigner(message, &message.Extradata)
	}
	sort.Slice(messages, func(i, j int) bool {
		return (messages)[i].CreatedTime > (messages)[j].CreatedTime
	})
	service.GetSession().Response(messages)
}

func (service *Service) DeleteMessagebox(request *Kernel.Request) {
	service.data.deleteMessagebox(request.Data.(string))
}
