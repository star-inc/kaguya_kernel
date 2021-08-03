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
	Kernel "gopkg.in/star-inc/kaguyakernel.v2"
)

type MessageboxSource struct {
	Base
	ClientID string
}

// NewMessageboxSource: create a new Source instance to connect RethinkDB Server for Messagebox.
func NewMessageboxSource(config Kernel.RethinkConfig) (*Base, error) {
	var err error
	instance := new(Base)
	instance.Term = rethinkdb.DB(config.DatabaseName)
	instance.Session, err = rethinkdb.Connect(config.ConnectConfig)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (s *MessageboxSource) GetFetchCursor() *rethinkdb.Cursor {
	return s.GetRawFetchCursor(s.ClientID)
}
