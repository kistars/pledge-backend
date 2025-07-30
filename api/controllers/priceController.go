package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kistars/pledge-backend/api/models/ws"
	"github.com/kistars/pledge-backend/log"
	"github.com/kistars/pledge-backend/utils"
)

type PriceController struct{}

func (p *PriceController) NewPrice(ctx *gin.Context) {
	defer func() {
		recoverRes := recover()
		if recoverRes != nil {
			log.Logger.Sugar().Error("new price recover", recoverRes)
		}
	}()

	// upgrade to websocket
	upgrader := &websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Logger.Sugar().Error("websocket request error", err)
	}

	randomId := ""
	remoteIP := ctx.RemoteIP()
	if remoteIP != "" {
		randomId = strings.Replace(remoteIP, ".", "_", -1) + "_" + utils.GetRandomString(23)
	} else {
		randomId = utils.GetRandomString(32)
	}

	server := &ws.Server{
		Id:       randomId,
		Socket:   conn,
		Send:     make(chan []byte, 800),
		LastTime: time.Now().Unix(),
	}

	go server.ReadAndWrite()
}
