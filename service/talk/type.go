package talk

import KaguyaKernel "github.com/star-inc/kaguya_kernel"

type ServiceInterface interface {
	KaguyaKernel.ServiceInterface
	fetchMessage()
	syncMessageBox()
	getMessageBox(KaguyaKernel.Request)
	getMessage(KaguyaKernel.Request)
	sendMessage(KaguyaKernel.Request)
}

type Message struct {
	ContentType int    `json:"contentType"`
	TargetType  int    `json:"targetType"`
	Origin      string `json:"origin"`
	Target      string `json:"target"`
	Content     []byte `json:"content"`
}
