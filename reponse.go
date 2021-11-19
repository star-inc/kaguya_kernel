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
	"encoding/json"
	"time"
)

// Response
// Method is a mark of the origin method,
// to declare where is the Response sent, the field is omitempty.
type Response struct {
	Data      []byte        `json:"data"`
	Signature string        `json:"signature"`
	Timestamp time.Duration `json:"timestamp"`
	Method    string        `json:"method,omitempty"`
}

// NewResponse will generate a Response.
func NewResponse(session *Session, method string, data interface{}) *Response {
	// Generate Response
	instance := new(Response)
	// Get Current Timestamp
	currentTimestamp := time.Now().UnixNano()
	// If data is nil, ignore to compress.
	if data != nil {
		// Encode data into JSON format.
		dataBytes, err := json.Marshal(data)
		if err != nil {
			session.RaiseError(ErrorJSONEncodingResponseData)
		}
		// Let dataBytes compressed by GZip.
		instance.Data = compress(dataBytes)
	} else {
		// Set return as nil.
		instance.Data = nil
	}
	instance.Method = method
	instance.Timestamp = time.Duration(currentTimestamp)
	signature := NewSignature(session, time.Duration(currentTimestamp), method, instance.Data)
	hashHex, err := signature.JSONHashHex()
	if err != nil {
		session.RaiseError(ErrorGenerateSignature)
	}
	instance.Signature = hashHex
	return instance
}

// JSON will stringify the response into JSON format.
func (r *Response) JSON() ([]byte, error) {
	val, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return val, nil
}
