package consts

const (
	APIVersion string = "v1"
)

const (
	SecretKey string = "dHepoR82!e#"
)

type RoleDefine struct {
	ADMIN string
	USER  string
}

// Role
var Role = RoleDefine{"admin", "user"}
