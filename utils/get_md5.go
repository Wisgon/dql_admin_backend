package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5(raw string) (result string) {
	h := md5.New()
	h.Write([]byte(raw)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
