package main

import (
	"github.com/kistars/pledge-backend/db"
	"github.com/kistars/pledge-backend/schedule/models"
)

func main() {
	// init mysql
	db.InitMysql()

	// init redis
	db.InitRedis()

	// create table
	models.InitTable()

	// pool task

}
