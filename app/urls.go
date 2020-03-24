package app

import (
	"github.com/rezwanul-haque/ID-Service/controllers/ping"
	"github.com/rezwanul-haque/Metadata-Service/controllers/users"
	"github.com/rezwanul-haque/Metadata-Service/utils/consts"
)

const (
	APIBase = "api/" + consts.APIVersion
)

func mapUrls() {
	router.GET(APIBase+"/ping", ping.Ping)

	router.POST(APIBase+"/users/meta", users.Create)
	router.PATCH(APIBase+"/users/:user_id/meta", users.Update)
	router.GET(APIBase+"/users", users.GetAllUsers)
	router.POST(APIBase + "/users/meta/search")
}