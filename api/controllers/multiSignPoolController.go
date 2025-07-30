package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models/request"
	"github.com/kistars/pledge-backend/api/models/response"
	"github.com/kistars/pledge-backend/api/services"
	"github.com/kistars/pledge-backend/api/validate"
	"github.com/kistars/pledge-backend/log"
)

type MultiSignPoolController struct{}

func (m *MultiSignPoolController) SetMultiSign(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.SetMultiSign{}
	log.Logger.Sugar().Info("SetMultiSign req ", req)

	errCode := validate.NewMutiSign().SetMultiSign(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, err := services.NewMultiSignService().SetMultiSign(&req)
	if errCode != statecode.CommonSuccess {
		log.Logger.Error(err.Error())
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, nil)
}

func (m *MultiSignPoolController) GetMultiSign(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.GetMultiSign{}
	result := response.MultiSign{}
	log.Logger.Sugar().Info("GetMultiSign req ", nil)

	errCode := validate.NewMutiSign().GetMultiSign(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, err := services.NewMultiSignService().GetMultiSign(&result, req.ChainId)
	if errCode != statecode.CommonSuccess {
		log.Logger.Error(err.Error())
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
}
