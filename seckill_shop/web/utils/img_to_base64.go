package utils

import (
	"os"
	"encoding/base64"
)

func Img2Base64(imgPath string) string {
	file,_ := os.Open(imgPath)

	defer file.Close()

	bufByte := make([]byte,100000)
	n,_ := file.Read(bufByte)

	imgBase64Str := base64.StdEncoding.EncodeToString(bufByte[:n])

	return imgBase64Str
	
}
