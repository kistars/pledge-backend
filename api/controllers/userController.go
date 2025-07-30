package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models/request"
	"github.com/kistars/pledge-backend/api/models/response"
	"github.com/kistars/pledge-backend/api/services"
	"github.com/kistars/pledge-backend/api/validate"
	"github.com/kistars/pledge-backend/db"
)

type UserController struct{}

func (c *UserController) Login(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.Login{}
	result := response.Login{}

	errCode := validate.NewUser().Login(ctx, &req) // 检查请求参数
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewUser().Login(&req, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
}

func (c *UserController) Logout(ctx *gin.Context) {
	res := response.Gin{Res: ctx}

	usernameIntf, _ := ctx.Get("username")

	//delete username in redis
	_, _ = db.RedisDelete(usernameIntf.(string))

	res.Response(ctx, statecode.CommonSuccess, nil)
}
