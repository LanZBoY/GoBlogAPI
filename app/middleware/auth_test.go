package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wentee/blog/app/config"
	"wentee/blog/app/schema/apperror/errcode"
	AuthSchema "wentee/blog/app/schema/auth"
	"wentee/blog/app/utils/reqcontext"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateToken(exp time.Time) string {
	claims := AuthSchema.JWTClaims{
		UserInfo: AuthSchema.JWTUserInfo{Id: "123"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.SERVICE_NAME,
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(config.JWT_SECRET))
	return tokenStr
}

func TestRequiredAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMw := AuthMiddleware{}

	router := gin.New()
	router.Use(ErrorHandler())
	router.Use(authMw.RequiredAuth())

	var gotUser AuthSchema.JWTUserInfo
	router.GET("/", func(c *gin.Context) {
		val, _ := reqcontext.GetUserInfo(c)
		gotUser = val
		c.Status(http.StatusOK)
	})

	t.Run("ValidToken", func(t *testing.T) {
		token := generateToken(time.Now().Add(time.Hour))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "123", gotUser.Id)
	})

	t.Run("MissingToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), errcode.INVALID_TOKEN)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		token := generateToken(time.Now().Add(-time.Hour))
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), errcode.EXPIRED_TOKEN)
	})
}
