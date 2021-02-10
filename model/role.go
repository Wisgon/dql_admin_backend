package model

import (
	"context"
	"dql_admin_backend/config"
	"dql_admin_backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
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

func GetAccessablePages(rolesString []string) (accessablePages map[string]map[string]bool, err error) {
	accessablePages = make(map[string]map[string]bool)
	filters := ""
	for _, roleString := range rolesString {
		filters += "eq(role_id,\"" + roleString + "\") or "
	}
	filterByte := []byte(filters)
	filters = string(filterByte[:len(filterByte)-4]) // 去掉最后一个 or
	query := fmt.Sprintf(`{
		roles(func:type("Role")) @filter(%s) {
			uid
			role_id
			accessable_pages
		}
	}`, filters)
	roles, err := getRoles(query)
	if err != nil {
		return
	}

	for _, role := range roles.Roles {
		accessablePages[role.RoleID] = make(map[string]bool)
		for _, ap := range role.AccessablePages {
			accessablePages[role.RoleID][ap] = true
		}
	}
	return
}

func DoEdit(role Role) error {
	var ctx = context.Background()
	setPagesMutation := ""
	for _, page := range role.AccessablePages {
		setPagesMutation += "<" + role.UID + "> " + "<accessable_pages> \"" + page + "\" .\n"
	}

	nowTime := strconv.Itoa(int(time.Now().Unix()))
	version := utils.GetMd5(nowTime)
	updatePermissionVersionMutation := "<" + config.SystemConfigNodeId + "> <permission_version> \"" + version + "\" .\n"

	resp, err := MutationSetWithUpsert(ctx, []string{setPagesMutation, updatePermissionVersionMutation}, "")
	if err != nil {
		log.Println("add accessable_pages error:" + err.Error() + " setPagesMutation:" + setPagesMutation)
		return err
	}
	fmt.Println("resp:", resp)
	// 成功的话，resp.Json是没有东西的
	_ = resp.Uids
	return nil
}

// ================below is useful function
func getRoles(query string) (roles RolesStru, err error) {
	resp, err := Query(context.Background(), query)
	if err != nil {
		log.Println("query roles error: " + err.Error() + "  query:" + query)
		return
	}
	// fmt.Println("resp:", string(resp.Json))
	err = json.Unmarshal(resp.Json, &roles)
	if err != nil {
		log.Println("parse roles json error:" + err.Error())
		return
	}
	if len(roles.Roles) == 0 {
		err = errors.New("role not found!")
		log.Println("no roles found, query:" + query)
		return
	}
	return
}
