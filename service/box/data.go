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
	session    *Rethink.Session
	database   Rethink.Term
	listenerID string
}

func newData(config Kernel.RethinkConfig, listenerID string) *Data {
	var err error
	data := new(Data)
	data.session, err = Rethink.Connect(config.ConnectConfig)
	if err != nil {
		log.Panicln(err)
	}
	data.database = Rethink.DB(config.DatabaseName)
	data.listenerID = listenerID
	return data
}

func (data *Data) fetchMessage(session *Kernel.Session) {
	cursor, err := data.database.Table(data.listenerID).Changes().Run(data.session)
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

func (data *Data) getHistoryMessagebox(timestamp int, count int) *[]Messagebox {
	messages := new([]Messagebox)
	cursor, err := data.database.Table(data.listenerID).
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
	if err != nil {
		log.Panicln(err)
	}
	return messages
}

func (data *Data) getMessagebox(target string) *Messagebox {
	messagebox := new(Messagebox)
	cursor, err := data.database.Table(data.listenerID).Get(target).Run(data.session)
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

func (data *Data) deleteMessagebox(target string) {
	err := data.database.Table(data.listenerID).Get(target).Delete().Exec(data.session)
	if err != nil {
		log.Panicln(err)
	}
}
