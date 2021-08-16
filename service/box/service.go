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
	"log"
)

// Service: this is the struct of Messagebox Service.
type Service struct {
	Kernel.Service
	source                *KernelSource.MessageboxSource
	syncExtraDataAssigner func(syncMessagebox data.SyncMessagebox) interface{}
}

// NewServiceInterface: create service interface of Messagebox.
func NewServiceInterface(source KernelSource.Interface, syncExtraDataAssigner func(syncMessagebox data.SyncMessagebox) interface{}) ServiceInterface {
	service := new(Service)
	service.source = source.(*KernelSource.MessageboxSource)
	service.syncExtraDataAssigner = syncExtraDataAssigner
	return service
}

// CheckPermission: check the permission of client.
func (service *Service) CheckPermission() bool {
	return service.GetGuard().Permission(service.source.ClientID)
}

// Fetch: do the fetch for data, if there is a change in database, it will throw the event out.
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

// SyncMessagebox: get the history messageboxes for client.
func (service *Service) SyncMessagebox(request *Kernel.Request) {
	query := request.Data.(map[string]interface{})
	timestamp := int64(query["timestamp"].(float64))
	limit := int64(query["count"].(float64))
	syncMessageboxes := data.FetchSyncMessageboxesByTimestamp(service.source, timestamp, limit)
	for i, syncMessagebox := range syncMessageboxes {
		syncMessageboxes[i].ExtraData = service.syncExtraDataAssigner(syncMessagebox)
	}
	service.GetSession().Response(syncMessageboxes)
}

// DeleteMessagebox: delete a messagebox by the request of client.
func (service *Service) DeleteMessagebox(request *Kernel.Request) {
	messagebox := data.NewMessagebox()
	err := messagebox.Load(service.source, request.Data.(string))
	if err != nil {
		service.GetSession().RaiseError(err.Error())
	}
	err = messagebox.Destroy(service.source)
	if err != nil {
		service.GetSession().RaiseError(err.Error())
	}
}
