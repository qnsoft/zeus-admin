package middleware

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
	"zeus/pkg/api/dto"
	"zeus/pkg/api/log"
	"zeus/pkg/api/model"
)

//todo : 用单独的claims model去掉user model
func JwtAuth() *jwt.GinJWTMiddleware {
					jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
					Realm:            "zeus jwt",
					SigningAlgorithm: "RS256",
					PubKeyFile:       viper.GetString("jwt.key.public"),
					PrivKeyFile:      viper.GetString("jwt.key.private"),
					Timeout:          time.Hour * 24,
					MaxRefresh:       time.Hour * 24 * 90,
					IdentityKey:      "id",
					PayloadFunc: func(data interface{}) jwt.MapClaims {
					if v, ok := data.(*model.User); ok {
					return jwt.MapClaims{
					"id":   v.Id,
					"name": v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				UserName: claims["name"].(string),
				Id:       int(claims["id"].(float64)),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginDto dto.LoginDto
			if err := c.ShouldBind(&loginDto); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginDto.Username
			password := loginDto.Password

			//todo : 实现登陆验证
			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &model.User{
					Id:       250,
					UserName: loginDto.Username,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*model.User); ok && v.UserName == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Error(err.Error())
	}
	return jwtMiddleware
}
