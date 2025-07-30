package services

import (
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models"
	"github.com/kistars/pledge-backend/api/models/request"
)

type TokenList struct{}

func NewTokenList() *TokenList {
	return new(TokenList)
}

func (t *TokenList) DebtTokenList(req *request.TokenList) (int, []models.TokenInfo) {
	res, err := models.NewTokenInfo().GetTokenInfo(req) // read data from db
	if err != nil {
		return statecode.CommonErrServerErr, nil
	}
	return statecode.CommonSuccess, res
}

func (t *TokenList) GetTokenList(req *request.TokenList) (int, []models.TokenList) {
	res, err := models.NewTokenInfo().GetTokenList(req) // read data from db
	if err != nil {
		return statecode.CommonErrServerErr, nil
	}
	return statecode.CommonSuccess, res
}
