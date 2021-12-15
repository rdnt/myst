package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/app/client/rest/generated"
)

func Error(c *gin.Context, code int, message interface{}) {
	msg := "unknown error"

	switch message.(type) {
	case string:
		msg = message.(string)
	case fmt.Stringer:
		msg = message.(fmt.Stringer).String()
	case error:
		msg = message.(error).Error()
	}

	c.JSON(
		code, generated.Error{
			Code:    http.StatusText(code),
			Message: msg,
		},
	)
}
