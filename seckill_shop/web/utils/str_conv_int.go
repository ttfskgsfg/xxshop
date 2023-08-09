package utils

import (
	"unsafe"
	"strconv"
)

func StrToInt(str string) int32  {

	value_int64,_ := strconv.ParseInt(str,10,64)
	value_int := *(*int32)(unsafe.Pointer(&value_int64))
	return value_int
}


func StrToFloat32(str string) float32  {

	value_float64,_ := strconv.ParseFloat(str,64)
	return float32(value_float64)
}
