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
	"crypto/sha256"
	"encoding/json"
	"gopkg.in/olahol/melody.v1"
	"log"
	"time"
)

const (
	ErrorJSONEncodingResponse = "JSON_encoding_response_error"
	ErrorGenerateSignature    = "Generate_signature_error"
	ErrorSessionClosed        = "Session_closed_error"
)

type Session struct {
	socketSession *melody.Session
}

func NewSession(socketSession *melody.Session) *Session {
	session := new(Session)
	session.socketSession = socketSession
	return session
}

func (session *Session) Response(data interface{}) {
	response := new(Response)
	response.Data = data
	response.Timestamp = time.Now().UnixNano()
	hashString, err := json.Marshal(data)
	if err != nil {
		session.RaiseError(ErrorGenerateSignature)
		return
	}
	response.Signature = sha256.Sum256(hashString)
	responseString, err := json.Marshal(response)
	if err != nil {
		session.RaiseError(ErrorJSONEncodingResponse)
		return
	}
	err = session.socketSession.Write(responseString)
	if err != nil {
		log.Println(ErrorSessionClosed)
		return
	}
}

func (session *Session) RaiseError(message string) {
	session.Response(&ErrorReport{
		Timestamp: time.Now().UnixNano(),
		Error:     message,
	})
}
