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
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
)

// SyncContainer: this is the container to wrap client a messagebox with extra data.
type SyncContainer struct {
	Messagebox
	ExtraData interface{} `json:"extraData"`
}

// FetchSyncContainersByTimestamp: ToDo
func FetchSyncContainersByTimestamp(source *RethinkSource, timestamp int, limit int) []SyncContainer {
	containers := make([]SyncContainer, limit)
	cursor, err := source.Term.Table(source.Table).
		OrderBy(Rethink.Desc("createdTime")).
		Filter(Rethink.Row.Field("createdTime").Lt(timestamp)).
		Limit(limit).
		Run(source.Session)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	err = cursor.All(&containers)
	if err == Rethink.ErrEmptyResult {
		return nil
	}
	if err != nil {
		log.Panicln(err)
	}
	return containers
}
