package model

import (
	"context"
	"dql_admin_backend/config"
	"dql_admin_backend/utils"
	"errors"
	"fmt"
	"log"

	"github.com/buger/jsonparser"
)

type User struct {
	ID          string
	UserName    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone"`
	Avatar      string `json:"avatar"` //头像图片地址
	Roles       []Role `json:"roles"`
}

type Role struct {
	RoleID   string `json:"role_id"`
	RoleName string `json:"role_name"`
}

func (u User) CreateUser() error {
	var ctx = context.Background()
	query := fmt.Sprintf(`
	query{
		q1(func: eq(user_name, "%s")) {
			un as uid
		}
		q2(func: eq(phone_number, "%s")) {
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
	_:new_user <user_name> "%s" .
	_:new_user <pwd> "%s" .
	_:new_user <phone_number> "%s" .
	_:new_user <avatar> "%s" .
	_:new_user <dgraph.type> "User" .
	_:new_user <create_time> "%s" .
	_:new_user <role> <%s> .
	`, u.UserName, u.Password, u.PhoneNumber, u.Avatar, nowTime, config.NormalRoleId)
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
			userInfo(func: uid(%s)) {
				user_name
				phone_number
				avatar
				role{
					role_id
					name
				}
			}
		}`, u.ID)
	case "user_name":
		query = fmt.Sprintf(`{
			userInfo(func: eq(user_name, "%s")) {
				user_name
				phone_number
				avatar
				role{
					role_id
					name
				}
			}
		}`, u.UserName)
	case "phone_number":
		query = fmt.Sprintf(`{
			userInfo(func: eq(phone_number, "%s")) {
				user_name
				phone_number
				avatar
				role{
					role_id
					name
				}
			}
		}`, u.PhoneNumber)
	default:
		err := errors.New("no such user attr: " + condition)
		return err
	}

	resp, err := Query(context.Background(), query)
	if err != nil {
		log.Println("query userinfo error: " + err.Error())
		return err
	}
	// fmt.Println("resp:", resp)

	_, err = jsonparser.ArrayEach(resp.Json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		username, newErr := jsonparser.GetString(value, "user_name")
		if newErr != nil {
			log.Println("get username error: " + newErr.Error())
		}
		phoneNumber, newErr := jsonparser.GetString(value, "phone_number")
		if newErr != nil {
			log.Println("get phone_number error: " + newErr.Error())
		}
		avatar, newErr := jsonparser.GetString(value, "avatar")
		if newErr != nil {
			log.Println("get avatar error: " + newErr.Error())
		}
		u.UserName, u.Avatar, u.PhoneNumber = username, avatar, phoneNumber
		_, insideErr := jsonparser.ArrayEach(value, func(childValue []byte, dataType jsonparser.ValueType, offset int, err error) {
			roleId, newErr := jsonparser.GetString(childValue, "role_id")
			if newErr != nil {
				log.Println("get role id error:" + newErr.Error())
			}
			roleName, newErr := jsonparser.GetString(childValue, "name")
			if newErr != nil {
				log.Println("get role name error:" + newErr.Error())
			}
			role := Role{
				RoleID:   roleId,
				RoleName: roleName,
			}
			u.Roles = append(u.Roles, role)
		}, "role")
		if insideErr != nil {
			log.Println("parse role error:", insideErr.Error())
		}
	}, "userInfo")

	if err != nil {
		log.Println("array each error: " + err.Error())
		return err
	}
	return nil
}

func (u *User) VerifyPwd() (result bool, err error) {
	query := fmt.Sprintf(`{
		verify(func: eq(user_name, "%s")) {
			uid
			pwd
			role {
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

	_, err = jsonparser.ArrayEach(resp.Json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// theResult, newErr := jsonparser.GetBoolean(value, "checkpwd(pwd)")
		// if newErr != nil {
		// 	log.Println("get checkpwd error: " + newErr.Error())
		// }
		var theResult bool
		pwd, newErr := jsonparser.GetString(value, "pwd")
		if newErr != nil {
			log.Println("get pwd error: " + newErr.Error())
		}
		md5Pwd := utils.GetMd5(pwd)
		if u.Password == md5Pwd {
			theResult = true
		} else {
			theResult = false
		}

		uid, newErr := jsonparser.GetString(value, "uid")
		if newErr != nil {
			log.Println("get uid error: " + newErr.Error())
		}

		_, insideErr := jsonparser.ArrayEach(value, func(childValue []byte, dataType jsonparser.ValueType, offset int, err error) {
			roleId, newErr := jsonparser.GetString(childValue, "role_id")
			if newErr != nil {
				log.Println("get role id error:" + newErr.Error())
			}
			role := Role{
				RoleID:   roleId,
			}
			u.Roles = append(u.Roles, role)
		}, "role")
		if insideErr != nil {
			log.Println("verify pwd parse role error:", insideErr.Error())
		}

		result = theResult
		u.ID = uid

	}, "verify")
	if err != nil {
		log.Println("get result error: " + err.Error())
		return false, err
	}
	// fmt.Println("result: ", result)
	return
}
