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

// sign will Generate a Signature as hex string by SHA256.
// Due to the data has been turned into compressed bytes,
// there will be no JSON ordering problem while doing the verification.
func sign(session *Session, currentTimestamp int64, method string, dataBytes []byte) string {
	instance := new(Signature)
	instance.Data = dataBytes
	instance.Method = method
	instance.Salt = session.requestSalt
	instance.Timestamp = currentTimestamp
	signatureString, err := json.Marshal(instance)
	if err == nil {
		signatureHash := sha256.Sum256(signatureString)
		return fmt.Sprintf("%x", signatureHash)
	} else {
		session.RaiseError(ErrorGenerateSignature)
		return ErrorGenerateSignature
	}
}
