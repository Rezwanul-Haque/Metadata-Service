package lms_db

import (
	"database/sql"
	"fmt"
	"github.com/rezwanul-haque/Metadata-Service/src/utils/helpers"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Client *sql.DB

	username = helpers.GoDotEnvVariable("MYSQL_LMS_USERNAME")
	password = helpers.GoDotEnvVariable("MYSQL_LMS_PASSWORD")
	host     = helpers.GoDotEnvVariable("MYSQL_LMS_HOST")
	port     = helpers.GoDotEnvVariable("MYSQL_LMS_PORT")
	schema   = helpers.GoDotEnvVariable("MYSQL_LMS_SCHEMA")
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		username, password, host, port, schema,
	)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("database successfully configured")
}
