package KaguyaKernel

import (
	"encoding/json"
	"time"
)

// Response
// Method is a mark of the origin method,
// to declare where is the Response sent, the field is omitempty.
type Response struct {
	Data      []byte `json:"data"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
	Method    string `json:"method,omitempty"`
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
	instance.Timestamp = currentTimestamp
	signature := NewSignature(session, currentTimestamp, method, instance.Data)
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
