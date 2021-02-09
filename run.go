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

import (
	"encoding/json"
	"gopkg.in/olahol/melody.v1"
	"reflect"
)

func Run(service ServiceInterface) *melody.Melody {
	worker := melody.New()
	worker.HandleConnect(func(socketSession *melody.Session) {
		service.SetSession(NewSession(socketSession))
		go service.Fetch()
	})
	worker.HandleMessage(func(socketSession *melody.Session, message []byte) {
		request := new(Request)
		err := json.Unmarshal(message, request)
		if err != nil {
			panic(err)
		}
		method := reflect.ValueOf(service).MethodByName(request.Type)
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(request)})
		}
	})
	return worker
}
