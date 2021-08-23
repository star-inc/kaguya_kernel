package KaguyaKernel

import (
	"encoding/json"
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
func NewResponse(session *Session, currentTimestamp int64, method string, dataBytes []byte) *Response {
	// Generate Response
	instance := new(Response)
	instance.Data = dataBytes
	instance.Method = method
	instance.Timestamp = currentTimestamp
	signature := NewSignature(session, currentTimestamp, method, dataBytes)
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
