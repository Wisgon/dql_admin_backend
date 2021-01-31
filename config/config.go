package config

type DgraphInfo struct {
	URL string `json:"url"`
}

var DB = DgraphInfo{
	URL: "localhost:9080",
}

var JwtSignString = "Rf9REe9feFe98ReY"

var Root = "E:/Documents/my_projects/dql_admin_backend"
