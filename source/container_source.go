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

package source

import (
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
)

type ContainerSource struct {
	Base
	RelationID string
}

// NewContainerSource will create a new Source instance to connect rethinkdbDB Server for Container.
func NewContainerSource(config rethinkdb.ConnectOpts, databaseName string) (Interface, error) {
	var err error
	instance := new(ContainerSource)
	instance.Term = rethinkdb.DB(databaseName)
	instance.Session, err = rethinkdb.Connect(config)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (s *ContainerSource) CheckReady() bool {
	return s.RelationID != "" && s.checkTerm()
}

func (s *ContainerSource) GetFetchCursor() *rethinkdb.Cursor {
	return s.GetRawFetchCursor(s.RelationID)
}

func (s *ContainerSource) checkTerm() bool {
	cursor, err := s.
		GetTerm().
		TableList().
		Contains(s.RelationID).
		Run(s.GetSession())
	if err != nil {
		log.Panicln(err)
	}
	var status bool
	if err = cursor.One(&status); err != nil {
		log.Panicln(err)
	}
	return status
}

func (s *ContainerSource) CreateTerm() error {
	return s.
		GetTerm().
		TableCreate(
			s.RelationID,
			rethinkdb.TableCreateOpts{PrimaryKey: "id"},
		).
		IndexCreate("origin").
		IndexCreate("createdTime").
		Exec(s.GetSession())
}

func (s *ContainerSource) DropTerm() error {
	return s.
		GetTerm().
		TableDrop(s.RelationID).
		Exec(s.GetSession())
}
