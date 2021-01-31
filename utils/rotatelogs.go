package utils

import (
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func UseRotateLog(logPath string) error {
	logf, err := rotatelogs.New(
		logPath+"/log.%Y%m%d%H%M",                      // 如果文件夹不存在则自动创建
		rotatelogs.WithLinkName(logPath+"/newest.log"), // 生成软链，指向最新日志文件
		//MaxAge and RotationCount cannot be both set  两者不能同时设置
		rotatelogs.WithMaxAge(48*time.Hour), //clear 最小分钟为单位, 也就是这段时间后，会清理第一个log文件，重头开始
		//rotatelogs.WithRotationCount(5),        //number 默认7份 大于7份 或到了清理时间 开始清理
		rotatelogs.WithRotationTime(time.Hour), //rotate 最小为1分钟轮询。默认60s  低于1分钟就按1分钟来
	)
	if err != nil {
		log.Printf("failed to create rotatelogs: %s", err)
		return err
	}
	log.SetOutput(logf)

	return nil
}
