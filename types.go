/*
Package KaguyaKernel : The kernel for Kaguya

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
package KaguyaKernel

type Request struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type Response struct {
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	Signature [32]byte    `json:"signature"`
}

type ErrorReport struct {
	Timestamp int64  `json:"timestamp"`
	Error     string `json:"error"`
}

type User struct {
	Identity    string `json:"identity"`
	DisplayName string `json:"displayName"`
	Username    string `json:"username"`
}
