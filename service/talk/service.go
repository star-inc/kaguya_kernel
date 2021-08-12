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

package talk

import (
	"context"
	"github.com/mitchellh/mapstructure"
	Kernel "gopkg.in/star-inc/kaguyakernel.v2"
	"gopkg.in/star-inc/kaguyakernel.v2/data"
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"log"
	"strings"
)

const (
	ErrorEmptyContent   = "Content_is_empty"
	ErrorInvalidContent = "Content_is_invalid"
	ErrorOriginNotEmpty = "Origin_is_not_empty"
)

type Service struct {
	Kernel.Service
	source           *KernelSource.ContainerSource
	contentValidator func(contentType int, content string) bool
}

func NewServiceInterface(source KernelSource.Interface, contentValidator func(contentType int, content string) bool) ServiceInterface {
	service := new(Service)
	service.source = source.(*KernelSource.ContainerSource)
	service.contentValidator = contentValidator
	return service
}

func (service *Service) CheckPermission() bool {
	return false
}

func (service *Service) Fetch(ctx context.Context) {
	cursor := service.source.GetFetchCursor()
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	var row interface{}
	for cursor.Next(&row) {
		select {
		case <-ctx.Done():
			return
		default:
			service.GetSession().Response(row)
		}
	}
	if err := cursor.Err(); err != nil {
		service.GetSession().RaiseError(err.Error())
	}
}

func (service *Service) GetHistoryMessages(request *Kernel.Request) {
	query := request.Data.(map[string]interface{})
	timestamp := int64(query["timestamp"].(float64))
	limit := int64(query["count"].(float64))
	containers := data.FetchContainersByTimestamp(service.source, timestamp, limit)
	service.GetSession().Response(containers)
}

func (service *Service) GetMessage(request *Kernel.Request) {
	container := new(data.Container)
	err := container.Load(service.source, request.Data.(string))
	if err == nil {
		service.GetSession().Response(container)
	} else {
		service.GetSession().RaiseError(err.Error())
	}
}

func (service *Service) SendMessage(request *Kernel.Request) {
	message := new(data.Message)
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
	container := data.NewContainer(message)
	err = container.Create(service.source)
	if err != nil {
		service.GetSession().RaiseError(err.Error())
	}
}

func (service *Service) CancelSentMessage(request *Kernel.Request) {
	container := new(data.Container)
	err := container.Load(service.source, request.Data.(string))
	if err != nil {
		service.GetSession().RaiseError(err.Error())
	}
	err = container.Destroy(service.source)
	if err != nil {
		service.GetSession().RaiseError(err.Error())
	}
}
