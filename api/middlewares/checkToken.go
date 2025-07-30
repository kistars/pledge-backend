package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/kistars/pledge-backend/api/common/statecode"
	"github.com/kistars/pledge-backend/api/models/response"
	"github.com/kistars/pledge-backend/config"
	"github.com/kistars/pledge-backend/db"
	"github.com/kistars/pledge-backend/utils"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := response.Gin{Res: c}
		token := c.Request.Header.Get("authCode")

		username, err := utils.ParseToken(token, config.Config.Jwt.SecretKey)
		if err != nil {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}

		if username != config.Config.DefaultAdmin.Username {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}

		// Judge whether the user logout
		resByteArr, err := db.RedisGet(username)
		if string(resByteArr) != `"login_ok"` {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}

		c.Set("username", username)

		c.Next()
	}
}
