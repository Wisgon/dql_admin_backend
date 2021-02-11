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

func CreateRole(role Role) error {
	var ctx = context.Background()
	nowTime := utils.ChangeTimeFormat("normal2dql", utils.GetTimeString("date_and_time"))
	mutation := fmt.Sprintf(`
		_:r <create_time> "%s" .
		_:r <name> "%s" .
		_:r <role_id> "%s" .
		_:r <dgraph.type> "Role" .
	`, nowTime, role.RoleName, role.RoleID)
	for _, page := range role.AccessablePages {
		mutation += "\n_:r <accessable_pages> \"" + page + "\" .\n"
	}
	resp, err := MutationSet(ctx, mutation)
	if err != nil {
		return err
	}
	// fmt.Println("resp:", resp)
	_ = resp.Uids
	return nil
}

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

func EditRole(role Role) error {
	// 由于dgraph无法一次性设置好数组，所以这里的策略是对于这种整个数组的更改，我们删掉原来的整个数组，然后重新建立数组，删除用 <0xyyy> <accessable_pages> * .
	var ctx = context.Background()

	deleteMutation := utils.CombineNQuad(role.UID, "accessable_pages", "*", "other")

	nowTime := utils.ChangeTimeFormat("normal2dql", utils.GetTimeString("date_and_time"))
	setPagesMutation := utils.CombineNQuad(role.UID, "update_time", nowTime, "string")
	if role.RoleName != "" {
		setPagesMutation += utils.CombineNQuad(role.UID, "name", role.RoleName, "string")
	}
	for _, page := range role.AccessablePages {
		setPagesMutation += utils.CombineNQuad(role.UID, "accessable_pages", page, "string")
	}

	nowUnix := strconv.Itoa(int(time.Now().Unix()))
	version := utils.GetMd5(nowUnix)
	updatePermissionVersionMutation := utils.CombineNQuad(config.SystemConfigNodeId, "permission_version", version, "string")

	dmArray, smArray := []string{deleteMutation}, []string{setPagesMutation, updatePermissionVersionMutation}

	resp, err := MutationDeleteAndSetWithUpsert(ctx, dmArray, smArray, "")
	if err != nil {
		log.Println("add accessable_pages error:" + err.Error() + " setPagesMutation:" + setPagesMutation)
		return err
	}
	// fmt.Println("resp:", resp)
	// 成功的话，resp.Json是没有东西的
	_ = resp.Uids
	return nil
}

/**
example:
upsert{
    query{
		find_nodes(func: uid(0x4e22)){
			~roles{
				linked_uid as uid
			}
		}
	}

	mutation{
		delete{
			uid(linked_uid) <roles> <0x4e22> .
      		<0x4e22> * * .
		}
	}
}
*/
func DeleteRole(role Role) error {
	// 要删除role，删除后要用reverse找到所有与这个node关联的uid然后删除
	var ctx = context.Background()
	// 首先，找到所有与这个节点有关的uid
	query := fmt.Sprintf(`
	query {
		find_linked(func: uid(%s)) {
			~roles {
				linked_nodes as uid
			}
		}
	}
	`, role.UID)

	//然后，组装delete语句，先删edge后删node
	deleteEdgeMutation := "uid(linked_nodes) <roles> <" + role.UID + "> .\n"
	deleteNodeMutation := utils.CombineNQuad(role.UID, "", "", "deleteAll")
	dmArray := []string{deleteEdgeMutation, deleteNodeMutation}

	resp, err := MutationDeleteWithUpsert(ctx, dmArray, query)
	if err != nil {
		log.Printf("delete role error:"+err.Error()+" query:%+v\n", dmArray)
		return err
	}
	// fmt.Println("resp:", resp)
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
