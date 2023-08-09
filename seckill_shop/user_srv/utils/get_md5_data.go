package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5pwd(pwd string) string {
	h := md5.New()
	md5code := "12ffjkjru98hsjq%*"
	h.Write([]byte(pwd + md5code))
	return hex.EncodeToString(h.Sum(nil))
}


