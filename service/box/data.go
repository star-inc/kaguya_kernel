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
package box

import (
	Kernel "github.com/star-inc/kaguya_kernel"
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
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

func (data Data) getHistoryMessages(timestamp int, count int) *[]Messagebox {
	messages := new([]Messagebox)
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

func (data Data) getMessagebox(origin string) *Messagebox {
	messagebox := new(Messagebox)
	cursor, err := data.database.Table(data.tableName).Get(origin).Run(data.session)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	err = cursor.One(messagebox)
	if err != nil {
		log.Panicln(err)
	}
	return messagebox
}

func (data Data) replaceMessagebox(messagebox *Messagebox) {
	err := data.database.Table(data.tableName).Replace(messagebox).Exec(data.session)
	if err != nil {
		log.Panicln(err)
	}
}

func (data Data) deleteMessagebox(origin string) {
	err := data.database.Table(data.tableName).Get(origin).Delete().Exec(data.session)
	if err != nil {
		log.Panicln(err)
	}
}