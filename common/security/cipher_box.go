package security

import (
	"crypto/rsa"
	"github.com/zyl0501/go-push/tools/config"
)

var (
	CipherBoxIns = CipherBox{AesKeyLength: config.Aes_key_length}
)

type CipherBox struct {
	AesKeyLength int
	PrivateKey   rsa.PrivateKey
	PublicKey    rsa.PublicKey
}

func RandomAESKey()([]byte){
	return nil
}
