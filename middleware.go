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

import "sync"

// MiddlewareInterface is the interface to get middlewares for kernel.
type MiddlewareInterface interface {
	OnRequestBefore() []func(session *Session, wg *sync.WaitGroup, request *Request)
	OnRequestAfter() []func(session *Session, wg *sync.WaitGroup, request *Request)
	OnResponseBefore() []func(session *Session, wg *sync.WaitGroup, method string, data interface{})
	OnResponseAfter() []func(session *Session, wg *sync.WaitGroup, method string, response *Response)
}

func doMiddlewareBeforeRequest(session *Session, request *Request) {
	if middlewares := session.middlewares.OnRequestBefore(); middlewares != nil {
		wg := new(sync.WaitGroup)
		wg.Add(len(middlewares))
		for _, middleware := range middlewares {
			go middleware(session, wg, request)
		}
		wg.Wait()
	}
}

func doMiddlewareAfterRequest(session *Session, request *Request) {
	if middlewares := session.middlewares.OnRequestAfter(); middlewares != nil {
		wg := new(sync.WaitGroup)
		wg.Add(len(middlewares))
		for _, middleware := range middlewares {
			go middleware(session, wg, request)
		}
		wg.Wait()
	}
}

func doMiddlewareBeforeRespond(session *Session, method string, data interface{}) {
	if middlewares := session.middlewares.OnResponseBefore(); middlewares != nil {
		wg := new(sync.WaitGroup)
		wg.Add(len(middlewares))
		for _, middleware := range middlewares {
			go middleware(session, wg, method, data)
		}
		wg.Wait()
	}
}

func doMiddlewareAfterRespond(session *Session, method string, response *Response) {
	if middlewares := session.middlewares.OnResponseAfter(); middlewares != nil {
		wg := new(sync.WaitGroup)
		wg.Add(len(middlewares))
		for _, middleware := range middlewares {
			go middleware(session, wg, method, response)
		}
		wg.Wait()
	}
}
