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
	Kernel "gopkg.in/star-inc/kaguyakernel.v1"
)

type ServiceInterface interface {
	Kernel.ServiceInterface
	SyncMessagebox(*Kernel.Request)
	DeleteMessagebox(*Kernel.Request)
}

type SyncMessagebox struct {
	Messagebox
	ExtraData interface{} `json:"extraData"`
}

type Messagebox struct {
	CreatedTime int64       `rethinkdb:"createdTime" json:"createdTime"`
	LastSeen    int64       `rethinkdb:"lastSeen" json:"lastSeen"`
	Metadata    interface{} `rethinkdb:"metadata" json:"metadata"`
	Origin      string      `rethinkdb:"origin" json:"origin"`
	Target      string      `rethinkdb:"target" json:"target"`
}

type MessageboxForNew struct {
	CreatedTime int64       `rethinkdb:"createdTime"`
	Metadata    interface{} `rethinkdb:"metadata"`
	Origin      string      `rethinkdb:"origin"`
	Target      string      `rethinkdb:"target"`
}

type MessageboxForUpdateSeen struct {
	LastSeen int64  `rethinkdb:"lastSeen"`
	Target   string `rethinkdb:"target"`
}
