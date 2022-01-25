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
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"gopkg.in/star-inc/kaguyakernel.v2/time"
	"log"
)

type Messagebox struct {
	CreatedTime time.NanoTime `rethinkdb:"createdTime" json:"createdTime"`
	LastSeen    time.NanoTime `rethinkdb:"lastSeen" json:"lastSeen"`
	Metadata    interface{}   `rethinkdb:"metadata" json:"metadata"`
	Origin      string        `rethinkdb:"origin" json:"origin"`
	Target      string        `rethinkdb:"target" json:"target"`
}

// NewMessagebox ToDo
func NewMessagebox() Interface {
	instance := new(Messagebox)
	return instance
}

// CheckReady will check model is ready.
func (m *Messagebox) CheckReady() bool {
	return m != nil && m.Origin != "" && m.Target != ""
}

// Load ToDo
func (m *Messagebox) Load(source KernelSource.Interface, filter ...interface{}) error {
	sourceInstance := source.(*KernelSource.MessageboxSource)
	cursor, err := source.GetTerm().Table(sourceInstance.ClientID).Get(filter[0].(string)).Run(source.GetSession())
	if err != nil {
		return err
	}
	defer func() {
		if err := cursor.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	return cursor.One(m)
}

// Reload ToDo
func (m *Messagebox) Reload(source KernelSource.Interface) error {
	return m.Load(source, m.Target)
}

// Create ToDo
func (m *Messagebox) Create(source KernelSource.Interface) error {
	sourceInstance := source.(*KernelSource.MessageboxSource)
	return source.GetTerm().Table(sourceInstance.ClientID).Insert(m).Exec(source.GetSession())
}

// Replace ToDo
func (m *Messagebox) Replace(source KernelSource.Interface) error {
	sourceInstance := source.(*KernelSource.MessageboxSource)
	return source.GetTerm().Table(sourceInstance.ClientID).Replace(m).Exec(source.GetSession())
}

// Destroy ToDo
func (m *Messagebox) Destroy(source KernelSource.Interface) error {
	sourceInstance := source.(*KernelSource.MessageboxSource)
	return source.GetTerm().Table(sourceInstance.ClientID).Get(m.Target).Delete().Exec(source.GetSession())
}
