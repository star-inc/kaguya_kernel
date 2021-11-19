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
	"errors"
	"github.com/google/uuid"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"gopkg.in/star-inc/kaguyakernel.v2/time"
	"log"
)

// Container is the data structure, only can be modified by server only, to include a message into database.
type Container struct {
	UUID        string        `rethinkdb:"id,omitempty" json:"uuid"`
	Message     *Message      `rethinkdb:"message" json:"message"`
	CreatedTime time.NanoTime `rethinkdb:"createdTime" json:"createdTime"`
	Canceled    bool          `rethinkdb:"canceled" json:"canceled"`
}

// NewContainer will include a message automatically, the function will fill the information required for Container.
func NewContainer(message *Message) Interface {
	instance := new(Container)
	instance.UUID = uuid.New().String()
	instance.Message = message
	instance.CreatedTime = time.Now()
	instance.Canceled = false
	return instance
}

// CheckReady will check model is ready.
func (c *Container) CheckReady() bool {
	return c != nil && c.UUID != ""
}

// Load will load a message from database, filter is the message ID.
func (c *Container) Load(source KernelSource.Interface, filter ...interface{}) error {
	sourceInstance := source.(*KernelSource.ContainerSource)
	cursor, err := source.GetTerm().Table(sourceInstance.RelationID).Get(filter[0].(string)).Run(source.GetSession())
	if err != nil {
		return err
	}
	defer func() {
		if err := cursor.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	return cursor.One(c)
}

// Reload will reload a message from database,
func (c *Container) Reload(source KernelSource.Interface) error {
	return c.Load(source, c.UUID)
}

// Create will create a new message to database.
func (c *Container) Create(source KernelSource.Interface) error {
	sourceInstance := source.(*KernelSource.ContainerSource)
	return source.GetTerm().Table(sourceInstance.RelationID).Insert(c).Exec(source.GetSession())
}

// Replace will update a message context in database.
func (c *Container) Replace(source KernelSource.Interface) error {
	sourceInstance := source.(*KernelSource.ContainerSource)
	return source.GetTerm().Table(sourceInstance.RelationID).Replace(c).Exec(source.GetSession())
}

// Destroy is the method can not be called.
func (c *Container) Destroy(_ KernelSource.Interface) error {
	return errors.New(ErrorBadMethodCallException)
}

// FetchContainersByTimestamp ToDo
func FetchContainersByTimestamp(source *KernelSource.ContainerSource, timestamp time.NanoTime, limit int64) []Container {
	containers := make([]Container, limit)
	cursor, err := source.GetTerm().Table(source.RelationID).
		OrderBy(rethinkdb.Desc("createdTime")).
		Filter(rethinkdb.Row.Field("createdTime").Lt(timestamp)).
		Limit(limit).
		OrderBy(rethinkdb.Asc("createdTime")).
		Run(source.GetSession())
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

// CountContainersByTimestamp ToDo
func CountContainersByTimestamp(source *KernelSource.ContainerSource, timestamp time.NanoTime) int {
	cursor, err := source.GetTerm().Table(source.RelationID).
		Filter(rethinkdb.Row.Field("createdTime").Gt(timestamp)).
		Count().
		Run(source.GetSession())
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		if err := cursor.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	var count int
	if err = cursor.One(&count); err != nil {
		log.Panicln(err)
	}
	return count
}
