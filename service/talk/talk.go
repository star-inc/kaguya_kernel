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
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	Kernel "github.com/star-inc/kaguya_kernel"
	"log"
	"sort"
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
	readMessagesHook func(*DatabaseMessage)
	sendMessageHook  func(*DatabaseMessage)
}

func NewServiceInterface(
	dbConfig Kernel.RethinkConfig,
	chatRoomID string,
	contentValidator func(int, string) bool,
	readMessagesHook func(*DatabaseMessage),
	sendMessageHook func(*DatabaseMessage),
) ServiceInterface {
	service := new(Service)
	service.data = newData(dbConfig, chatRoomID)
	service.contentValidator = contentValidator
	service.readMessagesHook = readMessagesHook
	service.sendMessageHook = sendMessageHook
	return service
}

func (service *Service) CheckPermission() bool {
	if !service.GetGuard().Permission(service.data.chatRoomID) {
		return false
	}
	return true
}

func (service *Service) Fetch(ctx context.Context) {
	cursor := service.data.getFetchCursor()
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	var row interface{}
	for cursor.Next(&row) {
		select {
		case <-ctx.Done():
			log.Println("Stop Fetching")
			return
		default:
			service.GetSession().Response(row)
			row := row.(map[string]interface{})
			message := new(DatabaseMessage)
			err := mapstructure.Decode(row["new_val"], message)
			if err != nil {
				log.Println(err)
				return
			}
			service.readMessagesHook(message)
		}
	}
	if err := cursor.Err(); err != nil {
		service.GetSession().RaiseError(err.Error())
	}
}

func (service *Service) GetHistoryMessages(request *Kernel.Request) {
	data := request.Data.(map[string]interface{})
	messages := *service.data.getHistoryMessages(
		int(data["timestamp"].(float64)),
		int(data["count"].(float64)),
	)
	sort.Slice(messages, func(i, j int) bool {
		return (messages)[i].CreatedTime < (messages)[j].CreatedTime
	})
	if messages != nil && len(messages) != 0 {
		service.readMessagesHook(&messages[len(messages)-1])
	}
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
	if message.Message.Origin != service.GetGuard().Me() {
		service.GetSession().RaiseError(Kernel.ErrorForbidden)
		return
	}
	message.Message.Content = fmt.Sprint(time.Now().UnixNano())
	message.Canceled = true
	service.data.updateMessage(message)
}
