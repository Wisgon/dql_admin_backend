package config

type DgraphInfo struct {
	URL string `json:"url"`
}

var DB = DgraphInfo{
	URL: "localhost:9080",
}

var JwtSignString = "Rf9REe9feFe98ReY"

var STATUS = map[string]int{
	"OK":            0,
	"NotFound":      1,
	"AuthForbidden": 2,
	"InternalError": 3,
	"ParseError":    4,
	"InvalidParam":  5,
}

// 除这些之外，因为不同环境这些设置不同，所以放置到了local_config.go里并gitignore了
// local_config里面必须包含下面这些变量设置：
/*


var NormalRoleId = "0x4e25" // 普通用户的角色id，创建用户时的默认选项
var SystemConfigNodeId = "0x2222"  //
var NormalRoleId = "0x22221"

根目录
var Root = "/home/zhilong/Documents/my_projects/dql_admin_backend"
*/
