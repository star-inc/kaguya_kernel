package data

type SyncMessagebox struct {
	Messagebox
	ExtraData interface{} `json:"extraData"`
}
