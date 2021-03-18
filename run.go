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
package KaguyaKernel

import (
	"context"
	"encoding/json"
	"gopkg.in/olahol/melody.v1"
	"log"
	"reflect"
)

const (
	ErrorJSONDecodingRequest = "JSON_decoding_request_error"
	ErrorInvalidRequestType  = "Request_type_is_invalid"
)

func Run(service ServiceInterface, guard AuthorizeInterface, requestSalt string) *melody.Melody {
	worker := melody.New()
	ctx, cancel := context.WithCancel(context.Background())
	worker.HandleConnect(func(socketSession *melody.Session) {
		service.SetSession(NewSession(socketSession, requestSalt))
		service.SetGuard(guard)
		if !service.CheckPermission() {
			defer func() {
				err := socketSession.Close()
				if err != nil {
					log.Println(err)
				}
			}()
			service.GetSession().RaiseError(ErrorForbidden)
			return
		}
		go service.Fetch(ctx)
	})
	worker.HandleMessage(func(socketSession *melody.Session, message []byte) {
		request := new(Request)
		err := json.Unmarshal(message, request)
		if err != nil {
			service.GetSession().RaiseError(ErrorJSONDecodingRequest)
			return
		}
		method := reflect.ValueOf(service).MethodByName(request.Type)
		if !service.CheckRequestType(method) {
			service.GetSession().RaiseError(ErrorInvalidRequestType)
			return
		}
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(request)})
		}
	})
	worker.HandleDisconnect(func(_ *melody.Session) {
		cancel()
	})
	return worker
}
