package validate

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models/request"
)

type User struct{}

func NewUser() *User {
	return new(User)
}

func (u *User) Login(c *gin.Context, req *request.Login) int {
	err := c.ShouldBind(req)

	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			if e.Field() == "Name" && e.Tag() == "required" {
				return statecode.PNameEmpty
			}
			if e.Field() == "Password" && e.Tag() == "required" {
				return statecode.PNameEmpty
			}
		}
		return statecode.CommonErrServerErr
	}

	return statecode.CommonSuccess
}
