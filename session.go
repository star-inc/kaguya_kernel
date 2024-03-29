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
	"bytes"
	"compress/gzip"
	"errors"
	"gopkg.in/olahol/melody.v1"
	"gopkg.in/star-inc/kaguyakernel.v2/time"
	"log"
	"runtime"
)

var (
	ErrorJSONEncodingResponseData = errors.New("json_encoding_response_data_error")
	ErrorJSONEncodingResponse     = errors.New("json_encoding_response_error")
	ErrorGenerateSignature        = errors.New("generate_signature_error")
	ErrorSessionClosed            = errors.New("session_closed_error")
)

type Session struct {
	socketSession *melody.Session
	requestSalt   string
	middlewares   MiddlewareInterface
}

// NewSession will start a new session.
func NewSession(socketSession *melody.Session, middlewares MiddlewareInterface, requestSalt string) *Session {
	session := new(Session)
	session.socketSession = socketSession
	session.requestSalt = requestSalt
	session.middlewares = middlewares
	return session
}

// Respond will send data to client.
func (session *Session) Respond(data interface{}) {
	// Find original method from Caller.
	skip := 1
	if _, ok := data.(*ErrorReport); ok {
		skip = 2
	}
	pc, _, _, _ := runtime.Caller(skip)
	method := runtime.FuncForPC(pc).Name()
	// Do middlewares [before]
	doMiddlewareBeforeRespond(session, method, data)
	// Create a new Response object.
	response := NewResponse(session, method, data)
	// Encode the response into JSON format.
	// JSON package will convert bytes to base64 automatically,
	// so dataBytes with compressed will be encoded into Base64 format.
	responseBytes, err := response.JSON()
	if err != nil {
		session.RaiseError(ErrorJSONEncodingResponse)
		return
	}
	// Flush the response.
	if err = session.socketSession.Write(responseBytes); err != nil {
		log.Println(ErrorSessionClosed)
		return
	}
	// Do middlewares [after]
	doMiddlewareAfterRespond(session, method, response)
}

// compress bytes by GZip.
func compress(raw []byte) []byte {
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	if _, err := gz.Write(raw); err != nil {
		log.Panicln(err)
	}
	if err := gz.Flush(); err != nil {
		log.Panicln(err)
	}
	if err := gz.Close(); err != nil {
		log.Panicln(err)
	}
	return compressed.Bytes()
}

// RaiseError will throw an error to client.
func (session *Session) RaiseError(message error) {
	pc, _, _, _ := runtime.Caller(1)
	method := runtime.FuncForPC(pc).Name()
	log.Printf("[%s] %s\n", method, message)
	session.Respond(&ErrorReport{
		Timestamp: time.Now(),
		Error:     message.Error(),
	})
}
