package config

import (
	"fmt"
	"github.com/kreditplus/scorepro/config/credential"
	"github.com/kreditplus/scorepro/utils/monitoring"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"os"
	"strconv"
)

func NewDbMssql() *gorm.DB {

	var (
		dbHost   = credential.DbHost
		port     = credential.DbPort
		user     = credential.DbUsername
		password = credential.DbPassword
		database = credential.DbName
	)

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		user, password, dbHost, port, database)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		extra := map[string]interface{}{
			"message":err.Error(),
		}
		monitoring.SendToSentry(nil, extra, "DATABASE")
		panic(err)
	}

	idle, _ := strconv.Atoi(os.Getenv("SET_MAX_IDLE_CONN"))
	open, _ := strconv.Atoi(os.Getenv("SET_MAX_OPEN_CONN"))

	pool, err := db.DB()
	pool.SetMaxIdleConns(idle)
	pool.SetMaxOpenConns(open)

	return db.Debug()
}

func NewDbMssqlScorepro() *gorm.DB {

	var (
		dbHost   = credential.DbHost
		port     = credential.DbPort
		user     = credential.DbUsername
		password = credential.DbPassword
		database = credential.DbNameScorepro
	)

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		user, password, dbHost, port, database)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		tags := map[string]string{
			"app.pkg":"main",
			"app.func":"main",
		}
		extra := map[string]interface{}{
			"message":err.Error(),
		}
		monitoring.SendToSentry(tags, extra, "DATABASE")
		panic(err)
	}

	idle, _ := strconv.Atoi(os.Getenv("SET_MAX_IDLE_CONN"))
	open, _ := strconv.Atoi(os.Getenv("SET_MAX_OPEN_CONN"))

	pool, err := db.DB()
	pool.SetMaxIdleConns(idle)
	pool.SetMaxOpenConns(open)

	return db.Debug()
}
