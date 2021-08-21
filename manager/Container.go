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
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"log"
)

type Container struct {
	source *KernelSource.ContainerSource
}

func NewContainerManager() Interface {
	manager := new(Container)
	return manager
}

func (c Container) Check() bool {
	cursor, err := c.source.
		GetTerm().
		TableList().
		Contains(c.source.RelationID).
		Run(c.source.GetSession())
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

func (c Container) Create() error {
	return c.source.
		GetTerm().
		TableCreate(
			c.source.RelationID,
			rethinkdb.TableCreateOpts{PrimaryKey: "id"},
		).
		IndexCreate("origin").
		IndexCreate("createdTime").
		Exec(c.source.GetSession())
}

func (c Container) Drop() error {
	return c.source.
		GetTerm().
		TableDrop(c.source.RelationID).
		Exec(c.source.GetSession())
}
