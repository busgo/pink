package handler

import (
	"github.com/busgo/pink/http/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	OK                = 0
	BusinessErrorCode = 100000
	ParamErrorCode    = 100001
)

func WriteOK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, &model.Result{
		Code: OK,
		Data: data,
	})
}

func WriteErrorWithCode(c echo.Context, code int32, message string) error {
	return c.JSON(http.StatusOK, &model.Result{
		Code:    code,
		Message: message,
	})
}

func WriteParamError(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, &model.Result{
		Code:    ParamErrorCode,
		Message: message,
	})
}

func WriteBusinessError(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, &model.Result{
		Code:    BusinessErrorCode,
		Message: message,
	})
}
