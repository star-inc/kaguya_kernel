package talk

import Kernel "github.com/star-inc/kaguya_kernel"

type ServiceInterface interface {
	Kernel.ServiceInterface
	SyncMessageBox()
	GetMessageBox(*Kernel.Request)
	GetMessage(*Kernel.Request)
	SendMessage(*Kernel.Request)
}

type Message struct {
	ContentType int    `json:"contentType"`
	TargetType  int    `json:"targetType"`
	Origin      string `json:"origin"`
	Target      string `json:"target"`
	Content     []byte `json:"content"`
}

type DatabaseMessage struct {
	UUID      string   `rethinkdb:"id,omitempty"`
	Message   *Message `rethinkdb:"message"`
	Timestamp int64    `rethinkdb:"timestamp"`
}
