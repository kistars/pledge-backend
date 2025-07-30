package services

import (
	"encoding/json"

	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models"
	"github.com/kistars/pledge-backend/api/models/request"
	"github.com/kistars/pledge-backend/api/models/response"
)

type MultiSignService struct {
}

func NewMultiSignService() *MultiSignService {
	return new(MultiSignService)
}

func (m *MultiSignService) SetMultiSign(multiSign *request.SetMultiSign) (int, error) {
	//db set
	err := models.NewMultiSign().Set(multiSign)
	if err != nil {
		return statecode.CommonErrServerErr, err
	}
	return statecode.CommonSuccess, nil
}

func (m *MultiSignService) GetMultiSign(multiSign *response.MultiSign, chainId int) (int, error) {
	multiSignModel := models.NewMultiSign()
	err := multiSignModel.Get(chainId)
	if err != nil {
		return statecode.CommonErrServerErr, err
	}

	var multiSignAccount []string
	_ = json.Unmarshal([]byte(multiSignModel.MultiSignAccount), &multiSignAccount)

	multiSign.SpName = multiSignModel.SpName
	multiSign.SpToken = multiSignModel.SpToken
	multiSign.JpName = multiSignModel.JpName
	multiSign.JpToken = multiSignModel.JpToken
	multiSign.SpAddress = multiSignModel.SpAddress
	multiSign.JpAddress = multiSignModel.JpAddress
	multiSign.SpHash = multiSignModel.SpHash
	multiSign.JpHash = multiSignModel.JpHash
	multiSign.MultiSignAccount = multiSignAccount
	return statecode.CommonSuccess, nil
}
