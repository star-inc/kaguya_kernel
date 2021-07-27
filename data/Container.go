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
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

// Container: Container is the data structure, only can be modified by server only, to include a message into database.
type Container struct {
	UUID        string   `rethinkdb:"id,omitempty" json:"uuid"`
	Message     *Message `rethinkdb:"message" json:"message"`
	CreatedTime int64    `rethinkdb:"createdTime" json:"createdTime"`
	Canceled    bool     `rethinkdb:"canceled" json:"canceled"`
}

// NewContainer: include a message automatically, the function will fill the information required for Container.
func NewContainer(source *RethinkSource, message *Message) Interface {
	instance := new(Container)
	instance.UUID = uuid.New().String()
	instance.Message = message
	instance.CreatedTime = time.Now().UnixNano()
	instance.Canceled = false
	return instance
}

// Load: load a message from database, filter is the message ID.
func (c *Container) Load(source *RethinkSource, filter ...interface{}) error {
	cursor, err := source.Term.Table(source.Table).Get(filter[0].(string)).Run(source.Session)
	if err != nil {
		return err
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	return cursor.One(c)
}

// Create: create a new message to database.
func (c *Container) Create(source *RethinkSource) error {
	return source.Term.Table(source.Table).Insert(c).Exec(source.Session)
}

// Update: update a message context in database.
func (c *Container) Update(source *RethinkSource) error {
	return source.Term.Table(source.Table).Update(c).Exec(source.Session)
}

// Replace: update a message context in database.
func (c *Container) Replace(source *RethinkSource) error {
	return source.Term.Table(source.Table).Replace(c).Exec(source.Session)
}

// Destroy: the method can not be called.
func (c *Container) Destroy(_ *RethinkSource) error {
	return errors.New(ErrorBadMethodCallException)
}
