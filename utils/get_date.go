package utils

import (
	"strings"
	"time"
)

func GetTimeString(mode string) string {
	// 整个字符串是这样的：2021-02-05 16:33:26.278697399 +0800 CST m=+0.000334755
	switch mode {
	case "date_and_time":
		return string([]byte(time.Now().String())[:19])
	case "date":
		return string([]byte(time.Now().String())[:10])
	case "time":
		return string([]byte(time.Now().String())[11:19])
	default:
		return ""
	}
}

func ChangeTimeFormat(mode string, timeStr string) string {
	// 在普通日期时间格式和dql要求的日期时间格式之间转换
	// 普通格式： 2020-06-03 12:33:11
	// dql要求的格式： 2020-06-03T12:33:11Z
	switch mode {
	case "dql2normal":
		noT := strings.Replace(timeStr, "T", " ", 1)
		noZ := strings.Replace(noT, "Z", "", 1)
		return noZ
	case "normal2dql":
		return strings.Replace(timeStr, " ", "T", 1)
	default:
		return ""
	}
}
