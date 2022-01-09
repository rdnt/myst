package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"myst/internal/client/api/http/generated"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, code int, v interface{}) {
	msg := "unknown error"

	switch v.(type) {
	case string:
		msg = v.(string)
	case fmt.Stringer:
		msg = v.(fmt.Stringer).String()
	case error:
		msg = v.(error).Error()
	default:
		b, err := json.Marshal(v)
		if err == nil {
			msg = string(b)
		}
	}

	c.JSON(
		code, generated.Error{
			Code:    http.StatusText(code),
			Message: msg,
		},
	)
}

func Success(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, v)
}
