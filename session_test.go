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
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
)

func Test_compress(t *testing.T) {
	data := []byte("test")
	compressed := b64Compress(data)
	if target := []byte("H4sIAAAAAAAA/ypJLS4BAAAA//8BAAD//wx+f9gEAAAA"); !bytes.Equal(compressed, target) {
		t.Fatalf("\n%s\nis not equal to\n%s", compressed, target)
	}
}

func Test_NewSignature(t *testing.T) {
	signature := NewSignature(&Session{}, 1629552882143314889, "Test_sign", []byte("test"))
	target := "fd159afc02f3a88985ff1da2600d9c4a5b28a8fa792d9d0607e62936e8faae34"
	if hex, _ := signature.JSONHashHex(); hex != target {
		t.Fatalf("\n%s\nis not equal to\n%s", hex, target)
	}
}

func Test_NewResponse(t *testing.T) {
	session := &Session{}
	method := "Test_NewResponse"
	data := []byte("test")
	// Generated
	response := NewResponse(session, method, data)
	// Raw
	jsonData, _ := json.Marshal(data)
	compressed := compress(jsonData)
	raw := &Response{Data: compressed, Signature: response.Signature, Timestamp: response.Timestamp, Method: method}
	if !reflect.DeepEqual(response, raw) {
		t.Fatalf("\n%#v\nis not equal to\n%#v", response, raw)
	}
}

func b64Compress(raw []byte) []byte {
	compressedBytes := compress(raw)
	compressedResult := make([]byte, base64.StdEncoding.EncodedLen(len(compressedBytes)))
	base64.StdEncoding.Encode(compressedResult, compressedBytes)
	return compressedResult
}
