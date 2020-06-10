package app

import (
	"github.com/rezwanul-haque/Metadata-Service/src/controllers/ping"
	"github.com/rezwanul-haque/Metadata-Service/src/controllers/users"
	"github.com/rezwanul-haque/Metadata-Service/src/utils/consts"
)

const (
	APIBase = "api/" + consts.APIVersion
)

func mapUrls() {
	router.GET(APIBase+"/ping", ping.Ping)

	router.POST(APIBase+"/users/meta", users.Create)
	router.PATCH(APIBase+"/users/:user_id/meta", users.Update)
	router.GET(APIBase+"/users", users.Get)
	router.POST(APIBase+"/users/meta/search", users.Search)
}
