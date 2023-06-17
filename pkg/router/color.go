package router

import (
	"net/http"

	"myst/pkg/logger"
)

func StatusColor(status int) logger.Color {
	switch {
	case status >= http.StatusOK && status < http.StatusMultipleChoices:
		return logger.GreenBg | logger.Black
	case status >= http.StatusMultipleChoices && status < http.StatusBadRequest:
		return logger.WhiteBg | logger.Black
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		return logger.YellowBg | logger.Black
	default:
		return logger.RedBg | logger.Black
	}
}

func MethodColor(method string) logger.Color {
	switch method {
	case http.MethodGet:
		return logger.GreenBg | logger.Black
	case http.MethodPost:
		return logger.BlueBg | logger.Black
	case http.MethodPut:
		return logger.CyanBg | logger.Black
	case http.MethodPatch:
		return logger.YellowBg | logger.Black
	case http.MethodDelete:
		return logger.RedBg | logger.Black
	default:
		return logger.MagentaBg | logger.Black
	}
}
