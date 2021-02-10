package model

import (
	"context"
	"dql_admin_backend/config"
	"dql_admin_backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// ==============below is some method of struct User
func (u User) CreateUser() error {
	var ctx = context.Background()
	query := fmt.Sprintf(`
	query{
		q1(func: eq(username, "%s")) {
			un as uid
		}
		q2(func: eq(phone, "%s")) {
			ph as uid
		}
	}
	`, u.UserName, u.PhoneNumber)

	var mutationStrings []map[string]string
	mutationString := make(map[string]string)

	if u.Avatar == "" {
		u.Avatar = "https://i.gtimg.cn/club/item/face/img/2/15922_100.gif"
	}

	// 设置cond
	cond := "@if(eq(len(un), 0) AND eq(len(ph), 0))"
	mutationString["cond"] = cond
	nowTime := utils.ChangeTimeFormat("normal2dql", utils.GetTimeString("date_and_time"))
	// 组装mutation
	mu := fmt.Sprintf(`
	_:new_user <username> "%s" .
	_:new_user <password> "%s" .
	_:new_user <phone> "%s" .
	_:new_user <avatar> "%s" .
	_:new_user <dgraph.type> "User" .
	_:new_user <create_time> "%s" .
	_:new_user <update_time> "%s" .
	_:new_user <roles> <%s> .
	`, u.UserName, u.Password, u.PhoneNumber, u.Avatar, nowTime, nowTime, config.NormalRoleId)
	// 上面的config.NormalRoleId，是当前的普通role的再数据库的id，如果重新建库可能会有不同，要去config.NormalRoleId修改
	mutationString["mutation"] = mu

	// fmt.Printf("mu: %+v\n\n", mutationString)

	mutationStrings = append(mutationStrings, mutationString)

	resp, err := MutationSetWithConditionUpsert(ctx, mutationStrings, query)
	if err != nil {
		log.Println("mutation with upsert error: " + err.Error())
		return err
	}
	// fmt.Printf("resp: %+v /n", resp)
	q1_count, err := utils.CountJsonArray(resp.Json, "q1")
	if err != nil {
		log.Println("parse resp q1 json error: " + err.Error())
		return err
	}
	q2_count, err := utils.CountJsonArray(resp.Json, "q2")
	if err != nil {
		log.Println("parse resp q2 json error: " + err.Error())
		return err
	}

	if q1_count > 0 {
		return errors.New("用户名已被注册")
	}

	if q2_count > 0 {
		return errors.New("手机已被注册")
	}

	return nil
}

func (u *User) GetUserInfo(condition string) error {
	var query string
	switch condition {
	case "id":
		query = fmt.Sprintf(`{
			users(func: uid(%s)) {
				username
				phone
				avatar
				roles{
					role_id
					name
				}
			}
		}`, u.UID)
	case "username":
		query = fmt.Sprintf(`{
			users(func: eq(username, "%s")) {
				username
				phone
				avatar
				roles{
					role_id
					name
				}
			}
		}`, u.UserName)
	case "phone":
		query = fmt.Sprintf(`{
			users(func: eq(phone, "%s")) {
				username
				phone
				avatar
				roles{
					role_id
					name
				}
			}
		}`, u.PhoneNumber)
	default:
		err := errors.New("no such user attr: " + condition)
		return err
	}

	// 用user的struct代替了下面的jsonparser
	users, err := getUsers(query)
	if err != nil {
		return err
	}
	user := users.Users[0] // 第一个就是
	u.UserName, u.Avatar, u.PhoneNumber, u.Roles = user.UserName, user.Avatar, user.PhoneNumber, user.Roles

	// _, err = jsonparser.ArrayEach(resp.Json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 	username, newErr := jsonparser.GetString(value, "username")
	// 	if newErr != nil {
	// 		log.Println("get username error: " + newErr.Error())
	// 	}
	// 	phoneNumber, newErr := jsonparser.GetString(value, "phone")
	// 	if newErr != nil {
	// 		log.Println("get phone error: " + newErr.Error())
	// 	}
	// 	avatar, newErr := jsonparser.GetString(value, "avatar")
	// 	if newErr != nil {
	// 		log.Println("get avatar error: " + newErr.Error())
	// 	}
	// 	u.UserName, u.Avatar, u.PhoneNumber = username, avatar, phoneNumber
	// 	_, insideErr := jsonparser.ArrayEach(value, func(childValue []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 		roleId, newErr := jsonparser.GetString(childValue, "role_id")
	// 		if newErr != nil {
	// 			log.Println("get role id error:" + newErr.Error())
	// 		}
	// 		roleName, newErr := jsonparser.GetString(childValue, "name")
	// 		if newErr != nil {
	// 			log.Println("get role name error:" + newErr.Error())
	// 		}
	// 		role := Role{
	// 			RoleID:   roleId,
	// 			RoleName: roleName,
	// 		}
	// 		u.Roles = append(u.Roles, role)
	// 	}, "roles")
	// 	if insideErr != nil {
	// 		log.Println("parse role error:", insideErr.Error())
	// 	}
	// }, "userInfo")

	// if err != nil {
	// 	log.Println("array each error: " + err.Error())
	// 	return err
	// }
	return nil
}

func (u *User) VerifyPwd() (result bool, err error) {
	query := fmt.Sprintf(`{
		users(func: eq(username, "%s")) {
			uid
			password
			roles {
				role_id
			}
		}
	}
	`, u.UserName)
	resp, err := Query(context.Background(), query)
	if err != nil {
		log.Println("query error: " + err.Error())
		return false, err
	}
	// fmt.Printf("resp: %+v\n", resp)
	users := UsersStru{}
	err = json.Unmarshal(resp.Json, &users)
	if err != nil {
		log.Println("parse users json error:" + err.Error())
		//return err
	}
	if len(users.Users) == 0 {
		err = errors.New("user not found!")
		return false, err
	}
	user := users.Users[0]

	if u.Password == utils.GetMd5(user.Password) {
		result = true
	} else {
		result = false
	}

	u.UID, u.Roles = user.UID, user.Roles

	// _, err = jsonparser.ArrayEach(resp.Json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 	// theResult, newErr := jsonparser.GetBoolean(value, "checkpwd(pwd)")
	// 	// if newErr != nil {
	// 	// 	log.Println("get checkpwd error: " + newErr.Error())
	// 	// }
	// 	var theResult bool
	// 	pwd, newErr := jsonparser.GetString(value, "password")
	// 	if newErr != nil {
	// 		log.Println("get password error: " + newErr.Error())
	// 	}
	// 	md5Pwd := utils.GetMd5(pwd)
	// 	if u.Password == md5Pwd {
	// 		theResult = true
	// 	} else {
	// 		theResult = false
	// 	}

	// 	uid, newErr := jsonparser.GetString(value, "uid")
	// 	if newErr != nil {
	// 		log.Println("get uid error: " + newErr.Error())
	// 	}

	// 	_, insideErr := jsonparser.ArrayEach(value, func(childValue []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 		roleId, newErr := jsonparser.GetString(childValue, "role_id")
	// 		if newErr != nil {
	// 			log.Println("get role id error:" + newErr.Error())
	// 		}
	// 		role := Role{
	// 			RoleID: roleId,
	// 		}
	// 		u.Roles = append(u.Roles, role)
	// 	}, "roles")
	// 	if insideErr != nil {
	// 		log.Println("verify pwd parse role error:", insideErr.Error())
	// 	}

	// 	result = theResult
	// 	u.UID = uid

	// }, "verify")
	// if err != nil {
	// 	log.Println("get result error: " + err.Error())
	// 	return false, err
	// }
	// fmt.Println("result: ", result)
	return
}

// ============below is not method but is outside function

func GetUserList(pageSize int, pageNo int, username string, fuzz bool) (users UsersStru, err error) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	var query = ""
	if username == "" {
		query = fmt.Sprintf(`{
			users(func: type(User), first:%d, offset:%d) {
				uid
				username
				phone
				update_time
				roles{
					uid
					role_id
					name
				}
			}
		}
		`, pageSize, pageSize*(pageNo-1))
	} else {
		if fuzz == true {
			query = fmt.Sprintf(`{
				users(func: type(User), first:%d, offset:%d) @filter(regexp(username, /%s/)) {
					uid
					username
					phone
					update_time
					roles{
						uid
						role_id
						name
					}
				}
			}`, pageSize, pageSize*(pageNo-1), username)
		} else {
			query = fmt.Sprintf(`{
				users(func: type(User), first:%d, offset:%d) @filter(anyofterms(username, %s)) {
					uid
					username
					phone
					update_time
					roles{
						uid
						role_id
						name
					}
				}
			}`, pageSize, pageSize*(pageNo-1), username)
		}
	}

	users, err = getUsers(query)
	for i, _ := range users.Users {
		// todo: email暂时写死
		users.Users[i].Email = "xxxx@qq.com"
		permissions := []string{}
		for _, role := range users.Users[i].Roles {
			permissions = append(permissions, role.UID)
		}
		users.Users[i].Permissions = permissions
		users.Users[i].UpdateTime = utils.ChangeTimeFormat("dql2normal", users.Users[i].UpdateTime)
	}
	return
}

