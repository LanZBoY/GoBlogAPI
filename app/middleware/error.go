package middleware

import (
	"errors"
	"net/http"
	"wentee/blog/app/schema/apperror"

	"github.com/gin-gonic/gin"
)

type HandlerFuncWithError func(c *gin.Context) error

func ErrorHandlerWrapper(handlerFunc HandlerFuncWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handlerFunc(c); err != nil {
			var appError *apperror.AppError

			if errors.As(err, &appError) {
				c.JSON(appError.Status, gin.H{"code": appError.Code, "message": appError.Message})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}

	}
}
