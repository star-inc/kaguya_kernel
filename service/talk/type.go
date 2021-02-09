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
	ContentType int    `rethinkdb:"content_type" json:"contentType"`
	TargetType  int    `rethinkdb:"target_type" json:"targetType"`
	Origin      string `rethinkdb:"origin" json:"origin"`
	Target      string `rethinkdb:"target" json:"target"`
	Content     string `rethinkdb:"content" json:"content"`
}

type DatabaseMessage struct {
	UUID        string   `rethinkdb:"id,omitempty" json:"uuid"`
	Message     *Message `rethinkdb:"message" json:"message"`
	CreatedTime int64    `rethinkdb:"created_time" json:"created_time"`
}
