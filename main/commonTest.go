package main

import (
	"dql_admin_backend/model"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	// date := utils.GetTimeString("date_and_time")
	// fmt.Println("date:", date)

	// test json parse
	jsonStr := `
	{"users":[{"uid":"0x271b","password":"123456","roles":[{"role_id":"admin"},{"role_id":"editor"}]}]}
	`
	users := model.UsersStru{}
	err := json.Unmarshal([]byte(jsonStr), &users)
	if err != nil {
		log.Println("err:" + err.Error())
		return
	}
	fmt.Println("p:", users)
}
