package middleware

import (
	"fmt"
	"net/http"
	"wentee/blog/app/schema/apperror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HandlerFuncWithError func(c *gin.Context) error

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Run This Before

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case apperror.AppError:
				c.JSON(e.Status, gin.H{"Code": e.Code, "Message": e.GetMessage()})
			case validator.ValidationErrors:
				fieldErrs := make([]string, len(e))

				for i, fieldErr := range e {
					fieldErrs[i] = fmt.Sprintf("%v", fieldErr.Error())
				}

				c.JSON(http.StatusUnprocessableEntity, gin.H{"Message": fieldErrs})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"Message": "Internal Server Error"})
			}
		}
	}
}
