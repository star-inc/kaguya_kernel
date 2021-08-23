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
	KernelSource "gopkg.in/star-inc/kaguyakernel.v2/source"
	"time"
)

// RefreshMessageboxAfterSent will refresh Messagebox after sent a Message.
// target is the relation ID, used for getting the room, as known as chat room ID.
func RefreshMessageboxAfterSent(source *KernelSource.MessageboxSource, target string, message *data.Message, metadata string) {
	messagebox := new(data.Messagebox)
	if err := messagebox.Load(source, target); err != nil {
		panic(err)
	}
	if !messagebox.CheckReady() {
		messagebox.Target = target
	}
	messagebox.Origin = message.Origin
	messagebox.CreatedTime = time.Now().UnixNano()
	messagebox.Metadata = metadata
	if messagebox.CheckReady() {
		if err := messagebox.Replace(source); err != nil {
			panic(err)
		}
	} else {
		messagebox.Target = target
		if err := messagebox.Create(source); err != nil {
			panic(err)
		}
	}
}
