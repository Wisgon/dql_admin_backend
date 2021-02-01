package utils

import (
	"github.com/buger/jsonparser"
)

// 用jsonparser算字符串的json中的数组的元素个数
func CountJsonArray(jsonString []byte, arg ...string) (count int, err error) {
	_, err = jsonparser.ArrayEach(jsonString, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		count += 1
	}, arg...)
	return
}
