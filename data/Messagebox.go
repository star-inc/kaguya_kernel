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
	"log"
)

type Messagebox struct {
	CreatedTime int64       `rethinkdb:"createdTime" json:"createdTime"`
	LastSeen    int64       `rethinkdb:"lastSeen" json:"lastSeen"`
	Metadata    interface{} `rethinkdb:"metadata" json:"metadata"`
	Origin      string      `rethinkdb:"origin" json:"origin"`
	Target      string      `rethinkdb:"target" json:"target"`
}

// NewMessagebox: ToDO
func NewMessagebox() Interface {
	instance := new(Messagebox)
	return instance
}

// Load: ToDO
func (m *Messagebox) Load(source *RethinkSource, filter ...interface{}) error {
	cursor, err := source.Term.Table(source.Table).Get(filter[0].(string)).Run(source.Session)
	if err != nil {
		return err
	}
	defer func() {
		err := cursor.Close()
		log.Println(err)
	}()
	return cursor.One(m)
}

// Reload: ToDO
func (m *Messagebox) Reload(source *RethinkSource) error {
	return m.Load(source, m.Target)
}

// Create: ToDO
func (m *Messagebox) Create(source *RethinkSource) error {
	return source.Term.Table(source.Table).Insert(m).Exec(source.Session)
}

// Update: ToDO
func (m *Messagebox) Update(source *RethinkSource) error {
	return source.Term.Table(source.Table).Update(m).Exec(source.Session)
}

// Replace: ToDO
func (m *Messagebox) Replace(source *RethinkSource) error {
	return source.Term.Table(source.Table).Replace(m).Exec(source.Session)
}

// Destroy: ToDO
func (m *Messagebox) Destroy(source *RethinkSource) error {
	return source.Term.Table(source.Table).Get(m.Target).Delete().Exec(source.Session)
}
