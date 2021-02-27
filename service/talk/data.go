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
	"github.com/google/uuid"
	Kernel "github.com/star-inc/kaguya_kernel"
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
	"time"
)

type Data struct {
	session   *Rethink.Session
	database  Rethink.Term
	tableName string
}

func newData(config Kernel.RethinkConfig, tableName string) *Data {
	var err error
	data := new(Data)
	data.session, err = Rethink.Connect(config.ConnectConfig)
	if err != nil {
		log.Panicln(err)
	}
	data.database = Rethink.DB(config.DatabaseName)
	data.tableName = tableName
	return data
}

func newDatabaseMessage(rawMessage *Message) *DatabaseMessage {
	dbMessage := new(DatabaseMessage)
	dbMessage.UUID = uuid.New().String()
	dbMessage.CreatedTime = time.Now().UnixNano()
	dbMessage.Message = rawMessage
	return dbMessage
}

func (data Data) fetchMessage(session *Kernel.Session) {
	cursor, err := data.database.Table(data.tableName).Changes().Run(data.session)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	var row interface{}
	for cursor.Next(&row) {
		session.Response(row)
	}
	if err := cursor.Err(); err != nil {
		session.RaiseError(err.Error())
	}
}

func (data Data) getHistoryMessages(timestamp int, count int) *[]DatabaseMessage {
	messages := new([]DatabaseMessage)
	cursor, err := data.database.Table(data.tableName).
		OrderBy(Rethink.Asc("createdTime")).
		Filter(Rethink.Row.Field("createdTime").Ge(timestamp)).
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
	if err != nil {
		log.Panicln(err)
	}
	return messages
}

func (data Data) getMessage(messageID string) *DatabaseMessage {
	message := new(DatabaseMessage)
	cursor, err := data.database.Table(data.tableName).Get(messageID).Run(data.session)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	err = cursor.One(message)
	if err != nil {
		log.Panicln(err)
	}
	return message
}

func (data Data) saveMessage(rawMessage *Message) {
	message := newDatabaseMessage(rawMessage)
	err := data.database.Table(data.tableName).Insert(message).Exec(data.session)
	if err != nil {
		log.Panicln(err)
	}
}
