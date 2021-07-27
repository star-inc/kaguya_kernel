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
	data2 "github.com/star-inc/kaguya_kernel/data"
	"github.com/star-inc/kaguya_kernel/data/Message"
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
	data             *data2.Data
	contentValidator func(int, string) bool
	readMessagesHook func(*DatabaseMessage)
	sendMessageHook  func(*DatabaseMessage)
}

func NewServiceInterface(dbConfig Kernel.RethinkConfig, chatRoomID string, contentValidator func(int, string) bool, readMessagesHook func(*DatabaseMessage), sendMessageHook func(*DatabaseMessage)) ServiceInterface {
	return service
}

func (service *Service) CheckPermission() bool {
}

func (service *Service) Fetch(ctx context.Context) {
}

func (service *Service) GetHistoryMessages(request *Kernel.Request) {
}

func (service *Service) GetMessage(request *Kernel.Request) {
}

func (service *Service) SendMessage(request *Kernel.Request) {
}

func (service *Service) CancelSentMessage(request *Kernel.Request) {
}
