package model

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
