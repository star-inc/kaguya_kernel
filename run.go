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

// Run: to execute kernel with specific arguments.
func Run(service ServiceInterface, guard AuthorizeInterface, middlewares MiddlewareInterface, requestSalt string) *melody.Melody {
	worker := melody.New()
	fetchCtx, fetchCancel := context.WithCancel(context.Background())
	worker.HandleConnect(func(socketSession *melody.Session) {
		session := NewSession(socketSession, middlewares, requestSalt)
		service.SetGuard(guard)
		service.SetSession(session)
		connectHandler(service, fetchCtx)
	})
	worker.HandleMessage(func(_ *melody.Session, message []byte) {
		messageHandler(service, message)
	})
	worker.HandleDisconnect(func(_ *melody.Session) {
		disconnectHandler(fetchCancel)
	})
	return worker
}

func connectHandler(service ServiceInterface, fetchCtx context.Context) {
	if !service.CheckPermission() {
		defer func() {
			err := service.GetSession().socketSession.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		service.GetSession().RaiseError(ErrorForbidden)
		return
	}
	go service.Fetch(fetchCtx)
}

func messageHandler(service ServiceInterface, message []byte) {
	request := new(Request)
	// Decode JSON into Request.
	err := json.Unmarshal(message, request)
	if err != nil {
		service.GetSession().RaiseError(ErrorJSONDecodingRequest)
		return
	}
	// Check method requested is exists.
	method := reflect.ValueOf(service).MethodByName(request.Type)
	if !service.CheckRequestType(method) {
		service.GetSession().RaiseError(ErrorInvalidRequestType)
		return
	}
	// Check method requested is valid (can be requested by client).
	if method.IsValid() {
		// Do middlewares [before]
		doMiddlewareBeforeRequest(service.GetSession(), request)
		// Do main
		method.Call([]reflect.Value{reflect.ValueOf(request)})
		// Do middlewares [after]
		doMiddlewareAfterRequest(service.GetSession(), request)
	}
}

func disconnectHandler(fetchCancel context.CancelFunc) {
	fetchCancel()
}
