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
	ErrorResponseWriting      = "Response_writing_error"
)

type ResponseHandler struct {
	session *melody.Session
}

func NewResponseHandler(session *melody.Session) *ResponseHandler {
	handler := new(ResponseHandler)
	handler.session = session
	return handler
}

func (handler *ResponseHandler) Response(data interface{}) {
	responseContainer := new(Response)
	responseContainer.Timestamp = time.Now().Unix()
	responseContainer.Data = data
	hashString, err := json.Marshal(data)
	if err != nil {
		handler.ErrorRaise(ErrorGenerateSignature)
		return
	}
	responseContainer.Signature = sha256.Sum256(hashString)
	dataString, err := json.Marshal(data)
	if err != nil {
		handler.ErrorRaise(ErrorJSONEncodingResponse)
		return
	}
	err = handler.session.Write(dataString)
	if err != nil {
		log.Panicln(ErrorResponseWriting)
		return
	}
}

func (handler *ResponseHandler) ErrorRaise(message string) {
	handler.Response(&ErrorRaise{Error: message})
}