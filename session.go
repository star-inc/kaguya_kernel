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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/olahol/melody.v1"
	"log"
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

// NewSession: starts a new session.
func NewSession(socketSession *melody.Session, middlewares MiddlewareInterface, requestSalt string) *Session {
	session := new(Session)
	session.socketSession = socketSession
	session.requestSalt = requestSalt
	session.middlewares = middlewares
	return session
}

// Response: response a data to client.
func (session *Session) Response(data interface{}) {
	// Do middlewares [before]
	doMiddlewareBeforeResponse(session, data)
	// Encode data into JSON format.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		session.RaiseError(ErrorJSONEncodingResponseData)
		return
	}
	// Let dataBytes compressed by GZip, and encoded by Base64.
	dataBytes = compress(dataBytes)
	// Create a new Response object.
	now := time.Now().UnixNano()
	response := responseFactory(session, now, dataBytes)
	// Encode the response into JSON format.
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
	doMiddlewareAfterResponse(session, response)
}

// responseFactory: Generate a Response.
func responseFactory(session *Session, currentTimestamp int64, dataBytes []byte) *Response {
	// Generate Response
	instance := new(Response)
	instance.Data = dataBytes
	instance.Timestamp = currentTimestamp
	instance.Signature = sign(session, currentTimestamp, dataBytes)
	return instance
}

// sign: Generate a Signature as hex string by SHA256.
// Due to the data has been turned into compressed bytes,
// there will be no JSON ordering problem while doing the verification.
func sign(session *Session, currentTimestamp int64, dataBytes []byte) string {
	instance := new(Signature)
	instance.Data = dataBytes
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

// compress: compress bytes by GZip, and encode into Base64 format.
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
	compressedBytes := compressed.Bytes()
	compressedResult := make([]byte, base64.StdEncoding.EncodedLen(len(compressedBytes)))
	base64.StdEncoding.Encode(compressedResult, compressedBytes)
	return compressedResult
}

// RaiseError: throw an error to client.
func (session *Session) RaiseError(message string) {
	session.Response(&ErrorReport{
		Timestamp: time.Now().UnixNano(),
		Error:     message,
	})
}
