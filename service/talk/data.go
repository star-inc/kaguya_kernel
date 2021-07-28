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
	"log"
	"time"

	"github.com/google/uuid"
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	Kernel "gopkg.in/star-inc/kaguyakernel.v1"
)

type Data struct {
	session    *Rethink.Session
	database   Rethink.Term
	chatRoomID string
}

func newData(config Kernel.RethinkConfig, chatRoomID string) *Data {
	var err error
	data := new(Data)
	data.session, err = Rethink.Connect(config.ConnectConfig)
	if err != nil {
		log.Panicln(err)
	}
	data.database = Rethink.DB(config.DatabaseName)
	data.chatRoomID = chatRoomID
	return data
}

func newDatabaseMessage(rawMessage *Message) *DatabaseMessage {
	dbMessage := new(DatabaseMessage)
	dbMessage.UUID = uuid.New().String()
	dbMessage.CreatedTime = time.Now().UnixNano()
	dbMessage.Message = rawMessage
	dbMessage.Canceled = false
	return dbMessage
}

func (data *Data) getFetchCursor() *Rethink.Cursor {
	cursor, err := data.database.Table(data.chatRoomID).
		Changes().
		Run(data.session)
	if err != nil {
		log.Panicln(err)
	}
	return cursor
}

func (data *Data) getHistoryMessages(timestamp int, count int) *[]DatabaseMessage {
	messages := new([]DatabaseMessage)
	cursor, err := data.database.Table(data.chatRoomID).
		OrderBy(Rethink.Desc("createdTime")).
		Filter(Rethink.Row.Field("createdTime").Lt(timestamp)).
		Limit(count).
		Run(data.session)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	err = cursor.All(messages)
	if err == Rethink.ErrEmptyResult {
		return nil
	}
	if err != nil {
		log.Panicln(err)
	}
	return messages
}

func (data *Data) getMessage(messageID string) *DatabaseMessage {
	message := new(DatabaseMessage)
	cursor, err := data.database.Table(data.chatRoomID).Get(messageID).Run(data.session)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	err = cursor.One(message)
	if err == Rethink.ErrEmptyResult {
		return nil
	}
	if err != nil {
		log.Panicln(err)
	}
	return message
}

func (data *Data) insertMessage(rawMessage *Message) *DatabaseMessage {
	message := newDatabaseMessage(rawMessage)
	err := data.database.Table(data.chatRoomID).Insert(message).Exec(data.session)
	if err != nil {
		log.Panicln(err)
	}
	return message
}

func (data *Data) updateMessage(message *DatabaseMessage) {
	err := data.database.Table(data.chatRoomID).Replace(message).Exec(data.session)
	if err != nil {
		log.Panicln(err)
	}
}
