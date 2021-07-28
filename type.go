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

type Request struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

type Response struct {
	Data      interface{} `json:"data"`
	Signature string      `json:"signature"`
	Timestamp int64       `json:"timestamp"`
}

type ErrorReport struct {
	Error     string `json:"error"`
	Timestamp int64  `json:"timestamp"`
}

type Signature struct {
	Data      interface{} `json:"data"`
	Salt      string      `json:"salt"`
	Timestamp int64       `json:"timestamp"`
}
