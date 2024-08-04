package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type JSendResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, JSendResponse{
		Status: "success",
		Data:   data,
	})
}

func Fail(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusBadRequest, JSendResponse{
		Status: "fail",
		Data:   data,
	})
}

func Error(c echo.Context, message string, code int, data interface{}) error {
	return c.JSON(http.StatusInternalServerError, JSendResponse{
		Status:  "error",
		Message: message,
		Code:    code,
		Data:    data,
	})
}
