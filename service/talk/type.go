package talk

import Kernel "github.com/star-inc/kaguya_kernel"

type ServiceInterface interface {
	Kernel.ServiceInterface
	fetchMessage()
	syncMessageBox()
	getMessageBox(Kernel.Request)
	getMessage(Kernel.Request)
	sendMessage(Kernel.Request)
}

type Message struct {
	ContentType int    `json:"contentType"`
	TargetType  int    `json:"targetType"`
	Origin      string `json:"origin"`
	Target      string `json:"target"`
	Content     []byte `json:"content"`
}
