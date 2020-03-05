package database

import (
	"app/constant"
	"fmt"
	"os"
	"strconv"
	"time"

	"database/sql"
	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db  *sql.DB
	err error
)

func init() {
	var (
		dbHost     = constant.ServerConfig.Get(fmt.Sprint("MYSQL_RW_HOST"))
		dbPort, _  = strconv.Atoi(constant.ServerConfig.Get("MYSQL_RW_PORT"))
		dbUser     = constant.ServerConfig.Get("MYSQL_RW_USERNAME")
		dbPwd      = constant.ServerConfig.Get("MYSQL_RW_PASSWORD")
	)
	dataBaseName := constant.ServerConfig.Get("MYSQL_DATABASE")
	Db, err = manager.
		New(dataBaseName, dbUser, os.Getenv(dbPwd), os.Getenv(dbHost)).
		Set(
			manager.SetCharset("utf8"),
			manager.SetAllowCleartextPasswords(true),
			manager.SetInterpolateParams(true),
			manager.SetTimeout(1 * time.Second),
			manager.SetReadTimeout(1 * time.Second),
		).
		Port(dbPort).
		Open(true)
}