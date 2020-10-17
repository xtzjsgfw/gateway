package middleware

import (
	"gateway/extend/code"
	"gateway/extend/jwt"
	"gateway/extend/redis"
	"gateway/extend/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			utils.ResponseFormat(c, code.TokenIsNotExistError, nil)
			c.Abort()
			return
		}

		token = token[7:]
		jwtInstance := jwt.NewJWT()
		claims, err := jwtInstance.ParseToken(token)
		if err != nil {
			utils.ResponseFormat(c, code.TokenParseError, nil)
			c.Abort()
			return
		}
		rdb := redis.GetRedis()
		tokenCache, err := rdb.Get("TOKEN:" + claims.Mobile).Result()
		if err != nil {
			utils.ResponseFormat(c, code.TokenIsNotExistError, nil)
			c.Abort()
			return
		}

		if token != tokenCache {
			utils.ResponseFormat(c, code.TokenInvalid, nil)
			c.Abort()
			return
		}
	}
}
