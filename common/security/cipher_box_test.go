package security

import (
	"testing"
)

func TestRsa(t *testing.T) {
	cipher, err := NewRsaCipher()
	if err != nil {
		t.Error("New Failure")
		return
	}
	content := "abc123"
	result, err := cipher.Encrypt([]byte(content))
	if err != nil {
		t.Error("Encrypt Failure", err)
		return
	}
	decryptContent, err := cipher.Decrypt(result)
	if err != nil {
		t.Error("Decrypt Failure", err)
		return
	}
	if string(decryptContent) == content {
		t.Log("OK")
	} else {
		t.Error("Decrypt Failure Content")
	}
}

func TestAes(t *testing.T) {
	cipher := AesCipher{CipherBoxIns.RandomAESKey(), CipherBoxIns.RandomAESIV()}
	content := "abc123"
	result, err := cipher.Encrypt([]byte(content))
	if err != nil {
		t.Error("Encrypt Failure", err)
	}

	decryptContent, err := cipher.Decrypt(result)
	if err != nil {
		t.Error("Decrypt Failure")
	}
	if string(decryptContent) == content {
		t.Log("OK")
	} else {
		t.Error("Decrypt Failure Content")
	}
}
