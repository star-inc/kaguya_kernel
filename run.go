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
	"gopkg.in/olahol/melody.v1"
	"reflect"
)

func Run(service *ServiceInterface, request Request) *melody.Melody {
	worker := melody.New()
	worker.HandleMessage(func(s *melody.Session, msg []byte) {
		reflect.ValueOf(service).MethodByName(request.Type).Call([]reflect.Value{})
	})
	return worker
}
