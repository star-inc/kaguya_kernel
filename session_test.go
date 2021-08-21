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
	"reflect"
	"testing"
	"time"
)

func Test_responseFactory(t *testing.T) {
	session := &Session{}
	data := []byte("test")
	currentTimestamp := time.Now().UnixNano()
	response := responseFactory(session, currentTimestamp, data)
	raw := &Response{Data: data, Signature: response.Signature, Timestamp: currentTimestamp}
	if !reflect.DeepEqual(response, raw) {
		t.Fatalf("\n%#v\nis not equal to\n%#v", response, raw)
	}
}

func Test_sign(t *testing.T) {
	signature := sign(&Session{}, 1629552882143314889, []byte("test"))
	if target := "35e267e2d0fe80f8a3acfa02d5ecc648069c38d4a2d937f70a02ac63014610aa"; signature != target {
		t.Fatalf("\n%s\nis not equal to\n%s", signature, target)
	}
}

func Test_compress(t *testing.T) {
	data := []byte("test")
	compressed := compress(data)
	if target := []byte("H4sIAAAAAAAA/ypJLS4BAAAA//8BAAD//wx+f9gEAAAA"); bytes.Compare(compressed, target) != 0 {
		t.Fatalf("\n%s\nis not equal to\n%s", compressed, target)
	}
}
