package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func GetRolesList(pageSize int, pageNo int) (roles RolesStru, err error) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query := fmt.Sprintf(`{
		roles(func: type(Role), first:%d, offset:%d) {
			uid
			role_id
			name
		}
	}
		
	`, pageSize, pageSize*(pageNo-1))
	roles, err = getRoles(query)
	return
}

// ================below is useful function
func getRoles(query string) (roles RolesStru, err error) {
	resp, err := Query(context.Background(), query)
	// fmt.Printf("user is:%+v\n", users)
	if err != nil {
		log.Println("query users error: " + err.Error())
		return
	}
	//fmt.Println("resp:", string(resp.Json))
	err = json.Unmarshal(resp.Json, &roles)
	if err != nil {
		log.Println("parse users json error:" + err.Error())
		return
	}
	if len(roles.Roles) == 0 {
		err = errors.New("user not found!")
		log.Println("no roles found, query:" + query)
		return
	}
	return
}
