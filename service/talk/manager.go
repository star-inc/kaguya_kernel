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
	"github.com/star-inc/kaguya_kernel/data"
	Rethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
)

type Manager struct {
	data.Data
}

func NewManager(config Kernel.RethinkConfig, chatRoomID string) *Manager {
}

func (manager *Manager) Check() bool {
}

func (manager *Manager) Create() bool {
}

func (manager *Manager) Drop() bool {
}
