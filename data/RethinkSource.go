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
	Kernel "gopkg.in/star-inc/kaguyakernel.v2"
	"log"
)

type RethinkSource struct {
	Table   string
	Term    rethinkdb.Term
	Session *rethinkdb.Session
}

// NewRethinkSource: create a new RethinkSource instance to connect RethinkDB Server.
func NewRethinkSource(config Kernel.RethinkConfig, table string) (*RethinkSource, error) {
	var err error
	instance := new(RethinkSource)
	instance.Term = rethinkdb.DB(config.DatabaseName)
	instance.Session, err = rethinkdb.Connect(config.ConnectConfig)
	if err != nil {
		return nil, err
	}
	instance.Table = table
	return instance, nil
}

func (s *RethinkSource) GetFetchCursor() *rethinkdb.Cursor {
	cursor, err := s.Term.Table(s.Table).Changes().Run(s.Session)
	if err != nil {
		log.Panicln(err)
	}
	return cursor
}
