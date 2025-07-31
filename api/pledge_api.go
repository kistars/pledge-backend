package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kistars/pledge-backend/api/middlewares"
	"github.com/kistars/pledge-backend/api/models"
	"github.com/kistars/pledge-backend/api/models/kucoin"
	"github.com/kistars/pledge-backend/api/models/ws"
	"github.com/kistars/pledge-backend/api/routes"
	"github.com/kistars/pledge-backend/api/static"
	"github.com/kistars/pledge-backend/api/validate"
	"github.com/kistars/pledge-backend/config"
	"github.com/kistars/pledge-backend/db"
)

func main() {
	// init mysql
	db.InitMysql()

	// init redis
	db.InitRedis()
	models.InitTable()

	// validate
	validate.BindingValidator()

	// websocket server
	go ws.StartServer()

	// get plgr price from kucoin-exchange
	go kucoin.GetExchangePrice()

	// gin start
	gin.SetMode(gin.DebugMode)
	app := gin.Default()
	staticPath := static.GetCurrentAbPathByCaller()
	app.Static("/storage/", staticPath)
	app.Use(middlewares.Cors()) // 「 Cross domain Middleware 」
	routes.InitRoute(app)
	_ = app.Run(":" + config.Config.Env.Port)
}
