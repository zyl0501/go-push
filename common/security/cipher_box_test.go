package security

import "testing"

func TestRsa(t *testing.T) {
	cipher, err := NewRsaCipher()
	if err != nil {
		t.Error("New Failure")
	}
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
