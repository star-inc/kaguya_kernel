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
package data

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type DatabaseMessage struct {
	source      RethinkSource
	Canceled    bool     `rethinkdb:"canceled" json:"canceled"`
	CreatedTime int64    `rethinkdb:"createdTime" json:"createdTime"`
	Message     *Message `rethinkdb:"message" json:"message"`
	UUID        string   `rethinkdb:"id,omitempty" json:"uuid"`
}

func NewDatabaseMessage(message *Message) *DatabaseMessage {
	d := new(DatabaseMessage)
	d.UUID = uuid.New().String()
	d.CreatedTime = time.Now().UnixNano()
	d.Message = message
	d.Canceled = false
	return d
}

func (d *DatabaseMessage) Get(chatRoomID string, messageID string) error {
	cursor, err := d.source.Term.Table(chatRoomID).Get(messageID).Run(d.source.Session)
	if err != nil {
		return err
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	return cursor.One(d)
}

func (d *DatabaseMessage) Insert(chatRoomID string) error {
	return d.source.Term.Table(chatRoomID).Insert(d).Exec(d.source.Session)
}

func (d *DatabaseMessage) Update(chatRoomID string) error {
	return d.source.Term.Table(chatRoomID).Replace(d).Exec(d.source.Session)
}
