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

/**
RDF N-Quad format.
Each triple has the form:
<subject> <predicate> <object> .
*/
func CombineNQuad(uid string, predicate string, value string, valueType string) string {
	switch valueType {
	case "string":
		return "<" + uid + "> <" + predicate + "> \"" + value + "\" .\n"
	case "deleteAll":
		return "<" + uid + "> * * .\n"
	default:
		return "<" + uid + "> <" + predicate + "> " + value + " .\n"
	}
}
