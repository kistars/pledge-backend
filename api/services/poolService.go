package services

import (
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models"
	"github.com/kistars/pledge-backend/log"
)

type PoolService struct{}

func NewPoolService() *PoolService {
	return new(PoolService)
}

func (p *PoolService) PoolBaseInfo(chainId int, result *[]models.PoolBaseInfoRes) int {
	err := models.NewPoolBases().PoolBaseInfo(chainId, result) // read data from db
	if err != nil {
		log.Logger.Error(err.Error())
		return statecode.CommonErrServerErr
	}
	return statecode.CommonSuccess
}

func (p *PoolService) PoolDataInfo(chainId int, result *[]models.PoolDataInfoRes) int {
	err := models.NewPoolData().PoolDataInfo(chainId, result) // read data from db
	if err != nil {
		log.Logger.Error(err.Error())
		return statecode.CommonErrServerErr
	}
	return statecode.CommonSuccess
}
