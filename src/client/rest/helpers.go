package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/src/client/rest/generated"
)

func Error(c *gin.Context, code int, v interface{}) {
	msg := "unknown error"

	switch v := v.(type) {
	case string:
		msg = v
	case fmt.Stringer:
		msg = v.String()
	case error:
		msg = v.Error()
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
