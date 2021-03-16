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
package box

import (
	Kernel "github.com/star-inc/kaguya_kernel"
)

type ServiceInterface interface {
	Kernel.ServiceInterface
	GetHistoryMessagebox(*Kernel.Request)
	DeleteMessagebox(*Kernel.Request)
}

type Messagebox struct {
	CreatedTime  int64  `rethinkdb:"createdTime" json:"createdTime"`
	LastSeenTime int64  `rethinkdb:"lastSeenTime" json:"lastSeenTime"`
	Metadata     string `rethinkdb:"metadata" json:"metadata"`
	Origin       string `rethinkdb:"origin" json:"origin"`
	Target       string `rethinkdb:"target" json:"target"`
}
