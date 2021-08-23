package KaguyaKernel

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Signature struct {
	Data      []byte `json:"data"`
	Salt      string `json:"salt"`
	Timestamp int64  `json:"timestamp"`
	Method    string `json:"method,omitempty"`
}

// NewSignature will Generate a Signature as hex string by SHA256.
// Due to the data has been turned into compressed bytes,
// there will be no JSON ordering problem while doing the verification.
func NewSignature(session *Session, currentTimestamp int64, method string, dataBytes []byte) *Signature {
	instance := new(Signature)
	instance.Data = dataBytes
	instance.Method = method
	instance.Salt = session.requestSalt
	instance.Timestamp = currentTimestamp
	return instance
}

// JSONHashHex will generate a hex after hashing the json signature.
func (s *Signature) JSONHashHex() (string, error) {
	signatureString, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	signatureHash := sha256.Sum256(signatureString)
	return fmt.Sprintf("%x", signatureHash), nil
}
