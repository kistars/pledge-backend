package validate

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models/request"
)

type Search struct{}

func NewSearch() *Search {
	return &Search{}
}

func (s *Search) Search(c *gin.Context, req *request.Search) int {

	err := c.ShouldBindJSON(req)
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			if e.Field() == "ChainID" && e.Tag() == "required" {
				return statecode.ChainIdEmpty
			}
		}
		return statecode.CommonErrServerErr
	}

	if req.ChainID != 97 && req.ChainID != 56 {
		return statecode.ChainIdErr
	}

	return statecode.CommonSuccess
}
