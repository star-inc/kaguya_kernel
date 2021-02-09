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
package talk

import (
	Kernel "github.com/star-inc/kaguya_kernel"
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
)

type Manage struct {
	Data
}

func NewManage(config Kernel.RethinkConfig, tableName string) *Data {
	var err error
	data := new(Data)
	data.session, err = Rethink.Connect(config.ConnectConfig)
	if err != nil {
		log.Fatalln(err)
	}
	data.database = Rethink.DB(config.DatabaseName)
	data.tableName = tableName
	return data
}

func (manage *Manage) Create(string) bool {
	err := manage.database.TableCreate(
		manage.tableName,
		Rethink.TableCreateOpts{PrimaryKey: "id"},
	).
		IndexCreate("origin").
		IndexCreate("createdTime").
		Exec(manage.session)
	if err != nil {
		return false
	}
	return true
}

func (manage *Manage) Drop(string) bool {
	err := manage.database.TableDrop(manage.tableName).Exec(manage.session)
	if err != nil {
		return false
	}
	return true
}
