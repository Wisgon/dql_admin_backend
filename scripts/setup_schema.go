//这个scripts包的作用是存放各种脚本
package main

import (
	"context"
	"dql_admin_backend/model"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	var Ctx = context.Background()
	content, err := ioutil.ReadFile("./schema")
	if err != nil {
		log.Println("read file error: " + err.Error())
	}
	fmt.Println(string(content))
	err = model.SetupSchema(Ctx, string(content))
	if err != nil {
		log.Println("setup error: " + err.Error())
	} else {
		log.Println("Set schema success")
	}
}
