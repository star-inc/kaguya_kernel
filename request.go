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

package KaguyaKernel

// Request the data structure for receiving from clients.
// Type is the method name that the client requested to server,
// kernel will to the reflection and find the method in the ServiceInterface,
// if the only argument of the method is Request, the method will be executed and returned,
// otherwise the request will be denied.
type Request struct {
	Processed bool
	Data      interface{} `json:"data"`
	Type      string      `json:"type"`
}
