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

package moonbox

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
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

func (service *Service) InitPublicKey(request *Kernel.Request) {
	moonpass := request.Data.(*data.Moonpass)
	blk := make([]byte, 32)
	if _, err := rand.Read(blk); err != nil {
		log.Panicln(err)
	}
	blkHash := sha256.Sum256(blk)
	blkHashHex := fmt.Sprintf("%x", blkHash)
	question := moonpass.Encrypt(blk)
	confirmMessage := data.NewMoonpassConfirm(moonpass, blkHashHex, question)
	service.GetSession().Respond(confirmMessage)
	request.Processed = true
}

func (service *Service) ConfirmPublicKey(request *Kernel.Request) {
	confirmMessage := request.Data.(*data.MoonpassConfirm)
	blkHash := sha256.Sum256(confirmMessage.Data)
	blkHashHex := fmt.Sprintf("%x", blkHash)
	question := confirmMessage.Moonpass.Encrypt(confirmMessage.Data)
	questionHash := sha256.Sum256(question)
	questionHashHex := fmt.Sprintf("%x", questionHash)
	if confirmMessage.Hash == blkHashHex && confirmMessage.Hash == questionHashHex {
		if err := confirmMessage.Moonpass.Create(service.source); err != nil {
			panic(err)
		}
		service.GetSession().Respond("")
	} else {
		service.GetSession().Respond("")
	}
	request.Processed = true
}

func (service *Service) GetPublicKey(request *Kernel.Request) {
	moonpass := new(data.Moonpass)
	if err := moonpass.Load(service.source, request.Data.(string)); err != nil {
		panic(err)
	} else {
		service.GetSession().Respond(moonpass)
	}
	request.Processed = true
}

func (service *Service) SyncPublicKeys(_ *Kernel.Request) {
	moonpasses := make([]data.Moonpass, 2)
	service.GetSession().Respond(moonpasses)
}
