package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(content []byte) string{
	h := md5.New()
	h.Write([]byte(content))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}