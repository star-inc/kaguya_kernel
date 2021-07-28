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

package box

import (
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	Kernel "gopkg.in/star-inc/kaguyakernel.v1"
	"log"
)

type Manager struct {
	Data
}

func NewManager(config Kernel.RethinkConfig, listenerID string) *Manager {
	var err error
	manager := new(Manager)
	manager.session, err = Rethink.Connect(config.ConnectConfig)
	if err != nil {
		log.Panicln(err)
	}
	manager.database = Rethink.DB(config.DatabaseName)
	manager.listenerID = listenerID
	return manager
}

func (manager *Manager) Check() bool {
	cursor, err := manager.database.TableList().
		Contains(manager.listenerID).
		Run(manager.session)
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

func (manager *Manager) Create() bool {
	err := manager.database.TableCreate(
		manager.listenerID,
		Rethink.TableCreateOpts{PrimaryKey: "target"},
	).
		IndexCreate("origin").
		IndexCreate("createdTime").
		Exec(manager.session)
	if err != nil {
		return false
	}
	return true
}

func (manager *Manager) Drop() bool {
	err := manager.database.
		TableDrop(manager.listenerID).
		Exec(manager.session)
	if err != nil {
		return false
	}
	return true
}
