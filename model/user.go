package model

import (
	"context"
	"dql_admin_backend/utils"
	"errors"
	"fmt"
	"log"
)

type User struct {
	ID          string
	UserName    string `json:"username"`
	Password    string `json:"password"`
	Port        int
	PhoneNumber string `json:"phone"`
	Email       string
	MemberSN    string //会员编号
	Avatar      string //头像图片地址
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
	// 组装mutation
	mu := fmt.Sprintf(`
	_:new_user <user_name> "%s" .
	_:new_user <pwd> "%s" .
	_:new_user <phone_number> "%s" .
	_:new_user <avatar> "%s" .
	_:new_user <dgraph.type> "User" .
	`, u.UserName, u.Password, u.PhoneNumber, u.Avatar)
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
		log.Println("parse resp json error: " + err.Error())
		return err
	}
	q2_count, err := utils.CountJsonArray(resp.Json, "q2")
	if err != nil {
		log.Println("parse resp json error: " + err.Error())
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

func (u User) UpgradeMember() error {
	// todo: 升级会员，需要哈希一个会员编号和分配一个port，分配port是难点，用conditional upsert
	return nil
}

func (u *User) GetUserInfo(condition string) error {
	return nil
}

func (u *User) VerifyPwd() (result bool, err error) {
	// query := fmt.Sprintf(`{
	// 	verify(func: eq(user_name, "%s")) {
	// 		uid
	// 		checkpwd(pwd, "%s")
	// 	}
	// }
	// `, u.UserName, u.Password)
	// resp, err := service.Query(context.Background(), query)
	// if err != nil {
	// 	log.Println("query error: " + err.Error())
	// 	return false, err
	// }
	// //fmt.Printf("resp: %+v\n", resp)

	// _, err = jsonparser.ArrayEach(resp.Json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 	theResult, newErr := jsonparser.GetBoolean(value, "checkpwd(pwd)")
	// 	if newErr != nil {
	// 		log.Println("get checkpwd error: " + newErr.Error())
	// 	}
	// 	uid, newErr := jsonparser.GetString(value, "uid")
	// 	if newErr != nil {
	// 		log.Println("get uid error: " + newErr.Error())
	// 	}
	// 	result = theResult
	// 	u.ID = uid
	// }, "verify")
	// if err != nil {
	// 	log.Println("get result error: " + err.Error())
	// 	return false, err
	// }
	// // fmt.Println("result: ", result)
	return
}
