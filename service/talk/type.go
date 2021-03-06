package talk

import Kernel "github.com/star-inc/kaguya_kernel"

type ServiceInterface interface {
	Kernel.ServiceInterface
	GetHistoryMessages(*Kernel.Request)
	GetMessage(*Kernel.Request)
	SendMessage(*Kernel.Request)
}

type Message struct {
	Origin      string `rethinkdb:"origin" json:"origin"`
	Content     string `rethinkdb:"content" json:"content"`
	ContentType int    `rethinkdb:"contentType" json:"contentType"`
}

type DatabaseMessage struct {
	UUID        string   `rethinkdb:"id,omitempty" json:"uuid"`
	Message     *Message `rethinkdb:"message" json:"message"`
	Canceled    bool     `rethinkdb:"canceled" json:"canceled"`
	CreatedTime int64    `rethinkdb:"createdTime" json:"createdTime"`
}
