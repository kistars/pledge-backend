package services

import (
	"fmt"

	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models"
	"github.com/kistars/pledge-backend/api/models/request"
	"github.com/kistars/pledge-backend/log"
)

type SearchService struct{}

func NewSearch() *SearchService {
	return new(SearchService)
}

func (s *SearchService) Search(req *request.Search) (int, int64, []models.Pool) {
	whereCondition := fmt.Sprintf(`chain_id='%v'`, req.ChainID)
	if req.LendTokenSymbol != "" {
		whereCondition += fmt.Sprintf(` and lend_token_symbol='%v'`, req.LendTokenSymbol)
	}
	if req.State != "" {
		whereCondition += fmt.Sprintf(` and state='%v'`, req.State)
	}

	total, data, err := models.NewPool().Pagination(req, whereCondition) // read data from db
	if err != nil {
		log.Logger.Error(err.Error())
		return statecode.CommonErrServerErr, 0, nil
	}

	return statecode.CommonSuccess, total, data
}
