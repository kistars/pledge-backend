package controllers

import (
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models"
	"github.com/kistars/pledge-backend/api/models/request"
	"github.com/kistars/pledge-backend/api/models/response"
	"github.com/kistars/pledge-backend/api/services"
	"github.com/kistars/pledge-backend/api/validate"
	"github.com/kistars/pledge-backend/config"
)

type PoolController struct{}

func (p *PoolController) PoolBaseInfo(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.PoolBaseInfo{}
	var result []models.PoolBaseInfoRes

	errCode := validate.NewPoolBaseInfo().PoolBaseInfo(ctx, &req) // 验证请求参数
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewPoolService().PoolBaseInfo(req.ChainId, &result) // read data from db
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
}

func (p *PoolController) PoolDataInfo(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.PoolDataInfo{}
	var result []models.PoolDataInfoRes

	errCode := validate.NewPoolDataInfo().PoolDataInfo(ctx, &req) // 验证请求参数
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewPoolService().PoolDataInfo(req.ChainId, &result) // read data from db
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
}

func (p *PoolController) TokenList(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.TokenList{}
	result := response.TokenList{}

	errCode := validate.NewTokenList().TokenList(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, map[string]string{
			"error": "chainId error",
		})
		return
	}

	errCode, data := services.NewTokenList().GetTokenList(&req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, map[string]string{
			"error": "chainId error",
		})
		return
	}

	var BaseUrl = p.GetBaseUrl()
	result.Name = "Pledge Token List"
	result.LogoURI = BaseUrl + "storage/img/Pledge-project-logo.png"
	result.Timestamp = time.Now()
	result.Version = response.Version{
		Major: 2,
		Minor: 16,
		Patch: 12,
	}
	for _, v := range data {
		result.Tokens = append(result.Tokens, response.Token{
			Name:     v.Symbol,
			Symbol:   v.Symbol,
			Decimals: v.Decimals,
			Address:  v.Token,
			ChainID:  v.ChainId,
			LogoURI:  v.Logo,
		})
	}

	res.Response(ctx, statecode.CommonSuccess, result)
}

func (p *PoolController) Search(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.Search{}
	result := response.Search{}

	errCode := validate.NewSearch().Search(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, count, pools := services.NewSearch().Search(&req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	result.Rows = pools
	result.Count = count
	res.Response(ctx, statecode.CommonSuccess, result)
}

func (p *PoolController) DebtTokenList(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.TokenList{}

	errCode := validate.NewTokenList().TokenList(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, result := services.NewTokenList().DebtTokenList(&req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
}

func (p *PoolController) GetBaseUrl() string {
	domainName := config.Config.Env.DomainName
	domainNameSlice := strings.Split(domainName, "")
	pattern := "\\d+"
	isNumber, _ := regexp.MatchString(pattern, domainNameSlice[0])
	if isNumber {
		return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + ":" + config.Config.Env.Port + "/"
	}
	return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + "/"
}
