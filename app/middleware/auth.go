package middleware

import (
	"net/http"
	"wentee/blog/app/config"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	AuthSchema "wentee/blog/app/schema/auth"
	"wentee/blog/app/utils/reqcontext"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	// userCollection *mongo.Collection
}

func (authMiddleware *AuthMiddleware) RequiredAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.Error(apperror.New(http.StatusForbidden, errcode.INVALID_TOKEN, nil))
			c.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &AuthSchema.JWTClaims{}, func(token *jwt.Token) (any, error) {

			return []byte(config.JWT_SECRET), nil
		})

		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		claim, ok := token.Claims.(*AuthSchema.JWTClaims)
		if !ok {
			c.Error(apperror.New(http.StatusForbidden, errcode.INVALID_TOKEN, nil))
			c.Abort()
			return
		}

		c.Set(reqcontext.USER_INFO, claim.UserInfo)
		c.Next()
	}
}
