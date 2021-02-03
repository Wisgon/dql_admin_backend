package main

import (
	"dql_admin_backend/config"
	"dql_admin_backend/route"
	"dql_admin_backend/utils"
	"log"
)

func main() {
	defer func() {
		errMessage := recover()
		if errMessage != nil {
			log.Println("server 异常关闭： " + errMessage.(string))
		}
	}()
	utils.UseRotateLog(config.Root + "/logs")
	log.Println("starting server~~~")
	r := route.Router
	route.Users()
	r.Run(":8063")
}
