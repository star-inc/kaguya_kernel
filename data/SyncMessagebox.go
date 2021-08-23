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

package data

import (
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"log"
)

// SyncMessagebox this is the wrapper to wrap client a messagebox with extra data.
type SyncMessagebox struct {
	Messagebox
	ExtraData interface{} `json:"extraData"`
}

// FetchSyncMessageboxesByTimestamp ToDo
func FetchSyncMessageboxesByTimestamp(source *KernelSource.MessageboxSource, timestamp int64, limit int64) []SyncMessagebox {
	containers := make([]SyncMessagebox, limit)
	cursor, err := source.Term.Table(source.ClientID).
		OrderBy(rethinkdb.Desc("createdTime")).
		Filter(rethinkdb.Row.Field("createdTime").Lt(timestamp)).
		Limit(limit).
		Run(source.Session)
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		if err := cursor.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	err = cursor.All(&containers)
	if err == rethinkdb.ErrEmptyResult {
		return nil
	}
	if err != nil {
		log.Panicln(err)
	}
	return containers
}
