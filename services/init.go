package services

import "dql_admin_backend/model"

type RegistForm struct {
	model.User
	PhoneCode string `form:"phoneCode"` // 验证码
}

type Pagination struct {
	PageSize int `json:"pageSize"`
	PageNo   int `json:"pageNo"`
}
