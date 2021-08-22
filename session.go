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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gopkg.in/olahol/melody.v1"
	"log"
	"runtime"
	"time"
)

const (
	ErrorJSONEncodingResponseData = "JSON_encoding_response_data_error"
	ErrorJSONEncodingResponse     = "JSON_encoding_response_error"
	ErrorGenerateSignature        = "Generate_signature_error"
	ErrorSessionClosed            = "Session_closed_error"
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

// Response will respond data to client.
func (session *Session) Response(data interface{}) {
	// Find original method from Caller.
	skip := 1
	if _, ok := data.(*ErrorReport); ok {
		skip = 2
	}
	pc, _, _, _ := runtime.Caller(skip)
	method := runtime.FuncForPC(pc).Name()
	// Do middlewares [before]
	doMiddlewareBeforeResponse(session, method, data)
	// Encode data into JSON format.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		session.RaiseError(ErrorJSONEncodingResponseData)
		return
	}
	// Let dataBytes compressed by GZip.
	dataBytes = compress(dataBytes)
	// Create a new Response object.
	now := time.Now().UnixNano()
	response := responseFactory(session, now, method, dataBytes)
	// Encode the response into JSON format.
	// JSON package will convert bytes to base64 automatically,
	// so dataBytes with compressed will be encoded into Base64 format.
	responseBytes, err := json.Marshal(response)
	if err != nil {
		session.RaiseError(ErrorJSONEncodingResponse)
		return
	}
	// Flush the response.
	err = session.socketSession.Write(responseBytes)
	if err != nil {
		log.Println(ErrorSessionClosed)
		return
	}
	// Do middlewares [after]
	doMiddlewareAfterResponse(session, method, response)
}

// responseFactory will Generate a Response.
func responseFactory(session *Session, currentTimestamp int64, method string, dataBytes []byte) *Response {
	// Generate Response
	instance := new(Response)
	instance.Data = dataBytes
	instance.Method = method
	instance.Timestamp = currentTimestamp
	instance.Signature = sign(session, currentTimestamp, method, dataBytes)
	return instance
}

// sign will Generate a Signature as hex string by SHA256.
// Due to the data has been turned into compressed bytes,
// there will be no JSON ordering problem while doing the verification.
func sign(session *Session, currentTimestamp int64, method string, dataBytes []byte) string {
	instance := new(Signature)
	instance.Data = dataBytes
	instance.Method = method
	instance.Salt = session.requestSalt
	instance.Timestamp = currentTimestamp
	signatureString, err := json.Marshal(instance)
	if err == nil {
		signatureHash := sha256.Sum256(signatureString)
		return fmt.Sprintf("%x", signatureHash)
	} else {
		session.RaiseError(ErrorGenerateSignature)
		return ErrorGenerateSignature
	}
}

// compress will compress bytes by GZip.
func compress(raw []byte) []byte {
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	if _, err := gz.Write(raw); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return compressed.Bytes()
}

// RaiseError will throw an error to client.
func (session *Session) RaiseError(message string) {
	pc, _, _, _ := runtime.Caller(1)
	method := runtime.FuncForPC(pc).Name()
	log.Printf("[%s] %s\n", method, message)
	session.Response(&ErrorReport{
		Timestamp: time.Now().UnixNano(),
		Error:     message,
	})
}
