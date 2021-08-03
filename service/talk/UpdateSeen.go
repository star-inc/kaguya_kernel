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

package talk

import "gopkg.in/star-inc/kaguyakernel.v2/data"

// UpdateSeen: refresh Messagebox by Container.
func UpdateSeen(source *data.RethinkSource, relatedID string, container *data.Container) {
	messagebox := new(data.Messagebox)
	err := messagebox.Load(source, relatedID)
	if err != nil {
		panic(err)
	}
	messagebox.Target = source.Table
	messagebox.Origin = container.Message.Origin
	messagebox.CreatedTime = container.CreatedTime
	messagebox.LastSeen = container.CreatedTime
	if messagebox.CheckReady() {
		err = messagebox.Replace(source)
		if err != nil {
			panic(err)
		}
	} else {
		err = messagebox.Create(source)
		if err != nil {
			panic(err)
		}
	}
}
