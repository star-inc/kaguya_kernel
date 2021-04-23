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
	ErrorJSONEncodingResponse = "JSON_encoding_response_error"
	ErrorGenerateSignature    = "Generate_signature_error"
	ErrorSessionClosed        = "Session_closed_error"
)

type Session struct {
	socketSession *melody.Session
	requestSalt   string
}

func NewSession(socketSession *melody.Session, requestSalt string) *Session {
	session := new(Session)
	session.socketSession = socketSession
	session.requestSalt = requestSalt
	return session
}

func (session *Session) Response(data interface{}) {
	// Export as a GZip string with base64
	if data != nil {
		dataString, err := json.Marshal(data)
		if err != nil {
			session.RaiseError(ErrorJSONEncodingResponse)
			return
		}
		data = compress(dataString)
	}
	// Generate Signature
	// Due to the data has been turned into a string,
	// there will be no JSON ordering problem while doing the verification.
	signature := new(Signature)
	signature.Data = data
	signature.Salt = session.requestSalt
	signature.Timestamp = time.Now().UnixNano()
	rawSignatureString, err := json.Marshal(signature)
	if err != nil {
		session.RaiseError(ErrorGenerateSignature)
		return
	}
	signatureString := sha256.Sum256(rawSignatureString)
	// Generate Response
	response := new(Response)
	response.Data = data
	response.Timestamp = signature.Timestamp
	response.Signature = fmt.Sprintf("%x", signatureString)
	responseString, err := json.Marshal(response)
	// Write
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

func compress(raw []byte) string {
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	defer func() {
		if err := gz.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := gz.Write(raw); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(
		compressed.Bytes(),
	)
}

func (session *Session) RaiseError(message string) {
	session.Response(&ErrorReport{
		Timestamp: time.Now().UnixNano(),
		Error:     message,
	})
}
