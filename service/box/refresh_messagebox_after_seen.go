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
	"gopkg.in/star-inc/kaguyakernel.v2/time"
	"log"
)

// RefreshMessageboxAfterSeen will refresh Messagebox by fetching a Container or getting the history.
// target is the relation ID, used for getting the room, as known as chat room ID.
// origin is the ID of the Message sender.
func RefreshMessageboxAfterSeen(source *KernelSource.MessageboxSource, target string, origin string) {
	if target == "" || origin == "" {
		log.Panicf("target or origin is not specified. %s %s\n", target, origin)
	}
	messagebox := new(data.Messagebox)
	if err := messagebox.Load(source, target); err != nil {
		log.Panicln(err)
	}
	messagebox.LastSeen = time.Now()
	if messagebox.CheckReady() {
		if err := messagebox.Replace(source); err != nil {
			log.Panicln(err)
		}
	} else {
		messagebox.Target = target
		messagebox.Origin = origin
		messagebox.CreatedTime = time.Now()
		if err := messagebox.Create(source); err != nil {
			log.Panicln(err)
		}
	}
}
