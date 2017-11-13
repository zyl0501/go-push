package security

import (
	"crypto/aes"
	"crypto/cipher"
)

type AesCipher struct {
	Key []byte
	Iv []byte
}
func (ac *AesCipher)Decrypt(data []byte) ([]byte, error){
	var err error
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	key := ac.Key
	iv := ac.Iv
	decrypted := make([]byte, len(data))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, data)
	return decrypted, nil
}
func (ac *AesCipher)Encrypt(data []byte) ([]byte, error){
	key := ac.Key
	iv := ac.Iv
	encrypted := make([]byte, len(data))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(data))
	return encrypted, nil
}