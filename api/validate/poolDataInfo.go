package validate

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models/request"
)

type PoolDataInfo struct{}

func NewPoolDataInfo() *PoolDataInfo {
	return &PoolDataInfo{}
}

func (v *PoolDataInfo) PoolDataInfo(c *gin.Context, req *request.PoolDataInfo) int {
	err := c.ShouldBind(req)
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			if e.Field() == "ChainId" && e.Tag() == "required" {
				return statecode.ChainIdEmpty
			}
		}
		return statecode.CommonErrServerErr
	}

	if req.ChainId != 97 && req.ChainId != 56 {
		return statecode.ChainIdErr
	}

	return statecode.CommonSuccess
}
