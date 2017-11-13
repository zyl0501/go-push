package security

import (
	"crypto/rsa"
	"crypto/rand"
)

type RsaCipher struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (cipher *RsaCipher) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, cipher.privateKey, data)
}
func (cipher *RsaCipher) Encrypt(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, cipher.publicKey, data)
}

func NewRsaCipher() (*RsaCipher, error) {
	pub, err := CipherBoxIns.PublicKey()
	if err != nil {
		return nil, err
	}
	pri, err := CipherBoxIns.PrivateKey()
	if err != nil {
		return nil, err
	}
	return &RsaCipher{privateKey: pri, publicKey: pub}, nil
}
