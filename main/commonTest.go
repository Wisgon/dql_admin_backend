package main

import (
	"dql_admin_backend/utils"
	"fmt"
)

func main() {
	date := utils.GetTimeString("date_and_time")
	fmt.Println("date:", date)
}
