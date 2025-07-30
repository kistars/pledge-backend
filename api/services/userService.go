package services

import (
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models/request"
	"github.com/kistars/pledge-backend/api/models/response"
	"github.com/kistars/pledge-backend/config"
	"github.com/kistars/pledge-backend/db"
	"github.com/kistars/pledge-backend/log"
	"github.com/kistars/pledge-backend/utils"
)

type UserService struct {
}

func NewUser() *UserService {
	return new(UserService)
}

func (u *UserService) Login(req *request.Login, result *response.Login) int {
	log.Logger.Sugar().Info("contractService", req)
	if req.Name == "admin" && req.Password == "password" {
		token, err := utils.CreateToken(req.Name) // 创建jwt token
		if err != nil {
			log.Logger.Error("CreateToken" + err.Error())
			return statecode.CommonErrServerErr
		}

		result.TokenId = token
		_ = db.RedisSet(req.Name, "login_ok", config.Config.Jwt.ExpireTime)
		return statecode.CommonSuccess
	} else {
		return statecode.NameOrPasswordErr
	}
}
