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
	"github.com/star-inc/kaguya_kernel/service/talk"
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
)

type Hook struct {
	session                   *Rethink.Session
	database                  Rethink.Term
	chatRoomID                string
	getRelation               func(string) []string
	metadataGenerator         func(*talk.DatabaseMessage) string
	messageboxNotFoundHandler func(string) bool
}

func NewHook(
	config Kernel.RethinkConfig,
	chatRoomID string,
	getRelation func(string) []string,
	metadataGenerator func(*talk.DatabaseMessage) string,
	messageboxNotFoundHandler func(string) bool,
) *Hook {
	var err error
	hook := new(Hook)
	hook.session, err = Rethink.Connect(config.ConnectConfig)
	if err != nil {
		log.Panicln(err)
	}
	hook.database = Rethink.DB(config.DatabaseName)
	hook.chatRoomID = chatRoomID
	hook.getRelation = getRelation
	hook.metadataGenerator = metadataGenerator
	hook.messageboxNotFoundHandler = messageboxNotFoundHandler
	return hook
}

func (hook *Hook) MessageTrigger(message *talk.DatabaseMessage) {
	messagebox := new(Messagebox)
	messagebox.Target = hook.chatRoomID
	messagebox.Origin = message.Message.Origin
	messagebox.CreatedTime = message.CreatedTime
	messagebox.Metadata = hook.metadataGenerator(message)
	for _, relatedID := range hook.getRelation(hook.chatRoomID) {
		if !hook.checkMessagebox(relatedID) &&
			hook.messageboxNotFoundHandler(relatedID) {
			continue
		}
		hook.newMessagebox(relatedID, messagebox)
	}
}

func (hook *Hook) SeenTrigger(message *talk.DatabaseMessage) {
	messagebox := new(Messagebox)
	messagebox.Target = hook.chatRoomID
	messagebox.LastSeen = message.CreatedTime
	for _, relatedID := range hook.getRelation(hook.chatRoomID) {
		if !hook.checkMessagebox(relatedID) &&
			hook.messageboxNotFoundHandler(relatedID) {
			continue
		}
		hook.newMessagebox(relatedID, messagebox)
	}
}

func (hook Hook) checkMessagebox(relatedID string) bool {
	cursor, err := hook.database.TableList().
		Contains(relatedID).
		Run(hook.session)
	if err != nil {
		log.Panicln(err)
	}
	var status bool
	err = cursor.One(&status)
	if err != nil {
		log.Panicln(err)
	}
	return status
}

func (hook Hook) newMessagebox(relatedID string, messagebox *Messagebox) {
	err := hook.database.Table(relatedID).
		Insert(
			messagebox,
			Rethink.InsertOpts{
				Conflict: "update",
			},
		).
		Exec(hook.session)
	if err != nil {
		log.Panicln(err)
	}
}
