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

package box

import (
	"context"
	Kernel "gopkg.in/star-inc/kaguyakernel.v2"
	"gopkg.in/star-inc/kaguyakernel.v2/data"
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"gopkg.in/star-inc/kaguyakernel.v2/time"
	"log"
)

type SyncExtraDataAssigner func(syncMessagebox data.SyncMessagebox) interface{}

// Service is the data structure of Messagebox Service.
type Service struct {
	Kernel.Service
	source                *KernelSource.MessageboxSource
	syncExtraDataAssigner SyncExtraDataAssigner
}

// NewServiceInterface will create service interface of Messagebox.
func NewServiceInterface(source KernelSource.Interface, syncExtraDataAssigner SyncExtraDataAssigner) ServiceInterface {
	service := new(Service)
	service.source = source.(*KernelSource.MessageboxSource)
	service.syncExtraDataAssigner = syncExtraDataAssigner
	return service
}

// CheckPermission will check the permission of client.
func (service *Service) CheckPermission() bool {
	return service.GetGuard().Permission(service.source.ClientID)
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
		service.GetSession().RaiseError(err)
	}
}

// SyncMessagebox will get the history messageboxes for client.
func (service *Service) SyncMessagebox(request *Kernel.Request) {
	query := request.Data.(map[string]interface{})
	timestamp := time.NanoTime(query["timestamp"].(float64))
	limit := int64(query["count"].(float64))
	syncMessageboxes := data.FetchSyncMessageboxesByTimestamp(service.source, timestamp, limit)
	if syncMessageboxes != nil && len(syncMessageboxes) != 0 {
		for i, syncMessagebox := range syncMessageboxes {
			syncMessageboxes[i].ExtraData = service.syncExtraDataAssigner(syncMessagebox)
		}
	}
	service.GetSession().Respond(syncMessageboxes)
	request.Processed = true
}

// DeleteMessagebox will delete a messagebox by the request from client.
func (service *Service) DeleteMessagebox(request *Kernel.Request) {
	messagebox := data.NewMessagebox()
	if err := messagebox.Load(service.source, request.Data.(string)); err != nil {
		service.GetSession().RaiseError(err)
	}
	if err := messagebox.Destroy(service.source); err != nil {
		service.GetSession().RaiseError(err)
	}
	request.Processed = true
}
