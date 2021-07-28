package talk

import Kernel "gopkg.in/star-inc/kaguyakernel.v1"

type ServiceInterface interface {
	Kernel.ServiceInterface
	GetHistoryMessages(*Kernel.Request)
	GetMessage(*Kernel.Request)
	SendMessage(*Kernel.Request)
	CancelSentMessage(*Kernel.Request)
}

type Message struct {
	Content     string `rethinkdb:"content" json:"content"`
	ContentType int    `rethinkdb:"contentType" json:"contentType"`
	Origin      string `rethinkdb:"origin" json:"origin"`
}

type DatabaseMessage struct {
	Canceled    bool     `rethinkdb:"canceled" json:"canceled"`
	CreatedTime int64    `rethinkdb:"createdTime" json:"createdTime"`
	Message     *Message `rethinkdb:"message" json:"message"`
	UUID        string   `rethinkdb:"id,omitempty" json:"uuid"`
}
