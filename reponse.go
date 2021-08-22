package KaguyaKernel

// Response
// Method is a mark of the origin method,
// to declare where is the Response sent, the field is omitempty.
type Response struct {
	Data      []byte `json:"data"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
	Method    string `json:"method,omitempty"`
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
