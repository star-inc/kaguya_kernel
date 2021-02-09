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
	Content     string `json:"content"`
}

type DatabaseMessage struct {
	UUID        string `rethinkdb:"id,omitempty" json:"uuid"`
	Message     []byte `rethinkdb:"message" json:"message"`
	CreatedTime int64  `rethinkdb:"created_time" json:"created_time"`
}
