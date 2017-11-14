package security

import (
	"crypto/rsa"
	"github.com/zyl0501/go-push/tools/config"
	"math/rand"
	"math"
	"github.com/zyl0501/go-push/tools/utils"
	"encoding/pem"
	"errors"
	"crypto/x509"
)

var (
	CipherBoxIns = CipherBox{AesKeyLength: config.AesKeyLength}
)

type CipherBox struct {
	AesKeyLength int
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
}

func (cb *CipherBox) PublicKey() (*rsa.PublicKey, error) {
	if cb.publicKey == nil {
		block, _ := pem.Decode([]byte(config.PublicKey))
		if block == nil {
			return nil, errors.New("public key error")
		}
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		cb.publicKey = pubInterface.(*rsa.PublicKey)
	}
	return cb.publicKey, nil
}

func (cb *CipherBox) PrivateKey() (*rsa.PrivateKey, error) {
	if cb.privateKey == nil {
		block, _ := pem.Decode([]byte(config.PrivateKey))
		if block == nil {
			return nil, errors.New("private key error!")
		}
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		cb.privateKey = priv
	}
	return cb.privateKey, nil
}

func (cb *CipherBox) RandomAESKey() ([]byte) {
	length := cb.AesKeyLength
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte(rand.Intn(math.MaxInt8))
	}
	return result
}
func (cb *CipherBox) RandomAESIV() ([]byte) {
	length := cb.AesKeyLength
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte(rand.Intn(math.MaxInt8))
	}
	return result
}
func (cb *CipherBox) MixKey(clientKey []byte, serverKey []byte) ([]byte) {
	length := cb.AesKeyLength
	sessionKey := make([]byte, length)
	for i := 0; i < length; i++ {
		a := clientKey[i]
		b := serverKey[i]
		sum := utils.AbsInt(int(a + b))
		var c int
		if sum%2 == 0 {
			c = int(a) ^ int(b)
		} else {
			c = int(b) ^ int(a)
		}
		sessionKey[i] = byte(c);
	}
	return sessionKey;
}
