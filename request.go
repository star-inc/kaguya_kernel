package KaguyaKernel

// Request the data structure for receiving from clients.
// Type is the method name that the client requested to server,
// kernel will to the reflection and find the method in the ServiceInterface,
// if the only argument of the method is Request, the method will be executed and returned,
// otherwise the request will be denied.
type Request struct {
	Processed bool
	Data      interface{} `json:"data"`
	Type      string      `json:"type"`
}
