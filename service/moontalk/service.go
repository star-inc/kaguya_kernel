// Copyright 2021 Star Inc.(https://starinc.xyz)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package moontalk

import (
	"context"
	"github.com/mitchellh/mapstructure"
	Kernel "gopkg.in/star-inc/kaguyakernel.v2"
	"gopkg.in/star-inc/kaguyakernel.v2/data"
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"log"
)

// Service is the data structure of Moonpass Service.
type Service struct {
	Kernel.Service
	source *KernelSource.ContainerSource
}

// NewServiceInterface will create service interface of Moonpass.
func NewServiceInterface(source KernelSource.Interface) ServiceInterface {
	service := new(Service)
	service.source = source.(*KernelSource.ContainerSource)
	return service
}

// CheckPermission will check the permission of client.
func (service *Service) CheckPermission() bool {
	return service.GetGuard().Permission(service.source.RelationID)
}

// Fetch will do the fetch for data, if there is a change in database, it will throw the event out.
func (service *Service) Fetch(ctx context.Context) {
	cursor := service.source.GetFetchCursor()
	defer func() {
		if err := cursor.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	var row interface{}
	for cursor.Next(&row) {
		select {
		case <-ctx.Done():
			return
		default:
			service.GetSession().Respond(row)
		}
	}
	if err := cursor.Err(); err != nil {
		service.GetSession().RaiseError(err.Error())
	}
}

// GetHistoryMessages will get the history messages for client.
func (service *Service) GetHistoryMessages(request *Kernel.Request) {
	query := request.Data.(map[string]interface{})
	timestamp := int64(query["timestamp"].(float64))
	limit := int64(query["count"].(float64))
	containers := data.FetchContainersByTimestamp(service.source, timestamp, limit)
	service.GetSession().Respond(containers)
	request.Processed = true
}

// GetMessage will get the message specific for client.
func (service *Service) GetMessage(request *Kernel.Request) {
	container := new(data.Container)
	err := container.Load(service.source, request.Data.(string))
	if err == nil {
		service.GetSession().Respond(container)
	} else {
		service.GetSession().RaiseError(err.Error())
	}
	request.Processed = true
}

// SendMessage will send a message by the request from client.
func (service *Service) SendMessage(request *Kernel.Request) {
	message := new(data.Message)
	if err := mapstructure.Decode(request.Data, message); err != nil {
		log.Panicln(err)
		return
	}
	message.Origin = service.GetGuard().Me()
	container := data.NewContainer(message)
	if err := container.Create(service.source); err != nil {
		service.GetSession().RaiseError(err.Error())
	}
	request.Processed = true
}

// CancelSentMessage will cancel a message delivery by the request from client.
func (service *Service) CancelSentMessage(request *Kernel.Request) {
	container := new(data.Container)
	if err := container.Load(service.source, request.Data.(string)); err != nil {
		service.GetSession().RaiseError(err.Error())
	}
	container.Canceled = true
	container.Message.Content = ""
	if err := container.Replace(service.source); err != nil {
		service.GetSession().RaiseError(err.Error())
	}
	request.Processed = true
}
