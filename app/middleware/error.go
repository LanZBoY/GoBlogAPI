package middleware

import (
	"errors"
	"wentee/blog/app/schema/apperror"

	"github.com/gin-gonic/gin"
)

type HandlerFuncWithError func(c *gin.Context) error

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Run This Before

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appError *apperror.AppError
			if errors.As(err, &appError) {
				c.JSON(appError.Status, gin.H{"Code": appError.Code, "Message": appError.GetMessage()})
				return
			}
		}
	}
}
