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
	"reflect"
)

type Service struct {
	authorize AuthorizeInterface
	session   *Session
}

type ServiceInterface interface {
	Fetch()
	CheckPermission() bool
	GetGuard() AuthorizeInterface
	SetGuard(authorization AuthorizeInterface)
	GetSession() *Session
	SetSession(session *Session)
	CheckRequestType(method reflect.Value) bool
}

func (service *Service) GetGuard() AuthorizeInterface {
	return service.authorize
}

func (service *Service) SetGuard(authorization AuthorizeInterface) {
	service.authorize = authorization
}

func (service *Service) GetSession() *Session {
	return service.session
}

func (service *Service) SetSession(session *Session) {
	service.session = session
}

func (service *Service) CheckRequestType(method reflect.Value) bool {
	return true
}
