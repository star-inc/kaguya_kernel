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

package manager

import (
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"log"
)

type Messagebox struct {
	source *KernelSource.MessageboxSource
}

func NewMessageboxManager() Interface {
	m := new(Messagebox)
	return m
}

func (m *Messagebox) Check() bool {
	cursor, err := m.source.
		GetTerm().
		TableList().
		Contains(m.source.ClientID).
		Run(m.source.GetSession())
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

func (m *Messagebox) Create() error {
	return m.source.
		GetTerm().
		TableCreate(
			m.source.ClientID,
			Rethink.TableCreateOpts{PrimaryKey: "target"},
		).
		IndexCreate("origin").
		IndexCreate("createdTime").
		Exec(m.source.GetSession())
}

func (m *Messagebox) Drop() error {
	return m.source.
		GetTerm().
		TableDrop(m.source.ClientID).
		Exec(m.source.GetSession())
}