func UpdateUser(user User) error {
	var ctx = context.Background()
	// mutationSet := ""
	// mutationDelete := ""
	// uid, err := jsonparser.GetString(updateData, "update_data", "uid")
	// if err != nil {
	// 	log.Println("update user get uid error:" + err.Error())
	// 	return err
	// }
	// err = jsonparser.ObjectEach(updateData, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
	// 	switch string(key) {
	// 	case "username":
	// 		mutationSet += "<" + uid + "> <username> \"" + string(value) + "\" .\n"
	// 	case "password":
	// 		mutationSet += "<" + uid + "> <password> \"" + string(value) + "\" .\n"
	// 	case "permissions":
	// 		// permissions是由两个元素"add"和"delete"组成的对象，而每个元素又是对应的要添加或删除的permission的role的uid数组
	// 		_, newErr := jsonparser.ArrayEach(value, func(perValue []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 			mutationSet += "<" + uid + "> <roles> <" + string(perValue) + "> .\n"
	// 		}, "add")
	// 		if newErr != nil {
	// 			log.Println("parse permissions add error:" + newErr.Error())
	// 			return newErr
	// 		}
	// 		_, newErr2 := jsonparser.ArrayEach(value, func(perValue []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 			mutationDelete += "<" + uid + "> <roles> <" + string(perValue) + "> .\n"
	// 		}, "delete")
	// 		if newErr2 != nil {
	// 			log.Println("parse permissions add error:" + newErr2.Error())
	// 			return newErr2
	// 		}
	// 	case "uid":
	// 		return nil
	// 	default:
	// 		log.Println("update user no key name:" + string(key))
	// 	}
	// 	return nil
	// }, "update_data")
	// if err != nil {
	// 	log.Println("update user parse json error:" + err.Error())
	// 	return err
	// }
	// nowTime := utils.ChangeTimeFormat("normal2dql", utils.GetTimeString("date_and_time"))
	// mutationSet += "<" + uid + "> <update_time> \"" + nowTime + "\" .\n"
	// //fmt.Println("mutation:", mutationSet, "delete:", mutationDelete)

	// // 接下来执行数据库操作
	// msArray, mdArray := []string{}, []string{}
	// if mutationSet != "" {
	// 	msArray = append(msArray, mutationSet)
	// }
	// if mutationDelete != "" {
	// 	mdArray = append(mdArray, mutationDelete)
	// }

	mutationDelete := ""
	mutationSet := ""

	if user.UserName != "" {
		mutationSet += "<" + user.UID + "> <username> \"" + user.UserName + "\" .\n"
	}
	if user.Password != "" {
		mutationSet += "< " + user.UID + "> <password> \"" + user.Password + "\" .\n"
	}
	if len(user.Permissions) != 0 {
		// 这里的每个permission就是每个role的uid
		mutationDelete += "<" + user.UID + "> <roles> * .\n"
		for _, permission := range user.Permissions {
			mutationSet += "<" + user.UID + "> <roles> <" + permission + "> .\n"
		}
	}

	msArray, mdArray := []string{}, []string{}
	if mutationDelete != "" {
		mdArray = append(mdArray, mutationDelete)
	}
	if mutationSet != "" {
		msArray = append(msArray, mutationSet)
	}

	resp, err := MutationDeleteAndSetWithUpsert(ctx, mdArray, msArray, "") //先delete后set
	if err != nil {
		log.Println("update user mutation set error:" + err.Error())
		return err
	}
	// fmt.Println("resp:", len(resp.Json)) // if success len(resp.Json) is 0
	if len(resp.Json) != 0 {
		err = errors.New("some error happen")
		log.Println("some error happen, len(resp.Json) is not 0, resp.Json:" + string(resp.Json))
		return err
	}
	return nil
}

func GetSystemConfig() (sc SystemConfigStru, err error) {
	ctx := context.Background()
	query := fmt.Sprintf(`{
		system_configs(func: uid(%s)) {
			uid
			permission_version
		}
	}`, config.SystemConfigNodeId)
	resp, err := Query(ctx, query)
	if err != nil {
		log.Println("query system config error:" + err.Error())
		return
	}
	err = json.Unmarshal(resp.Json, &sc)
	if err != nil {
		log.Println("parse sc error:" + err.Error())
		return
	}
	return
}

// ===========below is some useful function

func getUsers(query string) (users UsersStru, err error) {
	ctx := context.Background()
	resp, err := Query(ctx, query)
	// fmt.Printf("user is:%+v\n", users)
	if err != nil {
		log.Println("query users error: " + err.Error())
		return
	}
	//fmt.Println("resp:", string(resp.Json))
	err = json.Unmarshal(resp.Json, &users)
	if err != nil {
		log.Println("parse users json error:" + err.Error())
		return
	}
	// if len(users.Users) == 0 {
	// 	err = errors.New("user not found!")
	// 	log.Println("no user found, query:" + query)
	// 	return
	// }
	return
}
