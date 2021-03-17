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

type Seen struct {
	Data
}

func NewSeen(config Kernel.RethinkConfig) *Seen {
	var err error
	seen := new(Seen)
	seen.session, err = Rethink.Connect(config.ConnectConfig)
	if err != nil {
		log.Panicln(err)
	}
	seen.database = Rethink.DB(config.DatabaseName)
	return seen
}

func (seen *Seen) CountUnreadMessages(chatRoomID string, timestamp int64) int {
	cursor, err := seen.database.Table(chatRoomID).
		OrderBy(Rethink.Desc("createdTime")).
		Filter(Rethink.Row.Field("createdTime").Gt(timestamp)).
		Count().
		Run(seen.session)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	var count int
	err = cursor.One(&count)
	if err != nil {
		log.Panicln(err)
	}
	return count
}
