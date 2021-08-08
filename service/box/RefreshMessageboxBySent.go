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
	"gopkg.in/star-inc/kaguyakernel.v2/data"
	"gopkg.in/star-inc/kaguyakernel.v2/source"
)

// RefreshMessageboxBySent: refresh Messagebox by sent a message.
// target: target is the Table name of talk service(UserID), to update the specified row with a container from messageboxes.
func RefreshMessageboxBySent(source *source.MessageboxSource, target string, container *data.Container, relatedMessagebox []string, metadata string) {
	for _, relatedID := range relatedMessagebox {
		source.ClientID = relatedID
		messagebox := new(data.Messagebox)
		err := messagebox.Load(source, target)
		if err != nil {
			panic(err)
		}
		messagebox.Origin = container.Message.Origin
		messagebox.CreatedTime = container.CreatedTime
		messagebox.Metadata = metadata
		if messagebox.CheckReady() {
			err = messagebox.Replace(source)
			if err != nil {
				panic(err)
			}
		} else {
			messagebox.Target = target
			err = messagebox.Create(source)
			if err != nil {
				panic(err)
			}
		}
	}
}
