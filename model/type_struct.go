//这是schema的type的结构体，方便json操作
package model

type User struct {
	UID         string   `json:"uid"`
	UserName    string   `json:"username"`
	Password    string   `json:"password"`
	PhoneNumber string   `json:"phone"`
	Avatar      string   `json:"avatar"` //头像图片地址
	Roles       []Role   `json:"roles"`  // 这里不能用Roles RolesStru，json解析时会出错
	CreateTime  string   `json:"create_time"`
	UpdateTime  string   `json:"update_time"`
	Email       string   `json:"email,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type Role struct {
	UID             string   `json:"uid"`
	RoleID          string   `json:"role_id"`
	RoleName        string   `json:"name"`
	AccessablePages []string `json:"accessable_pages"`
}

type SystemConfig struct {
	UID               string `json:"uid"`
	PermissionVersion string `json:"permission_version"`
}

type UsersStru struct {
	Users []User `json:"users"`
}

type RolesStru struct {
	Roles []Role `json:"roles"`
}

// 虽然SystemConfig只有一个节点，但是query的时候还是拿到的数组
type SystemConfigStru struct {
	SystemConfigs []SystemConfig `json:"system_configs"`
}
